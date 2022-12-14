package servers

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"

	"cache-server/caches"
	"cache-server/helpers"
	"github.com/FishGoddess/vex"
)

const (
	// getCommand 是 get 命令。
	getCommand = byte(1)

	// setCommand 是 set 命令。
	setCommand = byte(2)

	// deleteCommand 是 delete 命令。
	deleteCommand = byte(3)

	// statusCommand 是 status 命令。
	statusCommand = byte(4)

	// nodesCommand 是 nodes 命令。
	nodesCommand = byte(5)
)

var (
	// commandNeedsMoreArgumentsErr 是命令需要更多参数的错误。
	commandNeedsMoreArgumentsErr = errors.New("command needs more arguments")

	// notFoundErr 是找不到的错误。
	notFoundErr = errors.New("not found")
)

// TCPServer 是 TCP 类型的服务器。
type TCPServer struct {

	// node 是内部用于记录集群信息的实例。
    // 这里使用了 Go 语言的组合模式，这样当前服务器实例也可以说成是节点实例，方法都可以互通。
	*node

	// cache 是内部用于存储数据的缓存组件。
	cache *caches.Cache

	// server 是内部真正用于服务的服务器。
	server *vex.Server

	// options 存储着这个服务器的选项配置。
	options *Options
}

// NewTCPServer 返回新的 TCP 服务器。
func NewTCPServer(cache *caches.Cache, options *Options) (*TCPServer, error) {

	n, err := newNode(options)
	if err != nil {
		return nil, err
	}

	return &TCPServer{
		node: n,
		cache:   cache,
		server:  vex.NewServer(),
		options: options,
	}, nil
}

// Run 运行这个 TCP 服务器。
func (ts *TCPServer) Run() error {
    // 注册几种命令的处理器
	ts.server.RegisterHandler(getCommand, ts.getHandler)
	ts.server.RegisterHandler(setCommand, ts.setHandler)
	ts.server.RegisterHandler(deleteCommand, ts.deleteHandler)
	ts.server.RegisterHandler(statusCommand, ts.statusHandler)
    
    // 新增的 nodes 命令，用于获取集群所有节点的名称。
	ts.server.RegisterHandler(nodesCommand, ts.nodesHandler)
	return ts.server.ListenAndServe("tcp", helpers.JoinAddressAndPort(ts.options.Address, ts.options.Port))
}

// Close 用于关闭服务器。
func (ts *TCPServer) Close() error {
	return ts.server.Close()
}

// =======================================================================

// getHandler 是处理 get 命令的的处理器。
func (ts *TCPServer) getHandler(args [][]byte) (body []byte, err error) {
    
    // 检查参数个数是否足够
	if len(args) < 1 {
		return nil, commandNeedsMoreArgumentsErr
	}

    // 使用一致性哈希选择出这个 key 所属的物理节点
	key := string(args[0])
	node, err := ts.selectNode(key)
	if err != nil {
		return nil, err
	}

    // 判断这个 key 所属的物理节点是否是当前节点，如果不是，需要响应重定向信息给客户端，并告知正确的节点地址
	if !ts.isCurrentNode(node) {
		return nil, fmt.Errorf("redirect to node %s", node)
	}

    // 调用缓存的 Get 方法，如果不存在就返回 notFoundErr 错误
	value, ok := ts.cache.Get(key)
	if !ok {
		return value, notFoundErr
	}
	return value, nil
}

// setHandler 是处理 set 命令的处理器。
func (ts *TCPServer) setHandler(args [][]byte) (body []byte, err error) {
    
    // 检查参数个数是否足够
	if len(args) < 3 {
		return nil, commandNeedsMoreArgumentsErr
	}

    // 使用一致性哈希选择出这个 key 所属的物理节点
	key := string(args[1])
	node, err := ts.selectNode(key)
	if err != nil {
		return nil, err
	}

    // 判断这个 key 所属的物理节点是否是当前节点，如果不是，需要响应重定向信息给客户端，并告知正确的节点地址
	if !ts.isCurrentNode(node) {
		return nil, fmt.Errorf("redirect to node %s", node)
	}

    // 读取 ttl，注意这里使用大端的方式读取，所以要求客户端也以大端的方式进行存储
	ttl := int64(binary.BigEndian.Uint64(args[0]))
	err = ts.cache.SetWithTTL(key, args[2], ttl)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

// deleteHandler 是处理 delete 命令的处理器。
func (ts *TCPServer) deleteHandler(args [][]byte) (body []byte, err error) {
    
    // 检查参数个数是否足够
	if len(args) < 1 {
		return nil, commandNeedsMoreArgumentsErr
	}

    // 使用一致性哈希选择出这个 key 所属的物理节点
	key := string(args[0])
	node, err := ts.selectNode(key)
	if err != nil {
		return nil, err
	}

    // 判断这个 key 所属的物理节点是否是当前节点，如果不是，需要响应重定向信息给客户端，并告知正确的节点地址
	if !ts.isCurrentNode(node) {
		return nil, fmt.Errorf("redirect to node %s", node)
	}

    // 删除指定的数据
	err = ts.cache.Delete(key)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

// statusHandler 是返回缓存状态的处理器。
func (ts *TCPServer) statusHandler(args [][]byte) (body []byte, err error) {
	return json.Marshal(ts.cache.Status())
}

// nodesHandler 是返回集群所有节点名称的处理器。
func (ts *TCPServer) nodesHandler(args [][]byte) (body []byte, err error) {
	return json.Marshal(ts.nodes())
}
