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
	temp := list.head.Next

	for temp != nil {
		if temp == i {
			temp.Prev.Next = temp.Next
			temp.Next.Prev = temp.Prev.Next
			list.length--
			return
		}
		temp = temp.Next
	}
}

// Переместить элемент в начало.
func (list *List) MoveToFront(i *ListItem) {
	if i == list.head {
		return
	}

	temp := list.head.Next // temp - второй элемент в списке

	for temp != nil {
		if temp == i {
			temp.Prev.Next = temp.Next
			if temp.Next != nil {
				temp.Next.Prev = temp.Prev
			}
			if temp == list.tail {
				list.tail = temp.Prev
			}

			list.head.Prev = temp
			temp.Next = list.head
			list.head = temp
			temp.Prev = nil

			return
		}
		temp = temp.Next
	}
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
