package engine

import (
	"gce/pkg/chess"
)

// AnalysisByDepth returns the evaluation of the board by analyzing it to a certain depth.
func AnalysisByDepth(board chess.Board, depth uint) float64 {
	return minimax(board, depth)
}

// EvaluatePosition returns the evaluation of the current board without doing any moves.
func EvaluatePosition(board chess.Board) float64 {
	if board.IsMated() {
		if board.Ctx.WhiteTurn {
			return 10000
		} else {
			return -10000
		}
	} else if board.IsDraw() {
		return 0
	}

	evaluation := float64(board.MaterialValueBalance())
	evaluation += BoardEvaluationByPieceSquareTable(board)

	// Simulates "tempo"
	if board.Ctx.WhiteTurn {
		evaluation += 0.15
	} else {
		evaluation -= 0.15
	}

	return evaluation
}
