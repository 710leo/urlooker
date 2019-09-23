package judge

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/710leo/urlooker/dataobj"
	"github.com/710leo/urlooker/modules/alarm/backend"
	"github.com/710leo/urlooker/modules/alarm/cache"
	"github.com/710leo/urlooker/modules/alarm/g"
)

func Judge(L *SafeLinkedList, item *dataobj.ItemStatus, now int64) {
	strategy, exists := cache.StrategyMap.Get(item.Sid)
	if !exists {
		return
	}
	event := &dataobj.Event{
		EventId:    fmt.Sprintf("e_%d_%s", item.Id, item.PK()),
		Url:        strategy.Url,
		StrategyId: item.Sid,
		Ip:         item.Ip,
		EventTime:  item.PushTime,
		RespCode:   item.RespCode,
		RespTime:   item.RespTime,
		Result:     item.Result,
	}

	historyData, isTriggered := compute(L, item, strategy)

	sendEventIfNeed(historyData, isTriggered, now, event, strategy.MaxStep)
}

func compute(L *SafeLinkedList, item *dataobj.ItemStatus, strategy dataobj.Strategy) (historyData []*HistoryData, isTriggered bool) {
	historyData, isEnough := L.HistoryData(strategy.Times)
	if !isEnough {
		return
	}

	num := 0 //策略触发次数
	for i := 0; i < strategy.Times; i++ {
		if historyData[i].Value != 0 {
			num++
		}
	}

	if num < strategy.Times {
		return
	}

	isTriggered = true
	return
}

func sendEventIfNeed(historyData []*HistoryData, isTriggered bool, now int64, event *dataobj.Event, maxStep int) {
	lastEvent, exists := cache.LastEvents.Get(event.EventId)
	if isTriggered {
		event.Status = "PROBLEM"
		if !exists || lastEvent.Status[0] == 'O' {
			// 本次触发了阈值，之前又没报过警，得产生一个报警Event
			event.CurrentStep = 1

			// 但是有些用户把最大报警次数配置成了0，相当于屏蔽了，要检查一下
			if maxStep == 0 {
				return
			}

			sendEvent(event)
			return
		}

		// 逻辑走到这里，说明之前Event是PROBLEM状态
		if lastEvent.CurrentStep >= maxStep {
			// 报警次数已经足够多，到达了最多报警次数了，不再报警
			return
		}

		if historyData[len(historyData)-1].Timestamp <= lastEvent.EventTime {
			// 产生过报警的点，就不能再使用来判断了，否则容易出现一分钟报一次的情况
			return
		}

		if now-lastEvent.EventTime < g.Config.Alarm.MinInterval {
			// 报警不能太频繁，间隔至少MinIntervals
			return
		}

		event.CurrentStep = lastEvent.CurrentStep + 1
		sendEvent(event)
	} else {
		// 如果LastEvent是Problem，报OK，否则啥都不做
		if exists && lastEvent.Status[0] == 'P' {
			event.Status = "OK"
			event.CurrentStep = 1
			sendEvent(event)
		}
	}
}

func sendEvent(event *dataobj.Event) {
	cache.LastEvents.Set(event.EventId, event)
	saveEvent(event)

	bs, err := json.Marshal(event)
	if err != nil {
		log.Printf("json marshal event %v fail: %v", event, err)
	}

	redisKey := g.Config.Alarm.QueuePattern
	rc := g.RedisConnPool.Get()
	defer rc.Close()

	_, err = rc.Do("LPUSH", redisKey, string(bs))
	if err != nil {
		log.Println(err)
	}
}

func saveEvent(event *dataobj.Event) {
	var reply string
	err := backend.CallRpc("Web.SaveEvent", event, &reply)
	if err != nil {
		log.Println("[ERROR] Web.SaveEvent:", err)
		return
	}
	if reply != "" {
		log.Println("reply:", reply)
	}
}
