package main

import (
	"context"
	"database/sql"

	"github.com/heroiclabs/nakama-common/runtime"
)

const MODULE_NAME = "tictactoe"
const RPC_JOIN_OR_CREATE_NAME = "tictactoe_match"
const RPC_KILL_MATCH_NAME = "tictactoe_kill"

func InitModule(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, initializer runtime.Initializer) error {
	logger.Info("setting up tic-tac-toe...")

	err := initializer.RegisterMatch(MODULE_NAME, newMatch)
	if err != nil {
		logger.Error("[RegisterMatch] error: ", err.Error())
		return err
	}

	if err := initializer.RegisterRpc(RPC_JOIN_OR_CREATE_NAME, TicTacToeMatchRPC); err != nil {
		logger.Error("Unable to register: %v", err)
		return err
	}

	/*if err := initializer.RegisterRpc(RPC_KILL_MATCH_NAME, KillTicTacToeMatchesRPC); err != nil {
		logger.Error("Unable to register: %v", err)
		return err
	}*/

	return nil
}
