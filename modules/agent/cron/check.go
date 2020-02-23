package cron

import (
	"fmt"
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
		items, err := GetItem()
		if err != nil {
			log.Println("[ERROR] ", err)
		}

		for _, item := range items {
			g.WorkerChan <- 1
			go utils.CheckTargetStatus(item)
		}
		<-t1.C
	}
}

func GetItem() ([]*dataobj.DetectedItem, error) {
	var resp dataobj.GetItemResponse
	log.Println(g.Config.IDC)
	err := backend.CallRpc("Web.GetItem", g.Config.IDC, &resp)
	if err != nil {
		return []*dataobj.DetectedItem{}, err
	}
	if resp.Message != "" {
		err := fmt.Errorf(resp.Message)
		return []*dataobj.DetectedItem{}, err
	}

	return resp.Data, err
}
