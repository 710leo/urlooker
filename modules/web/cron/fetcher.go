package cron

import (
	"log"
	"strings"
	"time"

	"github.com/710leo/urlooker/dataobj"
	"github.com/710leo/urlooker/modules/web/g"
	"github.com/710leo/urlooker/modules/web/model"
	"github.com/710leo/urlooker/modules/web/utils"
)

func GetDetectedItem() {
	t1 := time.NewTicker(time.Duration(60) * time.Second)
	for {
		err := getDetectedItem()
		if err != nil {
			time.Sleep(time.Second * 1)
			continue
		}
		<-t1.C
	}
}

func getDetectedItem() error {
	detectedItemMap := make(map[string][]*dataobj.DetectedItem)
	stras, err := model.GetAllStrategyByCron()
	if err != nil {
		log.Println("get strategies error:", err)
		return err
	}

	for _, s := range stras {
		_, domain, _, _ := utils.ParseUrl(s.Url)
		var ipIdcArr []dataobj.IpIdc
		if s.IP != "" {
			ips := strings.Split(s.IP, ",")
			for _, ip := range ips {
				var tmp dataobj.IpIdc
				tmp.Ip = ip
				tmp.Idc = "default"
				ipIdcArr = append(ipIdcArr, tmp)
			}
		} else {
			ipIdcArr = getIpAndIdc(domain)
		}

		for _, tmp := range ipIdcArr {
			detectedItem := newDetectedItem(s, tmp.Ip, tmp.Idc)
			key := utils.Getkey(tmp.Idc, int(detectedItem.Sid))

			if _, exists := detectedItemMap[key]; exists {
				detectedItemMap[key] = append(detectedItemMap[key], &detectedItem)
			} else {
				detectedItemMap[key] = []*dataobj.DetectedItem{&detectedItem}
			}
		}
	}

	for k, v := range detectedItemMap {
		log.Println(k)
		for _, i := range v {
			log.Println(i)
		}
	}

	g.DetectedItemMap.Set(detectedItemMap)
	return nil
}

func getIpAndIdc(domain string) []dataobj.IpIdc {

	//公司内部提供接口，拿到域名解析的ip和机房列表，获取方式写在InternalDns.CMD文件中
	if g.Config.InternalDns.Enable {
		ipIdcArr, err := utils.InternalDns(domain)
		if err != nil {
			log.Println(err)
		}
		return ipIdcArr
	}

	ipIdcArr := make([]dataobj.IpIdc, 0)

	if utils.IsIP(domain) {
		var tmp dataobj.IpIdc
		tmp.Ip = domain
		tmp.Idc = "default"
		ipIdcArr = append(ipIdcArr, tmp)
	} else {
		ips, _ := utils.LookupIP(domain, 5000)
		for _, ip := range ips {
			var tmp dataobj.IpIdc
			tmp.Ip = ip
			tmp.Idc = "default"
			ipIdcArr = append(ipIdcArr, tmp)
		}
	}

	return ipIdcArr
}

func newDetectedItem(s *model.Strategy, ip string, idc string) dataobj.DetectedItem {
	detectedItem := dataobj.DetectedItem{
		Ip:         ip,
		Idc:        idc,
		Creator:    s.Creator,
		Sid:        s.Id,
		Keywords:   s.Keywords,
		Data:       s.Data,
		Tag:        s.Tag,
		Endpoint:   s.Endpoint,
		ExpectCode: s.ExpectCode,
		Timeout:    s.Timeout,
		Header:     s.Header,
		PostData:   s.PostData,
		Method:     s.Method,
	}

	schema, domain, port, path := utils.ParseUrl(s.Url)
	if port == "" {
		detectedItem.Target = schema + "//" + ip + path
	} else {
		detectedItem.Target = schema + "//" + ip + ":" + port + path
	}

	detectedItem.Domain = domain

	return detectedItem
}
