package main

import (
	"fmt"
	"game-server/components"
	"game-server/gobblers"
	"math/rand/v2"
	"net/http"
	"strings"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/websocket"
)

var MODES = []components.HomeCard{
	{Name: "Gobblers", Desc: "TicTacToe but cooler", Route: "gobblers"},
	{Name: "Connect4", Desc: "Self explaining"},
	{Name: "Pente", Desc: "TicTacToe but 5 in a row and you can bully enemies pieces"},
}

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	fs := http.FileServer(http.Dir("public"))
	r.Handle("/public/*", http.StripPrefix("/public/", fs))
	r.Get("/", home)
	r.Get("/gobblers", gobblersHandler)
	r.HandleFunc("/gobblers/ws-{board_id}", gobblersWsHandler)
	http.ListenAndServe(":3000", r)
}

func render(w http.ResponseWriter, r *http.Request, comp templ.Component) {
	if r.Header.Get("hx-request") != "true" {
		templ.Handler(components.Layout("Game Server", comp)).ServeHTTP(w, r)
	} else {
		templ.Handler(comp).ServeHTTP(w, r)
	}
}

func home(w http.ResponseWriter, r *http.Request) {
	render(w, r, components.Home(MODES))
}

type WsMessage struct {
	Id   string
	Data interface{}
}

func gobblersWsHandler(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error while upgrading ws: ", err)
		return
	}
	defer conn.Close()

	fmt.Println("Client connected")
	g := gobblers.NewGame()
	conn.WriteMessage(websocket.TextMessage, []byte("turn"))

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Error while reading message:", err)
			break
		}
		println(string(msg))
		cmds := strings.Split(string(msg), ":")
		println(cmds)
		switch cmds[0] {
		case "move":
			parts := strings.Split(cmds[1], "-")
			for i := range 2 {
				n := parts[i][1] - '0'
				if parts[i][0] == 'n' {
					g.SelectNew(gobblers.Size(n))
				} else {
					g.SelectBoard(gobblers.Pos(n))
				}
			}
			_, move := g.Selection()
			if !g.DoMove(move) {
				conn.WriteMessage(websocket.TextMessage, []byte("stop"))
			}
			moves := g.Board().Moves()
			m := moves[rand.IntN(len(moves))]
			if !g.DoMove(m) {
				conn.WriteMessage(websocket.TextMessage, []byte("stop"))
			}
			from := ""
			to := "b" + string('0'+m.To)
			if m.New {
				from += "n" + string('0'+m.Size+3)
			} else {
				from += "b" + string('0'+m.From)
			}
			conn.WriteMessage(websocket.TextMessage, []byte("move:"+from+"-"+to))
			conn.WriteMessage(websocket.TextMessage, []byte("turn"))
		}
	}
}

func gobblersHandler(w http.ResponseWriter, r *http.Request) {
	render(w, r, components.Gobblers(gobblers.NewBoard()))
}
