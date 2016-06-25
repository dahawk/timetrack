//contains all functions and structs for session handling
package main

import (
	"fmt"
	"net/http"
	"time"
)

func getSessionID(w http.ResponseWriter, r *http.Request) string {
	Info.Println("getSessionID()")
	var userID string
	if userID = getUserID(r); userID == "" {
		http.Redirect(w, r, "/login", 303)
		return ""
	}
	if imp := getImpersonate(r); imp != "" {
		userID = imp
	}
	return userID
}

func setSession(userID string, w http.ResponseWriter) {
	value := map[string]string{
		"userID": userID,
	}
	placeCookie(value, w)
}

func setImpersonateSession(userID, impersonateID string, w http.ResponseWriter) {
	Info.Printf("setImpersonateSession(%s,%s)\n", userID, impersonateID)
	value := map[string]string{
		"userID":        userID,
		"impersonateID": impersonateID,
	}
	placeCookie(value, w)
}

func setDispalyInterval(from, to time.Time, w http.ResponseWriter, request *http.Request) {
	Info.Println("setDisplayInterval()")
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
	Info.Println("getDisplayInterval()")
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
		return time.Time{}, time.Time{}, err
	}
	to, err = time.Parse("2006-01-02 15:04:05.999999999 -0700 MST", cookieValue["to"])
	if err != nil {
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

func unimpersonateMe(w http.ResponseWriter, r *http.Request) {
	Info.Println("unimpersonate()")
	if cookie, err := r.Cookie("user-session"); err == nil {
		cookieValue := make(map[string]string)
		if err = cookieHandler.Decode("user-session", cookie.Value, &cookieValue); err == nil {
			cookieValue["impersonateID"] = ""
			if encoded, err := cookieHandler.Encode("user-session", cookieValue); err == nil {
				cookie := &http.Cookie{
					Name:  "user-session",
					Value: encoded,
					Path:  "/",
				}
				http.SetCookie(w, cookie)
			}
		}
	}
}

func getUserID(request *http.Request) (userID string) {
	Info.Println("getUserID()")
	if cookie, err := request.Cookie("user-session"); err == nil {
		cookieValue := make(map[string]string)
		if err = cookieHandler.Decode("user-session", cookie.Value, &cookieValue); err == nil {
			userID = cookieValue["userID"]
		}
	}
	return userID
}

func getImpersonate(request *http.Request) (userID string) {
	Info.Println("getImpersonate()")
	if cookie, err := request.Cookie("user-session"); err == nil {
		cookieValue := make(map[string]string)
		if err = cookieHandler.Decode("user-session", cookie.Value, &cookieValue); err == nil {
			userID = cookieValue["impersonateID"]
		}
	}
	return userID
}
