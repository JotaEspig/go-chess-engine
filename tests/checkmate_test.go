package tests

import (
	"gce/pkg/engine"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	schoolMate    = "r1bqkbnr/ppp2Qpp/2np4/4p3/2B1P3/8/PPPP1PPP/RNB1K1NR b KQkq - 0 4"
	foolsMate     = "rnb1kbnr/pppp1ppp/8/4p3/6Pq/5P2/PPPPP2P/RNBQKBNR w KQkq - 1 3"
	startPosition = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
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
