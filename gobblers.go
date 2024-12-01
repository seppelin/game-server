package main

import (
	"encoding/json"
	"fmt"
	"game-server/components"
	"game-server/gobblers"
	"log"
	"math/rand/v2"
	"net/http"

	"github.com/gorilla/websocket"
)

type PlayMessage struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

func (wm *PlayMessage) Set(t string, data interface{}) (err error) {
	wm.Type = t
	wm.Data, err = json.Marshal(data)
	fmt.Println("Msg type:", wm.Type, "data:", string(wm.Data))
	return
}

func GobblersHandler(w http.ResponseWriter, r *http.Request) {
	renderPage(w, r, components.Gobblers(gobblers.NewBoard(), 3, 0))
}

func GobblersPlayHandler(w http.ResponseWriter, r *http.Request) {
	// Init Websocket
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
	var msg PlayMessage

	// Start game
	g := gobblers.NewGame()
	msg.Set("turn", nil)
	conn.WriteJSON(msg)

	for {
		err = conn.ReadJSON(&msg)
		if err != nil {
			fmt.Println("Error while reading message:", err)
			break
		}
		switch msg.Type {
		case "move":
			var move gobblers.Move
			if err = json.Unmarshal(msg.Data, &move); err != nil {
				log.Printf("Error decoding user data: %v", err)
				continue
			}
			if !g.DoMove(move) {
				msg.Set("stop", nil)
				conn.WriteJSON(msg)
				continue
			}
			moves := g.Board().Moves()
			move = moves[rand.IntN(len(moves))]
			if !g.DoMove(move) {
				msg.Set("stop", nil)
				conn.WriteJSON(msg)
				continue
			}
			msg.Set("move", move)
			conn.WriteJSON(msg)
			msg.Set("turn", nil)
			conn.WriteJSON(msg)
		}
	}
}
