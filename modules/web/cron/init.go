package cron

func Init() {
	go GetDetectedItem()
	go DeleteOld()
}
