package lru

import (
	"math/rand"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCache(t *testing.T) {
	t.Run("empty cache", func(t *testing.T) {
		c := NewCache(10)

		_, ok := c.Get("aaa")
		require.False(t, ok)

		_, ok = c.Get("bbb")
		require.False(t, ok)
	})

	t.Run("simple", func(t *testing.T) {
		c := NewCache(5)

		wasInCache := c.Set("aaa", 100)
		require.False(t, wasInCache)

		wasInCache = c.Set("bbb", 200)
		require.False(t, wasInCache)

		val, ok := c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 100, val)

		val, ok = c.Get("bbb")
		require.True(t, ok)
		require.Equal(t, 200, val)

		wasInCache = c.Set("aaa", 300)
		require.True(t, wasInCache)

		val, ok = c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 300, val)

		val, ok = c.Get("ccc")
		require.False(t, ok)
		require.Nil(t, val)
	})

	// t.Run("purge logic", func(t *testing.T) {
	// 	// Write me
	// })
}

func TestCacheMultithreading(_ *testing.T) {
	// t.Skip() // Remove me if task with asterisk completed.

	c := NewCache(10)
	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Set(Key(strconv.Itoa(i)), i)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Get(Key(strconv.Itoa(rand.Intn(1_000_000))))
		}
	}()

	wg.Wait()
}

func TestExeedCacheLength(t *testing.T) {
	c := NewCache(3)

	c.Set("k1", 1)
	c.Set("k2", 2)
	c.Set("k3", 3)
	c.Set("k4", 4)
	require.Equal(t, 3, c.queue.length)
	val, ok := c.Get("k1")

	require.False(t, ok)
	require.Nil(t, val)

	val, ok = c.Get("k2")
	require.Equal(t, 2, val)
	require.True(t, ok)
	val, ok = c.Get("k3")
	require.True(t, ok)
	require.Equal(t, 3, val)
	val, ok = c.Get("k4")
	require.True(t, ok)
	require.Equal(t, 4, val)
}

func getQueueValues(c *LruCache) []pair {
	result := []pair{}

	temp := c.queue.head
	for temp != nil {
		p, ok := temp.Value.(pair)
		if !ok { // если не тот тип
			panic("Unexpected value type")
		}
		result = append(result, p)
		temp = temp.Next
	}

	return result
}

func TestCacheLogic(t *testing.T) {
	c := NewCache(3)

	c.Set("k1", 1) // k1
	c.Set("k2", 2) // k2 k1
	c.Set("k3", 3) // k3 k2 k1
	pairs := getQueueValues(c)
	require.EqualValues(t, []pair{
		{key: "k3", value: 3},
		{key: "k2", value: 2},
		{key: "k1", value: 1},
	},
		pairs)

	c.Set("k2", 22) // k2 k3 k1
	c.Set("k3", 33) // k3 k2 k1
	pairs = getQueueValues(c)
	require.EqualValues(t, []pair{
		{key: "k3", value: 33},
		{key: "k2", value: 22},
		{key: "k1", value: 1},
	},
		pairs)

	_, ok := c.Get("k1") // k1 k3 k2
	require.True(t, ok)
	pairs = getQueueValues(c)
	require.EqualValues(t, []pair{
		{key: "k1", value: 1},
		{key: "k3", value: 33},
		{key: "k2", value: 22},
	}, pairs)

	_, ok = c.Get("k3") // k3 k1 k2
	require.True(t, ok)
	pairs = getQueueValues(c)
	require.EqualValues(t, []pair{
		{key: "k3", value: 33},
		{key: "k1", value: 1},
		{key: "k2", value: 22},
	}, pairs)

	c.Set("k4", 4) // k4 k3 k1
	require.Equal(t, 3, c.queue.length)
	pairs = getQueueValues(c)
	require.EqualValues(t, []pair{
		{key: "k4", value: 4},
		{key: "k3", value: 33},
		{key: "k1", value: 1},
	}, pairs)

	c.Clear()
	require.Equal(t, 0, c.queue.length)
	require.Nil(t, c.queue.head)
	require.Nil(t, c.queue.tail)
	require.Empty(t, c.items)
}
