package cron

import (
	"time"

	"github.com/710leo/urlooker/modules/web/g"
	"github.com/710leo/urlooker/modules/web/model"

	log "github.com/sirupsen/logrus"
)

func DeleteOld() {
	t1 := time.NewTicker(time.Duration(60) * time.Second)
	for {
		<-t1.C
		err := model.RelSidIpRepo.DeleteOld(int64(g.Config.KeepDurationHour))
		if err != nil {
			log.Println("delete error:", err)
		}

		err = model.ItemStatusRepo.DeleteOld(int64(g.Config.KeepDurationHour))
		if err != nil {
			log.Println("delete error:", err)
		}
	}
}
