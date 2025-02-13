package tests

import (
	"gce/pkg/chess"
	"gce/pkg/engine"
	"testing"

	"github.com/stretchr/testify/assert"
)

// For now (09/02/2025), it generates more moves than the expected for depth > 2
func TestPerftInitialBoard(t *testing.T) {
	expectedNodes := []uint64{20, 400, 8902, 197_281, 4_865_609, 119_060_324}
	for depth := 1; depth <= 5; depth++ {
		initialBoard := chess.NewDefaultBoard()
		nodes, _ := engine.Perft(initialBoard, uint(depth))
		assert.Equal(t, expectedNodes[depth-1], nodes, "Depth %d", depth)
	}
}

func BenchmarkPerftInitialBoard(b *testing.B) {
	initialBoard := chess.NewDefaultBoard()
	for i := 0; i < b.N; i++ {
		engine.Perft(initialBoard, 4)
	}
}
