package g

import "fmt"

type HistoryData struct {
	Timestamp int64   `json:"timestamp"`
	Value     float64 `json:"value"`
}

type Sms struct {
	Tos     string `json:"tos"`
	Content string `json:"content"`
}

type Mail struct {
	Tos     string `json:"tos"`
	Subject string `json:"subject"`
	Content string `json:"content"`
}

type WeChat struct {
	Tos     string `json:"tos"`
	Content string `json:"content"`
}

type DingTalk struct {
	Tos     string `json:"tos"`
	Content string `json:"content"`
}

func (this *Sms) String() string {
	return fmt.Sprintf(
		"<Tos:%s, Content:%s>",
		this.Tos,
		this.Content,
	)
}

func (this *Mail) String() string {
	return fmt.Sprintf(
		"<Tos:%s, Subject:%s, Content:%s>",
		this.Tos,
		this.Subject,
		this.Content,
	)
}

func (this *WeChat) String() string {
	return fmt.Sprintf(
		"<Tos:%s, Content:%s>",
		this.Tos,
		this.Content,
	)
}

func (this *DingTalk) String() string {
	return fmt.Sprintf(
		"<Tos:%s, Content:%s>",
		this.Tos,
		this.Content,
	)
}
