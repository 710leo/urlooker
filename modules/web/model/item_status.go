package model

import (
	"crypto/md5"
	"fmt"
	"io"
	"time"

	"github.com/710leo/urlooker/modules/web/g"
	. "github.com/710leo/urlooker/modules/web/store"
)

type ItemStatus struct {
	Id       int64  `json:"id"`
	Sid      int64  `json:"sid"`
	Ip       string `json:"ip"`
	RespTime int    `json:"resp_time"`
	RespCode string `json:"resp_code"`
	PushTime int64  `json:"push_time"`
	Result   int64  `json:"result"`
}

var ItemStatusRepo *ItemStatus

func (this *ItemStatus) Save() error {
	sql := fmt.Sprintf("insert into item_status00 (ip, sid, resp_time, resp_code, push_time, result) value(?,?,?,?,?,?)")
	_, err := Orm.Exec(sql, this.Ip, this.Sid, this.RespTime, this.RespCode, this.PushTime, this.Result)
	return err
}

func (this *ItemStatus) GetByIpAndSid(ip string, sid int64) ([]*ItemStatus, error) {
	itemStatusArr := make([]*ItemStatus, 0)
	ts := time.Now().Unix() - int64(g.Config.ShowDurationMin*60)
	sql := fmt.Sprintf("select * from item_status00 where ip=? and sid=? and push_time > ?")

	err := Orm.Sql(sql, ip, sid, ts).Find(&itemStatusArr)
	return itemStatusArr, err
}

func (this *ItemStatus) PK() string {
	h := md5.New()
	io.WriteString(h, fmt.Sprintf("%s/%s", this.Sid, this.Ip))

	return fmt.Sprintf("%x", h.Sum(nil))
}

func (this *ItemStatus) DeleteOld(d int64) error {
	ts := time.Now().Unix() - 12*60*60
	sql := fmt.Sprintf("delete from item_status00 where push_time < ?")
	_, err := Orm.Exec(sql, ts)
	return err
}
