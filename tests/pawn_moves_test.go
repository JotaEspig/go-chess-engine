package tests

import (
	"gce/pkg/chess"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnPassant(t *testing.T) {
	b := chess.NewDefaultBoard()

	move := &chess.Move{OldPiecePos: e2, NewPiecePos: e4, PieceType: chess.PawnType}
	assert.True(t, b.MakeMove(move))
	assert.Equal(t, e3, b.Ctx.EnPassant)

	move = &chess.Move{OldPiecePos: e7, NewPiecePos: e6, PieceType: chess.PawnType}
	assert.True(t, b.MakeMove(move))
	assert.Equal(t, uint64(0), b.Ctx.EnPassant)

	move = &chess.Move{OldPiecePos: e4, NewPiecePos: e5, PieceType: chess.PawnType}
	assert.True(t, b.MakeMove(move))
	assert.Equal(t, uint64(0), b.Ctx.EnPassant)

	move = &chess.Move{OldPiecePos: d7, NewPiecePos: d5, PieceType: chess.PawnType}
	assert.True(t, b.MakeMove(move))
	assert.Equal(t, d6, b.Ctx.EnPassant)

	move = &chess.Move{OldPiecePos: e5, NewPiecePos: d6, PieceType: chess.PawnType, IsCapture: true}
	assert.True(t, b.MakeMove(move))
	assert.Equal(t, uint64(0), b.Ctx.EnPassant)
}
