package engine

// PartialBoard represents a board with only pieces of the same color on it.
type PartialBoard struct {
	Pawns   PiecesPosition
	Knights PiecesPosition
	Bishops PiecesPosition
	Rooks   PiecesPosition
	Queens  PiecesPosition
	King    PiecesPosition
}

func NewPartialBoard() PartialBoard {
	return PartialBoard{
		Pawns:   PiecesPosition{Type: PawnType},
		Knights: PiecesPosition{Type: KnightType},
		Bishops: PiecesPosition{Type: BishopType},
		Rooks:   PiecesPosition{Type: RookType},
		Queens:  PiecesPosition{Type: QueenType},
		King:    PiecesPosition{Type: KingType},
	}
}

// MaterialValue returns the total value of all the pieces on the board.
func (pb PartialBoard) MaterialValue() int64 {
	return pb.Pawns.Value() + pb.Knights.Value() + pb.Bishops.Value() + pb.Rooks.Value() + pb.Queens.Value()
}
