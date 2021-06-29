package utils

import (
	"crypto/tls"
	"io/ioutil"
	"log"
	"net"
	"strconv"
	"strings"
	"time"
    "fmt"
    "regexp"

	"github.com/710leo/urlooker/dataobj"
	"github.com/710leo/urlooker/modules/agent/g"

	"github.com/astaxie/beego/httplib"
)

const (
	NO_ERROR          = 0
	REQ_TIMEOUT       = 1
	INVALID_RESP_CODE = 2
	KEYWORD_UNMATCH   = 3
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
		RespTime: item.Timeout,
		Step:     int64(g.Config.Web.Interval),
		RespCode: "-",
	}
	reqStartTime := time.Now()

	defer func() {
		log.Printf("[detect]:sid:%d domain:%s result:%d\n", item.Sid, item.Domain, itemCheckResult.Status)
	}()

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
			if h.Key == "Host" {
				req.SetHost(h.Value)
			}
		}
	}
	resp, err := req.Response()
	itemCheckResult.PushTime = time.Now().Unix()

	if err != nil {
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

            contentsFormat := string(contents)
            for _, keywordsItem := range strings.Split(item.Keywords, "|||") {
                log.Printf("[keywords-item]:%s", keywordsItem)

                if strings.Contains(keywordsItem, ":::") {
                    keywordsItemArr := strings.Split(keywordsItem, ":::")
                    rule := keywordsItemArr[0]
                    if rule == "replace" {
                        contentsFormat = strings.Replace(contentsFormat, keywordsItemArr[1], keywordsItemArr[2], 1)
                        log.Printf("[replace]:%s to %s", keywordsItemArr[1], keywordsItemArr[2])
                        continue
                    }

                    if len(keywordsItemArr) > 2 {
                        // get regexp rule from keywordsItem(not:::regexp:::<regexp>)
                        rule = keywordsItemArr[1]
                    }

                    switch rule {
                        case "not":
                            // not contains
                            if strings.Contains(contentsFormat, keywordsItemArr[1]) {
                                itemCheckResult.Status = KEYWORD_UNMATCH
                                return
                            }
                        case "regexp":
                            // regexp contains
                            if rule == keywordsItemArr[0] {
                                reg := regexp.MustCompile(fmt.Sprintf(`%s`, keywordsItemArr[1]))
                                reg_arr := reg.FindAllStringSubmatch(contentsFormat, -1)
                                if reg_arr == nil {
                                    log.Printf("[regexp]:`%s` not match data: %s", keywordsItemArr[1], contentsFormat)
                                    itemCheckResult.Status = KEYWORD_UNMATCH
                                    return
                                }
                            } else {
                                // not regexp contains
                                reg := regexp.MustCompile(fmt.Sprintf(`%s`, keywordsItemArr[2]))
                                reg_arr := reg.FindAllStringSubmatch(contentsFormat, -1)
                                if reg_arr != nil {
                                    log.Printf("[regexp]:`%s` match data: %#v in %s", keywordsItemArr[2], reg_arr, contentsFormat)
                                    itemCheckResult.Status = KEYWORD_UNMATCH
                                    return
                                }
                            }
                        default:
                            log.Printf("[rule]: %s not found!", rule)
                    }

                    continue
                }

                if !strings.Contains(contentsFormat, keywordsItem) {
                    itemCheckResult.Status = KEYWORD_UNMATCH
                    return
                }
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

func IntranetIP() (ips []string, err error) {
	ips = make([]string, 0)

	ifaces, e := net.Interfaces()
	if e != nil {
		return ips, e
	}

	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}

		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}

		// ignore docker and warden bridge
		if strings.HasPrefix(iface.Name, "docker") || strings.HasPrefix(iface.Name, "w-") {
			continue
		}

		addrs, e := iface.Addrs()
		if e != nil {
			return ips, e
		}

		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			if ip == nil || ip.IsLoopback() {
				continue
			}

			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}

			ipStr := ip.String()
			if IsIntranet(ipStr) {
				ips = append(ips, ipStr)
			}
		}
	}

	return ips, nil
}

func IsIntranet(ipStr string) bool {
	if strings.HasPrefix(ipStr, "10.") {
		return true
	}

	if strings.HasPrefix(ipStr, "192.168.") {
		return true
	}

	if strings.HasPrefix(ipStr, "172.") {
		// 172.16.0.0-172.31.255.255
		arr := strings.Split(ipStr, ".")
		if len(arr) != 4 {
			return false
		}

		second, err := strconv.ParseInt(arr[1], 10, 64)
		if err != nil {
			return false
		}

		if second >= 16 && second <= 31 {
			return true
		}
	}

	return false
}
