package sender

import (
	"log"
	"os/exec"
	"path"
	"strings"
	"syscall"
	"time"

	"github.com/toolkits/file"
	"github.com/toolkits/sys"

	"github.com/710leo/urlooker/modules/alarm/g"
)

func ConsumeSms() {
	queue := g.Config.Queue.Sms
	for {
		L := PopAllSms(queue)
		if len(L) == 0 {
			time.Sleep(time.Millisecond * 200)
			continue
		}
		SendSmsList(L)
	}
}

func SendSmsList(L []*g.Sms) {
	for _, sms := range L {
		if sms.Tos == "" || sms.Content == "" {
			continue
		}

		toArr := strings.Split(sms.Tos, ",")
		log.Println("SmsCount", len(toArr))

		SmsWorkerChan <- 1
		go sendSms(sms.Tos, sms.Content)
	}
}

func sendSms(phone string, sms string) {
	defer func() {
		<-SmsWorkerChan
	}()

	sms_shell := path.Join(file.SelfDir(), "script", "send.sms.sh")
	if !file.IsExist(sms_shell) {
		log.Printf("%s not found", sms_shell)
		return
	}

	cmd := exec.Command(sms_shell, phone, "'"+sms+"'")
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	cmd.Start()
	err, isTimeout := sys.CmdRunWithTimeout(cmd, time.Second*10)
	log.Printf("%s %s %s", sms_shell, phone, sms)
	if err != nil {
		log.Printf("err: %v, isTimeout: %v", err, isTimeout)
	}

	return
}
