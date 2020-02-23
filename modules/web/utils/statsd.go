package utils

import (
	"log"
	"os"
	"time"

	"github.com/quipo/statsd"
)

var stats *statsd.StatsdBuffer

func InitStatsd(addr string) {
	prefix := "urlooker."
	statsdclient := statsd.NewStatsdClient("localhost:8125", prefix)
	err := statsdclient.CreateSocket()
	if nil != err {
		log.Println(err)
		os.Exit(1)
	}
	interval := time.Second * 2 // aggregate stats and flush every 2 seconds
	stats = statsd.NewStatsdBuffer(interval, statsdclient)
}

func PushStatsd(metric string, value int64) error {
	return stats.Gauge(metric, value)
}
