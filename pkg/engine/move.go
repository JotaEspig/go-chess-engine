package engine

// Move represents a move in the game.
type Move struct {
	OldPiecePos  uint64
	NewPiecePos  uint64
	IsCapture    bool
	PieceType    PieceType
	IsPromotion  bool
	NewPieceType PieceType
}

func (m Move) Is2SquarePawnMove() bool {
	if m.PieceType != PawnType {
		return false
	}
	return (m.OldPiecePos>>8 == m.NewPiecePos) || (m.OldPiecePos<<8 == m.NewPiecePos)
}
