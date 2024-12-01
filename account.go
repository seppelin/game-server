package main

import (
	"game-server/components"
	"game-server/db"
	"game-server/middlewares"
	"game-server/models"
	"log"
	"net/http"
	"time"

	"github.com/a-h/templ"
)

func UserFromRequest(r *http.Request) (models.User, bool) {
	session := r.Context().Value(middlewares.SessionNameKey{})
	return db.UserFromSession(session.(string))
}

func Account(w http.ResponseWriter, r *http.Request) {
	user, ok := UserFromRequest(r)
	renderPage(w, r, components.Account(user.Name, !ok))
}

func AccountSignIn(w http.ResponseWriter, r *http.Request) {
	templ.Handler(components.AccountSignInForm("", "", "")).ServeHTTP(w, r)
}

func AccountSignUp(w http.ResponseWriter, r *http.Request) {
	templ.Handler(components.AccountSignUpForm("", "", "")).ServeHTTP(w, r)
}

func AccountSignOut(w http.ResponseWriter, r *http.Request) {
	sessionName := r.Context().Value(middlewares.SessionNameKey{}).(string)
	err := db.ExpireSession(sessionName)
	if err != nil {
		log.Fatal(err)
	}
	templ.Handler(components.AccountOOB("Anon", true)).ServeHTTP(w, r)
}

func AccountPutSignUp(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	password := r.FormValue("password")
	valid := models.ValidUsername(name) && models.ValidPassword(password)
	if !valid {
		templ.Handler(components.AccountSignInForm(name, password, "Invalid username or password")).ServeHTTP(w, r)
	} else {
		userID, err := db.InsertUser(name, "", models.HashPassword(password))
		if err != nil {
			templ.Handler(components.AccountSignInForm(name, password, err.Error())).ServeHTTP(w, r)
			return
		}
		_ = db.ExpireSession(r.Context().Value(middlewares.SessionNameKey{}).(string))
		_, sessionName, err := db.SetSession(userID, getClientAddress(r), time.Now().Add(time.Hour*24*356))
		if err != nil {
			log.Fatal(err)
		}
		http.SetCookie(w, middlewares.GetCookie(sessionName))
		templ.Handler(components.AccountOOB(name, false)).ServeHTTP(w, r)
	}
}

func AccountPutSignIn(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	password := r.FormValue("password")

	user, err := db.UserByName(name)
	if err == nil {
		err = models.VerifyPassword(user.PwdHash, password)
		if err == nil {
			_ = db.ExpireSession(r.Context().Value(middlewares.SessionNameKey{}).(string))
			_, sessionName, err := db.SetSession(user.ID, getClientAddress(r), time.Now().Add(time.Hour*24*356))
			if err != nil {
				log.Fatal(err)
			}
			http.SetCookie(w, middlewares.GetCookie(sessionName))
			templ.Handler(components.AccountOOB(name, false)).ServeHTTP(w, r)
			return
		}
	}
	templ.Handler(components.AccountSignInForm(name, password, err.Error())).ServeHTTP(w, r)
}
