package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/securecookie"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var jsDate = "02.01.2006"
var jsTime = "15:04"
var jsDateTime = "02.01.2006 15:04"
var db *sqlx.DB
var cookieHandler = securecookie.New(
	securecookie.GenerateRandomKey(64),
	securecookie.GenerateRandomKey(32))

const salt = "<set secret salt here>"

type anonStruct struct {
	User     User
	UserList []User
	Error    string
	From     string
	To       string
}

func main() {
	var err error

	db, err = sqlx.Open("postgres", fmt.Sprintf("host=%s user=%s dbname=%s password=%s sslmode=disable",
		"<host>", "<username>", "<db name>", "<password>"))
	if err != nil {
		log.Fatalln(err)
	}

	mux := http.NewServeMux()
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	mux.HandleFunc("/", index)
	mux.HandleFunc("/addEntry", addEntry)
	mux.HandleFunc("/loadLogs", loadLogs)
	mux.HandleFunc("/edit", edit)
	mux.HandleFunc("/delete", delete)
	mux.HandleFunc("/login", login)
	mux.HandleFunc("/logout", logout)
	mux.HandleFunc("/activeEntry", activeEntry)
	mux.HandleFunc("/pdf", pdf)
	mux.HandleFunc("/admin", admin)
	mux.HandleFunc("/editUser", editUser)
	mux.HandleFunc("/storeUser", storeUser)
	mux.HandleFunc("/loadUsers", loadUsers)
	mux.HandleFunc("/toggleUser", toggleUser)

	http.ListenAndServe(":1234", mux)
}

func index(w http.ResponseWriter, r *http.Request) {
	var userID string
	if userID = getUserID(r); userID == "" {
		http.Redirect(w, r, "/login", 303)
		return
	}

	tpl, err := template.ParseFiles("tpl/index.tpl", "tpl/fragments.tpl")
	if checkErr(err, w) {
		return
	}
	trackingData, err := GetUser(userID)
	if checkErr(err, w) {
		return
	}

	dateFrom, dateTo, err := getDisplayInterval(r)
	if err != nil {
		dateFrom, dateTo = getDefaultDates()
	}

	data := anonStruct{
		User: trackingData,
		From: dateFrom.Format(jsDate),
		To:   dateTo.Format(jsDate),
	}

	err = tpl.ExecuteTemplate(w, "index.tpl", data)
	if checkErr(err, w) {
		return
	}
}

func login(w http.ResponseWriter, r *http.Request) {
	tpl, err := template.ParseFiles("tpl/login.tpl", "tpl/fragments.tpl")
	if checkErr(err, w) {
		return
	}

	if r.Method == "POST" {
		err = r.ParseForm()
		if checkErr(err, w) {
			return
		}
		username := r.Form.Get("user")
		password := r.Form.Get("password")
		id, msg := VerifyLogin(username, password)
		if msg != "" || id == "" {
			data := anonStruct{
				User: User{
					Username: "",
				},
				Error: msg,
			}

			tpl.ExecuteTemplate(w, "login.tpl", data)
			return
		}

		setSession(id, w)

		http.Redirect(w, r, "/", 307)
		return
	}

	data := anonStruct{
		User: User{
			Username: "",
		},
	}

	tpl.ExecuteTemplate(w, "login.tpl", data)
}

func setSession(userID string, w http.ResponseWriter) {
	value := map[string]string{
		"userID": userID,
	}

	if encoded, err := cookieHandler.Encode("user-session", value); err == nil {
		cookie := &http.Cookie{
			Name:  "user-session",
			Value: encoded,
			Path:  "/",
		}
		http.SetCookie(w, cookie)
	}
}

func getUserID(request *http.Request) (userID string) {
	if cookie, err := request.Cookie("user-session"); err == nil {
		cookieValue := make(map[string]string)
		if err = cookieHandler.Decode("user-session", cookie.Value, &cookieValue); err == nil {
			userID = cookieValue["userID"]
		}
	}
	return userID
}

func setDispalyInterval(from, to time.Time, w http.ResponseWriter, request *http.Request) {
	cookie, err := request.Cookie("user-session")
	if err != nil {
		fmt.Println(err)
		return
	}
	cookieValue := make(map[string]string)
	err = cookieHandler.Decode("user-session", cookie.Value, &cookieValue)
	if err != nil {
		fmt.Println(err)
		return
	}

	cookieValue["from"] = from.String()
	cookieValue["to"] = to.String()

	placeCookie(cookieValue, w)
}

func getDisplayInterval(request *http.Request) (from, to time.Time, err error) {
	cookie, err := request.Cookie("user-session")
	if err != nil {
		fmt.Println(err)
		return time.Time{}, time.Time{}, err
	}
	cookieValue := make(map[string]string)
	err = cookieHandler.Decode("user-session", cookie.Value, &cookieValue)
	if err != nil {
		fmt.Println(err)
		return time.Time{}, time.Time{}, err
	}

	from, err = time.Parse("2006-01-02 15:04:05.999999999 -0700 MST", cookieValue["from"])
	if err != nil {
		fmt.Println(err)
		return time.Time{}, time.Time{}, err
	}
	to, err = time.Parse("2006-01-02 15:04:05.999999999 -0700 MST", cookieValue["to"])
	if err != nil {
		fmt.Println(err)
		return time.Time{}, time.Time{}, err
	}

	return from, to, nil
}

func placeCookie(value map[string]string, w http.ResponseWriter) {
	if encoded, err := cookieHandler.Encode("user-session", value); err == nil {
		cookie := &http.Cookie{
			Name:  "user-session",
			Value: encoded,
			Path:  "/",
		}
		http.SetCookie(w, cookie)
	}
}

func logout(w http.ResponseWriter, r *http.Request) {
	tpl, err := template.ParseFiles("tpl/login.tpl", "tpl/fragments.tpl")
	if checkErr(err, w) {
		return
	}

	cookie := &http.Cookie{
		Name:   "user-session",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(w, cookie)

	data := anonStruct{
		User: User{
			Username: "",
		},
	}

	tpl.ExecuteTemplate(w, "login.tpl", data)
}

func activeEntry(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)
	if userID == "" {
		http.Error(w, "invalid session", http.StatusForbidden)
		return
	}

	entryID := r.URL.Query().Get("id")

	if entryID == "" {
		err := StoreEntry(userID, time.Time{}, time.Time{}, "Work time", "create", entryID, true)
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

	err = StoreEntry(userID, from, time.Time{}, entry.Type, "update", entryID, true)
	checkErr(err, w)
}

func checkErr(err error, w http.ResponseWriter) bool {
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return true
	}
	return false
}

func pdf(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)
	if userID == "" {
		http.Error(w, "invalid session", http.StatusForbidden)
		return
	}

	trackingData, err := GetUser(userID)
	if checkErr(err, w) {
		return
	}

	from := r.URL.Query().Get("dateFrom")
	dateFrom, err := time.Parse(jsDate, from)
	if checkErr(err, w) {
		return
	}

	to := r.URL.Query().Get("dateTo")
	dateTo, err := time.Parse(jsDate, to)
	if checkErr(err, w) {
		return
	}

	logs, err := GetLogsForUser(trackingData.UserID, dateFrom, dateTo, true)
	data := anonStruct{
		User: trackingData,
		From: from,
		To:   to,
	}
	if err == nil {
		generatePDF(w, logs, data)
	}
}
