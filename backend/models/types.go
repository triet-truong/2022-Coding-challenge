package models

import (
	"sync"

	"golang.org/x/exp/constraints"
)

type TaskQueue chan Incident

type ConcurrentMap[K constraints.Ordered, V any] interface {
	Get(K) V
	Set(K, V)
	Delete(K)
	Values() []V
}

var _ ConcurrentMap[int, any] = &ObjectMapByID[int, any]{}

type ObjectMapByID[K constraints.Ordered, V any] struct {
	rawMap map[K]V
	mutex  *sync.RWMutex
}

func NewObjectMapByID[K constraints.Ordered, V any]() ObjectMapByID[K, V] {
	return ObjectMapByID[K, V]{
		rawMap: make(map[K]V),
		mutex:  &sync.RWMutex{},
	}
}

func (m *ObjectMapByID[K, V]) Get(key K) V {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.rawMap[key]
}

func (m *ObjectMapByID[K, V]) Set(key K, value V) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.rawMap[key] = value
}

func (m *ObjectMapByID[K, V]) Copy() map[K]V {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	newMap := make(map[K]V)
	for k, v := range m.rawMap {
		newMap[k] = v
	}

	return newMap
}

func (m *ObjectMapByID[K, V]) Delete(key K) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	delete(m.rawMap, key)
}

func (m *ObjectMapByID[K, V]) Values() []V {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	res := make([]V, len(m.rawMap))
	i := 0
	for _, v := range m.rawMap {
		res[i] = v
		i++
	}
	return res
}
