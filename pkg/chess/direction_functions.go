package chess

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

func moveUp(piecePos uint64, multiplier int) uint64 {
	return piecePos << (8 * multiplier)
}

func moveDown(piecePos uint64, multiplier int) uint64 {
	return piecePos >> (8 * multiplier)
}

func moveLeft(piecePos uint64, multiplier int) uint64 {
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

func moveRight(piecePos uint64, multiplier int) uint64 {
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

func moveUpLeft(piecePos uint64, multiplier int) uint64 {
	return moveUp(moveLeft(piecePos, multiplier), multiplier)
}

func moveUpRight(piecePos uint64, multiplier int) uint64 {
	return moveUp(moveRight(piecePos, multiplier), multiplier)
}

func moveDownLeft(piecePos uint64, multiplier int) uint64 {
	return moveDown(moveLeft(piecePos, multiplier), multiplier)
}

func moveDownRight(piecePos uint64, multiplier int) uint64 {
	return moveDown(moveRight(piecePos, multiplier), multiplier)
}

func moveKnightL1(piecePos uint64) uint64 {
	return moveUp(moveLeft(piecePos, 1), 2)
}

func moveKnightL2(piecePos uint64) uint64 {
	return moveUp(moveRight(piecePos, 1), 2)
}

func moveKnightL3(piecePos uint64) uint64 {
	return moveDown(moveLeft(piecePos, 1), 2)
}

func moveKnightL4(piecePos uint64) uint64 {
	return moveDown(moveRight(piecePos, 1), 2)
}

func moveKnightL5(piecePos uint64) uint64 {
	return moveLeft(moveUp(piecePos, 1), 2)
}

func moveKnightL6(piecePos uint64) uint64 {
	return moveRight(moveUp(piecePos, 1), 2)
}

func moveKnightL7(piecePos uint64) uint64 {
	return moveLeft(moveDown(piecePos, 1), 2)
}

func moveKnightL8(piecePos uint64) uint64 {
	return moveRight(moveDown(piecePos, 1), 2)
}

func GetDirectionFunc(direction int) directionFunc {
	switch direction {
	case directionUp:
		return moveUp
	case directionDown:
		return moveDown
	case directionLeft:
		return moveLeft
	case directionRight:
		return moveRight
	case directionUpLeft:
		return moveUpLeft
	case directionUpRight:
		return moveUpRight
	case directionDownLeft:
		return moveDownLeft
	case directionDownRight:
		return moveDownRight
	default:
		return nil
	}
}
