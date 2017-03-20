package cmd

type request struct {
	New    *requestNew    `json:",omitempty"`
	Play   *requestPlay   `json:",omitempty"`
	Render *requestRender `json:",omitempty"`
}

type requestNew struct {
	Players int `json:"players"`
}

type requestPlay struct {
	Player  int         `json:"player"`
	Command string      `json:"command"`
	Names   []string    `json:"names"`
	Game    interface{} `json:"game"`
}

type requestRender struct {
	Player *int        `json:"player"`
	Game   interface{} `json:"game"`
}

type response struct {
	New         *responseNew         `json:",omitempty"`
	Play        *responsePlay        `json:",omitempty"`
	Render      *responseRender      `json:",omitempty"`
	UserError   *responseUserError   `json:",omitempty"`
	SystemError *responseSystemError `json:",omitempty"`
}

type responsePlay struct {
	Game             gameResponse `json:"game"`
	Logs             []log        `json:"logs"`
	RemainingCommand string       `json:"remaining_command"`
}

type responseNew struct {
	Game gameResponse `json:"game"`
	Logs []log        `json:"logs"`
}

type responseRender struct {
	Render string `json:"render"`
}

type responseUserError struct {
	Message string `json:"message"`
}

type responseSystemError struct {
	Message string `json:"message"`
}

type gameResponse struct {
	State  string             `json:"state"`
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
