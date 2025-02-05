package engine

import "gce/pkg/chess"

type Engine struct {
}

// AnalysisByDepth returns the evaluation of the board by analyzing it to a certain depth.
func AnalysisByDepth(board chess.Board) float32 {
	return 0
}

// Evaluate returns the evaluation of the current board without doing any moves.
func Evaluate(board chess.Board) float32 {
	return float32(board.MaterialValueBalance())
}
