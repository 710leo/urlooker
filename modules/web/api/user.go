package api

import (
	"log"
	"strconv"
	"strings"

	"github.com/710leo/urlooker/dataobj"
	"github.com/710leo/urlooker/modules/web/model"
)

type UsersResponse struct {
	Message string
	Data    []*dataobj.User
}

func (this *Web) GetUsersByTeam(req string, reply *UsersResponse) error {
	tids := strings.Split(req, ",")
	if len(tids) < 1 || tids[0] == "" {
		reply.Message = "user no exists!"
		return nil
	}
	allUsers := make([]*dataobj.User, 0)
	for _, tid := range tids {
		id, err := strconv.ParseInt(tid, 10, 64)
		if err != nil {
			log.Println("tid error:", err)
			continue
		}
		users, err := model.UsersInfoOfTeam(id)
		if err != nil {
			reply.Message = err.Error()
		}

		for _, user := range users {
			u := &dataobj.User{
				Id:     user.Id,
				Name:   user.Name,
				Cnname: user.Cnname,
				Phone:  user.Phone,
				Wechat: user.Wechat,
				Role:   user.Role,
			}
			allUsers = append(allUsers, u)
		}
	}

	reply.Data = allUsers
	return nil
}
