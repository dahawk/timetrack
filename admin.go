package main

import (
	"html/template"
	"net/http"
)

func admin(w http.ResponseWriter, r *http.Request) {
	var userID string
	if userID = getUserID(r); userID == "" {
		http.Redirect(w, r, "/login", 303)
		return
	}

	user, err := GetUser(userID)
	if checkErr(err, w) {
		return
	}

	if !user.Admin {
		http.Redirect(w, r, "/", 303)
		return
	}

	users, err := GetUserList()
	if checkErr(err, w) {
		return
	}
	u, err := GetUser(userID)
	if checkErr(err, w) {
		return
	}

	data := anonStruct{
		User:     u,
		UserList: users,
	}

	tpl, err := template.ParseFiles("tpl/admin.tpl", "tpl/fragments.tpl")
	if checkErr(err, w) {
		return
	}

	err = tpl.ExecuteTemplate(w, "admin.tpl", data)
	if checkErr(err, w) {
		return
	}
}
