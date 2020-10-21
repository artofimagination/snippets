package main

import (
	"fmt"
	"net/http"

	"mongodb-go-interface/mongodb"
)

func sayHello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hi! I am Server!")
}

func insertUser(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Inserting user")
	names, ok := r.URL.Query()["name"]
	if !ok || len(names[0]) < 1 {
		fmt.Fprintln(w, "Url Param 'name' is missing")
		return
	}

	name := names[0]
	emails, ok := r.URL.Query()["email"]
	if !ok || len(emails[0]) < 1 {
		fmt.Fprintln(w, "Url Param 'email' is missing")
		return
	}

	email := emails[0]

	passwords, ok := r.URL.Query()["password"]
	if !ok || len(emails[0]) < 1 {
		fmt.Fprintln(w, "Url Param 'password' is missing")
		return
	}

	password := passwords[0]

	result, err := mongodb.AddUser(name, email, password)
	if err != nil {
		fmt.Fprintln(w, err.Error())
	} else {
		fmt.Fprintln(w, result)
	}
}

func getUser(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Getting user")
	emails, ok := r.URL.Query()["email"]
	if !ok || len(emails[0]) < 1 {
		fmt.Fprintln(w, "Url Param 'email' is missing")
		return
	}

	email := emails[0]

	result, err := mongodb.GetUserByEmail(email)
	if err != nil {
		fmt.Fprintln(w, err.Error())
	} else {
		fmt.Fprintln(w, result)
	}
}

func getSettings(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Getting user settings")
	emails, ok := r.URL.Query()["email"]
	if !ok || len(emails[0]) < 1 {
		fmt.Fprintln(w, "Url Param 'email' is missing")
		return
	}

	email := emails[0]

	resultUser, err := mongodb.GetUserByEmail(email)
	if err != nil {
		fmt.Fprintln(w, err.Error())
	} else {
		fmt.Fprintln(w, resultUser)
	}

	resultSettings, err := mongodb.GetSettings(&resultUser.SettingsID)
	if err != nil {
		fmt.Fprintln(w, err.Error())
	} else {
		fmt.Fprintln(w, resultSettings)
	}
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Deleting user")
	emails, ok := r.URL.Query()["email"]
	if !ok || len(emails[0]) < 1 {
		fmt.Fprintln(w, "Url Param 'email' is missing")
		return
	}

	email := emails[0]

	err := mongodb.DeleteUser(email)
	if err != nil {
		fmt.Fprintln(w, err.Error())
	}
}

func deleteSettings(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Deleting user settings")
	emails, ok := r.URL.Query()["email"]
	if !ok || len(emails[0]) < 1 {
		fmt.Fprintln(w, "Url Param 'email' is missing")
		return
	}

	email := emails[0]

	resultUser, err := mongodb.GetUserByEmail(email)
	if err != nil {
		fmt.Fprintln(w, err.Error())
	} else {
		fmt.Fprintln(w, resultUser)
	}

	resultSettings, err := mongodb.DeleteSettings(&resultUser.SettingsID)
	if err != nil {
		fmt.Fprintln(w, err.Error())
	} else {
		fmt.Fprintln(w, resultSettings)
	}
}

func checkUserPass(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Check user pass")
	emails, ok := r.URL.Query()["email"]
	if !ok || len(emails[0]) < 1 {
		fmt.Fprintln(w, "Url Param 'email' is missing")
		return
	}

	email := emails[0]
	passwords, ok := r.URL.Query()["password"]
	if !ok || len(emails[0]) < 1 {
		fmt.Fprintln(w, "Url Param 'password' is missing")
		return
	}

	password := passwords[0]

	err := mongodb.CheckEmailAndPassword(email, password)
	if err != nil {
		fmt.Fprintln(w, err.Error())
	} else {
		fmt.Fprintln(w, "User pass matched!")
	}
}

func main() {
	http.HandleFunc("/", sayHello)
	http.HandleFunc("/insert", insertUser)
	http.HandleFunc("/get", getUser)
	http.HandleFunc("/delete", deleteUser)
	http.HandleFunc("/check", checkUserPass)
	http.HandleFunc("/get-settings", getSettings)
	http.HandleFunc("/delete-settings", deleteSettings)

	// Start HTTP server that accepts requests from the offer process to exchange SDP and Candidates
	panic(http.ListenAndServe(":8080", nil))
}
