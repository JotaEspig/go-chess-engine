package tests

import (
	"gce/pkg/engine"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	schoolMate = "r1bqkbnr/ppp2Qpp/2np4/4p3/2B1P3/8/PPPP1PPP/RNB1K1NR b KQkq - 0 4"
	foolsMate  = "rnb1kbnr/pppp1ppp/8/4p3/6Pq/5P2/PPPPP2P/RNBQKBNR w KQkq - 1 3"
)

func TestWhiteKingInCheck(t *testing.T) {
	// Fen where white king is in check
	b := engine.FenToBoard(foolsMate)
	assert.True(t, b.IsKingInCheck(b.White.King.Board, true))

	// Fen where white king is not in check
	b = engine.FenToBoard(schoolMate)
	assert.False(t, b.IsKingInCheck(b.White.King.Board, true))
}

func TestBlackKingInCheck(t *testing.T) {
	// Fen where black king is in check
	b := engine.FenToBoard(schoolMate)
	assert.True(t, b.IsKingInCheck(b.Black.King.Board, false))

	// Fen where black king is not in check
	b = engine.FenToBoard(foolsMate)
	assert.False(t, b.IsKingInCheck(b.Black.King.Board, false))
}

func TestIsWhiteMated(t *testing.T) {
	// In mate
	b := engine.FenToBoard(foolsMate)
	assert.True(t, b.IsWhiteMated())

	// Not in mate
	b = engine.FenToBoard(schoolMate)
	assert.False(t, b.IsWhiteMated())
}

func TestIsBlackMated(t *testing.T) {
	// Not in mate
	b := engine.FenToBoard(foolsMate)
	assert.False(t, b.IsBlackMated())

	b = engine.FenToBoard(schoolMate)
	assert.True(t, b.IsBlackMated())
}
