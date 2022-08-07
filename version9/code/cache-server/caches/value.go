package caches

import (
	"sync/atomic"
	"time"

	"cache-server/helpers"
)

const (
	// NeverDie 是一个常量，我们设计的时候规定如果 ttl 为 0，那就是永不过期，相当于灵丹妙药。
	NeverDie = 0
)

// value 是一个包装了数据的结构体。
type value struct {

	// data 存储着真正的数据。
	Data []byte

	// ttl 代表这个数据的寿命。
	// 这个值的单位是秒。
	Ttl int64

	// ctime 代表这个数据的创建时间。
	Ctime int64
}

// newValue 返回一个包装之后的数据。
func newValue(data []byte, ttl int64) *value {
	return &value{
        // 这里使用 Copy 是为了让这个数据和外界没有任何联系
        // 实际上也可以不复制，性能还更高，只要保证外界不会更改就可以了
        // 这里我们进行复制是因为从设计上来说最好这么干，但是后面为了性能可能就会去除掉这个复制步骤了
        // 所以从这也能看出，有时候需要做一些反设计的设计
		Data:  helpers.Copy(data),
		Ttl:   ttl,
		Ctime: time.Now().Unix(),
	}
}

// alive 返回这个数据是否存活。
func (v *value) alive() bool {
    // 首先判断是否有过期时间，然后判断当前时间是否超过了这个数据的死期
	return v.Ttl == NeverDie || time.Now().Unix()-v.Ctime < v.Ttl
}

// visit 返回这个数据的实际存储数据。
func (v *value) visit() []byte {
    // 这一步是为了实现 LRU 过期机制而加的
    // 在访问数据的时候，将创建时间更新为访问时间，这样就相当于最近访问的数据过期时间会延长
    // 因为获取数据一般都在读取操作中进行，读取操作使用的是读锁，尽可能保证并发的性能
    // 使用读锁就意味着没有保证写的并发安全，所以我们需要自己去处理并发安全的问题
    // 一般会使用锁来处理，但是这里牵扯到读取性能，如果使用锁，就会显得非常臃肿和拉跨
    // 于是就使用了 atomic 包轻量化地去处理，这里直接使用了交换更新数据，而没有使用 CAS 的方式
    // 后交换成功的会把先交换成功的时间改掉，所以这里不保证交换的时间一定是更加新的时间
    // 有兴趣的童鞋可以尝试使用 CAS 的方式去更新，注意 CAS 的重试次数限制，防止高并发的时候 CPU 浪费严重
	atomic.SwapInt64(&v.Ctime, time.Now().Unix())
	return v.Data
}
