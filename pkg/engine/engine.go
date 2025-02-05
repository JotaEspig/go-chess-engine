package engine

import (
	"gce/pkg/chess"
	"gce/pkg/utils"
)

type Engine struct {
}

func analysisByDepth(board *chess.Board, depth uint) float64 {
	return 0
}

// AnalysisByDepth returns the evaluation of the board by analyzing it to a certain depth.
func AnalysisByDepth(board *chess.Board, depth uint) float64 {
	copyBoard := board.Copy()
	return analysisByDepth(&copyBoard, depth)
}

// EvaluatePosition returns the evaluation of the current board without doing any moves.
func EvaluatePosition(board *chess.Board) float64 {
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
	evaluation += BoardEvaluationByPieceSquareTable(*board)

	// Simulates "tempo"
	if board.Ctx.WhiteTurn {
		evaluation += 0.15
	} else {
		evaluation -= 0.15
	}

	whiteTurn := board.Ctx.WhiteTurn
	// Ratio between of many moves each color has
	board.Ctx.WhiteTurn = true
	whiteMoves := board.White.AllPossibleMoves(*board)
	board.Ctx.WhiteTurn = false
	blackMoves := board.Black.AllPossibleMoves(*board)
	// Filter out invalid moves
	whiteMoves = utils.Filter(whiteMoves, func(m chess.Move) bool {
		copyBoard := board.Copy()
		copyBoard.Ctx.WhiteTurn = true
		return copyBoard.MakeMove(m)
	})
	blackMoves = utils.Filter(blackMoves, func(m chess.Move) bool {
		copyBoard := board.Copy()
		copyBoard.Ctx.WhiteTurn = false
		return copyBoard.MakeMove(m)
	})

	// Restores whiteturn
	board.Ctx.WhiteTurn = whiteTurn

	if len(whiteMoves) > len(blackMoves) {
		evaluation += float64(len(whiteMoves)/len(blackMoves)) - 1
	} else {
		evaluation -= float64(len(blackMoves)/len(whiteMoves)) - 1
	}

	return evaluation
}
