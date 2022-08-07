/*
 * @Author       : STY
 * @Date         : 2020-09-13 15:07:33
 * @LastEditors  : STY
 * @LastEditTime : 2022-08-01 08:26:18
 * @FilePath     : \cache-server\servers\http.go
 * @Description  :
 * Copyright 2022 OBKoro1, All Rights Reserved.
 * 2020-09-13 15:07:33
 */
package servers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"cache-server/caches"

	"github.com/julienschmidt/httprouter"
)

// HTTPServer 是 HTTP 服务器结构。
type HTTPServer struct {
	// Cache 是底层存储的结构。
	cache *caches.Cache
}

// NewHTTPServer 返回一个关于 cache 的新 HTTP 服务器。
func NewHTTPServer(cache *caches.Cache) *HTTPServer {
	return &HTTPServer{
		cache: cache,
	}
}

// Run 启动服务器。
func (hs *HTTPServer) Run(address string) error {
	return http.ListenAndServe(address, hs.routerHandler())
}

// =======================================================================

// routerHandler 返回注册的路由处理器。
// key 都从 url 上获取，value 从请求体中获取
func (hs *HTTPServer) routerHandler() http.Handler {

	// httprouter.New() 创建一个 http 路由组件，包括各种请求方法的路由
	router := httprouter.New()

	//GET 请求方法就用于缓存的查询
	router.GET("/cache/:key", hs.getHandler)
	router.GET("/status", hs.statusHandler)

	//PUT 请求方法就用于缓存的新建
	router.PUT("/cache/:key", hs.setHandler)

	//DELETE 请求就用于缓存的删除
	router.DELETE("/cache/:key", hs.deleteHandler)

	return router
}

// getHandler 用于获取缓存数据。
// httprouter.Params应该是将url处理成结构体
func (hs *HTTPServer) getHandler(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {

	// 1. 从url中获取key
	key := params.ByName("key")

	// 2. 从内存数据库中获取该key对应的value
	value, ok := hs.cache.Get(key)
	if !ok {
		// 如果缓存中找不到数据，就返回 404 的状态码
		writer.WriteHeader(http.StatusNotFound)
		return
	}

	// 3. 将value写到发送部分
	writer.Write(value)
}

// setHandler 用于保存缓存数据。
func (hs *HTTPServer) setHandler(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {

	// 1. 从url中获取key
	key := params.ByName("key")

	// 2. 从请求体中获取value，为了简单操作，整个请求体都被当作 value
	value, err := ioutil.ReadAll(request.Body)
	if err != nil {
		// 如果读取请求体失败，就返回 500 的状态码
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	// 3. 将k-v存入内存数据库
	hs.cache.Set(key, value)
}

// deleteHandler 用于删除缓存数据。
func (hs *HTTPServer) deleteHandler(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {

	// 1. 从url中获取key
	key := params.ByName("key")

	// 2. 处理业务
	hs.cache.Delete(key)
}

// statusHandler 用于获取缓存键值对的个数
func (hs *HTTPServer) statusHandler(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {

	// 1. 将个数编码成 JSON 字符串
	status, err := json.Marshal(map[string]interface{}{
		"count": hs.cache.Count(),
	})
	if err != nil {
		// 如果编码失败，就返回 500 的状态码
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	// 2. 写到发送部分
	writer.Write(status)
}
