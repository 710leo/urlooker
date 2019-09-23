package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/710leo/urlooker/dataobj"
	"github.com/710leo/urlooker/modules/web/g"

	"github.com/astaxie/beego/httplib"
)

var transport *http.Transport = &http.Transport{}

type MetricValue struct {
	Endpoint  string      `json:"endpoint"`
	Metric    string      `json:"metric"`
	Tags      string      `json:"tags"`
	Value     interface{} `json:"value"`
	Timestamp int64       `json:"timestamp"`
	Type      string      `json:"counterType"`
	Step      int64       `json:"step"`
}

func PushFalcon(itemCheckedArray []*dataobj.CheckResult, hostname string) {

	pushDatas := make([]*MetricValue, 0)
	for _, itemChecked := range itemCheckedArray {
		tags := fmt.Sprintf("ip=%s,domain=%s,creator=%s,from=%s", itemChecked.Ip, itemChecked.Domain, itemChecked.Creator, hostname)
		if len(itemChecked.Tag) > 0 { //补充用户自定义tag
			tags += "," + itemChecked.Tag
		}

		//url 状态
		data := getMetric(itemChecked, "url_status", tags, itemChecked.Status)
		pushDatas = append(pushDatas, &data)

		//url 响应时间
		data2 := getMetric(itemChecked, "url_resp_time", tags, int64(itemChecked.RespTime))
		pushDatas = append(pushDatas, &data2)
	}

	err := push(pushDatas)
	if err != nil {
		log.Println("push error", err)
	}
}

func getMetric(item *dataobj.CheckResult, metric, tags string, value int64) MetricValue {
	var data MetricValue
	data.Endpoint = fmt.Sprintf("url_%d", item.Sid)
	if item.Endpoint != "" {
		data.Endpoint = item.Endpoint
	}

	data.Timestamp = item.PushTime
	data.Metric = metric
	data.Type = "GAUGE"
	data.Step = int64(g.Config.Falcon.Interval)
	data.Tags = tags
	data.Value = value
	return data
}

func push(data []*MetricValue) error {
	d, err := json.Marshal(data)
	if err != nil {
		return err
	}

	log.Println("to falcon: ", string(d))

	_, err = httplib.Post(g.Config.Falcon.Addr).Body(d).String()
	if err != nil {
		return err
	}

	return nil
}
