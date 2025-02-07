package chess

import "fmt"

// Move represents a move in the game.
type Move struct {
	OldPiecePos       uint64
	NewPiecePos       uint64
	IsCastling        bool
	IsCapture         bool
	IsPromotion       bool
	IsCheck           bool
	PieceType         PieceType
	NewPieceType      PieceType
	CapturedPieceType PieceType
	IsCheckFieldSet   bool
}

func (m Move) Is2SquarePawnMove() bool {
	if m.PieceType != PawnType {
		return false
	}
	return (m.OldPiecePos>>16 == m.NewPiecePos) || (m.OldPiecePos<<16 == m.NewPiecePos)
}

func (m Move) String() string {
	return fmt.Sprintf("Move{OldPiecePos: %d, NewPiecePos: %d, IsCastling: %t, IsCapture: %t, IsPromotion: %t, IsCheck: %t, PieceType: %s, NewPieceType: %s, CapturedPieceType: %s}",
		m.OldPiecePos, m.NewPiecePos, m.IsCastling, m.IsCapture, m.IsPromotion, m.IsCheck, m.PieceType.String(), m.NewPieceType.String(), m.CapturedPieceType.String())
}
