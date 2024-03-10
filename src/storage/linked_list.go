package storage

func NewLinkedList[T any]() LinkedList[T] {
	return LinkedList[T]{}
}

type LinkedList[T any] struct {
	Head *LinkedListNode[T]
	Tail *LinkedListNode[T]
}

func (list *LinkedList[T]) Append(data T) {
	node := NewLinkedListNode(data)
	if list.Head == nil && list.Tail == nil {
		list.Head = &node
		list.Tail = &node
	} else {
		node.Previous = list.Tail
		list.Tail.Next = &node
		list.Tail = &node
	}
}

func (list *LinkedList[T]) DeleteNode(node *LinkedListNode[T]) {
	if node != nil {
		if node.Previous != nil {
			node.Previous.Next = node.Next
		}
		if node.Next != nil {
			node.Next.Previous = node.Previous
		}
		node = nil
	}
}

type LinkedListNode[T any] struct {
	Next     *LinkedListNode[T]
	Previous *LinkedListNode[T]
	Data     T
}

func NewLinkedListNode[T any](data T) LinkedListNode[T] {
	return LinkedListNode[T]{
		Data: data,
	}
}
