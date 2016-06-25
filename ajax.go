//contains all functions and structs to handle ajax requests
package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"time"
)

func getAjaxSessionID(w http.ResponseWriter, r *http.Request) (string, string) {
	userID := getUserID(r)
	if userID == "" {
		http.Error(w, "invalid session", http.StatusForbidden)
		return "", ""
	}
	if imp := getImpersonate(r); imp != "" {
		return imp, userID
	}
	return userID, userID
}

func addEntry(w http.ResponseWriter, r *http.Request) {
	Info.Println("addEntry()")
	userID, impostor := getAjaxSessionID(w, r)
	if userID == "" {
		return
	}

	err := r.ParseForm()
	if checkErr(err, w) {
		return
	}

	entryType := r.Form.Get("type")
	fromformatString := jsDateTime
	toformatString := jsDateTime
	if len(r.Form.Get("begin")) < 12 {
		fromformatString = jsDate
	}
	if len(r.Form.Get("end")) < 12 {
		toformatString = jsDate
	}
	entryBegin, err := time.Parse(fromformatString, r.Form.Get("begin"))
	if checkErr(err, w) {
		return
	}
	entryEnd, err := time.Parse(toformatString, r.Form.Get("end"))
	if checkErr(err, w) {
		return
	}
	createType := r.Form.Get("create_type")

	entryID := ""
	//notRound := r.Form.Get("notRound")

	if createType == updateType {
		entryID = r.Form.Get("entry_id")
	}
	err = StoreEntry(userID, entryBegin, entryEnd, entryType, createType, entryID, impostor, false, false)
	if checkErr(err, w) {
		return
	}
}

func loadLogs(w http.ResponseWriter, r *http.Request) {
	Info.Println("loadLogs()")
	userID, _ := getAjaxSessionID(w, r)
	if userID == "" {
		return
	}

	err := r.ParseForm()
	if checkErr(err, w) {
		return
	}
	from := r.Form.Get("from_date")
	to := r.Form.Get("to_date")

	dateFrom, dateTo := getDefaultDates()
	if from != "" {
		dateFrom, err = time.Parse(jsDate, from)
		checkErr(err, w)
	}
	if to != "" {
		dateTo, err = time.Parse(jsDate, to)
		checkErr(err, w)
		dur, err := time.ParseDuration("23h59m59s")
		checkErr(err, w)
		dateTo = dateTo.Add(dur)
	}

	setDispalyInterval(dateFrom, dateTo, w, r)

	logs, err := GetLogsForUser(userID, dateFrom, dateTo, false)
	if checkErr(err, w) {
		return
	}

	tpl, err := template.ParseFiles("tpl/table.tpl")
	if checkErr(err, w) {
		return
	}

	err = tpl.ExecuteTemplate(w, "table.tpl", logs)
	if checkErr(err, w) {
		return
	}
}

func edit(w http.ResponseWriter, r *http.Request) {
	Info.Println("edit()")
	userID, _ := getAjaxSessionID(w, r)
	if userID == "" {
		return
	}

	entryID := r.URL.Query().Get("id")
	if entryID == "" {
		http.Error(w, "incomplete request", http.StatusBadRequest)
		return
	}

	entry, err := GetEntry(entryID)
	if checkErr(err, w) {
		return
	}

	encoder := json.NewEncoder(w)
	err = encoder.Encode(entry)
	if checkErr(err, w) {
		return
	}
}

func delete(w http.ResponseWriter, r *http.Request) {
	Info.Println("delete()")
	userID, impostor := getAjaxSessionID(w, r)
	if userID == "" {
		return
	}

	entryID := r.URL.Query().Get("id")
	if entryID == "" {
		http.Error(w, "incomplete request", http.StatusBadRequest)
		return
	}

	err := DeleteEntry(entryID, impostor)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func editUser(w http.ResponseWriter, r *http.Request) {
	Info.Println("editUser()")
	userID, _ := getAjaxSessionID(w, r)
	if userID == "" {
		return
	}

	editID := r.URL.Query().Get("id")
	if editID == "" {
		http.Error(w, "incomplete request", http.StatusBadRequest)
		return
	}

	user, err := GetUser(editID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	encoder := json.NewEncoder(w)
	err = encoder.Encode(user)
	if checkErr(err, w) {
		return
	}
}

func storeUser(w http.ResponseWriter, r *http.Request) {
	Info.Println("storeUser()")
	userID := getUserID(r)
	if userID == "" {
		http.Error(w, "invalid session", http.StatusForbidden)
		return
	}

	user, err := GetUser(userID)
	if checkErr(err, w) {
		return
	}

	if !user.Admin {
		http.Error(w, "no admin privlieges", http.StatusForbidden)
		return
	}

	err = r.ParseForm()
	if checkErr(err, w) {
		return
	}

	editID := r.Form.Get("user_id")
	editName := r.Form.Get("name")
	editUserName := r.Form.Get("username")
	editPassword := r.Form.Get("password")
	editRepeat := r.Form.Get("repeat")
	editType := r.Form.Get("type")

	err = StoreUser(editID, editUserName, editName, editPassword, editRepeat, editType)
	if checkErr(err, w) {
		http.Error(w, err.Error(), http.StatusConflict)
	}
}

func loadUsers(w http.ResponseWriter, r *http.Request) {
	Info.Println("loadUsers()")
	userID := getUserID(r)
	if userID == "" {
		http.Error(w, "invalid session", http.StatusForbidden)
		return
	}

	user, err := GetUser(userID)
	if checkErr(err, w) {
		return
	}

	if !user.Admin {
		http.Error(w, "no admin privlieges", http.StatusForbidden)
		return
	}

	users, err := GetUserList()
	if checkErr(err, w) {
		return
	}

	tpl, err := template.ParseFiles("tpl/userTable.tpl")
	if checkErr(err, w) {
		return
	}

	err = tpl.ExecuteTemplate(w, "userTable.tpl", users)
	if checkErr(err, w) {
		return
	}
}

func toggleUser(w http.ResponseWriter, r *http.Request) {
	Info.Println("toggleUser()")
	userID := getUserID(r)
	if userID == "" {
		http.Error(w, "invalid session", http.StatusForbidden)
		return
	}

	user, err := GetUser(userID)
	if checkErr(err, w) {
		return
	}

	if !user.Admin {
		http.Error(w, "no admin privlieges", http.StatusForbidden)
		return
	}

	toggle := r.URL.Query().Get("action")
	toggleUser := r.URL.Query().Get("id")
	enabled := false
	if toggle == "enable" {
		Info.Println("toggleUser(enabled)")
		enabled = true
	}

	err = UpdateEnabled(toggleUser, enabled)
	checkErr(err, w)
}

func activeEntry(w http.ResponseWriter, r *http.Request) {
	Info.Println("activeEntry()")
	userID, impostor := getAjaxSessionID(w, r)
	if userID == "" {
		return
	}

	entryID := r.URL.Query().Get("id")

	if entryID == "" {
		err := StoreEntry(userID, time.Time{}, time.Time{}, workTimeConst, createConst, entryID, impostor, true, true)
		if checkErr(err, w) {
			return
		}
		newEntry, err := ActiveEntry(userID)
		if checkErr(err, w) {
			return
		}
		w.Write([]byte(newEntry))
		return
	}
	entry, err := GetEntry(entryID)
	if checkErr(err, w) {
		return
	}

	from, err := time.Parse(jsDateTime, fmt.Sprintf("%s %s", entry.DateFrom, entry.TimeFrom))
	checkErr(err, w)

	err = StoreEntry(userID, from, time.Time{}, entry.Type, updateType, entryID, impostor, true, true)
	checkErr(err, w)
}

func stats(w http.ResponseWriter, r *http.Request) {
	Info.Println("stats()")

	userID, _ := getAjaxSessionID(w, r)
	if userID == "" {
		return
	}

	err := r.ParseForm()
	if checkErr(err, w) {
		return
	}
	from := r.Form.Get("from_date")
	to := r.Form.Get("to_date")
	dateFrom, dateTo := getDefaultDates()
	if from != "" {
		dateFrom, err = time.Parse(jsDate, from)
		checkErr(err, w)
	}
	if to != "" {
		dateTo, err = time.Parse(jsDate, to)
		checkErr(err, w)
		dur, err := time.ParseDuration("23h59m59s")
		checkErr(err, w)
		dateTo = dateTo.Add(dur)
	}

	logs, err := GetLogsForUser(userID, dateFrom, dateTo, false)
	if checkErr(err, w) {
		return
	}
	user, err := GetUser(userID)
	if checkErr(err, w) {
		return
	}
	data := anonStruct{
		User: user,
		from: dateFrom,
		to:   dateTo,
	}

	stats := calculateStats(logs, data)
	data.Stats = stats

	tpl, err := template.ParseFiles("tpl/fragments.tpl")
	if checkErr(err, w) {
		return
	}

	err = tpl.ExecuteTemplate(w, "stats", data)
	if checkErr(err, w) {
		return
	}
}
