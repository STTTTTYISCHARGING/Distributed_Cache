package servers

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/FishGoddess/cachego"
	"github.com/FishGoddess/vex"
	"cache-server/caches"
	"stathat.com/c/consistent"
)

const (
	// ttlOfClient 是客户端连接的有效期，单位是秒，所以这里是 15 分钟。
	ttlOfClient = 15 * 60

	// redirectPrefix 是重定向错误的前缀，用于判断某个错误是不是重定向错误。
	redirectPrefix = "redirect to node "

	// maxRedirectTimes 是最大的重定向次数，如果某次操作重定向了 5 次，说明集群节点的波动太大了，几乎可以认为是不可用的了。
	maxRedirectTimes = 5

	// updateCircleDuration 是更新节点信息的时间间隔，主要是用于更新一致性哈希的节点情况。
	updateCircleDuration = 5 * time.Minute
)

var (
	// noClientIsAvailableErr 意味着没有任何的客户端连接可以使用了。
	noClientIsAvailableErr = errors.New("no client is available")

	// reachMaxRetriedTimesErr 意味着重定向次数已经超过了最大限制，说明集群处于不可用状态。
	reachedMaxRetriedTimesErr = errors.New("reached max redirect times")
)

// TCPClient 是客户端结构。
type TCPClient struct {

	// clients 存储了所有的客户端连接，这是一个缓存结构。
	clients *cachego.Cache

	// circle 存储了当前集群的一致性哈希信息，用于避免重定向。
	circle *consistent.Consistent
}

// NewTCPClient 返回一个新创建的客户端实例。
// 由于服务端已经是集群了，这里填的 address 是集群中的一个节点地址。
func NewTCPClient(address string) (*TCPClient, error) {

    // 内部使用的是 vex 客户端
	client, err := vex.NewClient("tcp", address)
	if err != nil {
		return nil, err
	}

    // 创建一致性哈希环，并将虚拟节点设置为和服务端一致，否则节点的判断会发生误差
	circle := consistent.New()
	circle.NumberOfReplicas = 1024
	circle.Set([]string{address})

    // 创建缓存，设置过期数据清理的时间间隔是 10 分钟，并给所有的客户端连接设置 15 分钟的有效期
	clients := cachego.NewCache()
	clients.AutoGc(10 * time.Minute)
	clients.SetWithTTL(address, client, ttlOfClient)

	tc := &TCPClient{
		clients: clients,
		circle:  circle,
	}
    
    // 开启一个定时任务，定期更新一致性哈希信息
	tc.updateCircleAtFixedDuration(updateCircleDuration)
	return tc, tc.updateCircleAndClients()
}

// updateCircleAtFixedDuration 会开启一个定时任务，定期更新一致性哈希信息。
func (tc *TCPClient) updateCircleAtFixedDuration(duration time.Duration) {
	go func() {
		ticker := time.NewTicker(duration)
		for {
			select {
			case <-ticker.C:
                // 获取集群的节点信息，并更新到一致性哈希中
				nodes, err := tc.nodes()
				if err == nil {
					tc.circle.Set(nodes)
				}
			}
		}
	}()
}

// nodes 返回集群的节点信息。
func (tc *TCPClient) nodes() ([]string, error) {

    // 获取一致性哈希的成员，只要有一个即可
    // 首先拿到这个节点的客户端连接，然后查询集群的节点信息并返回
	nodes := tc.circle.Members()
	for _, node := range nodes {
		client, err := tc.getOrCreateClient(node)
		if err != nil {
			continue
		}
		body, err := client.Do(nodesCommand, nil)
		if err != nil {
			return nil, err
		}
		var nodes []string
		err = json.Unmarshal(body, &nodes)
		return nodes, err
	}
	return nil, noClientIsAvailableErr
}

// getOrCreateClient 从缓存中拿到某个节点的客户端连接。
func (tc *TCPClient) getOrCreateClient(node string) (*vex.Client, error) {

    // 如果拿不到，说明这个节点的客户端连接要么没有要么过期了，所以需要新创建一个
	client, ok := tc.clients.Get(node)
	if !ok {
		var err error
		client, err = vex.NewClient("tcp", node)
		if err != nil {
			return nil, err
		}
        
        // 注意新创建的连接需要设置有效性
		tc.clients.SetWithTTL(node, client, ttlOfClient)
	}
	return client.(*vex.Client), nil
}

// updateCircleAndClients 更新一致性哈希和客户端连接。
func (tc *TCPClient) updateCircleAndClients() error {

	nodes, err := tc.nodes()
	if err != nil {
		return err
	}

    // 更新一致性哈希，并根据节点信息更新客户端连接信息
	tc.circle.Set(nodes)
	for _, node := range nodes {
		tc.getOrCreateClient(node)
	}
	return nil
}

// clientOf 返回某个 key 的客户端连接。
func (tc *TCPClient) clientOf(key string) (*vex.Client, error) {
    
    // 使用一致性哈希环判断这个 key 属于哪一个节点，然后获取这个节点的客户端连接
    // 所以一致性哈希环的准确性直接关系到重定向问题的解决
	node, err := tc.circle.Get(key)
	if err != nil {
		return nil, err
	}
	return tc.getOrCreateClient(node)
}

// doCommand 使用 client 执行命令。
func (tc *TCPClient) doCommand(client *vex.Client, command byte, args [][]byte) (body []byte, err error) {

    // 因为可能存在重定向，所以使用循环，但是不能一直重定向，所以设置了一个最大的重定向次数
	for i := 0; i < maxRedirectTimes; i++ {
		body, err := client.Do(command, args)
        
        // 判断发生的错误是不是重定向错误，如果是，就从错误中获取正确的节点地址，并拿到这个节点的客户端连接，再次执行命令
		if err != nil && strings.HasPrefix(err.Error(), redirectPrefix) {
			node := strings.TrimPrefix(err.Error(), redirectPrefix)
			rightClient, err := tc.getOrCreateClient(node)
			if err != nil {
				continue
			}
			client = rightClient
			continue
		}

		// 如果错误不是重定向错误，而是这个连接关闭的错误，说明这个节点出现问题，很可能是节点信息已经不准了，需要更新集群的节点信息
		if err != nil && strings.HasSuffix(err.Error(), "closed by the remote host.") {
			nodes, err := tc.nodes()
			if err == nil {
				tc.circle.Set(nodes)
			}
		}
		return body, err
	}
    
    // 重定向达到了最大次数，直接返回错误
	return nil, reachedMaxRetriedTimesErr
}

// Get 获取指定 key 的数据。
func (tc *TCPClient) Get(key string) ([]byte, error) {
	client, err := tc.clientOf(key)
	if err != nil {
		return nil, err
	}
	return tc.doCommand(client, getCommand, [][]byte{[]byte(key)})
}

// Set 添加数据到缓存中。
func (tc *TCPClient) Set(key string, value []byte, ttl int64) error {

	client, err := tc.clientOf(key)
	if err != nil {
		return err
	}

	ttlBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(ttlBytes, uint64(ttl))
	_, err = tc.doCommand(client, setCommand, [][]byte{
		ttlBytes, []byte(key), value,
	})
	return err
}

// Delete 删除指定 key 的数据。
func (tc *TCPClient) Delete(key string) error {

	client, err := tc.clientOf(key)
	if err != nil {
		return err
	}

	_, err = tc.doCommand(client, deleteCommand, [][]byte{[]byte(key)})
	return err
}

// Status 返回缓存服务的状态。
func (tc *TCPClient) Status() (*caches.Status, error) {

    // 由于缓存服务可能是一个集群，所以这里需要获取所有节点的状态，然后做一个汇总
	totalStatus := caches.NewStatus()
	nodes := tc.circle.Members()
	for _, node := range nodes {
		client, err := tc.getOrCreateClient(node)
		if err != nil {
			continue
		}
		body, err := client.Do(statusCommand, nil)
		if err != nil {
			return nil, err
		}
		status := caches.NewStatus()
		err = json.Unmarshal(body, status)
		if err != nil {
			return nil, err
		}
		totalStatus.Count += status.Count
		totalStatus.KeySize += status.KeySize
		totalStatus.ValueSize += status.ValueSize
	}
	return totalStatus, nil
}

// Nodes 返回集群的节点信息。
func (tc *TCPClient) Nodes() ([]string, error) {
	return tc.nodes()
}

// Close 会关闭这个客户端。
func (tc *TCPClient) Close() (err error) {

    // 当然需要将每一个客户端连接都关闭掉
	nodes := tc.circle.Members()
	for _, node := range nodes {
		client, ok := tc.clients.Get(node)
		if ok {
			err = client.(*vex.Client).Close()
		}
	}
	tc.clients.RemoveAll()
	return err
}
