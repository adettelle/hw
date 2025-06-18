package hw04lrucache

import "fmt"

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type List struct {
	head   *ListItem
	tail   *ListItem
	length int
}

func NewList() *List {
	return &List{}
}

func (list *List) Len() int {
	return list.length
}

func (list *List) Front() *ListItem {
	return list.head
}

func (list *List) Back() *ListItem {
	return list.tail
}

// Добавить значение в начало.
func (list *List) PushFront(v interface{}) *ListItem {
	i := ListItem{Value: v}

	if list.head == nil {
		list.head = &i
		list.tail = &i
		i.Prev = nil
		i.Next = nil
	} else {
		list.head.Prev = &i
		i.Next = list.head
		list.head = &i
	}
	list.length++

	return &i
}

// Добавить значение в конец.
func (list *List) PushBack(v interface{}) *ListItem {
	i := ListItem{Value: v}

	if list.head == nil {
		list.head = &i
		list.tail = &i
		i.Prev = nil
		i.Next = nil
	} else {
		list.tail.Next = &i
		i.Prev = list.tail
		list.tail = &i
	}
	list.length++

	return &i
}

// Удалить элемент.
// Считатем, что метод вызывается только от существующих в списке элементов.
func (list *List) Remove(i *ListItem) {
	if list.head == nil {
		return
	}

	if list.head == i {
		list.head = list.head.Next
		list.head.Prev = nil
		list.length--
		return
	}
	if list.tail == i {
		list.tail = list.tail.Prev
		list.tail.Next = nil
		list.length--
		return
	}

	i.Prev.Next = i.Next
	i.Next.Prev = i.Prev
	list.length--
}

// Переместить элемент в начало.
// Считатем, что метод вызывается только от существующих в списке элементов.
func (list *List) MoveToFront(i *ListItem) {
	if i == list.head {
		return
	}

	i.Prev.Next = i.Next
	if i.Next != nil {
		i.Next.Prev = i.Prev
	}

	if list.tail == i {
		list.tail = i.Prev
	}

	i.Next = list.head
	list.head.Prev = i
	i.Prev = nil
	list.head = i
}

func (list *List) DeleteLinkedList() {
	current := list.head
	for current != nil {
		current.Prev = nil
		temp := current
		current = current.Next
		temp.Next = nil
	}
	list.head = nil
	list.length = 0

	list.tail = nil
}

func (list *List) printList() {
	if list.head == nil {
		fmt.Println("Empty Linked List")
	} else {
		temp := list.head
		for temp != nil {
			fmt.Printf("%v <-> ", temp.Value)
			temp = temp.Next
		}
	}
	fmt.Println()
}
