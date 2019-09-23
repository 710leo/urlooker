package cache

import (
	"sync"

	"github.com/710leo/urlooker/dataobj"
)

type SafeStrategyMap struct {
	sync.RWMutex
	M map[int64]dataobj.Strategy
}

var StrategyMap = &SafeStrategyMap{M: make(map[int64]dataobj.Strategy)}

func (this *SafeStrategyMap) ReInit(m map[int64]dataobj.Strategy) {
	this.Lock()
	defer this.Unlock()
	this.M = m
}

func (this *SafeStrategyMap) Get(key int64) (dataobj.Strategy, bool) {
	this.RLock()
	defer this.RUnlock()
	s, exists := this.M[key]
	return s, exists
}
