package tykki

import (
	js "encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
)

type Bot struct {
	Pos       Pos
	BotId, Hp int
}

type Pos struct {
	X int `json:"x"`
	Y int `json:"y"`
}

func (p Pos) Plus(p2 Pos) Pos {
	return Pos{p.X + p2.X, p.Y + p2.Y}
}

type Response struct {
	Kind   string `json:"type"`
	Config Config
	You    Team

	// only available when kind is events
	Round  int `json:"roundId"`
	Events []Event

	// available at connect
	TeamId int

	// available at end
	WinnerTeamId int
}

type Config struct {
	Bots, Move, FieldRadius, StartHp int
	MaxCount, Asteroids, LoopTime    int
	Cannon, Radar, See               int
	NoWait                           bool
}

type Team struct {
	Bots []Bot
}

type Event struct {
	Event         string
	BotId, Source int
	Pos           Pos
	Damage        int
}

func Run(teamName string, onJoin func(Config), play func([]Bot, []Event) []Action) {
	dialer := websocket.Dialer{}
	conn, _, err := dialer.Dial("ws://localhost:3000", nil)
	if err != nil {
		panic(err)
	}

	var round, team int

	json := Response{}

game:
	for {
		if err := conn.ReadJSON(&json); err != nil {
			panic(err)
		}

		switch json.Kind {
		case "connected":
			team = json.TeamId
			onJoin(json.Config)
			conn.WriteJSON(object{
				"type":     "join",
				"teamName": teamName,
			})
		case "start":
			// no-op
		case "events":
			round = json.Round
			conn.WriteJSON(object{
				"type":    "actions",
				"roundId": round,
				"actions": play(json.You.Bots, json.Events),
			})
			s, _ := js.Marshal(object{
				"type":    "actions",
				"roundId": round,
				"actions": play(json.You.Bots, json.Events),
			})
			fmt.Println(string(s))
		case "end":
			if json.WinnerTeamId == team {
				fmt.Println("I won!")
			}
			break game
		}
	}
}

type object map[string]interface{}
