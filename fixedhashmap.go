package fixedhashmap

import "hash/fnv"

// FixedHashMap is a hashmap with a fixed size.
type FixedHashMap interface {
	// Set stores the given key/value pair in the hash map. Returns a boolean value indicating success / failure of the operation.
	Set(key string, value interface{}) bool

	// Get returns the value associated with the given key, or null if no value is set.
	Get(key string) interface{}

	// Delete deletes the value associated with the given key, returning the value on success or null if the key has no value.
	Delete(key string) interface{}

	// Load returns a float value representing the load factor (`(items in hash map)/(size of hash map)`) of the data structure. Since the size of the dat structure is fixed, this should never be greater than 1.
	Load() float64
}

type fixedHashMap struct {
	els      []*element
	n        uint32 // number of els stored
	capacity uint32
}

type element struct {
	key     string
	value   interface{}
	deleted bool
}

// New creates a new FixedHashMap.
func New(capacity uint32) FixedHashMap {
	return &fixedHashMap{
		els:      make([]*element, capacity),
		capacity: capacity,
	}
}

func (m *fixedHashMap) Set(key string, value interface{}) bool {
	// check max load
	if m.n == m.capacity {
		return false
	}

	// probe
	pos := m.find(key)
	if pos == -1 {
		// no space for this key. should not be reached because we check for max load
		return false
	}

	el := m.els[pos]
	if el != nil {
		// replace value
		el.value = value
		el.deleted = false
		return true
	}

	// otherwise we are setting a new element
	m.els[pos] = &element{
		key:   key,
		value: value,
	}
	m.n++
	return true
}

func (m *fixedHashMap) Get(key string) interface{} {
	// probe
	pos := m.find(key)
	if pos == -1 {
		// does not exist
		return nil
	}

	// check if probe was not successful
	el := m.els[pos]
	if el == nil {
		return nil
	}

	// check if we have a deleted flag
	if el.deleted {
		return nil
	}

	// return found probe
	return el.value
}

func (m *fixedHashMap) Delete(key string) interface{} {
	// probe
	pos := m.find(key)
	// could not find
	if pos == -1 {
		return nil
	}

	el := m.els[pos]
	// probe terminated at empty cell
	if el == nil {
		return nil
	}

	// delete
	el.deleted = true
	// garbage collect the value
	el.value = nil
	m.n--
	return el.value
}

func (m *fixedHashMap) Load() float64 {
	return float64(m.n) / float64(m.capacity)
}

// find finds the position key should be in the map.
func (m *fixedHashMap) find(key string) int {
	h := hash(key)
	idx := h % m.capacity
	start := idx

	// Look for first empty probe
	for m.els[idx] != nil {
		// Check if the key is the same
		if m.els[idx].key == key {
			return int(idx)
		}

		// linear probing
		idx++

		// wraparound
		idx %= m.capacity

		// ensure we don't probe the same index (cyclic)
		if idx == start {
			// key is not present and there are no empty values in the list
			return -1
		}
	}

	return int(idx)
}

func hash(s string) uint32 {
	// fnv1a is a reasonably fast (non-cryptographic) uniformly distributed hash function.
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}
