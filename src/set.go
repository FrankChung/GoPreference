package main

import (
	"reflect"
	"sync"
)

type set struct {
	m []observer
	*sync.RWMutex
}

func newSet() *set {
	return &set{m: make([]observer, 0), RWMutex: &sync.RWMutex{}}
}

func (s *set) add(item observer) {
	s.Lock()
	defer s.Unlock()
	for _, v := range s.m {
		if reflect.DeepEqual(v, item) {
			return
		}
	}
	s.m = append(s.m, item)
}

func (s *set) remove(item observer) {
	s.Lock()
	defer s.Unlock()
	for i, v := range s.m {
		if reflect.DeepEqual(v, item) {
			s.m, s.m[len(s.m)-1] = append(s.m[:i], s.m[i+1:]...), nil
			break
		}
	}
}

func (s *set) has(item observer) bool {
	s.RLock()
	defer s.RUnlock()
	for _, v := range s.m {
		if reflect.DeepEqual(v, item) {
			return true
		}
	}
	return false
}

func (s *set) len() int {
	s.RLock()
	defer s.RUnlock()
	return len(s.m)
}

func (s *set) list() []observer {
	s.RLock()
	defer s.RUnlock()
	list := make([]observer, 0, len(s.m))
	for _, item := range s.m {
		list = append(list, item)
	}
	return list
}

func (s *set) clear() {
	s.Lock()
	defer s.Unlock()
	s.m = make([]observer, 0)
}

func (s *set) isEmpty() bool {
	return s.len() == 0
}
