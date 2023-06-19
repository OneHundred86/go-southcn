package session

import (
	"fmt"
	"github.com/OneHundred86/go-southcn/cache"
	"reflect"
	"sync"
	"time"
)

type CacheSession struct {
	Session

	id        string
	data      map[string]any
	lock      sync.RWMutex
	cache     cache.Cache
	keyPrefix string
	ttl       time.Duration
}

func NewCacheSession(cache cache.Cache, keyPrefix string, ttl time.Duration) Session {
	return &CacheSession{
		data:      make(map[string]any),
		lock:      sync.RWMutex{},
		cache:     cache,
		keyPrefix: keyPrefix,
		ttl:       ttl,
	}
}

func (t *CacheSession) SetSessionID(ID string) {
	t.id = ID
}

func (t *CacheSession) GetSessionID() string {
	return t.id
}

func (t *CacheSession) makeKey() string {
	return fmt.Sprintf("%s%s", t.keyPrefix, t.id)
}

func (t *CacheSession) Load() {
	t.lock.Lock()
	defer t.lock.Unlock()
	t.cache.Get(t.makeKey(), &t.data)
}

func (t *CacheSession) Save() {
	if len(t.data) > 0 {
		t.cache.Set(t.makeKey(), t.data, t.ttl)
	}
}

func (t *CacheSession) Flush() {
	t.lock.Lock()
	defer t.lock.Unlock()

	t.data = nil
	t.cache.Delete(t.makeKey())
}

func (t *CacheSession) Set(key string, val any) {
	t.lock.Lock()
	defer t.lock.Unlock()

	t.data[key] = val
}

func (t *CacheSession) Get(key string, val any) (exists bool) {
	t.lock.RLock()
	defer t.lock.RUnlock()

	val0, exists := t.data[key]
	if !exists {
		return
	}

	// 反射
	rv0 := reflect.ValueOf(val) // rv0为指针（因为val为指针）
	rv := rv0.Elem()            // rv为rv0指针指向的值
	rv.Set(reflect.ValueOf(val0).Convert(rv.Type()))

	return
}

func (t *CacheSession) Del(key string) {
	t.lock.Lock()
	defer t.lock.Unlock()

	delete(t.data, key)
}

func (t *CacheSession) All() map[string]any {
	return t.data
}
