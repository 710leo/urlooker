package g

import (
	"github.com/toolkits/container/list"
)

var CheckResultQueue *list.SafeLinkedList
var WorkerChan chan int

func Init() {
	WorkerChan = make(chan int, Config.Worker)
	CheckResultQueue = list.NewSafeLinkedList()
}
