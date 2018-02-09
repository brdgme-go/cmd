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
	case request.PlayerCounts != nil:
		handlePlayerCounts(game, *request.PlayerCounts, encoder)
	case request.New != nil:
		handleNew(game, *request.New, encoder)
	case request.Status != nil:
		handleStatus(game, *request.Status, encoder)
	case request.Play != nil:
		handlePlay(game, *request.Play, encoder)
	case request.PubRender != nil:
		handlePubRender(game, *request.PubRender, encoder)
	case request.PlayerRender != nil:
		handlePlayerRender(game, *request.PlayerRender, encoder)
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

func toGameResponse(game brdgme.Gamer) (gameResponse, error) {
	gameJSON, err := json.Marshal(game)
	return gameResponse{
		State:  string(gameJSON),
		Status: game.Status(),
	}, err
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
			At:      time.Now().Format("2006-01-02T15:04:05.999999999"),
			Public:  v.Public,
			To:      to,
		}
	}
	return l
}

func handlePlayerCounts(game brdgme.Gamer, request requestPlayerCounts, out *json.Encoder) {
	out.Encode(response{
		PlayerCounts: &responsePlayerCounts{
			PlayerCounts: game.PlayerCounts(),
		},
	})
}

func handleNew(game brdgme.Gamer, request requestNew, out *json.Encoder) {
	logs, err := game.New(request.Players)
	if err == nil {
		gameResponse, err := toGameResponse(game)
		if err == nil {
			out.Encode(response{
				New: &responseNew{
					Game: gameResponse,
					Logs: toResponseLogs(logs),
				},
			})
		} else {
			out.Encode(response{
				SystemError: &responseSystemError{
					Message: fmt.Sprintf("Unable to create game response struct, %s", err),
				},
			})
		}
	} else {
		// Most likely due to incorrect player counts.
		out.Encode(response{
			UserError: &responseUserError{
				Message: fmt.Sprintf("Unable to start game, %s", err),
			},
		})
	}
}

func handleStatus(game brdgme.Gamer, request requestStatus, out *json.Encoder) {
	if err := unmarshalGame(request.Game, game); err != nil {
		out.Encode(response{
			SystemError: &responseSystemError{
				Message: fmt.Sprintf("Could not unmarshal game: %s", err),
			},
		})
		return
	}
	gameResp, err := toGameResponse(game)
	if err != nil {
		out.Encode(response{
			SystemError: &responseSystemError{
				Message: fmt.Sprintf("Could not get game response: %s", err),
			},
		})
		return
	}
	playerRs := []string{}
	for i, pc := 0, game.PlayerCount(); i < pc; i++ {
		playerRs = append(playerRs, game.PlayerRender(i))
	}
	out.Encode(response{
		Status: &responseStatus{
			Game:          gameResp,
			PublicRender:  game.PubRender(),
			PlayerRenders: playerRs,
		},
	})
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
		commandResponse, err := game.Command(request.Player, remainingCommand, request.Names)
		newRemainingCommand := strings.TrimSpace(commandResponse.Remaining)
		commandSucceeded = commandSucceeded || err == nil
		logs = append(logs, commandResponse.Logs...)
		if err != nil || newRemainingCommand == "" || remainingCommand == newRemainingCommand {
			if commandSucceeded {
				// Something has already worked, so we'll stay quiet
				gameResponse, err := toGameResponse(game)
				if err == nil {
					out.Encode(response{
						Play: &responsePlay{
							Game:             gameResponse,
							Logs:             toResponseLogs(logs),
							CanUndo:          commandResponse.CanUndo,
							RemainingCommand: newRemainingCommand,
						},
					})
				} else {
					out.Encode(response{
						SystemError: &responseSystemError{
							Message: fmt.Sprintf("Unable to create game response struct, %s", err),
						},
					})
				}
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

func toPubRender(game brdgme.Gamer) (pubRender, error) {
	pubState, err := json.Marshal(game.PubState())
	if err != nil {
		return pubRender{}, fmt.Errorf("could not marshal pub state, %v", err)
	}
	return pubRender{
		PubState: string(pubState),
		Render:   game.PubRender(),
	}, nil
}

func handlePubRender(game brdgme.Gamer, request requestPubRender, out *json.Encoder) {
	if err := unmarshalGame(request.Game, game); err != nil {
		out.Encode(response{
			SystemError: &responseSystemError{
				Message: fmt.Sprintf("could not unmarshal game: %s", err),
			},
		})
		return
	}
	pr, err := toPubRender(game)
	if err != nil {
		out.Encode(response{
			SystemError: &responseSystemError{
				Message: fmt.Sprintf("could not generate pub render: %s", err),
			},
		})
		return
	}
	out.Encode(response{
		PubRender: &responsePubRender{
			Render: pr,
		},
	})
}

func toPlayerRender(game brdgme.Gamer, player int) (playerRender, error) {
	playerState, err := json.Marshal(game.PlayerState(player))
	if err != nil {
		return playerRender{}, fmt.Errorf("could not marshal player state, %v", err)
	}
	return playerRender{
		PlayerState: string(playerState),
		Render:      game.PlayerRender(player),
	}, nil
}

func handlePlayerRender(game brdgme.Gamer, request requestPlayerRender, out *json.Encoder) {
	if err := unmarshalGame(request.Game, game); err != nil {
		out.Encode(response{
			SystemError: &responseSystemError{
				Message: fmt.Sprintf("Could not unmarshal game: %s", err),
			},
		})
		return
	}
	pr, err := toPlayerRender(game, request.Player)
	if err != nil {
		out.Encode(response{
			SystemError: &responseSystemError{
				Message: fmt.Sprintf("could not generate player render: %s", err),
			},
		})
		return
	}
	out.Encode(response{
		PlayerRender: &responsePlayerRender{
			Render: pr,
		},
	})
}
