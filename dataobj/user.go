package dataobj

import "time"

type User struct {
	Id       int64     `json:"id"`
	Name     string    `json:"name"`
	Cnname   string    `json:"cnname"`
	Password string    `json:"-"`
	Email    string    `json:"email"`
	Phone    string    `json:"phone"`
	Wechat   string    `json:"wechat"`
	Role     int       `json:"role"`
	Created  time.Time `json:"-" xorm:"<-"`
}
