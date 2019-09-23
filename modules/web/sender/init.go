package sender

import (
	"github.com/710leo/urlooker/modules/web/g"
	"github.com/710leo/urlooker/modules/web/utils"

	"github.com/toolkits/container/list"
)

var (
	NodeRing   *ConsistentHashNodeRing          // 服务节点的一致性哈希环
	SendQueues map[string]*list.SafeListLimited // 发送缓存队列,减少发起连接次数
)

func initRing() {
	NodeRing = NewConsistentHashNodeRing(g.Config.Alarm.Replicas, utils.KeysOfMap(g.Config.Alarm.Cluster))
}

func initSendQueues() {
	SendQueues = make(map[string]*list.SafeListLimited)

	for node, _ := range g.Config.Alarm.Cluster {
		Q := list.NewSafeListLimited(10240)
		SendQueues[node] = Q
	}
}

func startSendTasks() {
	for node, _ := range g.Config.Alarm.Cluster {
		queue := SendQueues[node]
		go SendToAlarm(queue, node)
	}
}

func Init() {
	initRing()
	initSendQueues()

	startSendTasks()
}
