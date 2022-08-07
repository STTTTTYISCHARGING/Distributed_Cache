// Author: FishGoddess
// Email: fishgoddess@qq.com
// Created at 2020/10/01 01:01:14

package main

import (
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
)

func main() {

	// 连接 redis 服务
	// 这边使用 redis 协议链接的方式进行连接
	conn, err := redis.DialURL("redis://192.168.32.128:6379")
	if err != nil {
		panic(err)
	}

	// 注意一定要关闭连接释放资源
	defer conn.Close()

	// 添加一个键值对，key 和 value，这个键值对是不会过期的
	reply, err := redis.String(conn.Do("set", "key", "value"))
	fmt.Println(reply, err)

	// 添加一个键值对，key1 和 value1，这个键值对的有效期是 3 秒
	reply, err = redis.String(conn.Do("set", "key1", "value1", "ex", 3))
	fmt.Println(reply, err)

	// 查询添加的两个键值对
	reply, err = redis.String(conn.Do("get", "key"))
	fmt.Println(reply, err)

	reply, err = redis.String(conn.Do("get", "key1"))
	fmt.Println(reply, err)

	time.Sleep(3 * time.Second)

	// 等待 3 秒之后再查看一次，发现 key 还健在，但是 key1 已经过期了
	reply, err = redis.String(conn.Do("get", "key"))
	fmt.Println(reply, err)

	reply, err = redis.String(conn.Do("get", "key1"))
	fmt.Println(reply, err)

	// 改变一个键值对的值
	reply, err = redis.String(conn.Do("set", "key", "newValue"))
	fmt.Println(reply, err)

	// 查看是否被修改
	reply, err = redis.String(conn.Do("get", "key"))
	fmt.Println(reply, err)

	// 删除一个 key
	n, err := redis.Int64(conn.Do("del", "key"))
	fmt.Println(n, err)

	// 查看是否被删除
	reply, err = redis.String(conn.Do("get", "key"))
	fmt.Println(reply, err)
}
