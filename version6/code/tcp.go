// Author: FishGoddess
// Email: fishgoddess@qq.com
// Created at 2020/10/18 00:58:30

package main

import (
	"fmt"
	"net"
	"strings"
	"sync"
	"time"
)

func main() {
	// 先启动服务器，然后启动客户端
	go runServer()
	time.Sleep(time.Second)
	runClient()
}

func runServer() {

	// 监听 tcp 网络，端口使用 5837
	listener, err := net.Listen("tcp", ":5837")
	if err != nil {
		panic(err)
	}
	// 最后记得关闭监听器
	defer listener.Close()

	// 使用 WaitGroup 记录着所有的连接，等待所有连接处理完
	wg := &sync.WaitGroup{}
	for {
		// 等待客户端连接
		conn, err := listener.Accept()
		if err != nil {
			// 如果错误是 ErrNetClosing，就说明这个 listener 已经关闭了
			// 具体的说明可以参考源码 src/internal/poll/fd.go 的 ErrNetClosing 说明
			if strings.Contains(err.Error(), "use of closed network connection") {
				break
			}
			continue
		}

		// 记录一个连接，然后开启一个 goroutine 去处理这个连接
		wg.Add(1)
		go func() {
			defer wg.Done()
			handleConn(conn)
		}()
	}
	wg.Wait()
}

// 处理连接
func handleConn(conn net.Conn) {

	defer conn.Close()

	// 读取数据，这里使用了 1024 个字节的内存存储
	// 正常应该是一个循环去读，知道读取完毕
	// 这里就简单演示一下
	request := make([]byte, 1024)
	n, err := conn.Read(request)
	if err != nil {
		return
	}

	// 收到客户端数据之后原样发送回去
	fmt.Println("服务端收到消息：", string(request[:n]))
	_, err = conn.Write(request[:n])
	if err != nil {
		return
	}
}

func runClient() {

	// 连接到服务端
	conn, err := net.Dial("tcp", "127.0.0.1:5837")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	// 向服务端发送数据
	_, err = conn.Write([]byte("客户端请求..."))
	if err != nil {
		panic(err)
	}

	// 读取服务端发送过来的数据
	response := make([]byte, 1024)
	n, err := conn.Read(response)
	if err != nil {
		return
	}
	fmt.Println("客户端收到消息：", string(response[:n]))
}
