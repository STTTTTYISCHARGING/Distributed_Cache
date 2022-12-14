package main

import (
	"flag"

	"cache-server/caches"
	"cache-server/servers"
)

func main() {

	// 解析所有的 flag
    // 因为我们加入了一些用户配置，options 还记得吗？所以我们需要给用户一个入口去修改这些配置
    // 这是修改监听地址的参数，默认是 5837 端口
	address := flag.String("address", ":5837", "The address used to listen, such as 127.0.0.1:5837.")

    // 创建一个默认配置
    // 下面几个 flag 就是修改对应的配置参数用的
	options := caches.DefaultOptions()
	flag.Int64Var(&options.MaxEntrySize, "maxEntrySize", options.MaxEntrySize, "The max memory size that entries can use. The unit is GB.")
	flag.IntVar(&options.MaxGcCount, "maxGcCount", options.MaxGcCount, "The max count of entries that gc will clean.")
	flag.Int64Var(&options.GcDuration, "gcDuration", options.GcDuration, "The duration between two gc tasks. The unit is Minute.")
    
    // 千万要记得调用这个方法，否则参数不会被解析
	flag.Parse()

	// 初始化缓存和服务器
	cache := caches.NewCacheWith(options)
	cache.AutoGc()
	err := servers.NewHTTPServer(cache).Run(*address)
	if err != nil {
		panic(err)
	}
}
