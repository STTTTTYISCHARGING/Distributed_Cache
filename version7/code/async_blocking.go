// Author: FishGoddess
// Email: fishgoddess@qq.com
// Created at 2020/10/22 00:57:15

package main

import "fmt"

func main() {

	// 方法的两个参数
	a := 1
	b := 2

	// 从管道中接收结果，这一步是阻塞的，因为在等待结果的产出
	sum := <-addAsync(a, b)
	fmt.Println(sum)
}

func addAsync(a int, b int) <-chan int {
	// 使用管道接收结果，注意需要设置一个缓冲位，否则没有取结果的话这个 goroutine 会被阻塞
	resultChan := make(chan int, 1)
	go func() {
		// 在新的 goroutine 中计算结果，并将结果发送到管道
		resultChan <- a + b
	}()
	return resultChan
}
