package utils

import (
	"fmt"

	"github.com/710leo/urlooker/dataobj"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
)

func PushPromethues(addr string, itemsChecked []*dataobj.CheckResult, ip string) {
	addr = "http://" + addr
	for _, itemChecked := range itemsChecked {
		completionTime := prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "urlooker_api_status",
		})
		completionTime.SetToCurrentTime()
		completionTime.Set(float64(itemChecked.Status))

		if err := push.New(addr, "urlooker").
			Collector(completionTime).
			Grouping("domain", itemChecked.Domain).
			Grouping("creator", itemChecked.Creator).
			Grouping("from", ip).
			Push(); err != nil {
			fmt.Println("Could not push completion time to Pushgateway:", err)
		}
	}
}
