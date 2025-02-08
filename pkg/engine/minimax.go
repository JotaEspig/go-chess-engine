package engine

import (
	"gce/pkg/chess"
	"math"
	"sort"
)

func minimax(board chess.Board, depth uint, returnCh chan AnalysisReport, nodesCountch chan struct{}) {
	var analysisReport AnalysisReport
	if board.Ctx.WhiteTurn {
		bestBoard, eval := alphaBetaMax(board, -math.MaxFloat64, math.MaxFloat64, depth, nodesCountch)
		analysisReport = AnalysisReport{BestBoard: bestBoard, Evaluation: eval}
	} else {
		bestBoard, eval := alphaBetaMin(board, -math.MaxFloat64, math.MaxFloat64, depth, nodesCountch)
		analysisReport = AnalysisReport{BestBoard: bestBoard, Evaluation: eval}
	}
	returnCh <- analysisReport
}

func alphaBetaMax(board chess.Board, alpha, beta float64, depth uint, nodesCount chan struct{}) (chess.Board, float64) {
	if board.IsMated() || board.IsDraw() || depth == 0 {
		board.MoveDone = nil
		return board, EvaluatePosition(board)
	}

	moves := MoveSlice(board.AllLegalMoves())
	sort.Sort(moves)
	bestValue := -math.MaxFloat64
	bestBoard := chess.Board{}
	for _, move := range moves {
		nodesCount <- struct{}{} // Increment nodes count
		board.MoveDone = nil     // Resets
		board.MakeMove(move)
		newBoard, val := alphaBetaMin(board, alpha, beta, depth-1, nodesCount)
		board = *board.PrevBoard // Restore board
		if val > bestValue {
			bestValue = val
			bestBoard = newBoard
			if val > alpha {
				alpha = val
			}
		}
		if val >= beta {
			return chess.Board{}, val
		}
	}
	return bestBoard, bestValue
}

func alphaBetaMin(board chess.Board, alpha, beta float64, depth uint, nodesCount chan struct{}) (chess.Board, float64) {
	if board.IsMated() || board.IsDraw() || depth == 0 {
		return board, EvaluatePosition(board)
	}

	moves := MoveSlice(board.AllLegalMoves())
	sort.Sort(moves)
	bestValue := math.MaxFloat64
	bestBoard := chess.Board{}
	for _, move := range moves {
		nodesCount <- struct{}{} // Increment nodes count
		board.MoveDone = nil     // Resets
		board.MakeMove(move)
		newBoard, val := alphaBetaMax(board, alpha, beta, depth-1, nodesCount)
		board = *board.PrevBoard // Restore board
		if val < bestValue {
			bestValue = val
			bestBoard = newBoard
			if val < beta {
				beta = val
			}
		}
		if val <= alpha {
			return chess.Board{}, val
		}
	}
	return bestBoard, bestValue
}
