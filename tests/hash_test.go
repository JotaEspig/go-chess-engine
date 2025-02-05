package tests

import (
	"fmt"
	"gce/pkg/chess"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestThreesoldHash(t *testing.T) {
	board := chess.NewDefaultBoard()
	hash := board.Hash()
	fmt.Println(hash)

	move, err := board.ParseMove("e4")
	assert.Nil(t, err)
	assert.True(t, board.MakeMove(move))

	hash = board.Hash()
	fmt.Println(hash)
}
