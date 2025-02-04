package tests

import (
	"gce/pkg/engine"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKingInCheck(t *testing.T) {
	// Fen where white king is in check
	b := engine.FenToBoard(foolsMate)
	assert.True(t, b.IsKingInCheck())

	// Fen where no king is in check
	b = engine.FenToBoard(startPosition)
	assert.False(t, b.IsKingInCheck())
}
func TestIsMated(t *testing.T) {
	// In mate
	b := engine.FenToBoard(foolsMate)
	assert.True(t, b.IsMated())

	// Not in mate
	b = engine.FenToBoard(startPosition)
	assert.False(t, b.IsMated())
}
