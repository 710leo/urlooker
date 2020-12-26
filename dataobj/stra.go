package dataobj

type Strategy struct {
	Id          int64  `json:"id"`
	Idc         string `json:"idc"`
	Method      string `json:"method"`
	Url         string `json:"url"`
	Enable      int    `json:"enable"`
	IP          string `json:"ip" xorm:"ip"`
	Keywords    string `json:"keywords"`
	Timeout     int    `json:"timeout"`
	Creator     string `json:"creator"`
	ExpectCode  string `json:"expect_code"`
	Note        string `json:"note"`
	Data        string `json:"data"`
	Endpoint    string `json:"endpoint"`
	Header      string `json:"header"`
	PostData    string `json:"post_data"`
	Tag         string `json:"tag"`
	MaxStep     int    `json:"max_step"`
	Times       int    `json:"times"`
	Teams       string `json:"teams"`
	DingWebhook string `json:"ding_webhook"`
}

type StrategyResponse struct {
	Message string
	Data    []*Strategy
}
