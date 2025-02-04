package tests

import (
	"gce/pkg/engine"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnPassant(t *testing.T) {
	b := engine.NewDefaultBoard()
	var e2 uint64 = 2048
	var e3 uint64 = e2 << 8 // EnPassant Possition
	var e4 uint64 = e3 << 8
	var e5 uint64 = e4 << 8
	var e7 uint64 = 2_251_799_813_685_248
	var e6 uint64 = e7 >> 8
	var d7 uint64 = 4_503_599_627_370_496
	var d6 uint64 = d7 >> 8
	var d5 uint64 = d6 >> 8
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
