package main

import "github.com/heroiclabs/nakama-common/runtime"

type Match struct{}

type MatchLabel struct {
	Open      int `json:"open"`
	TicTacToe int `json:"tictactoe"`
}

type MatchState struct {
	// match lifecycle related
	playing         bool
	label           *MatchLabel
	joinsInProgress int

	// user maps
	presences map[string]*runtime.Presence

	board      []Mark
	marks      map[string]Mark
	nextToPlay []string
	winner     string
}

type MatchRpcBody struct {
	MatchIds []string `json:"matchIds"`
}

////

type Mark int

// game opcodes must be POSITIVE integers
const (
	MarkEmpty Mark = iota
	MarkX
	MarkO
)
