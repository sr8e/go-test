package handlers

import (
	"html/template"
	"net/http"

	"github.com/sr8e/mellow-ir/db"
)

func MyPage(w http.ResponseWriter, r *http.Request, u *db.User) {
	tmp, err := template.ParseFiles("templates/mypage.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	tmp.Execute(w, u)
}
