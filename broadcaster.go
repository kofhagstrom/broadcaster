package broadcaster

import (
	"context"
	"errors"
	"sync"
	"time"
)

type Broadcaster[T any] struct {
	sync.RWMutex
	broadcastChannel chan T
	timeOut          time.Duration
	channelList      List[T]
}

func NewBroadcaster[T any](
	timeOut time.Duration,
) *Broadcaster[T] {
	b := Broadcaster[T]{
		broadcastChannel: make(chan T),
		timeOut:          timeOut,
		channelList:      NewDoublyLinkedList[T](),
	}
	b.Start()
	return &b
}

func (b *Broadcaster[T]) Start() {
	go func() {
		for {
			b.broadcast(<-b.broadcastChannel)
		}
	}()
}

func (b *Broadcaster[T]) Broadcast(msgs ...T) {
	for _, msg := range msgs {
		b.broadcastChannel <- msg
	}
}

func NewOnReceiveEffect[T any](
	checkFunc func(msg T) (bool, error),
	timeoutMsg string,
) OnReceiveEffect[T] {
	return OnReceiveEffect[T]{
		checkFunc:  checkFunc,
		timeoutMsg: timeoutMsg,
	}
}

type OnReceiveEffect[T any] struct {
	checkFunc  func(msg T) (bool, error)
	timeoutMsg string
}

func (o *OnReceiveEffect[T]) Error() error {
	return errors.New(o.timeoutMsg)
}

func (o *OnReceiveEffect[T]) CheckEffect(msg T) (bool, error) {
	return o.checkFunc(msg)
}

func (b *Broadcaster[T]) WaitForEffects(
	ctx context.Context,
	f func(),
	onReceiveEffects ...OnReceiveEffect[T],
) error {
	ctx, cancel := context.WithTimeout(ctx, b.timeOut)
	defer cancel()

	errChan := make(chan error)
	doneChan := make(chan struct{})

	var wg sync.WaitGroup

	for _, effect := range onReceiveEffects {
		node := b.registerListener()
		wg.Add(1)
		go func(n *Node[T], e OnReceiveEffect[T]) {
			defer wg.Done()
			errChan <- b.receive(ctx, node, e)
		}(node, effect)
	}

	go func() {
		wg.Wait()
		doneChan <- struct{}{}
	}()

	f()

	var errs []error
	for {
		select {
		case err := <-errChan:
			if err != nil {
				errs = append(errs, err)
			}
		case <-doneChan:
			goto finish
		case <-ctx.Done():
			goto finish
		}
	}
finish:
	return errors.Join(errs...)
}

func (b *Broadcaster[T]) broadcast(msg T) {
	defer b.RUnlock()
	b.RLock()

	ok, element := b.channelList.Head()
	for ok {
		element.Receive(msg)
		ok, element = b.channelList.Iter(element)
	}
}

func (b *Broadcaster[T]) registerListener() *Node[T] {
	defer b.Unlock()
	b.Lock()
	return b.channelList.Add()
}

func (b *Broadcaster[T]) receive(
	ctx context.Context,
	n *Node[T],
	e OnReceiveEffect[T],
) (err error) {
	ctx, cancel := context.WithTimeout(ctx, b.timeOut)
	defer cancel()
	for {
		select {
		case msg := <-n.Wait():
			done, checkErr := e.CheckEffect(msg)
			if checkErr != nil {
				err = checkErr
				goto cleanUp
			}
			if done {
				goto cleanUp
			}
		case <-ctx.Done():
			err = e.Error()
			goto cleanUp
		}
	}
cleanUp:
	defer b.Unlock()
	b.Lock()
	b.channelList.Remove(n)
	return
}
