package routes

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/toolkits/str"

	"github.com/710leo/urlooker/modules/web/g"
	"github.com/710leo/urlooker/modules/web/http/cookie"
	"github.com/710leo/urlooker/modules/web/http/errors"
	"github.com/710leo/urlooker/modules/web/model"
)

func StraRequired(r *http.Request) *model.Strategy {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	errors.MaybePanic(err)

	obj, err := model.GetStrategyById(id)
	errors.MaybePanic(err)
	if obj == nil {
		panic(errors.BadRequestError("no such item"))
	}
	return obj
}

func BindJson(r *http.Request, obj interface{}) error {
	if r.Body == nil {
		return fmt.Errorf("Empty request body")
	}
	defer r.Body.Close()
	body, _ := ioutil.ReadAll(r.Body)

	err := json.Unmarshal(body, obj)
	if err != nil {
		return fmt.Errorf("unmarshal body %s err:%v", string(body), err)
	}
	return err
}

func HostnameRequired(r *http.Request) string {
	vars := mux.Vars(r)
	hostname := vars["hostname"]

	if str.HasDangerousCharacters(hostname) {
		errors.Panic("hostname不合法")
	}

	return hostname
}

func LoginRequired(w http.ResponseWriter, r *http.Request) (int64, string) {
	userid, username, found := cookie.ReadUser(r)
	if !found {
		panic(errors.NotLoginError())
	}

	return userid, username
}

func AdminRequired(id int64, name string) {
	user, err := model.GetUserById(id)
	if err != nil {
		panic(errors.InternalServerError(err.Error()))
	}

	if user == nil {
		panic(errors.NotLoginError())
	}

	for _, admin := range g.Config.Admins {
		if user.Name == admin {
			return
		}
	}

	panic(errors.NotLoginError())
	return
}

func MeRequired(id int64, name string) *model.User {
	user, err := model.GetUserById(id)
	if err != nil {
		panic(errors.InternalServerError(err.Error()))
	}

	if user == nil {
		panic(errors.NotLoginError())
	}

	return user
}

func TeamRequired(r *http.Request) *model.Team {
	vars := mux.Vars(r)
	tid, err := strconv.ParseInt(vars["tid"], 10, 64)
	errors.MaybePanic(err)

	team, err := model.GetTeamById(tid)
	errors.MaybePanic(err)
	if team == nil {
		panic(errors.BadRequestError("no such team"))
	}

	return team
}

func UserMustBeMemberOfTeam(uid, tid int64) {
	is, err := model.IsMemberOfTeam(uid, tid)
	errors.MaybePanic(err)
	if is {
		return
	}

	team, err := model.GetTeamById(tid)
	errors.MaybePanic(err)
	if team != nil && team.Creator == uid {
		return
	}

	panic(errors.BadRequestError("用户不是团队的成员"))
}

func IsAdmin(username string) bool {
	for _, admin := range g.Config.Admins {
		if username == admin {
			return true
		}
	}
	return false
}
