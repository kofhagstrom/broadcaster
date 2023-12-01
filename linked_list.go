package broadcaster

func NewDoublyLinkedList[T any]() *DoublyLinkedList[T] {
	return &DoublyLinkedList[T]{head: nil, tail: nil}
}

type DoublyLinkedList[T any] struct {
	head *Node[T]
	tail *Node[T]
}

var _ List[any] = (*DoublyLinkedList[any])(nil)

func (l *DoublyLinkedList[T]) Add() *Node[T] {
	switch l.head {
	case nil:
		newNode := NewNode[T](nil, nil)
		l.head = newNode
		l.tail = newNode
		return newNode
	default:
		newLast := NewNode[T](l.tail, nil)
		l.tail.Next = newLast
		l.tail = newLast
		return newLast
	}
}

func (l *DoublyLinkedList[T]) Remove(n *Node[T]) {
	switch n.Prev {
	case nil:
		l.head = n.Next
		if n.Next != nil {
			n.Next.Prev = nil
		}
	default:
		n.Prev.Next = n.Next
		if n.Next != nil {
			n.Next.Prev = n.Prev
		}
	}

	if n.Next == nil {
		l.tail = n.Prev
	}

	close(n.Value)
}

func (l *DoublyLinkedList[T]) Head() (bool, *Node[T]) {
	ok := l.head != nil
	return ok, l.head
}

func (l *DoublyLinkedList[T]) Tail() (bool, *Node[T]) {
	ok := l.tail != nil
	return ok, l.tail
}

func (l *DoublyLinkedList[T]) Iter(node *Node[T]) (bool, *Node[T]) {
	ok := node.Next != nil
	return ok, node.Next
}
