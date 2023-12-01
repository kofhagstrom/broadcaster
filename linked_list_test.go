package broadcaster

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func (l *DoublyLinkedList[T]) HeadT(t *testing.T) *Node[T] {
	ok, head := l.Head()
	require.True(t, ok)
	return head
}

func (l *DoublyLinkedList[T]) TailT(t *testing.T) *Node[T] {
	ok, tail := l.Tail()
	require.True(t, ok)
	return tail
}

func (l *DoublyLinkedList[T]) IterT(t *testing.T, n *Node[T]) *Node[T] {
	ok, next := l.Iter(n)
	require.True(t, ok)
	return next
}

func TestAddOneNode(t *testing.T) {

	l := NewDoublyLinkedList[int]()

	node := l.Add()
	require.Equal(t, node, l.HeadT(t))
	require.Equal(t, node, l.TailT(t))

	require.NotNil(t, l.HeadT(t))
	require.Nil(t, l.head.Prev)
	require.Nil(t, l.head.Next)

	require.NotNil(t, l.TailT(t))
	require.Nil(t, l.TailT(t).Prev)
	require.Nil(t, l.TailT(t).Next)

	require.Equal(t, l.head, l.TailT(t))
}

func TestAddTwoNodes(t *testing.T) {

	l := NewDoublyLinkedList[int]()

	node1 := l.Add()
	node2 := l.Add()
	require.Equal(t, node1, l.HeadT(t))
	require.Equal(t, node2, l.TailT(t))

	require.NotNil(t, l.HeadT(t))
	require.Nil(t, l.head.Prev)

	require.NotNil(t, l.TailT(t))
	require.Nil(t, l.TailT(t).Next)

	require.NotEqual(t, l.head, l.TailT(t))
	require.Equal(t, l.head, l.TailT(t).Prev)
	require.Equal(t, l.head.Next, l.TailT(t))
}

func TestAddThreeNodes(t *testing.T) {

	l := NewDoublyLinkedList[int]()

	node1 := l.Add()
	node2 := l.Add()
	node3 := l.Add()

	require.Equal(t, node1, l.HeadT(t))
	require.Equal(t, node3, l.TailT(t))

	require.NotNil(t, l.HeadT(t))
	require.Nil(t, l.head.Prev)
	require.Equal(t, node1, l.HeadT(t))

	require.Equal(t, node2, l.head.Next)
	require.Equal(t, l.head, node2.Prev)
	require.Equal(t, l.TailT(t), node2.Next)

	require.NotNil(t, l.TailT(t))
	require.Nil(t, l.TailT(t).Next)
	require.Equal(t, node3, l.TailT(t))

	require.NotEqual(t, l.head, l.TailT(t))
	require.Equal(t, node2, l.TailT(t).Prev)
	require.Equal(t, node2.Next, l.TailT(t))
}

func TestAddThreeNodesAndRemoveHead(t *testing.T) {

	l := NewDoublyLinkedList[int]()

	node1 := l.Add()
	node2 := l.Add()
	node3 := l.Add()

	require.Equal(t, node1, l.HeadT(t))
	require.Equal(t, node3, l.TailT(t))

	l.Remove(node1)

	require.Equal(t, node2, l.HeadT(t))
	require.Equal(t, node3, l.TailT(t))

	require.Equal(t, node2, l.HeadT(t))
	require.Nil(t, l.head.Prev)

	require.Equal(t, l.head.Next, l.TailT(t))
	require.Equal(t, l.TailT(t).Prev, l.HeadT(t))

	require.Nil(t, l.TailT(t).Next)

}

func TestAddThreeNodesAndRemoveSecond(t *testing.T) {

	l := NewDoublyLinkedList[int]()

	node1 := l.Add()
	node2 := l.Add()
	node3 := l.Add()

	require.Equal(t, node1, l.HeadT(t))
	require.Equal(t, node3, l.TailT(t))

	l.Remove(node2)

	require.Equal(t, node1, l.HeadT(t))
	require.Equal(t, node3, l.TailT(t))

	require.Equal(t, node1, l.HeadT(t))
	require.Equal(t, node3, l.TailT(t))
	require.Nil(t, l.head.Prev)
	require.Nil(t, l.TailT(t).Next)

	require.NotEqual(t, l.head, l.TailT(t))
	require.Equal(t, l.head, l.TailT(t).Prev)
	require.Equal(t, l.head.Next, l.TailT(t))

}

func TestAddThreeNodesAndRemoveTail(t *testing.T) {

	l := NewDoublyLinkedList[int]()

	node1 := l.Add()
	node2 := l.Add()
	node3 := l.Add()

	require.Equal(t, node1, l.HeadT(t))
	require.Equal(t, node3, l.TailT(t))

	l.Remove(node3)

	require.Equal(t, node1, l.HeadT(t))
	require.Equal(t, node2, l.TailT(t))

	require.Equal(t, node1, l.HeadT(t))
	require.Equal(t, node2, l.TailT(t))
	require.Nil(t, l.head.Prev)
	require.Nil(t, l.TailT(t).Next)

	require.NotEqual(t, l.head, l.TailT(t))
	require.Equal(t, l.head, l.TailT(t).Prev)
	require.Equal(t, l.head.Next, l.TailT(t))

}

func TestAddFourNodesAndRemoveHead(t *testing.T) {

	l := NewDoublyLinkedList[int]()

	node1 := l.Add()
	node2 := l.Add()
	node3 := l.Add()
	node4 := l.Add()

	require.Equal(t, node1, l.HeadT(t))
	require.Equal(t, node4, l.TailT(t))

	l.Remove(node1)

	require.Equal(t, node2, l.HeadT(t))
	require.Equal(t, node4, l.TailT(t))

	require.Equal(t, node2, l.HeadT(t))
	require.Equal(t, node4, l.TailT(t))
	require.Nil(t, l.head.Prev)
	require.Nil(t, l.TailT(t).Next)

	require.NotEqual(t, l.head, l.TailT(t))
	require.Equal(t, node2.Next, node3)
	require.Equal(t, node3.Prev, node2)
	require.Equal(t, node3.Next, node4)
	require.Equal(t, node4.Prev, node3)

}

func TestAddFourNodesAndRemoveSecond(t *testing.T) {
	l := NewDoublyLinkedList[int]()

	node1 := l.Add()
	node2 := l.Add()
	node3 := l.Add()
	node4 := l.Add()

	require.Equal(t, node1, l.HeadT(t))
	require.Equal(t, node4, l.TailT(t))

	l.Remove(node2)

	require.Equal(t, node1, l.HeadT(t))
	require.Equal(t, node4, l.TailT(t))

	require.Equal(t, node1, l.HeadT(t))
	require.Equal(t, node4, l.TailT(t))
	require.Nil(t, node1.Prev)
	require.Nil(t, node4.Next)
	require.Nil(t, l.head.Prev)
	require.Nil(t, l.TailT(t).Next)

	require.Equal(t, node1.Next, node3)
	require.Equal(t, node3.Prev, node1)
	require.Equal(t, node3.Next, node4)
	require.Equal(t, node4.Prev, node3)
}

func TestAddFourNodesAndRemoveThird(t *testing.T) {
	l := NewDoublyLinkedList[int]()

	node1 := l.Add()
	node2 := l.Add()
	node3 := l.Add()
	node4 := l.Add()

	require.Equal(t, node1, l.HeadT(t))
	require.Equal(t, node4, l.TailT(t))

	l.Remove(node3)

	require.Equal(t, node1, l.HeadT(t))
	require.Equal(t, node4, l.TailT(t))

	require.Equal(t, node1, l.HeadT(t))
	require.Equal(t, node4, l.TailT(t))
	require.Nil(t, l.head.Prev)
	require.Nil(t, l.TailT(t).Next)

	require.Equal(t, node1.Next, node2)
	require.Equal(t, node2.Prev, node1)
	require.Equal(t, node2.Next, node4)
	require.Equal(t, node4.Prev, node2)
}

func TestAddFourNodesAndRemoveTail(t *testing.T) {

	l := NewDoublyLinkedList[int]()

	node1 := l.Add()
	node2 := l.Add()
	node3 := l.Add()
	node4 := l.Add()

	require.Equal(t, node1, l.HeadT(t))
	require.Equal(t, node4, l.TailT(t))

	l.Remove(node4)

	require.Equal(t, node1, l.HeadT(t))
	require.Equal(t, node3, l.TailT(t))

	require.Equal(t, node1, l.HeadT(t))
	require.Equal(t, node3, l.TailT(t))
	require.Nil(t, l.head.Prev)
	require.Nil(t, l.TailT(t).Next)

	require.Equal(t, node1.Next, node2)
	require.Equal(t, node2.Prev, node1)
	require.Equal(t, node2.Next, node3)
	require.Equal(t, node3.Prev, node2)

}

func TestIter(t *testing.T) {
	l := NewDoublyLinkedList[int]()

	node1 := l.Add()
	node2 := l.Add()
	node3 := l.Add()
	node4 := l.Add()

	node := l.HeadT(t)
	require.Equal(t, node1, node)

	node = l.IterT(t, node)
	require.Equal(t, node2, node)

	node = l.IterT(t, node)
	require.Equal(t, node3, node)

	node = l.IterT(t, node)
	require.Equal(t, node4, node)

	ok, node := l.Iter(node)
	require.False(t, ok)
	require.Nil(t, node)
}
