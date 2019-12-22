#smtp

## demo

```go
package main

import (
	"log"

	"github.com/toolkits/smtp"
)

func main() {
	//s := smtp.New("smtp.exmail.qq.com:25", "notify@a.com", "password")
	s := smtp.NewSMTP("smtp.exmail.qq.com:25", "notify@a.com", "password",false,false,false)
	log.Println(s.SendMail("notify@a.com", "ulric@b.com;rain@c.com", "这是subject", "这是body,<font color=red>red</font>"))
}
```