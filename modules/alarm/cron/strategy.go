package cron

import (
	"log"
	"time"

	"github.com/710leo/urlooker/dataobj"
	"github.com/710leo/urlooker/modules/alarm/backend"
	"github.com/710leo/urlooker/modules/alarm/cache"
	"github.com/710leo/urlooker/modules/alarm/g"
)

func SyncStrategies() {
	t1 := time.NewTicker(time.Duration(g.Config.Web.Interval) * time.Second)
	for {
		syncStrategies()

		<-t1.C
	}

}

func syncStrategies() {

	var strategiesResponse dataobj.StrategyResponse
	err := backend.CallRpc("Web.GetStrategies", "", &strategiesResponse)
	if err != nil {
		log.Println("[ERROR] Web.GetStrategies:", strategiesResponse.Data, strategiesResponse.Message, err)
		return
	}

	rebuildStrategyMap(strategiesResponse.Data)
}

func rebuildStrategyMap(strategiesResponse []*dataobj.Strategy) {

	m := make(map[int64]dataobj.Strategy)
	for _, strategy := range strategiesResponse {
		m[strategy.Id] = *strategy
	}

	cache.StrategyMap.ReInit(m)
}
