package ploader

import "sync"

type ID = string
type Value = interface{}

type CacheMap struct {
	sync.RWMutex
	data map[ID]FutureValue
}

func NewCacheMap() *CacheMap {
	return &CacheMap{
		data: make(map[ID]FutureValue),
	}
}

func (cm *CacheMap) LoadOrElse(
	key ID,
	producer func() (Value, error),
) FutureValue {
	if value, ok := cm.Load(key); ok {
		return value
	}
	cm.Lock() //using write lock
	if value, ok := cm.data[key]; ok {
		return value
	}
	value := MakeFutureValue(producer)
	cm.data[key] = value
	cm.Unlock() //write unlock
	return value
}

func (cm *CacheMap) Load(key ID) (value FutureValue, ok bool) {
	cm.RLock()
	result, ok := cm.data[key]
	cm.RUnlock()
	return result, ok
}

func (cm *CacheMap) Delete(key ID) {
	cm.Lock()
	delete(cm.data, key)
	cm.Unlock()
}

func (cm *CacheMap) Store(key ID, value FutureValue) {
	cm.Lock()
	cm.data[key] = value
	cm.Unlock()
}
