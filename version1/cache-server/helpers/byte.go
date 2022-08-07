/*
 * @Author       : STY
 * @Date         : 2020-09-13 15:07:33
 * @LastEditors  : STY
 * @LastEditTime : 2022-08-01 08:27:45
 * @FilePath     : \cache-server\helpers\byte.go
 * @Description  :
 * Copyright 2022 OBKoro1, All Rights Reserved.
 * 2020-09-13 15:07:33
 */
package helpers

// Copy 复制 src 到新的 []byte 中并返回。
func Copy(src []byte) []byte {

	// 1. 声明一个字符切片
	dst := make([]byte, len(src))

	// 2. 复制，避免引用类型值传递，不安全
	copy(dst, src)
	return dst
}
