package main

/*
this file defines shared interfaces between server and client
*/

type OpCode int

// game opcodes must be POSITIVE integers
const (
	// outgoing (starting in 100)
	OpUpdate OpCode = iota + 100
	OpFeedback

	// incoming (starting in 200)
	OpMove = iota + 198
)

// outgoing

type UpdateBody struct {
	Playing    bool            `json:"playing"`
	Board      []Mark          `json:"board"`
	Marks      map[string]Mark `json:"marks"`
	NextToPlay string          `json:"nextToPlay"`
}

type FeedbackBody string

// incoming

type MoveBody = Mark
