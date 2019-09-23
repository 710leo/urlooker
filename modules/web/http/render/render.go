package render

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/710leo/urlooker/modules/web/g"
	"github.com/710leo/urlooker/modules/web/http/cookie"
	"github.com/710leo/urlooker/modules/web/http/helper"

	"github.com/gorilla/context"
	"github.com/unrolled/render"
)

var Render *render.Render
var funcMap = template.FuncMap{
	"Times1000":       helper.Times1000,
	"UsersOfTeam":     helper.UsersOfTeam,
	"TeamsOfStrategy": helper.TeamsOfStrategy,
	"HumenTime":       helper.HumenTime,
	"GetFirst":        helper.GetFirst,
}

func Init() {
	Render = render.New(render.Options{
		Directory:  "modules/web/views",
		Extensions: []string{".html"},
		Delims:     render.Delims{"{{", "}}"},
		Funcs:      []template.FuncMap{funcMap},
		IndentJSON: false,
	})
}

func Put(r *http.Request, key string, val interface{}) {
	m, ok := context.GetOk(r, "_DATA_MAP_")
	if ok {
		mm := m.(map[string]interface{})
		mm[key] = val
		context.Set(r, "_DATA_MAP_", mm)
	} else {
		context.Set(r, "_DATA_MAP_", map[string]interface{}{key: val})
	}
}

func HTML(r *http.Request, w http.ResponseWriter, name string, htmlOpt ...render.HTMLOptions) {
	userid, username, found := cookie.ReadUser(r)

	Put(r, "Debug", g.Config.Debug)
	Put(r, "HasLogin", found)
	Put(r, "UserId", userid)
	Put(r, "UserName", username)
	Render.HTML(w, http.StatusOK, name, context.Get(r, "_DATA_MAP_"), htmlOpt...)
}

func Text(w http.ResponseWriter, v string, codes ...int) {
	code := http.StatusOK
	if len(codes) > 0 {
		code = codes[0]
	}
	Render.Text(w, code, v)
}

func MaybeError(w http.ResponseWriter, err error) {
	msg := ""
	if err != nil {
		msg = err.Error()
	}

	Render.JSON(w, http.StatusOK, map[string]string{"msg": msg})
}

func Error(w http.ResponseWriter, err error) {
	msg := ""
	if err != nil {
		msg = err.Error()
	}

	Render.JSON(w, http.StatusOK, map[string]string{"msg": msg})
}

func Message(w http.ResponseWriter, format string, args ...interface{}) {
	Render.JSON(w, http.StatusOK, map[string]string{"msg": fmt.Sprintf(format, args...)})
}

func Data(w http.ResponseWriter, v interface{}, msg ...string) {
	m := ""
	if len(msg) > 0 {
		m = msg[0]
	}

	Render.JSON(w, http.StatusOK, map[string]interface{}{"msg": m, "data": v})
}
