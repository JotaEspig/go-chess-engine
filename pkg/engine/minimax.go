package engine

import (
	"gce/pkg/chess"
	"math"
	"sort"
)

func minimax(board chess.Board, depth uint, returnCh chan AnalysisReport, nodesCountch chan chess.Move) {
	var analysisReport AnalysisReport
	if board.Ctx.WhiteTurn {
		analysisReport = alphaBetaMax(board, -math.MaxFloat64, math.MaxFloat64, depth, nodesCountch)
	} else {
		analysisReport = alphaBetaMin(board, -math.MaxFloat64, math.MaxFloat64, depth, nodesCountch)
	}
	returnCh <- analysisReport
}

func alphaBetaMax(board chess.Board, alpha, beta float64, depth uint, nodesCount chan chess.Move) AnalysisReport {
	if entry, ok := transpositionalTable[board.HashBoardWithContext()]; ok {
		return entry
	}

	if board.IsMated() || board.IsDraw() || depth == 0 {
		nodesCount <- *board.PrevBoard.MoveDone // Increment nodes count
		report := AnalysisReport{board, EvaluatePosition(board)}
		transpositionalTable[board.HashBoardWithContext()] = report
		return report
	}

	moves := MoveSlice(board.AllLegalMoves())
	sort.Sort(moves)
	bestReport := AnalysisReport{
		Evaluation: -math.MaxFloat64,
	}
	for _, move := range moves {
		board.MoveDone = nil // Resets
		board.MakeMove(move)
		report := alphaBetaMin(board, alpha, beta, depth-1, nodesCount)
		board = *board.PrevBoard // Restore board
		if report.Evaluation > bestReport.Evaluation {
			bestReport = report
			if report.Evaluation > alpha {
				alpha = report.Evaluation
			}
		}
		if report.Evaluation >= beta {
			return AnalysisReport{board, report.Evaluation}
		}
	}
	return bestReport
}

func alphaBetaMin(board chess.Board, alpha, beta float64, depth uint, nodesCount chan chess.Move) AnalysisReport {
	if entry, ok := transpositionalTable[board.HashBoardWithContext()]; ok {
		return entry
	}

	if board.IsMated() || board.IsDraw() || depth == 0 {
		nodesCount <- *board.PrevBoard.MoveDone // Increment nodes count
		report := AnalysisReport{board, EvaluatePosition(board)}
		transpositionalTable[board.HashBoardWithContext()] = report
		return report
	}

	moves := MoveSlice(board.AllLegalMoves())
	sort.Sort(moves)
	bestReport := AnalysisReport{
		Evaluation: math.MaxFloat64,
	}
	for _, move := range moves {
		board.MoveDone = nil // Resets
		board.MakeMove(move)
		report := alphaBetaMax(board, alpha, beta, depth-1, nodesCount)
		board = *board.PrevBoard // Restore board
		if report.Evaluation < bestReport.Evaluation {
			bestReport = report
			if report.Evaluation < beta {
				beta = report.Evaluation
			}
		}
		if report.Evaluation <= alpha {
			return AnalysisReport{board, report.Evaluation}
		}
	}
	return bestReport
}
