package caches

// Status 是一个代表缓存信息的结构体。
// 因为这个结构体需要被序列化为 Json 字符串并通过网络传输，所以这里使用到了 Json 的标签。
type Status struct {

	// Count 记录着缓存中的数据个数。
	Count int `json:"count"`

	// KeySize 记录着 key 占用的空间大小。
	KeySize int64 `json:"keySize"`

	// ValueSize 记录着 value 占用的空间大小。
	ValueSize int64 `json:"valueSize"`
}

// newStatus 返回一个缓存信息对象指针。
func NewStatus() *Status {
	return &Status{
		Count:     0,
		KeySize:   0,
		ValueSize: 0,
	}
}

// addEntry 可以将 key 和 value 的信息记录起来。
func (s *Status) addEntry(key string, value []byte) {
    // 每添加一个键值对，count 就需要加 1，key 占用的空间就是 string 的长度。
    // 同理，value 占用的空间就是切片的长度。
	s.Count++
	s.KeySize += int64(len(key))
	s.ValueSize += int64(len(value))
}

// subEntry 可以将 key 和 value 的信息从 Status 中减去。
func (s *Status) subEntry(key string, value []byte) {
    // 每减少一个键值对，count 就需要减 1，key 和 value 占用的空间也需要减去相应的大小。
	s.Count--
	s.KeySize -= int64(len(key))
	s.ValueSize -= int64(len(value))
}

// entrySize 返回键值对占用的总大小。
func (s *Status) entrySize() int64 {
	return s.KeySize + s.ValueSize
}
