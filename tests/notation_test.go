package tests

import (
	"gce/pkg/chess"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNotation(t *testing.T) {

	b := chess.NewDefaultBoard()

	moveNotation := "d4"
	move, err := b.ParseMove(moveNotation)
	assert.Nil(t, err)
	assert.Equal(t, move.NewPiecePos, d4)
	b.MakePseudoLegalMove(move)

	moveNotation = "e5"
	move, err = b.ParseMove(moveNotation)
	assert.Nil(t, err)
	assert.Equal(t, move.NewPiecePos, e5)
	b.MakePseudoLegalMove(move)

	moveNotation = "Nd2"
	move, err = b.ParseMove(moveNotation)
	assert.Nil(t, err)
	assert.Equal(t, move.NewPiecePos, chess.PositionToUInt64(3, 1))
	b.MakePseudoLegalMove(move)

	moveNotation = "Nf6"
	move, err = b.ParseMove(moveNotation)
	assert.Nil(t, err)
	assert.Equal(t, move.NewPiecePos, chess.PositionToUInt64(5, 5))
	b.MakePseudoLegalMove(move)

	moveNotation = "Ngf3"
	move, err = b.ParseMove(moveNotation)
	assert.Nil(t, err)
	assert.Equal(t, move.NewPiecePos, chess.PositionToUInt64(5, 2))
	b.MakePseudoLegalMove(move)

	moveNotation = "d6"
	move, err = b.ParseMove(moveNotation)
	assert.Nil(t, err)
	assert.Equal(t, move.NewPiecePos, d6)
	b.MakePseudoLegalMove(move)

	moveNotation = "dxe5"
	move, err = b.ParseMove(moveNotation)
	assert.Nil(t, err)
	assert.Equal(t, move.NewPiecePos, e5)
	assert.True(t, move.IsCapture)
	b.MakePseudoLegalMove(move)

	moveNotation = "a6"
	move, err = b.ParseMove(moveNotation)
	assert.Nil(t, err)
	assert.Equal(t, move.NewPiecePos, d6<<3)
	b.MakePseudoLegalMove(move)

	moveNotation = "exd6"
	move, err = b.ParseMove(moveNotation)
	assert.Nil(t, err)
	assert.Equal(t, move.NewPiecePos, d6)
	assert.True(t, move.IsCapture)
	b.MakePseudoLegalMove(move)

	moveNotation = "b5"
	move, err = b.ParseMove(moveNotation)
	assert.Nil(t, err)
	assert.Equal(t, move.NewPiecePos, e5<<3)
	b.MakePseudoLegalMove(move)

	moveNotation = "dxc7"
	move, err = b.ParseMove(moveNotation)
	assert.Nil(t, err)
	assert.Equal(t, move.NewPiecePos, e7<<2)
	assert.True(t, move.IsCapture)
	b.MakePseudoLegalMove(move)

	moveNotation = "Ra7"
	move, err = b.ParseMove(moveNotation)
	assert.Nil(t, err)
	assert.Equal(t, move.NewPiecePos, e7<<4)
	b.MakePseudoLegalMove(move)

	moveNotation = "cxd8=Q+"
	move, err = b.ParseMove(moveNotation)
	assert.Nil(t, err)
	assert.Equal(t, move.NewPiecePos, e7<<9)
	assert.True(t, move.IsCapture)
	assert.True(t, move.IsPromotion)
	b.MakePseudoLegalMove(move)
}
