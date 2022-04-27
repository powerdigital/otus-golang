package lrucache

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	size    int
	head    *ListItem
	tail    *ListItem
	storage []*ListItem
}

func (l *list) Len() int {
	return len(l.storage)
}

func (l *list) Front() *ListItem {
	if l.size == 0 {
		return nil
	}

	return l.storage[0]
}

func (l *list) Back() *ListItem {
	if l.size == 0 {
		return nil
	}

	return l.tail
}

func (l *list) PushFront(val interface{}) *ListItem {
	var next *ListItem
	if l.size != 0 {
		next = l.storage[0]
	}

	item := ListItem{val, next, nil}
	l.storage = append([]*ListItem{&item}, l.storage...)
	l.head = &item
	l.size++

	return l.head
}

func (l *list) PushBack(v interface{}) *ListItem {
	index := l.size
	item := ListItem{v, nil, l.storage[index-1]}

	l.storage = append(l.storage, &item)
	l.storage[index-1].Next = &item
	l.tail = &item
	l.size++

	return l.storage[index]
}

func (l *list) Remove(item *ListItem) {
	l.MoveToFront(item)
	l.storage = l.storage[1:l.Len()]
	l.size--
}

func (l *list) MoveToFront(item *ListItem) {
	if l.size < 2 || item.Prev == nil {
		return
	}

	index := l.getIndex(item)

	deepCopy := *item
	deepCopy.Prev = nil
	deepCopy.Next = l.head

	prev := item.Prev
	next := item.Next

	if item.Prev != nil {
		item.Prev.Next = next
	}

	if item.Next != nil {
		item.Next.Prev = prev
	} else {
		l.tail = prev
	}

	sliceParts := l.storage[index+1 : l.size]
	sliceParts = append(l.storage[0:index], sliceParts...)

	l.storage = append([]*ListItem{&deepCopy}, sliceParts...)
	l.head = item
}

func (l *list) getIndex(item *ListItem) int {
	var index int
	for k, v := range l.storage {
		if v == item {
			index = k
			break
		}
	}

	return index
}

func NewList() List {
	return new(list)
}
