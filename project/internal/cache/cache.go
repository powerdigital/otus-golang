package lrucache

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
}

func (c lruCache) Set(key Key, value interface{}) bool {
	_, wasInCache := c.Get(key)

	listItem := c.queue.PushFront(value)
	c.items[key] = listItem

	if c.queue.Len() > c.capacity {
		c.queue.Remove(c.queue.Back())
	}

	return wasInCache
}

func (c lruCache) Get(key Key) (value interface{}, ok bool) {
	if _, ok = c.items[key]; ok != false {
		value = c.items[key].Value
	}

	return value, ok
}

func (c lruCache) Clear() {
	for k := range c.items {
		delete(c.items, k)
	}

	for i := c.queue.Front(); i != nil; i = i.Next {
		c.queue.Remove(i)
	}
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}
