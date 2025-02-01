package engine

// PiecesPosition is a bitboard representing the positions of pieces of a single type on the board.
// Each PiecesPosition corresponds to a specific type of piece.
type PiecesPosition struct {
	// Bitboard representing the positions of the pieces.
	Board int64
	Type  PieceType
}

func (pp PiecesPosition) Value() int64 {
	multiplier := pp.Type.Value()

	// Count the number of bits set in the bitboard.
	// This is the number of pieces of this type on the board.
	var count int64 = 0
	for pp.Board != 0 {
		count += pp.Board & 1
		pp.Board >>= 1
	}
	// Multiply the number of pieces by the value of the piece.
	return count * multiplier
}

// SetPieceAt sets the bit at the given column and row to 1.
// Column and Row starts at 0. col == 0 and row == 0 means A1
func (pp *PiecesPosition) SetPieceAt(col, row int) {
	pp.Board |= PositionToInt64(col, row)
}

// ClearPieceAt sets the bit at the given column and row to 0.
// Column and Row starts at 0. col == 0 and row == 0 means A1
func (pp *PiecesPosition) ClearPieceAt(col, row int) {
	pp.Board &= ^(PositionToInt64(col, row))
}
