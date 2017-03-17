package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/brdgme-go/brdgme"
)

// Cli creates a CLI interface to a game.
func Cli(game brdgme.Gamer, in io.Reader, out io.Writer) {
	var request request
	decoder := json.NewDecoder(in)
	encoder := json.NewEncoder(out)
	if err := decoder.Decode(&request); err != nil {
		encoder.Encode(response{
			SystemError: &responseSystemError{
				Message: fmt.Sprintf("Unable to decode request: %v", err),
			}})
		return
	}
	switch {
	case request.New != nil:
		handleNew(game, *request.New, encoder)
	case request.Play != nil:
		handlePlay(game, *request.Play, encoder)
	case request.Render != nil:
		handleRender(game, *request.Render, encoder)
	default:
		encoder.Encode(response{
			SystemError: &responseSystemError{
				Message: "Could not parse command from request",
			},
		})
	}
}

func unmarshalGame(game interface{}, into brdgme.Gamer) error {
	gameJSON, err := json.Marshal(game)
	if err != nil {
		return err
	}
	return json.Unmarshal(gameJSON, into)
}

func toGameResponse(game brdgme.Gamer) gameResponse {
	whoseTurn := game.WhoseTurn()
	if whoseTurn == nil {
		whoseTurn = []int{}
	}
	winners := game.Winners()
	if winners == nil {
		winners = []int{}
	}
	eliminated := []int{}
	if eGame, ok := game.(brdgme.Eliminator); ok {
		if gEliminated := eGame.Eliminated(); gEliminated != nil {
			eliminated = gEliminated
		}
	}
	return gameResponse{
		Game:       game,
		IsFinished: game.IsFinished(),
		WhoseTurn:  whoseTurn,
		Winners:    winners,
		Eliminated: eliminated,
	}
}

func toResponseLogs(logs []brdgme.Log) []log {
	l := make([]log, len(logs))
	for k, v := range logs {
		to := []int{}
		if v.To != nil {
			to = v.To
		}
		l[k] = log{
			Content: v.Message,
			At:      time.Now().Format(time.RFC3339),
			Public:  v.Public,
			To:      to,
		}
	}
	return l
}

func handleNew(game brdgme.Gamer, request requestNew, out *json.Encoder) {
	logs, err := game.Start(request.Players)
	if err == nil {
		out.Encode(response{
			New: &responseNew{
				Game: toGameResponse(game),
				Logs: toResponseLogs(logs),
			},
		})
	} else {
		// Most likely due to incorrect player counts.
		out.Encode(response{
			UserError: &responseUserError{
				Message: fmt.Sprintf("Unable to start game, %s", err),
			},
		})
	}
}

func handlePlay(game brdgme.Gamer, request requestPlay, out *json.Encoder) {
	if err := unmarshalGame(request.Game, game); err != nil {
		out.Encode(response{
			SystemError: &responseSystemError{
				Message: fmt.Sprintf("Could not unmarshal game: %s", err),
			},
		})
	}
	logs := []brdgme.Log{}
	remainingCommand := strings.TrimSpace(request.Command)
	commandSucceeded := false
	for {
		newLogs, newRemainingCommand, err := game.Command(request.Player, remainingCommand, request.Names)
		newRemainingCommand = strings.TrimSpace(newRemainingCommand)
		commandSucceeded = commandSucceeded || err == nil
		logs = append(logs, newLogs...)
		if err != nil || newRemainingCommand == "" || remainingCommand == newRemainingCommand {
			if commandSucceeded {
				// Something has already worked, so we'll stay quiet
				out.Encode(response{
					Play: &responsePlay{
						Game:             toGameResponse(game),
						Logs:             toResponseLogs(logs),
						RemainingCommand: newRemainingCommand,
					},
				})
			} else if err != nil {
				// We got an error so lets return it
				out.Encode(response{
					UserError: &responseUserError{
						Message: fmt.Sprintf("Command failed, %s", err),
					},
				})
			} else {
				// No commands were parsed for some reason
				out.Encode(response{
					UserError: &responseUserError{
						Message: "No command was executed",
					},
				})
			}
			return
		}
		remainingCommand = newRemainingCommand
	}
}

func handleRender(game brdgme.Gamer, request requestRender, out *json.Encoder) {
	if err := unmarshalGame(request.Game, game); err != nil {
		out.Encode(response{
			SystemError: &responseSystemError{
				Message: fmt.Sprintf("Could not unmarshal game: %s", err),
			},
		})
	}
	out.Encode(response{
		Render: &responseRender{
			Render: game.Render(request.Player),
		},
	})
}
