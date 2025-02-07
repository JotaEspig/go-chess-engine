package engine

import (
	"gce/pkg/chess"

	"github.com/charmbracelet/log"
)

type MoveSlice []*chess.Move

func MoveSortingScore(move *chess.Move) int {
	score := 0
	if move.IsPromotion {
		score += 100
	}
	if move.IsCheck {
		score += 2
	}
	if move.IsCapture {
		switch move.CapturedPieceType {
		case chess.PawnType:
			score += 1
		case chess.KnightType:
			score += 3
		case chess.BishopType:
			score += 3
		case chess.RookType:
			score += 5
		case chess.QueenType:
			score += 9
		default:
			log.Fatalf("Unknown piece type: %v", move.String())
		}
	}
	return score
}

func (ms MoveSlice) Len() int {
	return len(ms)
}

func (ms MoveSlice) Less(i, j int) bool {
	return MoveSortingScore(ms[i]) > MoveSortingScore(ms[j])
}

func (ms MoveSlice) Swap(i, j int) {
	ms[i], ms[j] = ms[j], ms[i]
}
