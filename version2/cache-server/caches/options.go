/*
 * @Author       : STY
 * @Date         : 2020-09-20 17:27:19
 * @LastEditors  : STY
 * @LastEditTime : 2022-08-01 09:09:02
 * @FilePath     : \cache-server\caches\options.go
 * @Description  :
 * Copyright 2022 OBKoro1, All Rights Reserved.
 * 2020-09-20 17:27:19
 */
package caches

// Options 是一些选项的结构体。实现写满保护机制
// 这边的一些成员是后面的功能需要的，我们可以先不管它们，后面功能涉及到再回来看。
type Options struct {

	// MaxEntrySize 是写满保护的一个阈值，当缓存中的键值对占用空间达到这个值，就会触发写满保护。
	// 这个值的单位是 GB。
	MaxEntrySize int64

	// MaxGcCount 是自动淘汰机制的一个阈值，当清理的数据达到了这个值后就会停止清理了。
	MaxGcCount int

	// GcDuration 是自动淘汰机制的时间间隔，每隔固定的 GcDuration 时间就会进行一次自动淘汰。
	// 这个值的单位是分钟。
	GcDuration int64
}

// DefaultOptions 返回一个默认的选项设置对象。
func DefaultOptions() Options {
	return Options{
		MaxEntrySize: int64(4), // 默认是 4 GB
		MaxGcCount:   1000,     // 默认是 1000 个
		GcDuration:   60,       // 默认是 1 小时
	}
}
