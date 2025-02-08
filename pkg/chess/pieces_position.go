package chess

import (
	"math/bits"

	"github.com/charmbracelet/log"
)

// PiecesPosition is a bitboard representing the positions of pieces of a single type on the board.
// Each PiecesPosition corresponds to a specific type of piece.
type PiecesPosition struct {
	// Bitboard representing the positions of the pieces.
	Board uint64
	Type  PieceType
}

func (pp PiecesPosition) Value() uint64 {
	multiplier := pp.Type.Value()

	// Count the number of bits set in the bitboard.
	// This is the number of pieces of this type on the board.
	count := uint64(bits.OnesCount64(pp.Board))
	// Multiply the number of pieces by the value of the piece.
	return count * multiplier
}

func (pp PiecesPosition) AllPossibleMoves(b Board) []*Move {
	var moves []*Move
	movesFn := GetMovesFunction(pp.Type)
	if movesFn == nil {
		log.Fatal("Invalid piece type")
	}

	bitboard := pp.Board
	for bitboard != 0 {
		i := bits.TrailingZeros64(bitboard)
		newMoves := movesFn(b, 1<<i)
		moves = append(moves, newMoves...)
		bitboard &= bitboard - 1 // Removes the LSB
	}
	return moves
}

// SetPieceAt sets the bit at the given column and row to 1.
// Column and Row starts at 0. col == 0 and row == 0 means A1
func (pp *PiecesPosition) SetPieceAt(col, row int) {
	pp.Board |= PositionToUInt64(col, row)
}

// ClearPieceAt sets the bit at the given column and row to 0.
// Column and Row starts at 0. col == 0 and row == 0 means A1
func (pp *PiecesPosition) ClearPieceAt(col, row int) {
	pp.Board &= ^(PositionToUInt64(col, row))
}
