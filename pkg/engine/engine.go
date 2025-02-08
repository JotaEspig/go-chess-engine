package engine

import (
	"fmt"
	"gce/pkg/chess"
	"strconv"
	"time"

	"github.com/charmbracelet/log"
)

func rjust(s, fill string, width int) string {
	for len(s) < width {
		s = fill + s
	}
	return s
}

// AnalysisByDepth returns the evaluation of the board by analyzing it to a certain depth.
func AnalysisByDepth(board chess.Board, depth uint, returnCh chan AnalysisReport, nodesCountch chan chess.Move) AnalysisReport {
	go minimax(board, depth, returnCh, nodesCountch)
	nodes := 0
	captures := 0
	checks := 0
	promotions := 0
	startTime := time.Now()
	for {
		diffToStart := time.Now().Sub(startTime)
		select {
		case move := <-nodesCountch:
			nodes++
			if move.IsCapture {
				captures++
			}
			if move.IsCheck {
				checks++
			}
			if move.IsPromotion {
				promotions++
			}

			toPrint := fmt.Sprintf("%.2f", float64(nodes)/diffToStart.Seconds())
			toPrint = rjust(toPrint, " ", 10)
			fmt.Printf("\rNodes per second: %s", toPrint)
		case analysisReport := <-returnCh:
			fmt.Println()
			log.Infof("Total amount of nodes: %d", nodes)
			log.Infof("Total amount of captures: %d", captures)
			log.Infof("Total amount of checks: %d", checks)
			log.Infof("Total amount of promotions: %d", promotions)
			log.Infof("Total time: %s", time.Now().Sub(startTime))
			return analysisReport
		}
	}
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

func GetEngineMove(start, end chess.Board) (*chess.Move, string) {
	if end.MoveDone == nil {
		return GetEngineMove(start, *end.PrevBoard)
	}

	if start.Hash() == end.Hash() {
		return end.MoveDone, end.MoveToNotation(end.MoveDone)
	}
	return GetEngineMove(start, *end.PrevBoard)
}

func GetEngineLine(start, end *chess.Board) string {
	// Means last position of the board, so it does not have a move
	if end.MoveDone == nil {
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
