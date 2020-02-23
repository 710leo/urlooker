package sender

import (
	"log"
	"os/exec"
	"path"
	"strings"
	"syscall"
	"time"

	"github.com/710leo/urlooker/dataobj"
	"github.com/710leo/urlooker/modules/alarm/backend"
	"github.com/710leo/urlooker/modules/alarm/cache"
	"github.com/710leo/urlooker/modules/alarm/g"
	"github.com/710leo/urlooker/modules/web/api"

	"github.com/toolkits/file"
	"github.com/toolkits/smtp"
	"github.com/toolkits/sys"
)

var (
	SmsWorkerChan  chan int
	MailWorkerChan chan int
)

func Init() {
	workerConfig := g.Config.Worker
	SmsWorkerChan = make(chan int, workerConfig.Sms)
	MailWorkerChan = make(chan int, workerConfig.Mail)
}

func SendEvent(event *dataobj.Event) {
	mail := make([]string, 0)
	sms := make([]string, 0)
	users := getUsers(event.StrategyId)

	mailContent := BuildMail(event)
	smsContent := BuildSms(event)
	for _, user := range users {
		mail = append(mail, user.Email)
		sms = append(sms, user.Phone)
	}

	WriteSms(sms, smsContent)
	WriteMail(mail, smsContent, mailContent)
}

func sendSms(phone string, sms string) {
	defer func() {
		<-SmsWorkerChan
	}()

	sms_shell := path.Join(file.SelfDir(), "script", "sms.sh")
	if !file.IsExist(sms_shell) {
		log.Printf("%s not found", sms_shell)
		return
	}

	cmd := exec.Command(sms_shell, phone, "'"+sms+"'")
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	err := cmd.Start()
	if err != nil {
		log.Printf("start cmd err: %v shell:%v", err, sms_shell)
	}
	err, isTimeout := sys.CmdRunWithTimeout(cmd, time.Second*10)
	log.Printf("%s %s %s", sms_shell, phone, sms)
	if err != nil {
		log.Printf("err: %v, isTimeout: %v", err, isTimeout)
	}

	return
}

func sendMail(mail *g.Mail) {
	defer func() {
		<-MailWorkerChan
	}()

	//s := smtp.New(g.Config.Smtp.Addr, g.Config.Smtp.Username, g.Config.Smtp.Password)
	s := smtp.NewSMTP(g.Config.Smtp.Addr, g.Config.Smtp.Username, g.Config.Smtp.Password, g.Config.Smtp.Tls, false, false)
	err := s.SendMail(g.Config.Smtp.From, strings.Replace(mail.Tos, ",", ";", -1), mail.Subject, mail.Content, "text")
	if err != nil {
		log.Println(err, "tos:", mail.Tos)
		//SendSmsToMaintainer("sender:" + err.Error())
	}

	if g.Config.Debug {
		resp := "ok"
		if err != nil {
			resp = err.Error()
		}
		log.Println("==mail==>>>>", mail)
		log.Println("<<<<==mail==", resp)
	}
}

func getUsers(sid int64) []*dataobj.User {
	var usersResp api.UsersResponse
	var users []*dataobj.User
	strategy, exists := cache.StrategyMap.Get(sid)
	if !exists {
		log.Printf("strategyId: %d not exists", sid)
		return users
	}

	err := backend.CallRpc("Web.GetUsersByTeam", strategy.Teams, &usersResp)
	if err != nil {
		log.Println("Web.GetUsersByTeam Error:", err)
		return users
	}

	if usersResp.Message != "" {
		log.Println("Web.GetUsersByTeam Error:", usersResp.Message)
		return users
	}
	users = usersResp.Data

	return users
}
