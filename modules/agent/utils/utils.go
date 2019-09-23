package utils

import (
	"crypto/tls"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego/httplib"

	"github.com/710leo/urlooker/dataobj"
	"github.com/710leo/urlooker/modules/agent/g"
)

const (
	NO_ERROR          = 0
	REQ_TIMEOUT       = 1
	INVALID_RESP_CODE = 2
	KEYWORD_UNMATCH   = 3
	DNS_ERROR         = 4
)

func CheckTargetStatus(item *dataobj.DetectedItem) {
	defer func() {
		<-g.WorkerChan
	}()

	checkResult := checkTargetStatus(item)
	g.CheckResultQueue.PushFront(checkResult)
}

func checkTargetStatus(item *dataobj.DetectedItem) (itemCheckResult *dataobj.CheckResult) {
	itemCheckResult = &dataobj.CheckResult{
		Sid:      item.Sid,
		Domain:   item.Domain,
		Creator:  item.Creator,
		Tag:      item.Tag,
		Endpoint: item.Endpoint,
		Target:   item.Target,
		Ip:       item.Ip,
		RespTime: item.Timeout,
		RespCode: "0",
	}
	reqStartTime := time.Now()

	req := httplib.Get(item.Target)

	if item.Method == "post" {
		req = httplib.Post(item.Target)
	} else if item.Method == "put" {
		req = httplib.Put(item.Target)
	}

	req.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	req.SetTimeout(3*time.Second, 10*time.Second)
	req.Header("Content-Type", "application/json")
	req.SetHost(item.Domain)
	if item.Data != "" {
		req.Header("Cookie", item.Data)
	}

	if item.PostData != "" {
		req.Body(item.PostData)
	}

	if item.Header != "" {
		headers := parseHeader(item.Header)
		for _, h := range headers {
			req.Header(h.Key, h.Value)
		}
	}
	resp, err := req.Response()
	itemCheckResult.PushTime = time.Now().Unix()

	if err != nil {
		log.Println("[ERROR]:", item.Sid, item.Domain, err)
		itemCheckResult.Status = REQ_TIMEOUT
		return
	}
	defer resp.Body.Close()

	respCode := strconv.Itoa(resp.StatusCode)
	itemCheckResult.RespCode = respCode

	respTime := int(time.Now().Sub(reqStartTime).Nanoseconds() / 1000000)
	itemCheckResult.RespTime = respTime

	if respTime > item.Timeout {
		itemCheckResult.Status = REQ_TIMEOUT
		return
	}

	if strings.Index(respCode, item.ExpectCode) == 0 || (len(item.ExpectCode) == 0 && respCode == "200") {
		if len(item.Keywords) > 0 {
			contents, _ := ioutil.ReadAll(resp.Body)
			if !strings.Contains(string(contents), item.Keywords) {
				itemCheckResult.Status = KEYWORD_UNMATCH
				return
			}
		}

		itemCheckResult.Status = NO_ERROR
		return

	} else {
		itemCheckResult.Status = INVALID_RESP_CODE
	}
	return
}

type header struct {
	Key   string
	Value string
}

func parseHeader(h string) []header {
	headers := []header{}
	kvs := strings.Split(h, "\n")
	for _, kv := range kvs {
		arr := strings.Split(kv, ":")
		if len(arr) == 2 {
			tmp := header{
				Key:   arr[0],
				Value: arr[1],
			}
			headers = append(headers, tmp)
		}
	}
	return headers
}
