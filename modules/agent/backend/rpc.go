package backend

import (
	"fmt"
	"log"
	"math/rand"
	"net/rpc"
	"strings"
	"sync"
	"time"

	"github.com/toolkits/net"
)

type BackendClients struct {
	sync.RWMutex
	Clients   map[string]*rpc.Client
	Addresses []string
}

func (this *BackendClients) GetAddresses() []string {
	return this.Addresses
}

func (this *BackendClients) InitAddresses(addresses []string) {
	this.Addresses = addresses
}

func (this *BackendClients) InitClients(clients map[string]*rpc.Client) {
	this.Lock()
	defer this.Unlock()
	this.Clients = clients
}

func (this *BackendClients) ReplaceClient(addr string, client *rpc.Client) {
	this.Lock()
	defer this.Unlock()

	old, has := this.Clients[addr]
	if has && old != nil {
		old.Close()
	}

	this.Clients[addr] = client
}

func (this *BackendClients) GetClient(addr string) (*rpc.Client, bool) {
	this.RLock()
	defer this.RUnlock()
	c, has := this.Clients[addr]
	return c, has
}

var Clients = &BackendClients{Clients: make(map[string]*rpc.Client)}

func InitClients(addresses []string) {
	Clients.InitAddresses(addresses)

	cs := make(map[string]*rpc.Client)
	for _, endpoint := range addresses {
		client, err := net.JsonRpcClient("tcp", endpoint, time.Second)
		if err != nil {
			log.Fatalln("cannot connect to", endpoint)
		}
		cs[endpoint] = client
	}

	Clients.InitClients(cs)
}

func CallRpc(method string, args, reply interface{}) error {
	addrs := Clients.GetAddresses()

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for _, i := range r.Perm(len(addrs)) {
		addr := addrs[i]
		client, has := Clients.GetClient(addr)
		if !has {
			log.Println(addr, "has no client")
			continue
		}

		err := client.Call(method, args, reply)
		if err == nil {
			return nil
		}

		if err == rpc.ErrShutdown || strings.Contains(err.Error(), "connection refused") {
			// 后端可能重启了以至于原来持有的连接关闭，或者后端挂了
			// 可以尝试再次建立连接，搞定重启的情况
			client, err = net.JsonRpcClient("tcp", addr, time.Second)
			if err != nil {
				log.Println(addr, "is dead")
				continue
			} else {
				// 重新建立了与该实例的连接
				Clients.ReplaceClient(addr, client)
				return client.Call(method, args, reply)
			}
		}

		// 刚开始后端没挂，但是仍然失败了，比如请求时间比较长，还没有结束，后端重启了，unexpected EOF
		// 不确定此时后端逻辑是否真的执行过了，防止后端逻辑不幂等，无法重试
		return err
	}

	return fmt.Errorf("all backends are dead")
}
