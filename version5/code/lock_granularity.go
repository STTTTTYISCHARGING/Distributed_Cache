// Author: FishGoddess
// Email: fishgoddess@qq.com
// Created at 2020/10/11 14:10:30

package code

import (
	"fmt"
	"sync"
	"time"
)

func main() {

	//bigGranularity()
	//smallGranularity()
}

func bigGranularity() {
	count := 0

	lock := sync.Mutex{}
	lock.Lock() // 上锁
	fmt.Println(time.Now())
	count++
	fmt.Println(time.Now())
	lock.Unlock() // 解锁
}

func smallGranularity() {
	count := 0

	lock := sync.Mutex{}
	fmt.Println(time.Now())
	lock.Lock() // 上锁
	count++
	lock.Unlock() // 解锁
	fmt.Println(time.Now())
}
