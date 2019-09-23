package routes

import (
	"net/http"

	"github.com/toolkits/str"
	"github.com/toolkits/web"

	"github.com/710leo/urlooker/modules/web/http/errors"
	"github.com/710leo/urlooker/modules/web/http/param"
	"github.com/710leo/urlooker/modules/web/http/render"
	"github.com/710leo/urlooker/modules/web/model"
)

func HomeIndex(w http.ResponseWriter, r *http.Request) {
	me := MeRequired(LoginRequired(w, r))
	username := me.Name
	mine := param.Int(r, "mine", 1)
	query := param.String(r, "q", "")
	if str.HasDangerousCharacters(query) {
		errors.Panic("查询字符不合法")
	}

	if IsAdmin(username) {
		mine = 0
	}

	limit := param.Int(r, "limit", 10)
	total, err := model.GetAllStrategyCount(mine, query, username)
	errors.MaybePanic(err)
	pager := web.NewPaginator(r, limit, total)

	strategies, err := model.GetAllStrategy(mine, limit, pager.Offset(), query, username)

	errors.MaybePanic(err)
	render.Put(r, "Strategies", strategies)
	render.Put(r, "Pager", pager)
	render.Put(r, "Mine", mine)
	render.Put(r, "Query", query)
	render.HTML(r, w, "home/index")
}
