package hw04lrucache

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
	Storage []*ListItem
}

func (l *list) Len() int {
	return len(l.Storage)
}

func (l *list) Front() *ListItem {
	if l.Len() == 0 {
		return nil
	}

	return l.Storage[0]
}

func (l *list) Back() *ListItem {
	if l.Len() == 0 {
		return nil
	}

	return l.Storage[l.Len()-1]
}

func (l *list) PushFront(val interface{}) *ListItem {
	var next *ListItem
	var item ListItem

	if l.Len() != 0 {
		next = l.Storage[0]
	}

	item = ListItem{val, next, nil}

	var initList []*ListItem
	initList = append(initList, &item)
	l.Storage = append(initList, l.Storage...)

	return l.Storage[0]
}

func (l *list) PushBack(v interface{}) *ListItem {
	index := l.Len()
	item := ListItem{v, nil, l.Storage[index-1]}
	l.Storage = append(l.Storage, &item)

	l.Storage[index-1].Next = &item

	return l.Storage[index]
}

func (l *list) Remove(item *ListItem) {
	l.MoveToFront(item)
	l.Storage = l.Storage[1:l.Len()]
}

func (l *list) MoveToFront(item *ListItem) {
	if l.Len() < 2 || item.Prev == nil {
		return
	}

	index := l.getIndex(item)

	deepCopy := *item
	deepCopy.Prev = nil
	deepCopy.Next = l.Storage[0]

	prev := item.Prev
	next := item.Next

	if item.Prev != nil {
		item.Prev.Next = next
	}

	if item.Next != nil {
		item.Next.Prev = prev
	}

	after := l.Storage[index+1 : l.Len()]
	before := l.Storage[0:index]
	before = append(before, after...)
	storage := []*ListItem{&deepCopy}
	storage = append(storage, before...)

	l.Storage = storage
}

func (l *list) getIndex(item *ListItem) int {
	for k, v := range l.Storage {
		if v == item {
			return k
		}
	}

	panic("Index not found")
}

func NewList() List {
	return new(list)
}
