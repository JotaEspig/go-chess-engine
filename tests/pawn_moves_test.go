package tests

import (
	"gce/pkg/engine"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnPassant(t *testing.T) {
	b := engine.NewDefaultBoard()

	move := engine.Move{OldPiecePos: e2, NewPiecePos: e4, PieceType: engine.PawnType}
	assert.True(t, b.MakeMove(move))
	assert.Equal(t, e3, b.Ctx.EnPassant)

	move = engine.Move{OldPiecePos: e7, NewPiecePos: e6, PieceType: engine.PawnType}
	assert.True(t, b.MakeMove(move))
	assert.Equal(t, uint64(0), b.Ctx.EnPassant)

	move = engine.Move{OldPiecePos: e4, NewPiecePos: e5, PieceType: engine.PawnType}
	assert.True(t, b.MakeMove(move))
	assert.Equal(t, uint64(0), b.Ctx.EnPassant)

	move = engine.Move{OldPiecePos: d7, NewPiecePos: d5, PieceType: engine.PawnType}
	assert.True(t, b.MakeMove(move))
	assert.Equal(t, d6, b.Ctx.EnPassant)

	move = engine.Move{OldPiecePos: e5, NewPiecePos: d6, PieceType: engine.PawnType, IsCapture: true}
	assert.True(t, b.MakeMove(move))
	assert.Equal(t, uint64(0), b.Ctx.EnPassant)
}
