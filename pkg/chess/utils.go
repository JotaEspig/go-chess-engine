package chess

import (
	"fmt"
	"log"
)

func Filter[T any](s []T, f func(T) bool) []T {
	var r []T
	for _, v := range s {
		if f(v) {
			r = append(r, v)
		}
	}
	return r
}

func Bla() {
	fmt.Println("bla")
}

func PositionToUInt64(col, row int) uint64 {
	if col < 0 || col > 7 || row < 0 || row > 7 {
		log.Fatalf("Invalid row or col: %v, %v", row, col)
	}
	return 1 << uint(row*8+(7-col))
}

// Int64toPositions converts an int64 to a slice of positions.
// Positions are represented as [2]int, where the first element is the column and the second element is the row.
func Int64toPositions(i uint64) [][2]int {
	// Find the first bit set in the int64.
	// This is the position of the piece.
	// The row is the bit index divided by 8.
	// The column is the bit index modulo 8.
	var positions [][2]int
	for j := 0; j < 64; j++ {
		if i&(1<<uint(j)) != 0 {
			positions = append(positions, [2]int{7 - j%8, j / 8})
		}
	}
	return positions
}
