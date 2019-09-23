package model

import (
	"fmt"
	"time"

	"github.com/710leo/urlooker/dataobj"
	. "github.com/710leo/urlooker/modules/web/store"
)

type Event dataobj.Event

var EventRepo *Event

func (this *Event) Insert() error {
	_, err := Orm.Insert(this)
	return err
}

func (this *Event) GetByStrategyId(strategyId int64, before int) ([]*Event, error) {
	events := make([]*Event, 0)
	ts := time.Now().Unix() - int64(before)
	err := Orm.Where("strategy_id = ? and event_time > ?", strategyId, ts).Desc("event_time").Find(&events)
	return events, err
}

func (this *Event) String() string {
	return fmt.Sprintf(
		"<Id:%s, EventId:,%s Ip:%s, Url:%s, EventTime:%v, StrategyId:%d, RespTime:%s, RespCode:%s, Status:%s, (%d/%d)>",
		this.Id,
		this.EventId,
		this.Ip,
		this.Url,
		this.EventTime,
		this.StrategyId,
		this.RespTime,
		this.RespCode,
		this.Status,
		this.CurrentStep,
		this.MaxStep,
	)
}
