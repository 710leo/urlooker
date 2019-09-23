package routes

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/toolkits/str"
	"github.com/toolkits/web"

	"github.com/710leo/urlooker/modules/web/http/errors"
	"github.com/710leo/urlooker/modules/web/http/param"
	"github.com/710leo/urlooker/modules/web/http/render"
	"github.com/710leo/urlooker/modules/web/model"
)

func TeamsPage(w http.ResponseWriter, r *http.Request) {
	me := MeRequired(LoginRequired(w, r))

	query := param.String(r, "q", "")
	if str.HasDangerousCharacters(query) {
		errors.Panic("查询字符不合法")
	}
	limit := param.Int(r, "limit", 10)

	total, err := model.TeamCountOfUser(query, me.Id)
	errors.MaybePanic(err)
	pager := web.NewPaginator(r, limit, total)
	teams, err := model.TeamsOfUser(query, me.Id, limit, pager.Offset())
	errors.MaybePanic(err)
	for _, team := range teams {
		user, err := model.GetUserById(team.Creator)
		if err == nil && user != nil {
			team.CreatorName = user.Name
		}
	}

	render.Put(r, "Teams", teams)
	render.Put(r, "Query", query)
	render.Put(r, "Pager", pager)
	render.Put(r, "Me", me)
	render.Put(r, "Title", "Team")
	render.HTML(r, w, "team/index")
}

func TeamsJson(w http.ResponseWriter, r *http.Request) {
	MeRequired(LoginRequired(w, r))

	query := param.String(r, "query", "")
	limit := param.Int(r, "limit", 10)

	if str.HasDangerousCharacters(query) {
		render.Data(w, fmt.Errorf("query invalid"))
		return
	}

	teams, err := model.QueryTeams(query, limit)
	errors.MaybePanic(err)

	render.Data(w, map[string]interface{}{"teams": teams})
}

func CreateTeamGet(w http.ResponseWriter, r *http.Request) {
	me := MeRequired(LoginRequired(w, r))

	render.Put(r, "Me", me)
	render.Put(r, "Title", "Team")
	render.HTML(r, w, "team/create")
}

func CreateTeamPost(w http.ResponseWriter, r *http.Request) {
	me := MeRequired(LoginRequired(w, r))

	name := param.MustString(r, "name")
	if str.HasDangerousCharacters(name) {
		errors.Panic("team名称不合法")
	}
	resume := param.String(r, "resume", "")
	if str.HasDangerousCharacters(resume) {
		errors.Panic("resume不合法")
	}

	uidsStr := param.String(r, "users", "")
	if str.HasDangerousCharacters(uidsStr) {
		errors.Panic("users不合法")
	}
	uidSlice := strings.Split(uidsStr, ",")

	isci := false
	uids := make([]int64, 0)
	for _, u := range uidSlice {
		if u == "" {
			continue
		}
		uid, err := strconv.ParseInt(u, 10, 64)
		errors.MaybePanic(err)
		uids = append(uids, uid)
		if uid == me.Id {
			isci = true
		}
	}
	if !isci {
		// creator is member of team
		uids = append(uids, me.Id)
	}

	_, err := model.AddTeam(name, resume, me.Id, uids)
	render.MaybeError(w, err)
}

func UpdateTeamGet(w http.ResponseWriter, r *http.Request) {
	team := TeamRequired(r)
	me := MeRequired(LoginRequired(w, r))
	if !IsAdmin(me.Name) {
		UserMustBeMemberOfTeam(me.Id, team.Id)
	}

	uids := make([]string, 0)
	users, err := model.UsersOfTeam(team.Id)
	errors.MaybePanic(err)
	for _, user := range users {
		uids = append(uids, strconv.FormatInt(user.Id, 10))
	}

	render.Put(r, "Team", team)
	render.Put(r, "Uids", strings.Join(uids, ","))
	render.Put(r, "Me", me)
	render.Put(r, "Title", "Team")
	render.HTML(r, w, "team/edit")
}

func UpdateTeamPost(w http.ResponseWriter, r *http.Request) {
	me := MeRequired(LoginRequired(w, r))
	team := TeamRequired(r)
	if !IsAdmin(me.Name) {
		UserMustBeMemberOfTeam(me.Id, team.Id)
	}

	team.Resume = param.String(r, "resume", "")
	if str.HasDangerousCharacters(team.Resume) {
		errors.Panic("resume不合法")
	}
	uidsStr := param.String(r, "users", "")
	if str.HasDangerousCharacters(uidsStr) {
		errors.Panic("users不合法")
	}
	uidsSlice := strings.Split(uidsStr, ",")
	uids := make([]int64, 0)
	for _, uidStr := range uidsSlice {
		if uidStr == "" {
			continue
		}
		uid, err := strconv.ParseInt(uidStr, 10, 64)
		errors.MaybePanic(err)
		uids = append(uids, uid)
	}

	render.Data(w, team.Update(uids))
}

func GetUsersOfTeam(w http.ResponseWriter, r *http.Request) {
	MeRequired(LoginRequired(w, r))
	team := TeamRequired(r)

	users, err := model.UsersOfTeam(team.Id)
	errors.MaybePanic(err)
	render.Data(w, users)
}
