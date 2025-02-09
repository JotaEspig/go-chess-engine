package tests

import (
	"gce/pkg/chess"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMakeMove(t *testing.T) {
	b := chess.FenToBoard(startPosition)
	m, err := b.ParseMove("e4")
	assert.Nil(t, err)
	b.MakePseudoLegalMove(m)

}

func TestUndoMove(t *testing.T) {

}
