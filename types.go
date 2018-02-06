package cmd

import "github.com/brdgme-go/brdgme"

type request struct {
	PlayerCounts *requestPlayerCounts `json:",omitempty"`
	New          *requestNew          `json:",omitempty"`
	Status       *requestStatus       `json:",omitempty"`
	Play         *requestPlay         `json:",omitempty"`
	PubRender    *requestPubRender    `json:",omitempty"`
	PlayerRender *requestPlayerRender `json:",omitempty"`
}

type requestPlayerCounts struct {
}

type requestNew struct {
	Players int `json:"players"`
}

type requestStatus struct {
	Game string `json:"game"`
}

type requestPlay struct {
	Player  int      `json:"player"`
	Command string   `json:"command"`
	Names   []string `json:"names"`
	Game    string   `json:"game"`
}

type requestPubRender struct {
	Game string `json:"game"`
}

type requestPlayerRender struct {
	Player int    `json:"player"`
	Game   string `json:"game"`
}

type response struct {
	PlayerCounts *responsePlayerCounts `json:",omitempty"`
	New          *responseNew          `json:",omitempty"`
	Status       *responseStatus       `json:",omitempty"`
	Play         *responsePlay         `json:",omitempty"`
	PubRender    *responseRender       `json:",omitempty"`
	PlayerRender *responseRender       `json:",omitempty"`
	UserError    *responseUserError    `json:",omitempty"`
	SystemError  *responseSystemError  `json:",omitempty"`
}

type responsePlayerCounts struct {
	PlayerCounts int `json:"player_counts"`
}

type responseNew struct {
	Game          gameResponse `json:"game"`
	Logs          []log        `json:"logs"`
	PublicRender  string       `json:"public_render"`
	PlayerRenders []string     `json:"player_renders"`
}

type responseStatus struct {
	Game          gameResponse `json:"game"`
	PublicRender  string       `json:"public_render"`
	PlayerRenders []string     `json:"player_renders"`
}

type responsePlay struct {
	Game             gameResponse `json:"game"`
	Logs             []log        `json:"logs"`
	CanUndo          bool         `json:"can_undo"`
	RemainingCommand string       `json:"remaining_command"`
	PublicRender     string       `json:"public_render"`
	PlayerRenders    []string     `json:"player_renders"`
}

type responsePubRender struct {
	Render pubRender `json:"render"`
}

type pubRender struct {
	PubState string `json:"pub_state"`
	Render   string `json:"render"`
}

type responsePlayerRender struct {
	Render playerRender `json:"render"`
}

type playerRender struct {
	PlayerState string       `json:"player_state"`
	Render      string       `json:"render"`
	CommandSpec *brdgme.Spec `json:"command_spec,omitempty"`
}

type responseUserError struct {
	Message string `json:"message"`
}

type responseSystemError struct {
	Message string `json:"message"`
}

type gameResponse struct {
	State  string             `json:"state"`
	Points []float32          `json:"points"`
	Status gameResponseStatus `json:"status"`
}

type gameResponseStatus struct {
	Active   *gameResponseStatusActive   `json:",omitempty"`
	Finished *gameResponseStatusFinished `json:",omitempty"`
}

type gameResponseStatusActive struct {
	WhoseTurn  []int `json:"whose_turn"`
	Eliminated []int `json:"eliminated"`
}

type gameResponseStatusFinished struct {
	Winners []int `json:"winners"`
}

type log struct {
	Content string `json:"content"`
	At      string `json:"at"`
	Public  bool   `json:"public"`
	To      []int  `json:"to"`
}
