package main

/*
0 1 2
3 4 5
6 7 8
*/

func getFilledPiecesCount(state *MatchState) int {
	count := 0

	for i := 0; i < 9; i++ {
		cell := state.board[i]
		if cell != MarkEmpty {
			count += 1
		}
	}

	return count
}

func nextMarkToPlay(state *MatchState) Mark {
	if getFilledPiecesCount(state)%2 == 0 {
		return MarkX
	}
	return MarkO
}

func hasWon(state *MatchState, mark Mark) bool {
	b := state.board

	var combs [8][3]int

	// horizontals
	combs[0] = [3]int{0, 1, 2}
	combs[1] = [3]int{3, 4, 5}
	combs[2] = [3]int{6, 7, 8}

	// verticals
	combs[3] = [3]int{0, 3, 6}
	combs[4] = [3]int{1, 4, 7}
	combs[5] = [3]int{2, 5, 8}

	// diagonals
	combs[6] = [3]int{0, 4, 8}
	combs[7] = [3]int{2, 4, 6}

	for _, comb := range combs {
		if b[comb[0]] == mark && b[comb[1]] == mark && b[comb[2]] == mark {
			return true
		}
	}
	return false
}

func isBoardFull(state *MatchState) bool {
	return getFilledPiecesCount(state) == 9
}
