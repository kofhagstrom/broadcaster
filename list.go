package broadcaster

type List[T any] interface {
	Add() *Node[T]
	Remove(node *Node[T])
	Head() (bool, *Node[T])
	Tail() (bool, *Node[T])
	Iter(node *Node[T]) (bool, *Node[T])
}

func NewNode[T any](prev *Node[T], next *Node[T]) *Node[T] {
	return &Node[T]{Value: make(chan T, 1), Next: next, Prev: prev}
}

type Node[T any] struct {
	Value chan T
	Next  *Node[T]
	Prev  *Node[T]
}

func (n *Node[T]) Receive(msg T) {
	n.Value <- msg
}

func (n *Node[T]) Wait() chan T {
	return n.Value
}
