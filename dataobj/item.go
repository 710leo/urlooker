package dataobj

import (
	"crypto/md5"
	"fmt"
	"io"
)

type IpIdc struct {
	Ip  string
	Idc string
}

//下发给agent的数据结构
type DetectedItem struct {
	Sid        int64  `json:"sid"`
	Method     string `json:"method"`
	Domain     string `json:"domain"`
	Target     string `json:"target"`
	Ip         string `json:"ip"`
	Keywords   string `json:"keywords"`
	Timeout    int    `json:"timeout"`
	Creator    string `json:"creator"`
	Data       string `json:"data"`
	Endpoint   string `json:"endpoint"`
	Tag        string `json:"tag"`
	ExpectCode string `json:"expect_code"`
	Idc        string `json:"idc"`
	Header     string `json:"header"`
	PostData   string `json:"post_data"`
}

//agent上报的数据结构
type CheckResult struct {
	Sid      int64  `json:"sid"`
	Domain   string `json:"domain"`
	Target   string `json:"target"`
	Creator  string `json:"creator"`
	Endpoint string `json:"endpoint"`
	Tag      string `json:"tag"`
	RespCode string `json:"resp_code"`
	RespTime int    `json:"resp_time"`
	Status   int64  `json:"status"`
	PushTime int64  `json:"push_time"`
	Ip       string `json:"ip"`
}

type ItemStatus struct {
	Id       int64  `json:"id"`
	Sid      int64  `json:"sid"`
	Ip       string `json:"ip"`
	RespTime int    `json:"resp_time"`
	RespCode string `json:"resp_code"`
	PushTime int64  `json:"push_time"`
	Result   int64  `json:"result"`
}

func (this *ItemStatus) PK() string {
	h := md5.New()
	io.WriteString(h, fmt.Sprintf("%s/%s", this.Sid, this.Ip))

	return fmt.Sprintf("%x", h.Sum(nil))
}

type SendResultReq struct {
	Hostname     string
	CheckResults []*CheckResult
}

type GetItemResponse struct {
	Message string
	Data    []*DetectedItem
}
