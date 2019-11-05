
package memory

import (
	"container/list"
	"errors"
)

// EvictCallback is used to get a callback when a cache entry is evicted
type EvictCallback func(key interface{}, value interface{})

// RAM : LRU implements a non-thread safe fixed size LRU cache
type RAM struct {
	size      int
	evictList *list.List
	items     map[interface{}]*list.Element
	onEvict   EvictCallback
}

// Memory : a thread-safe fixed size LRU Cache.
type Memory struct {
	ram  *RAM
	lock Mutex
}

// PageTable is used to hold a value in the evictList
type PageTable struct {
	key   interface{}
	value interface{}
}

// NewLRU constructs an LRU of the given size
func NewRAM(size int, onEvict EvictCallback) (*RAM, error) {
	if size <= 0 {
		return nil, errors.New("Must provide a positive size")
	}

	c := &RAM{
		size:      size,
		evictList: list.New(),
		items:     make(map[interface{}]*list.Element),
		onEvict:   onEvict,
	}
	return c, nil
}

// Purge is used to completely clear the cache.
func (ram *RAM) Purge() {
	for k, v := range ram.items {
		if ram.onEvict != nil {
			ram.onEvict(k, v.Value.(*PageTable).value)
		}
		delete(ram.items, k)
	}
	ram.evictList.Init()
}

// Add adds a value to the cache.  Returns true if an eviction occurred.
func (ram *RAM) Add(key, value interface{}) (evicted bool) {
	// Check for existing item
	if ent, ok := ram.items[key]; ok {
		ram.evictList.MoveToFront(ent)
		ent.Value.(*PageTable).value = value

		return false
	}

	// Add new item
	ent := &PageTable{key, value}
	pageTable := ram.evictList.PushFront(ent)
	ram.items[key] = pageTable

	evict := ram.evictList.Len() > ram.size
	// Verify size not exceeded
	if evict {
		ram.removeOldest()
	}
	return evict
}

// Get looks up a key's value from the cache.
func (ram *RAM) Get(key interface{}) (value interface{}, ok bool) {
	if ent, ok := ram.items[key]; ok {
		ram.evictList.MoveToFront(ent)
		if ent.Value.(*PageTable) == nil {
			return nil, false
		}
		return ent.Value.(*PageTable).value, true
	}

	return
}

// Contains checks if a key is in the cache, without updating the recent-ness
// or deleting it for being stale.
func (ram *RAM) Contains(key interface{}) (ok bool) {
	_, ok = ram.items[key]

	return ok
}

// Peek returns the key value (or undefined if not found) without updating
// the "recently used"-ness of the key.
func (ram *RAM) Peek(key interface{}) (value interface{}, ok bool) {
	var ent *list.Element

	if ent, ok = ram.items[key]; ok {
		return ent.Value.(*PageTable).value, true
	}

	return nil, ok
}

// Remove removes the provided key from the cache, returning if the
// key was contained.
func (ram *RAM) Remove(key interface{}) (present bool) {
	if ent, ok := ram.items[key]; ok {
		ram.removeElement(ent)

		return true
	}

	return false
}

// RemoveOldest removes the oldest item from the cache.
func (ram *RAM) RemoveOldest() (key interface{}, value interface{}, ok bool) {
	ent := ram.evictList.Back()
	if ent != nil {
		ram.removeElement(ent)
		kv := ent.Value.(*PageTable)

		return kv.key, kv.value, true
	}

	return nil, nil, false
}

// GetOldest returns the oldest entry
func (ram *RAM) GetOldest() (key interface{}, value interface{}, ok bool) {
	ent := ram.evictList.Back()
	if ent != nil {
		kv := ent.Value.(*PageTable)

		return kv.key, kv.value, true
	}

	return nil, nil, false
}

// Keys returns a slice of the keys in the cache, from oldest to newest.
func (ram *RAM) Keys() []interface{} {
	keys := make([]interface{}, len(ram.items))

	i := 0

	for ent := ram.evictList.Back(); ent != nil; ent = ent.Prev() {
		keys[i] = ent.Value.(*PageTable).key
		i++
	}

	return keys
}

// Len returns the number of items in the cache.
func (ram *RAM) Len() int {
	return ram.evictList.Len()
}

// Resize changes the cache size.
func (ram *RAM) Resize(size int) (evicted int) {
	diff := ram.Len() - size
	if diff < 0 {
		diff = 0
	}

	for i := 0; i < diff; i++ {
		ram.removeOldest()
	}

	ram.size = size

	return diff
}

// removeOldest removes the oldest item from the cache.
func (ram *RAM) removeOldest() {
	ent := ram.evictList.Back()

	if ent != nil {
		ram.removeElement(ent)
	}
}

// removeElement is used to remove a given list element from the cache
func (ram *RAM) removeElement(e *list.Element) {
	ram.evictList.Remove(e)
	kv := e.Value.(*PageTable)
	delete(ram.items, kv.key)
	if ram.onEvict != nil {
		ram.onEvict(kv.key, kv.value)
	}
}

// New creates an LRU of the given size.
func New(size int) (*Memory, error) {
	return NewWithEvict(size, nil)
}

// NewWithEvict constructs a fixed size cache with the given eviction
// callback.
func NewWithEvict(size int, onEvicted func(key interface{}, value interface{})) (*Memory, error) {
	ram, err := NewRAM(size, EvictCallback(onEvicted))
	if err != nil {
		return nil, err
	}
	c := &Memory{
		ram: ram,
	}
	return c, nil
}

// Purge is used to completely clear the cache.
func (m *Memory) Purge() {
	m.lock.Acquire()
	m.ram.Purge()
	m.lock.Release()
}

// Add adds a value to the cache.  Returns true if an eviction occurred.
func (m *Memory) Add(key, value interface{}) (evicted bool) {
	m.lock.Acquire()
	evicted = m.ram.Add(key, value)
	m.lock.Release()

	return evicted
}

// Get looks up a key's value from the cache.
func (m *Memory) Get(key interface{}) (value interface{}, ok bool) {
	m.lock.Acquire()
	value, ok = m.ram.Get(key)
	m.lock.Release()

	return value, ok
}

// Contains checks if a key is in the cache, without updating the
// recent-ness or deleting it for being stale.
func (m *Memory) Contains(key interface{}) bool {
	m.lock.Acquire()
	containKey := m.ram.Contains(key)
	m.lock.Release()
	return containKey
}

// Peek returns the key value (or undefined if not found) without updating
// the "recently used"-ness of the key.
func (m *Memory) Peek(key interface{}) (value interface{}, ok bool) {
	m.lock.Acquire()
	value, ok = m.ram.Peek(key)
	m.lock.Release()
	return value, ok
}

// ContainsOrAdd checks if a key is in the cache  without updating the
// recent-ness or deleting it for being stale,  and if not, adds the value.
// Returns whether found and whether an eviction occurred.
func (m *Memory) ContainsOrAdd(key, value interface{}) (ok, evicted bool) {
	m.lock.Acquire()
	defer m.lock.Release()

	if m.ram.Contains(key) {
		return true, false
	}
	evicted = m.ram.Add(key, value)
	return false, evicted
}

// Remove removes the provided key from the cache.
func (m *Memory) Remove(key interface{}) (present bool) {
	m.lock.Acquire()
	present = m.ram.Remove(key)
	m.lock.Release()
	return
}

// Resize changes the cache size.
func (m *Memory) Resize(size int) (evicted int) {
	m.lock.Acquire()
	evicted = m.ram.Resize(size)
	m.lock.Release()

	return evicted
}

// RemoveOldest removes the oldest item from the cache.
func (m *Memory) RemoveOldest() (key interface{}, value interface{}, ok bool) {
	m.lock.Acquire()
	key, value, ok = m.ram.RemoveOldest()
	m.lock.Release()

	return
}

// GetOldest returns the oldest entry
func (m *Memory) GetOldest() (key interface{}, value interface{}, ok bool) {
	m.lock.Acquire()
	key, value, ok = m.ram.GetOldest()
	m.lock.Acquire()

	return
}

// Keys returns a slice of the keys in the cache, from oldest to newest.
func (m *Memory) Keys() []interface{} {
	m.lock.Acquire()
	keys := m.ram.Keys()
	m.lock.Release()

	return keys
}

// Len returns the number of items in the cache.
func (m *Memory) Len() int {
	m.lock.Acquire()
	length := m.ram.Len()
	m.lock.Release()

	return length
}