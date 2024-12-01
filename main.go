package main

import (
	"game-server/components"
	"game-server/db"
	"game-server/middlewares"
	"log"
	"net/http"
	"strings"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var MODES = []components.HomeCard{
	{Name: "Gobblers", Desc: "TicTacToe but cooler", Route: "gobblers"},
	{Name: "Connect4", Desc: "Self explaining"},
	{Name: "Pente", Desc: "TicTacToe but 5 in a row and you can bully enemies pieces"},
}

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middlewares.SessionName)
	r.Handle("/public/*", http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))
	r.Get("/", home)
	r.Get("/gobblers", GobblersHandler)
	r.Get("/account", Account)
	r.Get("/account/sign-out", AccountSignOut)
	r.Get("/account/sign-up", AccountSignUp)
	r.Get("/account/sign-in", AccountSignIn)
	r.Put("/account/sign-up", AccountPutSignUp)
	r.Put("/account/sign-in", AccountPutSignIn)

	r.HandleFunc("/gobblers/play-{game_id}-{game_state}", GobblersPlayHandler)
	http.ListenAndServe(":3000", r)
}

func home(w http.ResponseWriter, r *http.Request) {
	renderPage(w, r, components.Home(MODES))
}

func renderPage(w http.ResponseWriter, r *http.Request, comp templ.Component) {
	if r.Header.Get("hx-request") != "true" {
		// On full page load set session accessed_at
		session, err := db.SessionByName(r.Context().Value(middlewares.SessionNameKey{}).(string))

		anon := true
		name := "Anon"
		if err == nil {
			log.Println(session)
			db.SetAccessedAt(session.Id)
			user, err := db.User(session.UserID)
			if err == nil {
				anon = false
				name = user.Name
			} else {
				log.Println("Get user page reload:", err)
			}
		} else {
			log.Println("Get session page reload:", err, r.Context().Value(middlewares.SessionNameKey{}).(string))
		}
		templ.Handler(components.Layout("Game Server", name, anon, comp)).ServeHTTP(w, r)
	} else {
		templ.Handler(comp).ServeHTTP(w, r)
	}
}

func getClientAddress(r *http.Request) string {
	// Check for the X-Forwarded-For header
	xForwardedFor := r.Header.Get("X-Forwarded-For")
	if xForwardedFor != "" {
		// The X-Forwarded-For header can contain multiple IPs, take the first one
		ips := strings.Split(xForwardedFor, ",")
		return strings.TrimSpace(ips[0]) // Return the first IP address
	}

	// Fallback to RemoteAddr
	return r.RemoteAddr // This contains the IP and port
}
