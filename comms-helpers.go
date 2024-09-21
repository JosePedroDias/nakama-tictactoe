package main

import (
	"encoding/json"

	"github.com/heroiclabs/nakama-common/runtime"
)

func getJustSender(state *MatchState, userId string) []runtime.Presence {
	destinations := make([]runtime.Presence, 0)
	destinations = append(destinations, *state.presences[userId])
	return destinations
}

//lint:ignore U1000 optional method
func getAllButSender(state *MatchState, userId string) []runtime.Presence {
	destinations := make([]runtime.Presence, 0)

	for k, v := range state.presences {
		if k != userId {
			destinations = append(destinations, *v)
		}
	}

	return destinations
}

////

func bcUpdate(dispatcher runtime.MatchDispatcher, state *MatchState, destinations []runtime.Presence) {
	body := &UpdateBody{
		Playing:    state.playing,
		Marks:      state.marks,
		Board:      state.board,
		NextToPlay: state.nextToPlay[0],
	}

	data, err := json.Marshal(body)
	if err == nil {
		dispatcher.BroadcastMessage(int64(OpUpdate), data, destinations, nil, true)
	}
}

func bcFeedback(dispatcher runtime.MatchDispatcher, body string, destinations []runtime.Presence) {
	data, err := json.Marshal(body)
	if err == nil {
		dispatcher.BroadcastMessage(int64(OpFeedback), data, destinations, nil, true)
	}
}
