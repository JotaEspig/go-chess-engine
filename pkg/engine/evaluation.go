package engine

// AnalysisByDepth returns the evaluation of the board by analyzing it to a certain depth.
func AnalysisByDepth(board Board) float32 {
	return 0
}

// Evaluate returns the evaluation of the current board without doing any moves.
func Evaluate(board Board) float32 {
	return float32(board.MaterialValueBalance())
}
