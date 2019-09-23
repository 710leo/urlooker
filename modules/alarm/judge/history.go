package judge

import (
	"container/list"
	"sync"

	"github.com/710leo/urlooker/dataobj"
)

type HistoryData struct {
	Timestamp int64 `json:"timestamp"`
	Value     int64 `json:"value"`
}

var HistoryBigMap = make(map[string]*SafeItemMap)

func InitHistoryBigMap() {
	arr := []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "a", "b", "c", "d", "e", "f"}
	for i := 0; i < 16; i++ {
		for j := 0; j < 16; j++ {
			HistoryBigMap[arr[i]+arr[j]] = NewSafeItemMap()
		}
	}
}

type SafeItemMap struct {
	sync.RWMutex
	M map[string]*SafeLinkedList
}

func NewSafeItemMap() *SafeItemMap {
	return &SafeItemMap{M: make(map[string]*SafeLinkedList)}
}

func (this *SafeItemMap) Get(key string) (*SafeLinkedList, bool) {
	this.RLock()
	defer this.RUnlock()
	val, ok := this.M[key]
	return val, ok
}

func (this *SafeItemMap) Set(key string, val *SafeLinkedList) {
	this.Lock()
	defer this.Unlock()
	this.M[key] = val
}

func (this *SafeItemMap) Len() int {
	this.RLock()
	defer this.RUnlock()
	return len(this.M)
}

func (this *SafeItemMap) Delete(key string) {
	this.Lock()
	defer this.Unlock()
	delete(this.M, key)
}

func (this *SafeItemMap) BatchDelete(keys []string) {
	count := len(keys)
	if count == 0 {
		return
	}

	this.Lock()
	defer this.Unlock()
	for i := 0; i < count; i++ {
		delete(this.M, keys[i])
	}
}

func (this *SafeItemMap) CleanStale(before int64) {
	keys := []string{}

	this.RLock()
	for key, L := range this.M {
		front := L.Front()
		if front == nil {
			continue
		}

		if front.Value.(*dataobj.ItemStatus).PushTime < before {
			keys = append(keys, key)
		}
	}
	this.RUnlock()

	this.BatchDelete(keys)
}

func (this *SafeItemMap) PushFrontAndMaintain(key string, val *dataobj.ItemStatus, maxCount int, now int64) {
	if linkedList, exists := this.Get(key); exists {
		needJudge := linkedList.PushFrontAndMaintain(val, maxCount)
		if needJudge {
			Judge(linkedList, val, now)
		}
	} else {
		linkedList := list.New()
		linkedList.PushFront(val)
		safeList := &SafeLinkedList{L: linkedList}
		this.Set(key, safeList)
		Judge(safeList, val, now)
	}
}
