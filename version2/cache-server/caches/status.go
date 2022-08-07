/*
 * @Author       : STY
 * @Date         : 2020-09-20 17:26:33
 * @LastEditors  : STY
 * @LastEditTime : 2022-08-01 09:11:18
 * @FilePath     : \cache-server\caches\status.go
 * @Description  :
 * Copyright 2022 OBKoro1, All Rights Reserved.
 * 2020-09-20 17:26:33
 */
package caches

// Status是一个代表缓存信息的结构体。
// Status 结构体将会在 Cache 结构体中使用，所以没有使用并发安全的机制，Cache 结构体带锁
type Status struct {

	// 结构体需要被序列化为Json字符串并通过网络传输，所以这里用了Json的标签。

	// Count 记录着缓存中的数据个数。
	Count int `json:"count"`

	// KeySize 记录着 key 占用的空间大小。
	KeySize int64 `json:"keySize"`

	// ValueSize 记录着 value 占用的空间大小。
	ValueSize int64 `json:"valueSize"`
}

// newStatus 返回一个缓存信息对象指针。
func newStatus() *Status {
	return &Status{
		Count:     0,
		KeySize:   0,
		ValueSize: 0,
	}
}

// addEntry 可以将 key 和 value 的信息记录起来。
func (s *Status) addEntry(key string, value []byte) {

	// 每添加一个键值对，count 就需要加 1
	s.Count++

	// key 占用的空间就是 string 的长度。
	s.KeySize += int64(len(key))

	// value 占用的空间就是切片的长度。
	s.ValueSize += int64(len(value))
}

// subEntry 可以将 key 和 value 的信息从 Status 中减去。
func (s *Status) subEntry(key string, value []byte) {

	// 每减少一个键值对，count 就需要减 1
	s.Count--

	// key 和 value 占用的空间也需要减去相应的大小。
	s.KeySize -= int64(len(key))
	s.ValueSize -= int64(len(value))
}

// entrySize 返回键值对占用的总大小。
func (s *Status) entrySize() int64 {
	return s.KeySize + s.ValueSize
}
