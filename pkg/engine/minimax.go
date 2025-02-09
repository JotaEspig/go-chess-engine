package engine

import (
	"gce/pkg/chess"
	"math"
	"sort"
)

func minimax(board *chess.Board, depth uint, returnCh chan AnalysisReport, nodesCountch chan struct{}) {
	var analysisReport AnalysisReport
	if board.Ctx.WhiteTurn {
		analysisReport = alphaBetaMax(board, -math.MaxFloat64, math.MaxFloat64, depth, nodesCountch)
	} else {
		analysisReport = alphaBetaMin(board, -math.MaxFloat64, math.MaxFloat64, depth, nodesCountch)
	}
	returnCh <- analysisReport
}

func alphaBetaMax(board *chess.Board, alpha, beta float64, depth uint, nodesCount chan struct{}) AnalysisReport {
	if board.IsMated() || board.IsDraw() || depth == 0 {
		nodesCount <- struct{}{} // Increment nodes count
		report := AnalysisReport{*board, EvaluatePosition(*board), []chess.Move{}}
		return report
	}

	moves := MoveSlice(board.AllLegalMoves())
	sort.Sort(moves)
	bestReport := AnalysisReport{
		Evaluation: -math.MaxFloat64,
	}
	bestMove := chess.Move{} // Used for saving engine line
	for _, move := range moves {
		board.MakeLegalMove(move)
		report := alphaBetaMin(board, alpha, beta, depth-1, nodesCount)
		board.UndoMove()
		if report.Evaluation > bestReport.Evaluation {
			bestReport = report
			bestMove = move
			if report.Evaluation > alpha {
				alpha = report.Evaluation
			}
		}
		if report.Evaluation >= beta {
			return AnalysisReport{*board, report.Evaluation, []chess.Move{}}
		}
	}
	bestReport.Moves = append([]chess.Move{bestMove}, bestReport.Moves...) // Inserts at the begginning
	return bestReport
}

func alphaBetaMin(board *chess.Board, alpha, beta float64, depth uint, nodesCount chan struct{}) AnalysisReport {
	if board.IsMated() || board.IsDraw() || depth == 0 {
		nodesCount <- struct{}{} // Increment nodes count
		report := AnalysisReport{*board, EvaluatePosition(*board), []chess.Move{}}
		return report
	}

	moves := MoveSlice(board.AllLegalMoves())
	sort.Sort(moves)
	bestReport := AnalysisReport{
		Evaluation: math.MaxFloat64,
	}
	bestMove := chess.Move{} // Used for saving engine line
	for _, move := range moves {
		board.MakeLegalMove(move)
		report := alphaBetaMax(board, alpha, beta, depth-1, nodesCount)
		board.UndoMove()
		if report.Evaluation < bestReport.Evaluation {
			bestReport = report
			bestMove = move
			if report.Evaluation < beta {
				beta = report.Evaluation
			}
		}
		if report.Evaluation <= alpha {
			return AnalysisReport{*board, report.Evaluation, []chess.Move{}}
		}
	}
	bestReport.Moves = append([]chess.Move{bestMove}, bestReport.Moves...) // Inserts at the begginning
	return bestReport
}
