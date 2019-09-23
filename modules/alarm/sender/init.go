package sender

import (
	"github.com/710leo/urlooker/modules/alarm/g"
)

var (
	SmsWorkerChan  chan int
	MailWorkerChan chan int
)

func Init() {
	workerConfig := g.Config.Worker
	SmsWorkerChan = make(chan int, workerConfig.Sms)
	MailWorkerChan = make(chan int, workerConfig.Mail)

	Consume()
}

func Consume() {
	go ConsumeMail()
	go ConsumeSms()
}
