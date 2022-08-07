// Author: FishGoddess
// Email: fishgoddess@qq.com
// Created at 2020/09/28 23:44:50

package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func main() {
	//getMethodDemo()
	//putEntry()
	getEntry()
	//getCacheStatus()
	//deleteEntry()
}

func getMethodDemo() {

	// 发送 GET 请求到 www.qq.com
	// 这个方法返回一个响应对象和一个错误
	// 如果错误是 nil 说明响应是可用的，不过，千万要记得关闭响应对象！
	response, err := http.Get("http://www.qq.com")
	if err != nil {
		panic(err)
	}

	// 使用 defer 关闭响应体，也相当于关闭了响应
	defer response.Body.Close()

	// 将响应打印出来看看
	fmt.Println(response)
}

func putEntry() {

	// 先创建我们的请求体
	// 这是我们要发送的数据，比如我们在缓存中存储的 value
	body := strings.NewReader("value")

	// 创建 Request 请求实例
	// 指定 PUT 请求方法，并将我们需要发送的数据添加到请求中
	request, err := http.NewRequest("PUT", "http://localhost:5837/v1/cache/key", body)
	if err != nil {
		panic(err)
	}

	// 发送请求
	// 执行完之后会得到一个响应和错误，别忘了需要关闭响应！
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()

	// 最后，将响应打印出来看看
	fmt.Println(response)
}

func getEntry() {

	// 创建请求，指定为 GET 请求方法，由于不需要请求体，所以这里置为 nil
	request, err := http.NewRequest("GET", "http://localhost:5837/v1/cache/key", nil)
	if err != nil {
		panic(err)
	}

	// 发送请求
	// 执行完之后会得到一个响应和错误，别忘了需要关闭响应！
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()

	// 这里就不仅仅要打印整个响应了，因为我们是需要响应体的内容，所以先去读取响应体数据
	// 直接使用 ReadAll 读取全部数据吧，这个 body 很小
	value, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	// 将得到的数据打印出来
	fmt.Println(string(value))
	fmt.Println(response)
}

func getCacheStatus() {

	// 创建请求，指定为 GET 请求方法，由于不需要请求体，所以这里置为 nil
	request, err := http.NewRequest("GET", "http://localhost:5837/v1/status", nil)
	if err != nil {
		panic(err)
	}

	// 发送请求
	// 执行完之后会得到一个响应和错误，别忘了需要关闭响应！
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()

	// 这里就不打印整个响应了，因为我们只是需要响应体的内容，所以只需要去读取响应体数据
	// 直接使用 ReadAll 读取全部数据吧，这个 body 很小
	value, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	// 将得到的数据打印出来
	fmt.Println(string(value))
}

func deleteEntry() {

	// 创建请求，指定为 GET 请求方法，由于不需要请求体，所以这里置为 nil
	request, err := http.NewRequest("DELETE", "http://localhost:5837/v1/cache/key", nil)
	if err != nil {
		panic(err)
	}

	// 发送请求
	// 执行完之后会得到一个响应和错误，别忘了需要关闭响应！
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()

	// 直接打印响应看看
	fmt.Println(response)
}
