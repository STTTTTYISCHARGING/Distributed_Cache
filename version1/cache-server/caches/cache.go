/*
 * @Author       : STY
 * @Date         : 2020-09-13 15:07:33
 * @LastEditors  : STY
 * @LastEditTime : 2022-08-01 08:30:54
 * @FilePath     : \cache-server\caches\cache.go
 * @Description  :
 * Copyright 2022 OBKoro1, All Rights Reserved.
 * 2020-09-13 15:07:33
 */
package caches

import (
	"sync"

	"cache-server/helpers"
)

// Cache 是一个结构体，用于封装缓存底层结构的。如果有任何扩充的需求，就改写这里的数据结构。可以替换成LRU，后面应该会改。
type Cache struct {

	// 1. data 是一个 map，存储了所有的数据。value类型使用的是 []byte，这是为了方便网络传输。
	data map[string][]byte

	// 2. count 记录着 data 中键值对的个数。
	// 冗余设计，使用len(data)也可以得到键值对的个数。使用 count 记录是为了更快地得到结果。
	count int64

	// 3. lock 用于保证并发安全。
	lock *sync.RWMutex
}

// NewCache 返回一个缓存对象。Go语言中自己创造构造函数。
func NewCache() *Cache {
	return &Cache{
		// 在创建 map 的时候指定了 256 的大小
		// 这个 256 不是说只能存 256 个 键值对，而是预先分配 256 个槽位
		// 尽可能少的避免后面因为容量不足导致 map 扩容，扩容会分配内存，影响性能
		// 还有一点就是槽位少了，哈希冲突的几率就大了，map 查找的性能就会下降
		// 这里的 256 也不一定是最佳值，需要根据实际情况而定
		data:  make(map[string][]byte, 256),
		count: 0,
		lock:  &sync.RWMutex{},
	}
}

// Get 返回指定 key 的 value，如果找不到就返回 false。
func (c *Cache) Get(key string) ([]byte, bool) {
	// 由于查询数据不会改变数据的状态，所以可以并发执行，这里就使用读锁，加快读取速度
	c.lock.RLock()
	defer c.lock.RUnlock()
	value, ok := c.data[key]
	return value, ok
}

// Set 保存 key 和 value 到缓存中。
func (c *Cache) Set(key string, value []byte) {
	// Set 操作会改变数据的状态，需要保证串行执行，这里使用写锁
	c.lock.Lock()
	defer c.lock.Unlock()
	_, ok := c.data[key]
	if !ok {
		c.count++
	}
	// 这个 Copy 方法会将 value 拷贝一份出来
	// 这样即使传进来的 value 被修改或者清空了也不会影响缓存里面的数据
	c.data[key] = helpers.Copy(value)
}

// Delete 删除指定 key 的键值对数据。
func (c *Cache) Delete(key string) {
	// Delete 操作会改变数据的状态，需要保证串行执行，这里使用写锁
	c.lock.Lock()
	defer c.lock.Unlock()
	_, ok := c.data[key]
	if ok {
		c.count--
		delete(c.data, key)
	}
}

// Count 返回键值对数据的个数。
func (c *Cache) Count() int64 {
	// 由于查询数据个数不会改变数据的状态，所以可以并发执行，这里就使用读锁，加快读取速度
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.count
}
