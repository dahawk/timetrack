// contains all functions and structs to handle the admin features
package main

import (
	"errors"
	"html/template"
	"net/http"
	"strconv"
	"time"
)

func admin(w http.ResponseWriter, r *http.Request) {
	Info.Println("admin()")
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

func worktime(w http.ResponseWriter, r *http.Request) {
	Info.Println("worktime()")
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

	if r.Method == "POST" {
		Info.Println("POST")
		err = r.ParseForm()
		if checkErr(err, w) {
			return
		}

		row := workTimeRow{
			StartDate: time.Now(),
		}
		wUserID := r.Form.Get("user")

		if r.Form.Get("fulltime") != "fulltime" {
			errOccurred := false
			mon, err := strconv.ParseFloat(r.Form.Get("mon-input"), 64)
			errOccurred = (errOccurred || (err != nil))
			tue, err := strconv.ParseFloat(r.Form.Get("tue-input"), 64)
			errOccurred = (errOccurred || (err != nil))
			wed, err := strconv.ParseFloat(r.Form.Get("wed-input"), 64)
			errOccurred = (errOccurred || (err != nil))
			thu, err := strconv.ParseFloat(r.Form.Get("thu-input"), 64)
			errOccurred = (errOccurred || (err != nil))
			fri, err := strconv.ParseFloat(r.Form.Get("fri-input"), 64)
			errOccurred = (errOccurred || (err != nil))

			if errOccurred {
				checkErr(errors.New("error parsing daily work time"), w)
				return
			}
			row.Mon = mon
			row.Tue = tue
			row.Wed = wed
			row.Thu = thu
			row.Fri = fri
		} else {
			row.FullTime = true
			row.Mon = 7.7
			row.Tue = 7.7
			row.Wed = 7.7
			row.Thu = 7.7
			row.Fri = 7.7
		}

		err = finishCurrentWorkTime(wUserID)
		if checkErr(err, w) {
			return
		}

		err = insertWorkTime(row, wUserID)
		if checkErr(err, w) {
			return
		}
		wUser, err := GetUser(wUserID)
		if checkErr(err, w) {
			return
		}
		workTime, err := getCurrentWorkTime(wUserID)
		if checkErr(err, w) {
			return
		}
		tpl, err := template.ParseFiles("tpl/admin.tpl", "tpl/fragments.tpl")
		if checkErr(err, w) {
			return
		}

		data := anonStruct{
			User:         user,
			WorkTimeUser: wUser,
			WorkTime:     workTime,
		}

		err = tpl.ExecuteTemplate(w, "admin.tpl", data)
		if checkErr(err, w) {
			return
		}
		return
	}
	Info.Println("GET")
	wUserID := r.URL.Query().Get("id")
	wUser, err := GetUser(wUserID)
	if checkErr(err, w) {
		return
	}

	workTime, err := getCurrentWorkTime(wUserID)
	if checkErr(err, w) {
		return
	}

	data := anonStruct{
		User:         user,
		WorkTimeUser: wUser,
		WorkTime:     workTime,
	}

	tpl, err := template.ParseFiles("tpl/worktime.tpl", "tpl/fragments.tpl")
	if checkErr(err, w) {
		return
	}

	err = tpl.ExecuteTemplate(w, "worktime.tpl", data)
	if checkErr(err, w) {
		return
	}
}
