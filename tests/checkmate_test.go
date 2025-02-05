package tests

import (
	"gce/pkg/chess"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKingInCheck(t *testing.T) {
	// Fen where white king is in check
	b := chess.FenToBoard(foolsMate)
	assert.True(t, b.IsKingInCheck())

	// Fen where no king is in check
	b = chess.FenToBoard(startPosition)
	assert.False(t, b.IsKingInCheck())
}
func TestIsMated(t *testing.T) {
	// In mate
	b := chess.FenToBoard(foolsMate)
	assert.True(t, b.IsMated())

	// Not in mate
	b = chess.FenToBoard(startPosition)
	assert.False(t, b.IsMated())
}
