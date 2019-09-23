package cron

import (
	"log"
	"time"

	"github.com/710leo/urlooker/dataobj"

	"github.com/710leo/urlooker/modules/agent/backend"
	"github.com/710leo/urlooker/modules/agent/g"
	"github.com/710leo/urlooker/modules/agent/utils"
)

func StartCheck() {
	t1 := time.NewTicker(time.Duration(g.Config.Web.Interval) * time.Second)
	for {
		items, _ := GetItem()

		for _, item := range items {
			g.WorkerChan <- 1
			go utils.CheckTargetStatus(item)
		}
		<-t1.C
	}
}

func GetItem() ([]*dataobj.DetectedItem, error) {
	hostname, _ := g.Hostname()

	var resp dataobj.GetItemResponse
	err := backend.CallRpc("Web.GetItem", hostname, &resp)
	if err != nil {
		log.Println(err)
	}
	if resp.Message != "" {
		log.Println(resp.Message)
	}

	return resp.Data, err
}
