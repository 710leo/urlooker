package http

import (
	"net/http"

	"github.com/710leo/urlooker/modules/web/http/routes"

	"github.com/gorilla/mux"
)

func ConfigRouter(r *mux.Router) {
	configStraRoutes(r)
	configStaticRoutes(r)
	configApiRoutes(r)
	configDomainRoutes(r)
	configAuthRoutes(r)
	configUserRoutes(r)
	configTeamRoutes(r)
	configProcRoutes(r)
}

func configDomainRoutes(r *mux.Router) {
	r.HandleFunc("/url", routes.UrlStatus).Methods("GET")
}

func configChartRoutes(r *mux.Router) {
	r.HandleFunc("/chart", routes.UrlStatus).Methods("GET")
}

func configApiRoutes(r *mux.Router) {
	r.HandleFunc("/api/item/{hostname}", routes.GetHostIpItem).Methods("GET")
}

func configStraRoutes(r *mux.Router) {
	r.HandleFunc("/", routes.HomeIndex).Methods("GET")
	r.HandleFunc("/strategy/add", routes.AddStrategyGet).Methods("GET")
	r.HandleFunc("/strategy/add", routes.AddStrategyPost).Methods("POST")
	r.HandleFunc("/strategy/{id:[0-9]+}", routes.GetStrategyById).Methods("POST")
	r.HandleFunc("/strategy/{id:[0-9]+}/delete", routes.DeleteStrategy).Methods("POST")
	r.HandleFunc("/strategy/{id:[0-9]+}/edit", routes.UpdateStrategyGet).Methods("GET")
	r.HandleFunc("/strategy/{id:[0-9]+}/edit", routes.UpdateStrategy).Methods("POST")
	r.HandleFunc("/strategy/{id:[0-9]+}/teams", routes.GetTeamsOfStrategy).Methods("GET")
}

func configAuthRoutes(r *mux.Router) {
	r.HandleFunc("/auth/register", routes.RegisterPage).Methods("GET")
	r.HandleFunc("/auth/register", routes.Register).Methods("POST")
	r.HandleFunc("/auth/logout", routes.Logout).Methods("GET")
	r.HandleFunc("/auth/login", routes.Login).Methods("POST")
	r.HandleFunc("/auth/login", routes.LoginPage).Methods("GET")
}

func configUserRoutes(r *mux.Router) {
	r.HandleFunc("/me.json", routes.MeJson).Methods("GET")
	r.HandleFunc("/me/profile", routes.UpdateMyProfile).Methods("POST")
	r.HandleFunc("/me/chpwd", routes.ChangeMyPasswd).Methods("POST")
	r.HandleFunc("/users/query", routes.UsersJson).Methods("GET")
}

func configTeamRoutes(r *mux.Router) {
	r.HandleFunc("/teams", routes.TeamsPage).Methods("GET")
	r.HandleFunc("/teams/query", routes.TeamsJson).Methods("GET")
	r.HandleFunc("/team/create", routes.CreateTeamGet).Methods("GET")
	r.HandleFunc("/team/create", routes.CreateTeamPost).Methods("POST")
	r.HandleFunc("/team/{tid:[0-9]+}/edit", routes.UpdateTeamGet).Methods("GET")
	r.HandleFunc("/team/{tid:[0-9]+}/edit", routes.UpdateTeamPost).Methods("POST")
	r.HandleFunc("/team/{tid:[0-9]+}/users", routes.GetUsersOfTeam).Methods("GET")
}

func configProcRoutes(r *mux.Router) {
	//r.HandleFunc("/log", routes.GetLog).Methods("GET")
	r.HandleFunc("/version", routes.Version).Methods("GET")
	r.HandleFunc("/post_test", routes.Health).Methods("POST")
	r.HandleFunc("/put_test", routes.Health).Methods("PUT")
}

func configStaticRoutes(r *mux.Router) {
	r.PathPrefix("/css").Handler(http.FileServer(http.Dir("./modules/web/static")))
	r.PathPrefix("/fonts").Handler(http.FileServer(http.Dir("./modules/web/static")))
	r.PathPrefix("/js").Handler(http.FileServer(http.Dir("./modules/web/static")))
	r.PathPrefix("/img").Handler(http.FileServer(http.Dir("./modules/web/static")))
	r.PathPrefix("/lib").Handler(http.FileServer(http.Dir("./modules/web/static")))
}
