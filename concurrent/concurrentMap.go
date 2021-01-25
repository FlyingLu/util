package concurrent

import (
	"mygo/util"
	"sync"
)

type KeyValuePair struct {
	Key, Value interface{}
}

type hashTable map[interface{}]interface{}

type segment struct {
	hashTable
	sync.RWMutex
}

func (s *segment) put(key interface{}, value interface{}) {
	s.Lock()
	defer s.Unlock()
	s.hashTable[key] = value
}

func (s *segment) get(key interface{}) (value interface{}, ok bool) {
	s.RLock()
	defer s.RUnlock()
	value, ok = s.hashTable[key]
	return
}

func (s *segment) replace(key, value, newValue interface{}) (ok bool) {
	s.Lock()
	defer s.Unlock()
	previewsValue, previewsExists := s.hashTable[key]
	if previewsExists && previewsValue == value {
		s.hashTable[key] = newValue
		return true
	}
	return false
}

func (s *segment) putIfAbsent(key interface{}, value interface{}) (previousValue interface{}, previousExists bool) {
	s.Lock()
	defer s.Unlock()
	previousValue, previousExists = s.hashTable[key]
	if !previousExists {
		s.hashTable[key] = value
	}
	return
}

func (s *segment) putIfPresent(key interface{}, value interface{}) (previousValue interface{}, previousExists bool) {
	s.Lock()
	defer s.Unlock()
	previousValue, previousExists = s.hashTable[key]
	if previousExists {
		s.hashTable[key] = value
	}
	return
}

func (s *segment) remove(key interface{}) (previousValue interface{}, previousExists bool) {
	s.Lock()
	defer s.Unlock()
	previousValue, previousExists = s.hashTable[key]
	delete(s.hashTable, key)
	return
}

func (s *segment) removeKeyValuePair(key, value interface{}) (ok bool) {
	s.Lock()
	defer s.Unlock()
	previousValue, previousExists := s.hashTable[key]
	if previousExists && previousValue == value {
		delete(s.hashTable, key)
		return true
	}
	return false
}

type Map struct {
	segments []*segment
}

func NewMap() *Map {
	segments := make([]*segment, 8)
	for i := range segments {
		segments[i] = new(segment)
	}
	return &Map{
		segments: segments,
	}
}

func (m *Map) Size() int {
	for _, seg := range m.segments {
		seg.RLock()
		defer seg.RUnlock()
	}
	sum := 0
	for _, seg := range m.segments {
		sum += len(seg.hashTable)
	}
	return sum
}

func (m *Map) IsEmpty() bool {
	return m.Size() == 0
}

func (m *Map) ContainsKey(key interface{}) bool {
	h := util.Hash(key)
	s := m.segments[h%uint32(len(m.segments))]
	_, ok := s.get(key)
	return ok
}

func (m *Map) ContainsValue(value interface{}) bool {
	for _, seg := range m.segments {
		seg.RLock()
		defer seg.RUnlock()
	}
	for _, seg := range m.segments {
		for _, v := range seg.hashTable {
			if v == value {
				return true
			}
		}
	}
	return false
}

func (m *Map) Get(key interface{}) (value interface{}, ok bool) {
	h := util.Hash(key)
	s := m.segments[h%uint32(len(m.segments))]
	value, ok = s.get(key)
	return
}

func (m *Map) Put(key, value interface{}) {
	h := util.Hash(key)
	s := m.segments[h%uint32(len(m.segments))]
	s.put(key, value)
}

func (m *Map) Remove(key interface{}) (previousValue interface{}, previousExists bool) {
	h := util.Hash(key)
	s := m.segments[h%uint32(len(m.segments))]
	return s.remove(key)
}

func (m *Map) PutAll(m2 *Map) {
	for _, seg := range m.segments {
		seg.Lock()
		defer seg.Unlock()
	}
	for _, seg := range m2.segments {
		seg.RLock()
		defer seg.RUnlock()
	}
	for _, seg := range m2.segments {
		for k, v := range seg.hashTable {
			h := util.Hash(k)
			m2.segments[h%uint32(len(m.segments))].hashTable[k] = v
		}
	}
}

func (m *Map) Clear() {
	for _, seg := range m.segments {
		seg.Lock()
		defer seg.Unlock()
	}
	for _, seg := range m.segments {
		seg.hashTable = make(hashTable)
	}
}

func (m *Map) Keys() []interface{} {
	for _, seg := range m.segments {
		seg.RLock()
		defer seg.RUnlock()
	}
	var k []interface{}
	for _, seg := range m.segments {
		for key := range seg.hashTable {
			k = append(k, key)
		}
	}
	return k
}

func (m *Map) Values() []interface{} {
	for _, seg := range m.segments {
		seg.RLock()
		defer seg.RUnlock()
	}
	var v []interface{}
	for _, seg := range m.segments {
		for _, value := range seg.hashTable {
			v = append(v, value)
		}
	}
	return v
}

func (m *Map) KeyValuePairs() []*KeyValuePair {
	for _, seg := range m.segments {
		seg.RLock()
		defer seg.RUnlock()
	}
	var p []*KeyValuePair
	for _, seg := range m.segments {
		for key, value := range seg.hashTable {
			p = append(p, &KeyValuePair{Key: key, Value: value})
		}
	}
	return p
}

func (m *Map) GetOrDefault(key, defaultValue interface{}) (value interface{}) {
	h := util.Hash(key)
	s := m.segments[h%uint32(len(m.segments))]
	value, ok := s.get(key)
	if !ok {
		value = defaultValue
	}
	return
}

func (m *Map) ForEach(action func(key, value interface{})) {
	for _, seg := range m.segments {
		seg.RLock()
		defer seg.RUnlock()
	}
	for _, seg := range m.segments {
		for key, value := range seg.hashTable {
			action(key, value)
		}
	}
}

func (m *Map) ReplaceAll(function func(key, value interface{}) interface{}) {
	for _, seg := range m.segments {
		seg.Lock()
		defer seg.Unlock()
	}
	for _, seg := range m.segments {
		for key, value := range seg.hashTable {
			newValue := function(key, value)
			seg.hashTable[key] = newValue
		}
	}
}

func (m *Map) PutIfAbsent(key, value interface{}) (previousValue interface{}, previousExists bool) {
	h := util.Hash(key)
	s := m.segments[h%uint32(len(m.segments))]
	return s.putIfAbsent(key, value)
}

func (m *Map) PutIfPresent(key, value interface{}) (previousValue interface{}, previousExists bool) {
	h := util.Hash(key)
	s := m.segments[h%uint32(len(m.segments))]
	return s.putIfPresent(key, value)
}

func (m *Map) RemoveKeyValuePair(key, value interface{}) (ok bool) {
	h := util.Hash(key)
	s := m.segments[h%uint32(len(m.segments))]
	return s.removeKeyValuePair(key, value)
}

func (m *Map) Replace(key, value, newValue interface{}) (ok bool) {
	h := util.Hash(key)
	s := m.segments[h%uint32(len(m.segments))]
	return s.replace(key, value, newValue)
}

func (m *Map) ComputeIfAbsent(key interface{}, mappingFunction func(key interface{}) (value interface{})) {
	h := util.Hash(key)
	s := m.segments[h%uint32(len(m.segments))]
	s.Lock()
	defer s.Unlock()
	previewsValue, previewsOk := s.hashTable[key]
	if !previewsOk || previewsValue == nil {
		newValue := mappingFunction(key)
		if newValue != nil {
			s.hashTable[key] = newValue
		}
	}
}

func (m *Map) ComputeIfPresent(key interface{}, mappingFunction func(key interface{}) (value interface{})) {
	h := util.Hash(key)
	s := m.segments[h%uint32(len(m.segments))]
	s.Lock()
	defer s.Unlock()
	previewsValue, previewsOk := s.hashTable[key]
	if previewsOk && previewsValue != nil {
		newValue := mappingFunction(key)
		if newValue != nil {
			s.hashTable[key] = newValue
		} else if previewsOk {
			delete(s.hashTable, key)
		}
	}
}

func (m *Map) Compute(key interface{}, mappingFunction func(key interface{}) (value interface{})) {
	h := util.Hash(key)
	s := m.segments[h%uint32(len(m.segments))]
	s.Lock()
	defer s.Unlock()
	newValue := mappingFunction(key)
	if newValue == nil {
		_, previewsOk := s.hashTable[key]
		if previewsOk {
			delete(s.hashTable, key)
		}
	} else {
		s.hashTable[key] = newValue
	}
}

func (m *Map) Merge(key, defaultValue interface{}, remappingFunction func(key, valueOrDefaultValue interface{}) (newValue interface{})) {
	h := util.Hash(key)
	s := m.segments[h%uint32(len(m.segments))]
	s.Lock()
	defer s.Unlock()
	previewsValue, previewsOk := s.hashTable[key]
	if !previewsOk || previewsValue == nil {
		s.hashTable[key] = defaultValue
		return
	}
	newValue := remappingFunction(key, previewsValue)
	if newValue == nil {
		if previewsOk {
			delete(s.hashTable, key)
		}
	} else {
		s.hashTable[key] = newValue
	}
}
