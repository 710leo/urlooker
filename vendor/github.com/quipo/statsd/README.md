# StatsD client (Golang)

[![Build Status](https://travis-ci.org/quipo/statsd.png?branch=master)](https://travis-ci.org/quipo/statsd) 
[![GoDoc](https://godoc.org/github.com/quipo/statsd?status.png)](http://godoc.org/github.com/quipo/statsd)

## Introduction

Go Client library for [StatsD](https://github.com/etsy/statsd/). Contains a direct and a buffered client.
The buffered version will hold and aggregate values for the same key in memory before flushing them at the defined frequency.

This client library was inspired by the one embedded in the [Bit.ly NSQ](https://github.com/bitly/nsq/blob/master/util/statsd_client.go) project, and extended to support some extra custom events used at DataSift.

## Installation

    go get github.com/quipo/statsd

## Supported event types

* `Increment` - Count occurrences per second/minute of a specific event
* `Decrement` - Count occurrences per second/minute of a specific event
* `Timing` - To track a duration event
* `PrecisionTiming` - To track a duration event
* `Gauge` (int) / `FGauge` (float) - Gauges are a constant data type. They are not subject to averaging, and they don’t change unless you change them. That is, once you set a gauge value, it will be a flat line on the graph until you change it again
* `GaugeDelta` (int) / `FGaugeDelta` (float) - Same as above, but as a delta change to the previous value rather than a new absolute value
* `Absolute` (int) / `FAbsolute` (float) - Absolute-valued metric (not averaged/aggregated)
* `Total` - Continously increasing value, e.g. read operations since boot


## Sample usage

```go
package main

import (
	"log"
	"os"
	"time"

	"github.com/quipo/statsd"
)

func main() {
	// init
	prefix := "myproject."
	statsdclient := statsd.NewStatsdClient("localhost:8125", prefix)
	err := statsdclient.CreateSocket()
	if nil != err {
		log.Println(err)
		os.Exit(1)
	}
	interval := time.Second * 2 // aggregate stats and flush every 2 seconds
	stats := statsd.NewStatsdBuffer(interval, statsdclient)
	defer stats.Close()

	// not buffered: send immediately
	statsdclient.Incr("mymetric", 4)

	// buffered: aggregate in memory before flushing
	stats.Incr("mymetric", 1)
	stats.Incr("mymetric", 3)
	stats.Incr("mymetric", 1)
	stats.Incr("mymetric", 1)
}
```

The string `%HOST%` in the metric name will automatically be replaced with the hostname of the server the event is sent from.


## [Changelog](https://github.com/quipo/statsd/releases)

* `HEAD`:

    *

* [`v.1.4.0`](https://github.com/quipo/statsd/releases/tag/1.4.0)

    * Fixed behaviour of Gauge with positive numbers: the previous behaviour was the same as GaugeDelta
      (FGauge already had the correct behaviour)
    * Added more tests
    * Small optimisation: replace string formatting with concatenation (thanks to @agnivade)

* [`v.1.3.0`](https://github.com/quipo/statsd/releases/tag/v.1.3.0):

    * Added stdout client ("echo" service for debugging)
    * Fixed [issue #23](https://github.com/quipo/statsd/issues/23): GaugeDelta event Stats() should not send an absolute value of 0
    * Fixed FGauge's collation in the buffered client to only preserve the last value in the batch (it mistakenly had the same implementation of FGaugeDelta's collation)
    * Fixed FGaugeDelta with negative value not to send a 0 value first (it mistakenly had the same implementation of FGauge)
    * Added many tests
    * Added compile-time checks that the default events implement the Event interface

* [`v.1.2.0`](https://github.com/quipo/statsd/releases/tag/1.2.0): Sample rate support (thanks to [Hongjian Zhu](https://github.com/hongjianzhu))
*  [`v.1.1.0`](https://github.com/quipo/statsd/releases/tag/1.1.0):

    * Added `SendEvents` function to `Statsd` interface;
    * Using interface in buffered client constructor;
    * Added/Fixed tests

* [`v.1.0.0`](https://github.com/quipo/statsd/releases/tag/1.0.0): First stable release
* `v.0.0.9`: Added memoization to reduce memory allocations
* `v.0.0.8`: Pre-release

## Author

Lorenzo Alberton

* Web: [http://alberton.info](http://alberton.info)
* Twitter: [@lorenzoalberton](https://twitter.com/lorenzoalberton)
* Linkedin: [/in/lorenzoalberton](https://www.linkedin.com/in/lorenzoalberton)


## Copyright

See [LICENSE](LICENSE) document
