package sender

import (
	"fmt"
	"strings"
	"time"

	"github.com/710leo/urlooker/dataobj"
	"github.com/710leo/urlooker/modules/alarm/cache"
	"github.com/710leo/urlooker/modules/alarm/g"
)

func BuildMail(event *dataobj.Event) string {
	strategy, _ := cache.StrategyMap.Get(event.StrategyId)
	respTime := fmt.Sprintf("%dms", event.RespTime)
	return fmt.Sprintf(
		"状态:%s\n结果:%s\nUrl:%s\nIP:%s\n返回状态码:%s\n响应时间:%s\n时间:%s\n报警次数:%d\n备注:%s\n",
		event.Status,
		g.EventStatus[event.Result],
		event.Url,
		event.Ip,
		event.RespCode,
		respTime,
		humanTime(event.EventTime),
		event.CurrentStep,
		strategy.Note,
	)
}

func BuildSms(event *dataobj.Event) string {
	respTime := fmt.Sprintf("%dms", event.RespTime)
	return fmt.Sprintf(
		"[%s][%s %s][%s][%s][%s][O%d]",
		event.Status,
		showSubString(event.Url, 100),
		event.Ip,
		event.RespCode,
		respTime,
		humanTime(event.EventTime),
		event.CurrentStep,
	)
}

func BuildDingWebhook(event *dataobj.Event) string {
	strategy, _ := cache.StrategyMap.Get(event.StrategyId)
	respTime := fmt.Sprintf("%dms", event.RespTime)
	return fmt.Sprintf(
		"- 告警: %s\n- 状态: %s\n- Url: %s\n- IP: %s\n- 返回状态码: %s\n- 响应时间: %s\n- 时间: %s\n- 报警次数: %d\n",
		strategy.Note,
		event.Status,
		event.Url,
		event.Ip,
		event.RespCode,
		respTime,
		humanTime(event.EventTime),
		event.CurrentStep,
	)
}

func BuildWeChat(event *dataobj.Event) string {
	strategy, _ := cache.StrategyMap.Get(event.StrategyId)
	respTime := fmt.Sprintf("%dms", event.RespTime)
	return fmt.Sprintf(
		"状态:%s\nUrl:%s\n备注:%s\nIP:%s\n返回状态码:%s\n响应时间:%s\n时间:%s\n报警次数:%d\n",
		event.Status,
		event.Url,
		strategy.Note,
		event.Ip,
		event.RespCode,
		respTime,
		humanTime(event.EventTime),
		event.CurrentStep,
	)
}

func humanTime(ts int64) string {
	return time.Unix(ts, 0).Format("2006-01-02 15:04:05")
}

func showSubString(str string, length int) string {
	runeStr := []rune(str)
	s := ""
	if length > len(runeStr) {
		length = len(runeStr)
	}

	for i := 0; i < length; i++ {
		s += string(runeStr[i])
	}
	return s
}

func WriteSms(tos []string, content string) {
	if !g.Config.SmsEnabled {
		return
	}

	if len(tos) == 0 {
		return
	}

	sms := &g.Sms{Tos: strings.Join(tos, ","), Content: content}
	SmsWorkerChan <- 1
	go sendSms(sms.Tos, sms.Content)
}

func WriteMail(tos []string, subject, content string) {
	if !g.Config.Smtp.Enabled {
		return
	}

	if len(tos) == 0 {
		return
	}

	mail := &g.Mail{Tos: strings.Join(tos, ","), Subject: subject, Content: content}
	MailWorkerChan <- 1
	go sendMail(mail)
}

func WriteWeChat(tos []string, content string) {
	if !g.Config.WeChat.Enabled {
		return
	}
	if len(tos) == 0 {
		return
	}
	weChat := &g.WeChat{Tos: strings.Join(tos, "|"), Content: content}
	WeChatWorkerChan <- 1
	go sendWeChat(weChat)
}
