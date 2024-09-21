# simple tictactoe nakama authoritative server logic

You can find the client counterpart to this here: https://github.com/JosePedroDias/nakama-tictactoe-client


## module contents

- `tictactoe_match` [RPC](rpcs.go) (responsible for querying for open matches of this kind, creating one if none were found)
- `tictactoe_kill` [RPC](rpcs.go) (auxiliary RPC which kills of matches of this kind) ~> disabled
- `tictactoe` (authoritative match implementation)[match.go]. expects 2 players to join


## abstract game logic

see [logic.go](logic.go)


## opcodes

see [messages.go](messages.go)


## reference

### nakama

- https://heroiclabs.com/docs/nakama/getting-started/
- https://heroiclabs.com/docs/nakama/server-framework/go-runtime/function-reference/

### go plugin

- https://github.com/heroiclabs/nakama/tree/master/sample_go_module
- https://github.com/heroiclabs/nakama/blob/v3.22.0/go.mod
