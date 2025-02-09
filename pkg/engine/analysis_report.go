package engine

import (
	"gce/pkg/chess"
	"strconv"
	"strings"
)

type AnalysisReport struct {
	BestBoard  chess.Board
	Evaluation float64
	Moves      []chess.Move
}

func (ar AnalysisReport) GetEngineLine() string {
	line := ""
	length := len(ar.Moves)
	turnStartNumber := len(ar.BestBoard.MovesDone) - length
	board := ar.BestBoard
	if board.IsMated() {
		line = "#"
	}

	for length > 0 {
		lastMove := ar.Moves[length-1]
		board.UndoMove()
		notation := board.MoveToNotation(lastMove)
		turnNumber := ""
		if board.Ctx.WhiteTurn {
			moveNumber := strconv.Itoa(turnStartNumber + length)
			turnNumber = moveNumber + "." + " "
		}
		line = " " + turnNumber + notation + line
		length--
	}
	if !board.Ctx.WhiteTurn {
		line = "..." + line
	}
	return strings.TrimSpace(line)
}
