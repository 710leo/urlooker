package dataobj

import "fmt"

type Event struct {
	Id          int64  `json:"id"`
	EventId     string `json:"event_id"`
	Status      string `json:"status"`
	Url         string `json:"url"`
	Ip          string `json:"ip"`
	EventTime   int64  `json:"event_time"`
	StrategyId  int64  `json:"strategy_id"`
	RespTime    int    `json:"resp_time"`
	RespCode    string `json:"resp_code"`
	Result      int64  `json:"result"`
	CurrentStep int    `json:"current_step"`
	MaxStep     int    `json:"max_step"`
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
