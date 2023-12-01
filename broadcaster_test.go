package broadcaster

import (
	"context"
	"errors"
	"github.com/stretchr/testify/require"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

func TestBroadcastSeparately(t *testing.T) {

	b := NewBroadcaster[uuid.UUID](100 * time.Millisecond)

	for i := 0; i < 1_000; i++ {
		id := uuid.New()
		require.NoError(t, b.WaitForEffects(
			context.Background(),
			func() {
				b.Broadcast(id)
			},
			OnReceiveEffect[uuid.UUID]{
				checkFunc: func(msg uuid.UUID) (bool, error) {
					return msg == id, nil
				},
				timeoutMsg: "Timed out waiting for message",
			},
		))
	}
}

func TestBroadcastReturnsAllErrors(t *testing.T) {

	b := NewBroadcaster[uuid.UUID](100 * time.Millisecond)

	id1 := uuid.New()
	id2 := uuid.New()

	err := b.WaitForEffects(
		context.Background(),
		func() {
			b.Broadcast(id1, id2)
		},
		OnReceiveEffect[uuid.UUID]{
			checkFunc: func(msg uuid.UUID) (bool, error) {
				return msg == id1, errors.New(id1.String())
			},
			timeoutMsg: "Timed out waiting for message",
		},
		OnReceiveEffect[uuid.UUID]{
			checkFunc: func(msg uuid.UUID) (bool, error) {
				return msg == id2, errors.New(id2.String())
			},
			timeoutMsg: "Timed out waiting for message",
		},
	)
	require.Error(t, err)
	require.Contains(t, err.Error(), id1.String())
	require.Contains(t, err.Error(), id2.String())
}

func TestBroadcastVariadically(t *testing.T) {

	b := NewBroadcaster[uuid.UUID](1000 * time.Millisecond)

	//nolint:prealloc
	var ids []uuid.UUID
	var effects []OnReceiveEffect[uuid.UUID]

	for i := 0; i < 1_000; i++ {
		id := uuid.New()
		ids = append(ids, id)
	}

	for _, id := range ids {
		effects = append(
			effects,
			OnReceiveEffect[uuid.UUID]{
				checkFunc: func(msg uuid.UUID) (bool, error) {
					return msg == id, nil
				},
				timeoutMsg: "Timed out waiting for message",
			},
		)
	}

	require.NoError(
		t,
		b.WaitForEffects(
			context.Background(),
			func() {
				b.Broadcast(ids...)
			},
			effects...,
		),
	)
}

func TestBroadcastWithLocks(t *testing.T) {

	b := NewBroadcaster[uuid.UUID](100 * time.Millisecond)

	var lock sync.RWMutex

	for i := 0; i < 1_000; i++ {
		var id *uuid.UUID

		waitForID := OnReceiveEffect[uuid.UUID]{
			checkFunc: func(msg uuid.UUID) (bool, error) {
				defer lock.RUnlock()
				lock.RLock()
				return msg == *id, nil
			},
			timeoutMsg: "Timed out waiting for message",
		}

		require.NoError(
			t,
			b.WaitForEffects(
				context.Background(),
				func() {
					defer lock.Unlock()
					lock.Lock()
					broadcastID := uuid.New()
					b.Broadcast(broadcastID)
					id = &broadcastID
				},
				waitForID,
			),
		)
	}
}

type Message struct {
	id       uuid.UUID
	delivery amqp.Delivery
}

func TestBroadcastWithLocks2(t *testing.T) {
	for i := 0; i < 1_000; i++ {
		b := NewBroadcaster[Message](100 * time.Millisecond)
		var id *uuid.UUID
		var lock sync.RWMutex

		waitForID1 := OnReceiveEffect[Message]{
			checkFunc: func(msg Message) (bool, error) {
				defer lock.RUnlock()
				lock.RLock()
				return msg.id == *id, nil
			},
			timeoutMsg: "Timed out waiting for message 1",
		}

		waitForID2 := OnReceiveEffect[Message]{
			checkFunc: func(msg Message) (bool, error) {
				defer lock.RUnlock()
				lock.RLock()
				if msg.delivery.Type != "test" {
					return false, nil
				}
				msgID, err := uuid.Parse(msg.delivery.ContentType)
				if err != nil {
					return false, err
				}
				return msgID == *id, nil
			},
			timeoutMsg: "Timed out waiting for message 2",
		}

		require.NoError(
			t,
			b.WaitForEffects(
				context.Background(),
				func() {
					defer lock.Unlock()
					lock.Lock()
					broadcastID := uuid.New()
					b.Broadcast(
						Message{
							id: broadcastID,
						},
						Message{
							delivery: amqp.Delivery{
								Type:        "test",
								ContentType: broadcastID.String(),
							},
						},
					)
					id = &broadcastID
				},
				waitForID1,
				waitForID2,
			),
		)
	}
}
