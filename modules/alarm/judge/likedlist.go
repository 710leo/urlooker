package judge

import (
	"container/list"
	"sync"

	"github.com/710leo/urlooker/dataobj"
)

type SafeLinkedList struct {
	sync.RWMutex
	L *list.List
}

func (this *SafeLinkedList) HistoryData(limit int) ([]*HistoryData, bool) {
	if limit < 1 {
		return []*HistoryData{}, false
	}

	size := this.Len()
	if size <= limit {
		return []*HistoryData{}, false
	}

	firstElement := this.Front()
	firstItem := firstElement.Value.(*dataobj.ItemStatus)

	vs := make([]*HistoryData, limit)
	vs[0] = &HistoryData{Timestamp: firstItem.PushTime, Value: firstItem.Result}

	currentElement := firstElement
	i := 1
	for i < limit {
		nextElement := currentElement.Next()
		vs[i] = &HistoryData{
			Timestamp: nextElement.Value.(*dataobj.ItemStatus).PushTime,
			Value:     nextElement.Value.(*dataobj.ItemStatus).Result,
		}
		i++
		currentElement = nextElement
	}

	return vs, true
}

func (this *SafeLinkedList) PushFrontAndMaintain(v *dataobj.ItemStatus, maxCount int) bool {
	this.Lock()
	defer this.Unlock()

	sz := this.L.Len()
	if sz > 0 {
		// 新push上来的数据有可能重复了，或者timestamp不对，这种数据要丢掉
		if v.PushTime <= this.L.Front().Value.(*dataobj.ItemStatus).PushTime || v.PushTime <= 0 {
			return false
		}
	}

	this.L.PushFront(v)

	sz++
	if sz <= maxCount {
		return true
	}

	del := sz - maxCount
	for i := 0; i < del; i++ {
		this.L.Remove(this.L.Back())
	}

	return true
}

func (this *SafeLinkedList) Front() *list.Element {
	this.RLock()
	defer this.RUnlock()
	return this.L.Front()
}

func (this *SafeLinkedList) Len() int {
	this.RLock()
	defer this.RUnlock()
	return this.L.Len()
}
