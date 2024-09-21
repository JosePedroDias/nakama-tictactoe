package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/heroiclabs/nakama-common/runtime"
)

const TICK_RATE = 2 // number of ticks the server runs per second
const MIN_PLAYERS = 2
const MAX_PLAYERS = 2

func newMatch(
	ctx context.Context,
	logger runtime.Logger,
	db *sql.DB,
	nk runtime.NakamaModule) (m runtime.Match, err error) {
	return &Match{}, nil
}

func (m *Match) MatchInit(
	ctx context.Context,
	logger runtime.Logger,
	db *sql.DB,
	nk runtime.NakamaModule,
	params map[string]interface{}) (interface{}, int, string) {
	state := &MatchState{
		playing:         false,
		label:           &MatchLabel{Open: 1, TicTacToe: 1},
		joinsInProgress: 0,

		presences: make(map[string]*runtime.Presence),

		board:      make([]Mark, 9),
		marks:      make(map[string]Mark),
		nextToPlay: make([]string, 0),
		winner:     "",
	}

	label := ""
	labelBytes, err := json.Marshal(state.label)
	if err == nil {
		label = string(labelBytes)
	}

	return state, TICK_RATE, label
}

// https://heroiclabs.com/docs/nakama/server-framework/go-runtime/function-reference/match-handler/#MatchJoinAttempt
func (m *Match) MatchJoinAttempt(
	ctx context.Context,
	logger runtime.Logger,
	db *sql.DB,
	nk runtime.NakamaModule,
	dispatcher runtime.MatchDispatcher,
	tick int64,
	state_ interface{},
	presence runtime.Presence,
	metadata map[string]string) (interface{}, bool, string) {
	state := state_.(*MatchState)

	// Check if match is full
	totalCount := getConnectedUsersCount(state) + state.joinsInProgress
	if totalCount >= MAX_PLAYERS {
		return state, false, "match full"
	}

	state.joinsInProgress++
	return state, true, ""
}

func (m *Match) MatchJoin(
	ctx context.Context,
	logger runtime.Logger,
	db *sql.DB,
	nk runtime.NakamaModule,
	dispatcher runtime.MatchDispatcher,
	tick int64,
	state_ interface{},
	presences []runtime.Presence) interface{} {
	state := state_.(*MatchState)

	for _, p := range presences {
		state.joinsInProgress--
		id := p.GetUserId()
		state.presences[id] = &p
		state.nextToPlay = append(state.nextToPlay, id)
	}

	// match got full
	if getConnectedUsersCount(state) >= MAX_PLAYERS && state.label.Open == 1 {
		state.label.Open = 0
		labelStr, err := json.Marshal(state.label)
		if err == nil {
			dispatcher.MatchLabelUpdate(string(labelStr))
		}

		state.marks[state.nextToPlay[0]] = MarkX
		state.marks[state.nextToPlay[1]] = MarkO

		state.playing = true
		bcFeedback(dispatcher, "game started!", nil)

		bcUpdate(dispatcher, state, nil)
	}

	return state
}

func (m *Match) MatchLeave(
	ctx context.Context,
	logger runtime.Logger,
	db *sql.DB,
	nk runtime.NakamaModule,
	dispatcher runtime.MatchDispatcher,
	tick int64,
	state_ interface{},
	presences []runtime.Presence) interface{} {
	state := state_.(*MatchState)

	state.playing = false

	for _, p := range presences {
		id := p.GetUserId()
		delete(state.presences, id)
	}

	if len(state.presences) == 0 {
		return nil
	} else {
		bcFeedback(dispatcher, "game stopped!", nil)
	}

	bcUpdate(dispatcher, state, nil)

	return state
}

func (m *Match) MatchLoop(
	ctx context.Context,
	logger runtime.Logger,
	db *sql.DB,
	nk runtime.NakamaModule,
	dispatcher runtime.MatchDispatcher,
	tick int64,
	state_ interface{},
	messages []runtime.MatchData) interface{} {
	state := state_.(*MatchState)

	//logger.Debug("Running match loop. Tick: %d", tick)

	if !state.playing {
		return state
	}

	for _, message := range messages {
		senderUserId := message.GetUserId()
		op := message.GetOpCode()
		data := message.GetData()

		logger.Debug("SENDER USER ID: %s | OPCODE: %d | DATA: %s", senderUserId, op, data)

		switch op {

		case OpMove:
			var moveBody MoveBody
			if err := json.Unmarshal(data, &moveBody); err != nil {
				logger.Error("error unmarshalling move body: %v", err)
				continue
			}

			nextMark := nextMarkToPlay(state)
			var feedbackContents string
			if state.marks[senderUserId] != nextMark {
				feedbackContents = "it is not your time to play!"
			} else if state.board[moveBody] != MarkEmpty {
				feedbackContents = "you must pick an empty cell!"
			} else {
				state.board[moveBody] = nextMark

				if hasWon(state, nextMark) {
					state.winner = senderUserId
					state.playing = false
					presPtr := state.presences[senderUserId]
					winnerUsername := (*presPtr).GetUsername()
					feedbackContents := fmt.Sprintf("%s won!", winnerUsername)
					bcFeedback(dispatcher, feedbackContents, nil)
				} else if isBoardFull(state) {
					state.playing = false
					bcFeedback(dispatcher, "it's a tie!", nil)
				} else {
					rotateNextToPlay(state)
				}

				bcUpdate(dispatcher, state, nil)
			}

			if len(feedbackContents) > 0 {
				logger.Error(feedbackContents)
				bcFeedback(dispatcher, feedbackContents, getJustSender(state, message.GetUserId()))
			}

			if !state.playing {
				return nil
			}

		default:
			feedbackContents := fmt.Sprintf("unsupported opcode received: (%d)", op)
			logger.Error(feedbackContents)
			bcFeedback(dispatcher, feedbackContents, getJustSender(state, message.GetUserId()))
		}
	}

	return state
}

// https://heroiclabs.com/docs/nakama/server-framework/go-runtime/function-reference/match-handler/#MatchTerminate
func (m *Match) MatchTerminate(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, dispatcher runtime.MatchDispatcher, tick int64, state_ interface{}, graceSeconds int) interface{} {
	state := state_.(*MatchState)

	// ADD LOGIC IF/WHEN NEEDED
	//message := "Server shutting down in " + strconv.Itoa(graceSeconds) + " seconds."
	//dispatcher.BroadcastMessage(2, []byte(message), []runtime.Presence{}, nil, true)

	return state
}

func (m *Match) MatchSignal(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, dispatcher runtime.MatchDispatcher, tick int64, state_ interface{}, data string) (interface{}, string) {
	state := state_.(*MatchState)

	if data == "kill" {
		return nil, "killing match due to rpc signal"
	}

	return state, ""
}

////

func getUserIds(state *MatchState) []string {
	result := make([]string, 0)
	for k := range state.presences {
		result = append(result, k)
	}
	return result
}

func getConnectedUsersCount(state *MatchState) int {
	return len(getUserIds(state))
}

////

func rotateNextToPlay(state *MatchState) {
	state.nextToPlay = append(state.nextToPlay[1:], state.nextToPlay[0])
}
