package util

import "sync"

// SafeMap 线程安全的map
type SafeMap struct {
	sync.RWMutex
	Map map[int32]interface{}
}

func NewSafeMap() *SafeMap {
	sm := new(SafeMap)
	sm.Map = make(map[int32]interface{})
	return sm
}

func (m *SafeMap) ReadSafeMap(key int32) (interface{}, bool) {
	m.RLock()
	value, ok := m.Map[key]
	m.RUnlock()
	return value, ok
}

func (m *SafeMap) WriteSafeMap(key int32, value interface{}) {
	m.Lock()
	m.Map[key] = value
	m.Unlock()
}

func (m *SafeMap) DeleteKey(key int32) {
	m.Lock()
	delete(m.Map, key)
	m.Unlock()
}
