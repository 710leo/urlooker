package cron

import (
	"log"
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
		detectedItem := newDetectedItem(s)
		idc := detectedItem.Idc
		if _, exists := detectedItemMap[idc]; exists {
			detectedItemMap[idc] = append(detectedItemMap[idc], &detectedItem)
		} else {
			detectedItemMap[idc] = []*dataobj.DetectedItem{&detectedItem}
		}
	}

	g.DetectedItemMap.Set(detectedItemMap)
	return nil
}

func newDetectedItem(s *model.Strategy) dataobj.DetectedItem {
	_, domain, _, _ := utils.ParseUrl(s.Url)
	idc := s.Idc
	if idc == "" {
		idc = g.Config.IDC[0]
	}
	detectedItem := dataobj.DetectedItem{
		Idc:        idc,
		Target:     s.Url,
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
		Domain:     domain,
	}

	return detectedItem
}
