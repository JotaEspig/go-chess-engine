package engine

import (
	"gce/pkg/chess"
)

func minimax(board chess.Board, depth uint) float64 {
	if board.IsMated() || board.IsDraw() || depth == 0 {
		return EvaluatePosition(board)
	}

	var best float64
	if board.Ctx.WhiteTurn {
		best = -10000
		for _, move := range board.AllLegalMoves() {
			board.MakeMove(move)
			val := minimax(board, depth-1)
			if val > best {
				best = val
			}
			board = *board.PrevBoard // Restore board
		}
	} else {
		best = 10000
		for _, move := range board.AllLegalMoves() {
			board.MakeMove(move)
			val := minimax(board, depth-1)
			if val < best {
				best = val
			}
			board = *board.PrevBoard // Restore board
		}
	}

	return best
}
