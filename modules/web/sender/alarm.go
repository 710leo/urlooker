package sender

import (
	"log"
	"time"

	"github.com/toolkits/container/list"

	"github.com/710leo/urlooker/modules/web/backend"
	"github.com/710leo/urlooker/modules/web/g"
)

func SendToAlarm(Q *list.SafeListLimited, node string) {
	cfg := g.Config
	batch := cfg.Alarm.Batch
	addr := cfg.Alarm.Cluster[node]

	//todo：rpc 当数据量增大时，rpc调用改为并行方式
	for {
		items := Q.PopBackBy(batch)
		count := len(items)
		if count == 0 {
			time.Sleep(time.Duration(cfg.Alarm.SleepTime) * time.Second)
			continue
		}

		var resp string
		var err error
		sendOk := false
		for i := 0; i < 3; i++ {
			rpcClient := backend.NewRpcClient(addr)
			err = rpcClient.Call("Alarm.Send", items, &resp)
			if err == nil {
				sendOk = true
				break
			}
			time.Sleep(1)
		}

		if !sendOk {
			log.Printf("send alarm %s:%s fail: %v", node, addr, err)
		}
		if cfg.Debug {
			log.Println("<=", resp)
		}
	}
}
