package web

import (
	"gitlab.com/systemz/tasktab/model"
	"net/http"
)

type CountLogPage struct {
	AuthOk         bool
	User           model.User
	Counters       []model.CounterSessionList
	CounterRunning bool
}

func CountLog(w http.ResponseWriter, r *http.Request) {
	var page CountLogPage
	authOk, user := CheckAuth(w, r)

	page.User = user
	page.AuthOk = authOk
	page.Counters = model.CounterLogList(user.Id)
	display.HTML(w, http.StatusOK, "count_log", page)
}
