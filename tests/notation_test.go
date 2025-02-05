package tests

import (
	"gce/pkg/engine"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNotation(t *testing.T) {

	b := engine.NewDefaultBoard()

	moveNotation := "d4"
	move, err := b.ParseMove(moveNotation)
	assert.Nil(t, err)
	assert.Equal(t, move.NewPiecePos, d4)
	assert.True(t, b.MakeMove(move))

	moveNotation = "e5"
	move, err = b.ParseMove(moveNotation)
	assert.Nil(t, err)
	assert.Equal(t, move.NewPiecePos, e5)
	assert.True(t, b.MakeMove(move))

	moveNotation = "Nd2"
	move, err = b.ParseMove(moveNotation)
	assert.Nil(t, err)
	assert.Equal(t, move.NewPiecePos, engine.PositionToUInt64(3, 1))
	assert.True(t, b.MakeMove(move))

	moveNotation = "Nf6"
	move, err = b.ParseMove(moveNotation)
	assert.Nil(t, err)
	assert.Equal(t, move.NewPiecePos, engine.PositionToUInt64(5, 5))
	assert.True(t, b.MakeMove(move))

	moveNotation = "Ngf3"
	move, err = b.ParseMove(moveNotation)
	assert.Nil(t, err)
	assert.Equal(t, move.NewPiecePos, engine.PositionToUInt64(5, 2))
	assert.True(t, b.MakeMove(move))

	moveNotation = "d6"
	move, err = b.ParseMove(moveNotation)
	assert.Nil(t, err)
	assert.Equal(t, move.NewPiecePos, d6)
	assert.True(t, b.MakeMove(move))

	moveNotation = "dxe5"
	move, err = b.ParseMove(moveNotation)
	assert.Nil(t, err)
	assert.Equal(t, move.NewPiecePos, e5)
	assert.True(t, move.IsCapture)
	assert.True(t, b.MakeMove(move))

	moveNotation = "a6"
	move, err = b.ParseMove(moveNotation)
	assert.Nil(t, err)
	assert.Equal(t, move.NewPiecePos, d6<<3)
	assert.True(t, b.MakeMove(move))

	moveNotation = "exd6"
	move, err = b.ParseMove(moveNotation)
	assert.Nil(t, err)
	assert.Equal(t, move.NewPiecePos, d6)
	assert.True(t, move.IsCapture)
	assert.True(t, b.MakeMove(move))

	moveNotation = "b5"
	move, err = b.ParseMove(moveNotation)
	assert.Nil(t, err)
	assert.Equal(t, move.NewPiecePos, e5<<3)
	assert.True(t, b.MakeMove(move))

	moveNotation = "dxc7"
	move, err = b.ParseMove(moveNotation)
	assert.Nil(t, err)
	assert.Equal(t, move.NewPiecePos, e7<<2)
	assert.True(t, move.IsCapture)
	assert.True(t, b.MakeMove(move))

	moveNotation = "Ra7"
	move, err = b.ParseMove(moveNotation)
	assert.Nil(t, err)
	assert.Equal(t, move.NewPiecePos, e7<<4)
	assert.True(t, b.MakeMove(move))

	moveNotation = "cxd8=Q+"
	move, err = b.ParseMove(moveNotation)
	assert.Nil(t, err)
	assert.Equal(t, move.NewPiecePos, e7<<9)
	assert.True(t, move.IsCapture)
	assert.True(t, move.IsPromotion)
	assert.True(t, b.MakeMove(move))
}
