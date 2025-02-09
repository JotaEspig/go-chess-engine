package engine

import (
	"fmt"
	"gce/pkg/chess"
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
func AnalysisByDepth(board *chess.Board, depth uint, returnCh chan AnalysisReport, nodesCountch chan struct{}) AnalysisReport {
	go minimax(board, depth, returnCh, nodesCountch)
	nodes := 0
	startTime := time.Now()
	for {
		diffToStart := time.Now().Sub(startTime)
		select {
		case <-nodesCountch:
			nodes++

			toPrint := fmt.Sprintf("%.2f", float64(nodes)/diffToStart.Seconds())
			toPrint = rjust(toPrint, " ", 10)
			fmt.Printf("\rNodes per second: %s", toPrint)
		case analysisReport := <-returnCh:
			fmt.Println()
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
	return evaluation
}
