package engine

const (
	invalid = iota
	directionUp
	directionDown
	directionLeft
	directionRight
	directionUpLeft
	directionUpRight
	directionDownLeft
	directionDownRight
)

type directionFunc func(uint64, int) uint64

func up(piecePos uint64, multiplier int) uint64 {
	return piecePos << (8 * multiplier)
}

func down(piecePos uint64, multiplier int) uint64 {
	return piecePos >> (8 * multiplier)
}

func left(piecePos uint64, multiplier int) uint64 {
	// separate the byte from the rest of the bits to represent the row.
	// So when we move beyond the byte bounds, the value will be 0.
	// If we don't do this, when we lsh the entire piecePos,
	// the piece will go to the other side of the board (e.g. Rook move: A1 - Lsh 1 -> H2)
	var _byte uint8 = 0
	piecePosCopy := piecePos
	count := -1
	for _byte == 0 && count < 8 {
		_byte = uint8(piecePosCopy & 0xFF)
		piecePosCopy >>= 8
		count++
	}
	if count == 8 {
		return 0
	}

	_byte <<= uint8(multiplier)
	return uint64(_byte) << (8 * count)
}

func right(piecePos uint64, multiplier int) uint64 {
	// separate the byte from the rest of the bits to represent the row.
	// So when we move beyond the byte bounds, the value will be 0.
	// If we don't do this, when we lsh the entire piecePos,
	// the piece will go to the other side of the board (e.g. Rook move: A1 - Lsh 1 -> H2)
	var _byte uint8 = 0
	piecePosCopy := piecePos
	count := -1
	for _byte == 0 && count < 8 {
		_byte = uint8(piecePosCopy & 0xFF)
		piecePosCopy >>= 8
		count++
	}
	if count == 8 {
		return 0
	}

	_byte >>= uint8(multiplier)
	return uint64(_byte) << (8 * count)
}

func upLeft(piecePos uint64, multiplier int) uint64 {
	return up(left(piecePos, multiplier), multiplier)
}

func upRight(piecePos uint64, multiplier int) uint64 {
	return up(right(piecePos, multiplier), multiplier)
}

func downLeft(piecePos uint64, multiplier int) uint64 {
	return down(left(piecePos, multiplier), multiplier)
}

func downRight(piecePos uint64, multiplier int) uint64 {
	return down(right(piecePos, multiplier), multiplier)
}

func knightL1(piecePos uint64) uint64 {
	return up(left(piecePos, 1), 2)
}

func knightL2(piecePos uint64) uint64 {
	return up(right(piecePos, 1), 2)
}

func knightL3(piecePos uint64) uint64 {
	return down(left(piecePos, 1), 2)
}

func knightL4(piecePos uint64) uint64 {
	return down(right(piecePos, 1), 2)
}

func knightL5(piecePos uint64) uint64 {
	return left(up(piecePos, 1), 2)
}

func knightL6(piecePos uint64) uint64 {
	return right(up(piecePos, 1), 2)
}

func knightL7(piecePos uint64) uint64 {
	return left(down(piecePos, 1), 2)
}

func knightL8(piecePos uint64) uint64 {
	return right(down(piecePos, 1), 2)
}

func GetDirectionFunc(direction int) directionFunc {
	switch direction {
	case directionUp:
		return up
	case directionDown:
		return down
	case directionLeft:
		return left
	case directionRight:
		return right
	case directionUpLeft:
		return upLeft
	case directionUpRight:
		return upRight
	case directionDownLeft:
		return downLeft
	case directionDownRight:
		return downRight
	default:
		return nil
	}
}
