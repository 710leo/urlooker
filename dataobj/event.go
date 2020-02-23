package dataobj

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

type EventQueue struct {
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
