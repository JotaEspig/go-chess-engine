package engine

import "github.com/charmbracelet/log"

// MovesFunction is a function that returns all possible new positions for a piece position in the complete Board.
type MovesFunction func(Board, uint64) []Move

func normalMoves(board Board, pieceBoard uint64, directions []int, pieceType PieceType) []Move {
	moves := make([]Move, 0)
	for _, direction := range directions {
		fn := GetDirectionFunc(direction)
		for i := 1; i < 8; i++ {
			if fn == nil {
				log.Fatal("Invalid direction")
			}

			newPieceBoard := fn(pieceBoard, i)
			if newPieceBoard == 0 {
				break
			}

			// check for collision
			var color PartialBoard
			var invertedColor PartialBoard
			if board.Ctx.WhiteToMove {
				color = board.White
				invertedColor = board.Black
			} else {
				color = board.Black
				invertedColor = board.White
			}

			allColorBoard := color.AllBoardMask() & ^pieceBoard // Removes the piece from the board
			if newPieceBoard&allColorBoard != 0 {
				break
			}

			isCapture := newPieceBoard&invertedColor.AllBoardMask() != 0
			move := Move{OldPiecePos: pieceBoard, NewPiecePos: newPieceBoard, IsCapture: isCapture, PieceType: pieceType}
			moves = append(moves, move)
			// Capture check
			if isCapture {
				break
			}
		}
	}

	return moves
}

func knightMove(board Board, pieceBoard uint64, fn func(uint64) uint64) []Move {
	moves := make([]Move, 0)
	newPieceBoard := fn(pieceBoard)
	if newPieceBoard != 0 {
		var color PartialBoard
		var invertedColor PartialBoard
		if board.Ctx.WhiteToMove {
			color = board.White
			invertedColor = board.Black
		} else {
			color = board.Black
			invertedColor = board.White
		}

		allColorBoard := color.AllBoardMask() & ^pieceBoard // Removes the piece from the board
		if newPieceBoard&allColorBoard == 0 {
			isCapture := newPieceBoard&invertedColor.AllBoardMask() != 0
			move := Move{OldPiecePos: pieceBoard, NewPiecePos: newPieceBoard, IsCapture: isCapture, PieceType: KnightType}
			moves = append(moves, move)
		}
	}
	return moves
}

// Includes En Passant and promotion
func PawnMoves(board Board, pieceBoard uint64) []Move {
	moves := make([]Move, 0)
	// Implement pawn moves logic here
	return moves
}

func KnightMoves(board Board, pieceBoard uint64) []Move {
	moves := make([]Move, 0)

	moves = append(moves, knightMove(board, pieceBoard, moveKnightL1)...)
	moves = append(moves, knightMove(board, pieceBoard, moveKnightL2)...)
	moves = append(moves, knightMove(board, pieceBoard, moveKnightL3)...)
	moves = append(moves, knightMove(board, pieceBoard, moveKnightL4)...)
	moves = append(moves, knightMove(board, pieceBoard, moveKnightL5)...)
	moves = append(moves, knightMove(board, pieceBoard, moveKnightL6)...)
	moves = append(moves, knightMove(board, pieceBoard, moveKnightL7)...)
	moves = append(moves, knightMove(board, pieceBoard, moveKnightL8)...)

	return moves
}

func BishopMoves(board Board, pieceBoard uint64) []Move {
	directions := []int{directionUpLeft, directionUpRight, directionDownLeft, directionDownRight}
	return normalMoves(board, pieceBoard, directions, BishopType)
}

func RookMoves(board Board, pieceBoard uint64) []Move {
	directions := []int{directionUp, directionDown, directionLeft, directionRight}
	return normalMoves(board, pieceBoard, directions, RookType)
}

func QueenMoves(board Board, pieceBoard uint64) []Move {
	directions := []int{directionUp, directionDown, directionLeft, directionRight, directionUpLeft, directionUpRight, directionDownLeft, directionDownRight}
	return normalMoves(board, pieceBoard, directions, QueenType)
}

func KingMoves(board Board, pieceBoard uint64) []Move {
	moves := make([]Move, 0)

	directions := []int{directionUp, directionDown, directionLeft, directionRight, directionUpLeft, directionUpRight, directionDownLeft, directionDownRight}
	for _, direction := range directions {
		fn := GetDirectionFunc(direction)
		if fn == nil {
			log.Fatal("Invalid direction")
		}

		newPieceBoard := fn(pieceBoard, 1)
		if newPieceBoard == 0 {
			continue
		}

		// check for collision
		var color PartialBoard
		var invertedColor PartialBoard
		if board.Ctx.WhiteToMove {
			color = board.White
			invertedColor = board.Black
		} else {
			color = board.Black
			invertedColor = board.White
		}

		allColorBoard := color.AllBoardMask() & ^pieceBoard // Removes the piece from the board
		if newPieceBoard&allColorBoard != 0 {
			continue
		}

		isCapture := newPieceBoard&invertedColor.AllBoardMask() != 0
		move := Move{OldPiecePos: pieceBoard, NewPiecePos: newPieceBoard, IsCapture: isCapture, PieceType: KingType}
		moves = append(moves, move)
	}

	return moves
}
