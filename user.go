package main

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

//User is the struct to hold user data fetched from the db
type User struct {
	Name        string
	Username    string
	UserID      string `db:"id"`
	Password    string
	ActiveEntry string
	Admin       bool
	Enabled     bool
}

//GetUser fetches a user from the db and returns it
func GetUser(userID string) (User, error) {
	var u User
	err := db.Get(&u, "select * from \"user\" where id=$1", userID)
	if err != nil {
		return User{}, err
	}
	entry, err := ActiveEntry(u.UserID)
	if err != nil {
		return User{}, err
	}
	u.ActiveEntry = entry

	return u, nil
}

//VerifyLogin checks username and password against the stored values in the db.
//returns userID for match and "" + possible error for mismatch or problem along the way
func VerifyLogin(username, password string) (string, string) {
	var u User
	err := db.Get(&u, "select * from \"user\" where username=$1", username)
	if err != nil {
		return "", "Username or Password doesn't match"
	}
	if !u.Enabled {
		return "", "User Account is disabled"
	}
	pwd := fmt.Sprintf("%s%s%s", salt, password, salt)

	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(pwd))
	if err != nil {
		return "", "Username or Password doesn't match"
	}

	return u.UserID, ""
}

//GetUserList fetches a list of all users in the db.
func GetUserList() ([]User, error) {
	var users []User
	err := db.Select(&users, "select * from \"user\"")
	if err != nil {
		return []User{}, err
	}
	return users, nil
}

//StoreUser takes passed user data and updates the corresponding db entry
func StoreUser(id, username, name, password, repeat, t string) error {
	//fmt.Println(id, username, name, password, repeat, t)
	if t == "edit" {
		return updateUser(id, username, name, password, repeat)
	} else if t == "create" {
		return createUser(username, name, password, repeat)
	}
	return nil
}

func createUser(username, name, password, repeat string) error {
	if username != "" &&
		name != "" &&
		password != "" &&
		repeat != "" &&
		password == repeat {
		saltedPassword := fmt.Sprintf("%s%s%s", salt, password, salt)
		passwordCrypt, err := bcrypt.GenerateFromPassword([]byte(saltedPassword), 12)
		if err != nil {
			return err
		}
		//fmt.Printf("insert into \"user\" (username,name,password) values (%s,%s,%s)\n", username, name, passwordCrypt)
		_, err = db.Exec("insert into \"user\" (username,name,password) values ($1,$2,$3)",
			username, name, string(passwordCrypt))
		if err != nil {
			return err
		}
	}
	return nil
}

func updateUser(id, username, name, password, repeat string) error {
	if id != "" &&
		username != "" &&
		name != "" {
		//fmt.Printf("update \"user\" set username=%s, name=%s where id=%s\n", username, name, id)
		_, err := db.Exec("update \"user\" set username=$1, name=$2 where id=$3",
			username, name, id)
		if err != nil {
			return err
		}
	}
	if id != "" &&
		password != "" &&
		repeat != "" &&
		password == repeat {
		saltedPassword := fmt.Sprintf("%s%s%s", salt, password, salt)
		passwordCrypt, err := bcrypt.GenerateFromPassword([]byte(saltedPassword), 12)
		if err != nil {
			return err
		}
		//fmt.Printf("update \"user\" set password=%s where id=%s\n", passwordCrypt, id)
		_, err = db.Exec("update \"user\" set password=$1 where id=$2", passwordCrypt, id)
		if err != nil {
			return err
		}
	}
	return nil
}

//UpdateEnabled updates a users enabled flag according to the input
func UpdateEnabled(userID string, enabled bool) error {
	_, err := db.Exec("update \"user\" set enabled=$1 where id=$2", enabled, userID)

	return err
}
