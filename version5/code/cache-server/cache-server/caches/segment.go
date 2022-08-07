package caches

import (
	"errors"
	"sync"
)

// segment 就是数据块结构体。
type segment struct {

	// Data 和 Status 是导出字段，因为要通过 gob 进行持久化。

	// Data 存储这个数据块的数据。
	Data map[string]*value

	// Status 记录着这个数据块的情况。
	Status *Status

	// options 是缓存的选项设置。
	options *Options

	// lock 用于保证这个数据块的并发安全。
	lock *sync.RWMutex
}

// newSegment 返回一个使用 options 初始化过的 segment 实例。
func newSegment(options *Options) *segment {
	return &segment{
		// 初始化 map 的时候给出初始大小，可以避免大量扩容带来的性能损耗
		Data:    make(map[string]*value, options.MapSizeOfSegment),
		Status:  newStatus(),
		options: options,
		lock:    &sync.RWMutex{},
	}
}

// get 返回指定 key 的数据。
// 这个方法和原来 cache 的方法一样，只是移动到 segment 这里。
func (s *segment) get(key string) ([]byte, bool) {
	// 1. 加锁解锁
	s.lock.RLock()
	defer s.lock.RUnlock()

	// 2. 根据key获取value
	value, ok := s.Data[key]
	if !ok {
		return nil, false
	}

	// 3. 判断该value是否活着
	if !value.alive() {
		s.lock.RUnlock()
		s.delete(key)
		s.lock.RLock()
		return nil, false
	}
	return value.visit(), true
}

// set 添加一个数据进 segment。
func (s *segment) set(key string, value []byte, ttl int64) error {

	// 1. 加锁解锁
	s.lock.Lock()
	defer s.lock.Unlock()

	// 2. 判断是否已经有该数据
	if oldValue, ok := s.Data[key]; ok {
		s.Status.subEntry(key, oldValue.Data)
	}

	// 3. 判断容量是否到了设置的上限
	if !s.checkEntrySize(key, value) {
		if oldValue, ok := s.Data[key]; ok {
			s.Status.addEntry(key, oldValue.Data)
		}
		return errors.New("the entry size will exceed if you set this entry")
	}

	s.Status.addEntry(key, value)
	s.Data[key] = newValue(value, ttl)
	return nil
}

// delete 从 segment 中删除指定 key 的数据。
func (s *segment) delete(key string) {
	s.lock.Lock()
	defer s.lock.Unlock()
	if oldValue, ok := s.Data[key]; ok {
		s.Status.subEntry(key, oldValue.Data)
		delete(s.Data, key)
	}
}

// Status 返回这个 segment 的情况。
func (s *segment) status() Status {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return *s.Status
}

// checkEntrySize 会判断数据容量是否已经达到了设定的上限。
// 因为这个配置是针对整个缓存的，而这边判断大小是针对单个 segment 的，所以需要算出单个 segment 的上限来判断。
func (s *segment) checkEntrySize(newKey string, newValue []byte) bool {
	return s.Status.entrySize()+int64(len(newKey))+int64(len(newValue)) <= int64((s.options.MaxEntrySize*1024*1024)/s.options.SegmentSize)
}

// gc 会清理 segment 中过期的数据。
func (s *segment) gc() {

	// 1. 加锁解锁
	s.lock.Lock()
	defer s.lock.Unlock()

	// 2. 循环遍历value中的数据是否过期
	count := 0
	for key, value := range s.Data {
		if !value.alive() {
			s.Status.subEntry(key, value.Data)
			delete(s.Data, key)
			count++
			if count >= s.options.MaxGcCount {
				break
			}
		}
	}
}
