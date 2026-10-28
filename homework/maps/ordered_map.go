package main

import "cmp"

type OrderedMap[K cmp.Ordered, V any] struct {
	root *orderedMapNode[K, V]
	size int
}

type orderedMapNode[K cmp.Ordered, V any] struct {
	key         K
	value       V
	left, right *orderedMapNode[K, V]
}

func NewOrderedMap[K cmp.Ordered, V any]() OrderedMap[K, V] {
	return OrderedMap[K, V]{}
}

func (m *OrderedMap[K, V]) Insert(key K, value V) {
	link := m.findLink(key)
	if *link != nil {
		(*link).value = value
		return
	}
	*link = &orderedMapNode[K, V]{key: key, value: value}
	m.size++
}

func (m *OrderedMap[K, V]) Erase(key K) {
	link := m.findLink(key)
	node := *link
	if node == nil {
		return
	}

	switch {
	case node.left == nil:
		*link = node.right
	case node.right == nil:
		*link = node.left
	default:
		successorLink := &node.right
		for (*successorLink).left != nil {
			successorLink = &(*successorLink).left
		}
		successor := *successorLink
		node.key, node.value = successor.key, successor.value
		*successorLink = successor.right
	}
	m.size--
}

func (m *OrderedMap[K, V]) Contains(key K) bool {
	return *m.findLink(key) != nil
}

func (m *OrderedMap[K, V]) Size() int {
	return m.size
}

func (m *OrderedMap[K, V]) ForEach(action func(K, V)) {
	var stack []*orderedMapNode[K, V]
	node := m.root
	for node != nil || len(stack) > 0 {
		for node != nil {
			stack = append(stack, node)
			node = node.left
		}
		node = stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		action(node.key, node.value)
		node = node.right
	}
}

func (m *OrderedMap[K, V]) findLink(key K) **orderedMapNode[K, V] {
	link := &m.root
	for *link != nil {
		switch cmp.Compare(key, (*link).key) {
		case -1:
			link = &(*link).left
		case 1:
			link = &(*link).right
		default:
			return link
		}
	}
	return link
}
