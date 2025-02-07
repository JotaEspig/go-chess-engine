package engine

import (
	"gce/pkg/chess"
	"strconv"
)

// AnalysisByDepth returns the evaluation of the board by analyzing it to a certain depth.
func AnalysisByDepth(board chess.Board, depth uint) (chess.Board, float64) {
	return minimax(board, depth)
}

// EvaluatePosition returns the evaluation of the current board without doing any moves.
func EvaluatePosition(board chess.Board) float64 {
	if board.IsMated() {
		if board.Ctx.WhiteTurn {
			return -10000
		} else {
			return 10000
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

func GetEngineLine(start, end *chess.Board) string {
	// Means last position of the board, so it does not have a move
	if end.MoveDone == (chess.Move{}) {
		mateSuffix := ""
		if end.IsMated() {
			mateSuffix = "#"
		}
		// If it is mate on this turn, so we need to add the suffix to the previous move
		return GetEngineLine(start, end.PrevBoard) + mateSuffix
	}

	moveNotation := end.MoveToNotation(end.MoveDone)
	if end.PrevBoard == nil {
		return "1. " + moveNotation
	}
	moveNumberIfNeeded := ""
	if end.Ctx.WhiteTurn {
		moveNumberInt := end.Ctx.MoveNumber
		moveNumberIfNeeded = strconv.Itoa(int(moveNumberInt)) + ". "
	}
	prev := ""
	if start.Hash() != end.Hash() {
		checkSuffix := ""
		if end.IsKingInCheck() {
			checkSuffix = "+"
		}
		// If it is check on this turn, so we need to add the suffix to the previous move
		prev = GetEngineLine(start, end.PrevBoard) + checkSuffix + " "
	} else if !end.Ctx.WhiteTurn {
		prev = "... "
	}

	return prev + moveNumberIfNeeded + moveNotation
}
