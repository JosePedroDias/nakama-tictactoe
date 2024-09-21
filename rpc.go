package main

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/heroiclabs/nakama-common/runtime"
)

func TicTacToeMatchRPC(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {
	var min_size *int

	max_size := new(int)
	*max_size = 2

	var err error
	reply := MatchRpcBody{
		MatchIds: make([]string, 0),
	}

	// 1) check if a match already exists
	limit := 10
	var label string
	// https://heroiclabs.com/docs/nakama/server-framework/go-runtime/function-reference/#MatchList
	if matches, err := nk.MatchList(ctx, limit, true, label, min_size, max_size, "+label.open:1 +label.tictactoe:1"); err != nil {
		logger.Error("[MatchList]: %s", err)
	} else {
		//logger.Warn("[MatchList]: matches: %#v", matches)
		if len(matches) > 0 {
			for _, match := range matches {
				reply.MatchIds = append(reply.MatchIds, match.MatchId)
			}
		}
	}

	if len(reply.MatchIds) == 0 {
		// 2) create a new match
		// https://heroiclabs.com/docs/nakama/server-framework/go-runtime/function-reference/#MatchCreate
		if matchId, err := nk.MatchCreate(ctx, MODULE_NAME, nil); err != nil {
			logger.Error("[MatchCreate]: %s", err)
			return "", err
		} else {
			//logger.Warn("[MatchCreate]: creating match %s", matchId)
			reply.MatchIds = append(reply.MatchIds, matchId)
		}
	}

	reply2, err := json.Marshal(reply)
	if err != nil {
		return "", err
	} else {
		//logger.Warn("[MatchList]: %#v", reply2)
		return string(reply2), nil
	}
}

/*func KillTicTacToeMatchesRPC(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {
	var minSize *int
	var maxSize *int
	var label string
	if matches, err := nk.MatchList(ctx, 20, true, label, minSize, maxSize, "+label.tictactoe:1"); err != nil {
		logger.Error("[MatchList]: %s", err)
	} else {
		//logger.Warn("[MatchList]: matches: %#v", matches)
		if len(matches) > 0 {
			for _, match := range matches {
				nk.MatchSignal(ctx, match.MatchId, "kill")
			}
		}
	}
	return "{}", nil
}*/

//// WIP

/*func wasItCalledFromServer(ctx context.Context, logger runtime.Logger) bool {
	userIdValue := ctx.Value(runtime.RUNTIME_CTX_USER_ID)

	userId, ok := userIdValue.(string)
	if !ok {
		logger.Error("failed to get user ID from context")
		return false
	}

	logger.Debug("userId: %#v", userId)
	//logger.Debug("ctx.UserId: %#v", ctx.UserId)

	return true
}*/
