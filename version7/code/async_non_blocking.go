// Author: FishGoddess
// Email: fishgoddess@qq.com
// Created at 2020/10/22 01:03:56

package main

import (
	"fmt"
	"time"
)

func main() {

	// 方法的两个参数
	a := 1
	b := 2

	// 调用方法的时候加上回调函数
	// 这个回调函数会在得到结果之后执行
	addWithCallback(a, b, func(sum int) {
		fmt.Println(sum)
	})

	// 防止 main goroutine 比异步任务的 goroutine 先退出
	time.Sleep(time.Second)
}

func addWithCallback(a int, b int, callback func(sum int)) {
	go func() {
		// 在新的 goroutine 中计算结果，并将结果传递给回调函数
		sum := a + b
		callback(sum)
	}()
}
