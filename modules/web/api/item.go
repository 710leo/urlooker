package api

import (
	"log"
	"time"

	"github.com/710leo/urlooker/dataobj"
	"github.com/710leo/urlooker/modules/web/g"
	"github.com/710leo/urlooker/modules/web/model"
	"github.com/710leo/urlooker/modules/web/sender"
	"github.com/710leo/urlooker/modules/web/utils"
)

func (this *Web) SendResult(req dataobj.SendResultReq, reply *string) error {
	for _, arg := range req.CheckResults {
		itemStatus := model.ItemStatus{
			Ip:       arg.Ip,
			Sid:      arg.Sid,
			RespTime: arg.RespTime,
			RespCode: arg.RespCode,
			PushTime: arg.PushTime,
			Result:   arg.Status,
		}

		relSidIp := model.RelSidIp{
			Sid: arg.Sid,
			Ip:  arg.Ip,
			Ts:  time.Now().Unix(),
		}

		err := relSidIp.Save()
		if err != nil {
			log.Println("save sid_ip error:", err)
			*reply = "save sid_ip error:" + err.Error()
			return nil
		}

		err = itemStatus.Save()
		if err != nil {
			log.Println("save item error:", err)
			*reply = "save item error:" + err.Error()
			return nil
		}

		if g.Config.Alarm.Enable {
			node, err := sender.NodeRing.GetNode(itemStatus.PK())
			if err != nil {
				log.Println("error:", err)
				*reply = "get node error:" + err.Error()
				return nil
			}

			Q := sender.SendQueues[node]
			isSuccess := Q.PushFront(itemStatus)
			if !isSuccess {
				log.Println("push itemStatus error:", itemStatus)
				*reply = "push itemStatus error"
				return nil
			}
		}

	}

	if g.Config.Falcon.Enable {
		if len(req.CheckResults) > 0 {
			utils.PushFalcon(req.CheckResults, req.Hostname)
		}
	}

	*reply = ""
	return nil
}

func (this *Web) GetItem(hostname string, resp *dataobj.GetItemResponse) error {
	items, exists := g.DetectedItemMap.Get(hostname)
	if !exists {
		resp.Message = "no found item assigned to " + hostname
	}
	resp.Data = items
	return nil
}
