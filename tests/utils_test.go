package tests

import (
	"gce/pkg/chess"
	"testing"
)

func TestInt64ToRowCol(t *testing.T) {
	// Single positions
	for row := 0; row <= 7; row++ {
		for col := 0; col <= 7; col++ {
			i := chess.PositionToUInt64(col, row)
			positions := chess.Int64toPositions(i)
			if len(positions) != 1 {
				t.Errorf("len(Int64toPositions(%v)) = %v", i, positions)
			}
			c := positions[0][0]
			r := positions[0][1]
			if c != col || r != row {
				t.Errorf("RowColToInt64(%v, %v) = %v, Int64ToRowCol(%v) = (%v, %v)", col, row, i, i, c, r)
			}
		}
	}

	// Multiple positions
	i := chess.PositionToUInt64(0, 0) | chess.PositionToUInt64(1, 1) | chess.PositionToUInt64(2, 2)
	positions := chess.Int64toPositions(i)
	if len(positions) != 3 {
		t.Errorf("len(Int64toPositions(%v)) = %v", i, positions)
	}
	isInvalid := positions[0][0] != 0 || positions[0][1] != 0 ||
		positions[1][0] != 1 || positions[1][1] != 1 ||
		positions[2][0] != 2 || positions[2][1] != 2
	if isInvalid {
		t.Errorf("Int64toPositions(%v) = %v", i, positions)
	}
}
