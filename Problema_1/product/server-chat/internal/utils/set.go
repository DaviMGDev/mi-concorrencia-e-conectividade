package utils

import "sync"

var void = Void{}

type Set[T comparable] struct {
	elements map[T]Void
	mutex    sync.Mutex
}

func NewSet[T comparable]() *Set[T] {
	return &Set[T]{
		elements: make(map[T]Void),
		mutex:    sync.Mutex{},
	}
}
