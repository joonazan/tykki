package tykki

type Action struct {
	BotId  int    `json:"botId"`
	Action string `json:"type"`
	Pos    Pos    `json:"pos"`
}

func (b Bot) Move(pos Pos) Action {
	return Action{b.BotId, "move", pos.Plus(b.Pos)}
}

func (b Bot) Shoot(pos Pos) Action {
	return Action{b.BotId, "cannon", pos}
}

func (b Bot) Scan(pos Pos) Action {
	return Action{b.BotId, "radar", pos}
}
