package lru

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestList(t *testing.T) {
	t.Run("empty list", func(t *testing.T) {
		l := NewList()

		require.Equal(t, 0, l.Len())
		require.Nil(t, l.Front())
		require.Nil(t, l.Back())
	})

	t.Run("complex", func(t *testing.T) {
		l := NewList()

		l.PushFront(10) // [10]
		l.PushBack(20)  // [10, 20]
		l.PushBack(30)  // [10, 20, 30]
		require.Equal(t, 3, l.Len())

		middle := l.Front().Next // 20
		l.Remove(middle)         // [10, 30]
		require.Equal(t, 2, l.Len())

		for i, v := range [...]int{40, 50, 60, 70, 80} {
			if i%2 == 0 {
				l.PushFront(v)
			} else {
				l.PushBack(v)
			}
		} // [80, 60, 40, 10, 30, 50, 70]

		require.Equal(t, 7, l.Len())
		require.Equal(t, 80, l.Front().Value)
		require.Equal(t, 70, l.Back().Value)

		l.MoveToFront(l.Front()) // [80, 60, 40, 10, 30, 50, 70]
		l.MoveToFront(l.Back())  // [70, 80, 60, 40, 10, 30, 50]

		elems := make([]int, 0, l.Len())
		for i := l.Front(); i != nil; i = i.Next {
			elems = append(elems, i.Value.(int))
		}
		require.Equal(t, []int{70, 80, 60, 40, 10, 30, 50}, elems)
	})
}

func TestList2(t *testing.T) {
	doublyLL := NewList()

	li1 := doublyLL.PushFront(1)
	li2 := doublyLL.PushFront(2)
	li3 := doublyLL.PushFront(3)
	li4 := doublyLL.PushBack(4) // 3 2 1 4
	require.Equal(t, 4, doublyLL.length)

	liFront := doublyLL.Front()
	require.Equal(t, 3, liFront.Value)

	liBack := doublyLL.Back()
	require.Equal(t, 4, liBack.Value)

	doublyLL.Remove(li3) // 2 1 4
	require.Equal(t, 3, doublyLL.length)
	liFront = doublyLL.Front()
	require.Equal(t, 2, liFront.Value)

	doublyLL.MoveToFront(li4) // 4 2 1
	liFront = doublyLL.Front()
	require.Equal(t, 4, liFront.Value)

	doublyLL.printList()

	doublyLL.DeleteLinkedList()
	require.Equal(t, 0, doublyLL.length)

	require.Nil(t, doublyLL.head)
	require.Nil(t, doublyLL.tail)

	require.Nil(t, li1.Next)
	require.Nil(t, li1.Prev)
	require.Nil(t, li2.Next)
	require.Nil(t, li2.Prev)
	require.Nil(t, li4.Next)
	require.Nil(t, li4.Prev)
}
