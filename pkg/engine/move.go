package engine

// Move represents a move in the game.
type Move struct {
	OldPiecePos  uint64
	NewPiecePos  uint64
	IsCastling   bool
	IsCapture    bool
	IsPromotion  bool
	PieceType    PieceType
	NewPieceType PieceType
}

func (m Move) Is2SquarePawnMove() bool {
	if m.PieceType != PawnType {
		return false
	}
	return (m.OldPiecePos>>16 == m.NewPiecePos) || (m.OldPiecePos<<16 == m.NewPiecePos)
}
