package sender

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
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
	SmsWorkerChan         chan int
	MailWorkerChan        chan int
	WeChatWorkerChan      chan int
	DingWebhookWorkerChan chan int
	requestError          = errors.New("request error,check url or network")
)

const (
	// 发送消息使用导的url
	sendUrl = "https://qyapi.weixin.qq.com/cgi-bin/message/send?access_token="
	// 获取token使用导的url
	getToken = "https://qyapi.weixin.qq.com/cgi-bin/gettoken?corpid="
)

func Init() {
	workerConfig := g.Config.Worker
	SmsWorkerChan = make(chan int, workerConfig.Sms)
	MailWorkerChan = make(chan int, workerConfig.Mail)
	WeChatWorkerChan = make(chan int, workerConfig.WeChat)
	DingWebhookWorkerChan = make(chan int, workerConfig.DingWebhook)
}

func SendEvent(event *dataobj.Event) {
	mail := make([]string, 0)
	sms := make([]string, 0)
	weChat := make([]string, 0)
	users := getUsers(event.StrategyId)

	mailContent := BuildMail(event)
	smsContent := BuildSms(event)
	weChatContent := BuildWeChat(event)
	for _, user := range users {
		mail = append(mail, user.Email)
		sms = append(sms, user.Phone)
		weChat = append(weChat, user.Wechat)
	}

	WriteSms(sms, smsContent)
	WriteMail(mail, smsContent, mailContent)
	WriteWeChat(weChat, weChatContent)

	if g.Config.DingWebhook.Enabled {
		strategy, exists := cache.StrategyMap.Get(event.StrategyId)
		if !exists {
			log.Printf("strategyId: %d not exists", event.StrategyId)
		} else {
			content := BuildDingWebhook(event)
			title := "[告警] " + strategy.Note

			if strategy.DingWebhook != "" {
				log.Println("sending Ding Webhook(strategy)", strategy.DingWebhook)
				DingWebhookWorkerChan <- 1
				go sendDingWebhook(strategy.DingWebhook, title, content)
			}

			if g.Config.DingWebhook.Addr != "" && g.Config.DingWebhook.Addr != strategy.DingWebhook {
				log.Println("sending Ding Webhook(global)", g.Config.DingWebhook.Addr)
				DingWebhookWorkerChan <- 1
				go sendDingWebhook(g.Config.DingWebhook.Addr, title, content)
			}
		}
	}
}

type Markdown struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}

type At struct {
	AtMobiles []string `json:"atMobiles"`
	IsAtAll   bool     `json:"isAtAll"`
}

type DingWebhookMsg struct {
	MsgType  string    `json:"msgtype"`
	Markdown *Markdown `json:"markdown"`
	At       *At       `jons:"at"`
}

func sendDingWebhook(webhook, title, content string) {
	defer func() {
		<-DingWebhookWorkerChan
	}()

	log.Printf("send ding_webhook(%s):\n %s", webhook, content)

	msg := &DingWebhookMsg{
		MsgType: "markdown",
		Markdown: &Markdown{
			Title: title,
			Text:  content,
		},
	}
	buf, err := json.Marshal(msg)

	if err != nil {
		log.Printf("marshal ding_webhook msg err:%q", err)
		return
	}

	err = SendMsg(webhook, buf)
	if err != nil {
		log.Printf("send ding webhook msg err: %q", err)
	} else {
		log.Printf("send ding webhook success")
	}
}

func sendSms(phone string, sms string) {
	defer func() {
		<-SmsWorkerChan
	}()

	smsShell := path.Join(file.SelfDir(), "script", "sms.sh")
	if !file.IsExist(smsShell) {
		log.Printf("%s not found", smsShell)
		return
	}

	cmd := exec.Command(smsShell, phone, "'"+sms+"'")
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	err := cmd.Start()
	if err != nil {
		log.Printf("start cmd err: %v shell:%v", err, smsShell)
	}
	err, isTimeout := sys.CmdRunWithTimeout(cmd, time.Second*10)
	log.Printf("%s %s %s", smsShell, phone, sms)
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
		log.Printf("send mail err:%v tos:%v\n", err, mail.Tos)
		return
	}

	if g.Config.Debug {
		log.Println("==mail==>>>>", mail)
	}
}

func sendWeChat(weChat *g.WeChat) {
	defer func() {
		<-WeChatWorkerChan
	}()

	var msg = weChatMsg{
		ToUser:  weChat.Tos,
		ToParty: g.Config.WeChat.ToParty,
		MsgType: "text",
		AgentId: g.Config.WeChat.AgentId,
		Text:    map[string]string{"content": weChat.Content},
	}
	token, err := GetToken(g.Config.WeChat.CorpId, g.Config.WeChat.CorpSecret)
	buf, err := json.Marshal(msg)
	if err != nil {
		log.Println(err, "get weChat token error")
	}
	url := sendUrl + token.AccessToken
	err = SendMsg(url, buf)
	if err != nil {
		log.Println(err, "send weChat")
	} else {
		log.Println("==weChat==>>>>", weChat)
		log.Println("<<<<==weChatMsg==", string(buf))
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

// 定义微信文本消息结构体
type weChatMsg struct {
	ToUser  string            `json:"touser"`
	ToParty string            `json:"toparty"`
	MsgType string            `json:"msgtype"`
	AgentId int               `json:"agentid"`
	Text    map[string]string `json:"text"`
	Safe    int               `json:"safe"`
}

// 定义微信错误返回结构体
// 钉钉Webhook返回结构体
type sendMsgError struct {
	ErrCode int    `json:"errcode`
	ErrMsg  string `json:"errmsg"`
}

// 定义token结构体
type accessToken struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

// 获取返回
func GetToken(corpId, corpSecret string) (at accessToken, err error) {
	resp, err := http.Get(getToken + corpId + "&corpSecret=" + corpSecret)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		err = requestError
		return
	}
	buf, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(buf, &at)
	if at.AccessToken == "" {
		err = errors.New("微信企业号中的标识或者应用Secret错误")
	}
	return
}

// 发送微信
func SendMsg(url string, msgBody []byte) error {
	body := bytes.NewBuffer(msgBody)
	resp, err := http.Post(url, "application/json", body)
	if resp.StatusCode != 200 {
		return requestError
	}
	buf, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	var e sendMsgError
	err = json.Unmarshal(buf, &e)
	if err != nil {
		return err
	}
	if e.ErrCode != 0 && e.ErrMsg != "ok" {
		return errors.New(string(buf))
	}
	return nil
}
