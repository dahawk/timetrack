package main

import (
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
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

const (
	updateType    string = "update"
	createConst   string = "create"
	salt          string = "secretsalt"
	workTimeConst string = "Work time"
	holidayConst  string = "holiday"
	sickConst     string = "Sick leave"
)

type anonStruct struct {
	User          User
	WorkTimeUser  User
	UserList      []User
	Error         string
	From          string
	To            string
	from          time.Time
	to            time.Time
	Impersonating bool
	Stats         timeywimey
	WorkTime      workTimeRow
	FullTime      bool
}

var (
	//Info is the INFO level logger
	Info *log.Logger
	//Warning is the WARN level logger
	Warning *log.Logger
	//Error is the ERROR level logger
	Error *log.Logger
)

func initLogging(
	infoHandle io.Writer,
	warningHandle io.Writer,
	errorHandle io.Writer) {

	Info = log.New(infoHandle,
		"INFO: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Warning = log.New(warningHandle,
		"WARNING: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Error = log.New(errorHandle,
		"ERROR: ",
		log.Ldate|log.Ltime|log.Lshortfile)
}

func main() {
	initLogging(ioutil.Discard, os.Stdout, os.Stdout)
	Info.Println("main()")
	var err error

	db, err = sqlx.Open("postgres", fmt.Sprintf("host=%s user=%s dbname=%s password=%s sslmode=disable",
		"<host>", "<username>", "<db name>", "<password>"))
	if err != nil {
		Error.Fatalln(err)
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
	mux.HandleFunc("/impersonate", impersonate)
	mux.HandleFunc("/unimpersonate", unimpersonate)
	mux.HandleFunc("/worktime", worktime)
	mux.HandleFunc("/stats", stats)

	http.ListenAndServe(":1234", mux)
}

func index(w http.ResponseWriter, r *http.Request) {
	Info.Println("index()")
	var userID string
	if userID = getSessionID(w, r); userID == "" {
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

	impersonating := (getImpersonate(r) != "")

	dateFrom, dateTo, err := getDisplayInterval(r)
	if err != nil {
		dateFrom, dateTo = getDefaultDates()
	}

	logs, err := GetLogsForUser(userID, dateFrom, dateTo, false)
	if checkErr(err, w) {
		return
	}

	data := anonStruct{
		User:          trackingData,
		From:          dateFrom.Format(jsDate),
		To:            dateTo.Format(jsDate),
		to:            dateTo,
		from:          dateFrom,
		Impersonating: impersonating,
	}

	stats := calculateStats(logs, data)
	data.Stats = stats

	err = tpl.ExecuteTemplate(w, "index.tpl", data)
	if checkErr(err, w) {
		return
	}
}

func login(w http.ResponseWriter, r *http.Request) {
	Info.Println("login()")
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

			err := tpl.ExecuteTemplate(w, "login.tpl", data)
			checkErr(err, w)
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

func logout(w http.ResponseWriter, r *http.Request) {
	Info.Println("logout()")
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

func pdf(w http.ResponseWriter, r *http.Request) {
	Info.Println("pdf()")
	var userID string
	if userID = getSessionID(w, r); userID == "" {
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
		from: dateFrom,
		to:   dateTo,
	}
	if err == nil {
		generatePDF(w, logs, data)
	}
}

func impersonate(w http.ResponseWriter, r *http.Request) {
	Info.Println("impersonate()")
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
		http.Redirect(w, r, "/", 307)
		return
	}

	impersonateID := r.URL.Query().Get("id")
	if impersonateID == "" {
		http.Redirect(w, r, "/admin", 307)
		return
	}

	setImpersonateSession(userID, impersonateID, w)
	http.Redirect(w, r, "/", 307)
}

func unimpersonate(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)
	if userID == "" {
		http.Error(w, "invalid session", http.StatusForbidden)
		return
	}
	unimpersonateMe(w, r)
	http.Redirect(w, r, "/", 307)
}

func checkErr(err error, w http.ResponseWriter) bool {
	if err != nil {
		Error.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return true
	}
	return false
}
