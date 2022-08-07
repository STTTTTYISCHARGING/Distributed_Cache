// Author: FishGoddess
// Email: fishgoddess@qq.com
// Created at 2020/09/28 00:00:31

package main

import "testing"

// go test -v add.go add_test.go
// 首先，测试文件的名字必须以 _test 结尾，并且和被测试文件在同一个目录下。
// 测试方法的命名也有讲究，必须以 Test 开头，并且参数是 *testing.T。
func TestAdd(t *testing.T) {

	// 测试案例
	a := 1
	b := 2
	sum := 3

	result := Add(a, b)
	if result != sum {
		// t.Fatalf 方法表示测试出错，并立即退出测试
		// t 中还有多个类似的方法，用法都差不多，只是含义不同
		t.Fatalf("The result should be %d, but it returns %d!", sum, result)
	}
}

// go test add.go add_test.go -bench=.
func BenchmarkAdd(b *testing.B) {

	//b.ReportAllocs() // 显示内存情况，想要显示就取消注释
	b.ResetTimer() // 重置计时器，一般在测试之前执行将计时器清零，测试的时间更加准确

	// 注意这个 b.N 是内置的数字，并不是固定的大小，而是会根据测试情况进行变动的数字
	for i := 0; i < b.N; i++ {
		Add(1, 1)
	}
}
