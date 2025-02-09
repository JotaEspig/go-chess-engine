package chess

import "math/bits"

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
	// Find the first nonzero byte
	count := bits.TrailingZeros64(piecePos) / 8

	// Extract the relevant byte
	_byte := uint8((piecePos >> (8 * count)) & 0xFF)
	_byte <<= uint8(multiplier)
	newPos := uint64(_byte) << (8 * count)
	return newPos
}

func moveRight(piecePos uint64, multiplier int) uint64 {
	// Find the first nonzero byte
	count := bits.TrailingZeros64(piecePos) / 8

	// Extract the relevant byte
	_byte := uint8((piecePos >> (8 * count)) & 0xFF)
	_byte >>= uint8(multiplier)
	newPos := uint64(_byte) << (8 * count)
	return newPos
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

func moveL1(piecePos uint64) uint64 {
	return moveUp(moveLeft(piecePos, 1), 2)
}

func moveL2(piecePos uint64) uint64 {
	return moveUp(moveRight(piecePos, 1), 2)
}

func moveL3(piecePos uint64) uint64 {
	return moveDown(moveLeft(piecePos, 1), 2)
}

func moveL4(piecePos uint64) uint64 {
	return moveDown(moveRight(piecePos, 1), 2)
}

func moveL5(piecePos uint64) uint64 {
	return moveLeft(moveUp(piecePos, 1), 2)
}

func moveL6(piecePos uint64) uint64 {
	return moveRight(moveUp(piecePos, 1), 2)
}

func moveL7(piecePos uint64) uint64 {
	return moveLeft(moveDown(piecePos, 1), 2)
}

func moveL8(piecePos uint64) uint64 {
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
