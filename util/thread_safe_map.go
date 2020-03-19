package util

import "sync"

// === 线程安全的map ===
type SafeMap struct {
	sync.RWMutex
	Map map[string]interface{}
}

func NewSafeMap() *SafeMap {
	sm := new(SafeMap)
	sm.Map = make(map[string]interface{})
	return sm
}

func (this *SafeMap) ReadSafeMap(key string) (interface{}, bool) {
	this.RLock()
	value, ok := this.Map[key]
	this.RUnlock()
	return value, ok
}

func (this *SafeMap) WriteSafeMap(key string, value interface{}) {
	this.Lock()
	this.Map[key] = value
	this.Unlock()
}

func (this *SafeMap) DeleteKey(key string) {
	this.Lock()
	delete(this.Map, key)
	this.Unlock()
}

// === 线程安全的 map End ===
