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

	// Color configs
	var dirFn func(uint64, int) uint64
	var isInPromotionRow func(uint64) bool
	var isInInitialRow func(uint64) bool
	if board.Ctx.WhiteToMove {
		dirFn = moveUp
		isInPromotionRow = func(pos uint64) bool {
			return pos>>56 != 0
		}
		isInInitialRow = func(pos uint64) bool {
			return pos<<48 != 0
		}
	} else {
		dirFn = moveDown
		isInPromotionRow = func(pos uint64) bool {
			return pos<<56 != 0
		}
		isInInitialRow = func(pos uint64) bool {
			return pos>>48 != 0
		}
	}

	// Move forward
	newPieceBoard := dirFn(pieceBoard, 1)
	if newPieceBoard == 0 {
		log.Fatal("Pawn is already at the last row, How?")
	}

	// check for collision
	whiteMask := board.White.AllBoardMask()
	blackMask := board.Black.AllBoardMask()
	allColorBoard := whiteMask | blackMask
	// If there's no collision
	if newPieceBoard&allColorBoard == 0 {
		isPromotion := isInPromotionRow(newPieceBoard)
		move := Move{OldPiecePos: pieceBoard, NewPiecePos: newPieceBoard, IsPromotion: isPromotion, PieceType: PawnType}
		if !isPromotion {
			moves = append(moves, move)
		} else {
			for _, pieceType := range []PieceType{QueenType, RookType, BishopType, KnightType} {
				move.NewPieceType = pieceType
				moves = append(moves, move)
			}
		}
	}

	// Double pawn move
	// Check if pawn is at 2 row
	if isInInitialRow(pieceBoard) {
		// Check if row in front is not blocked
		if newPieceBoard != 0 {
			newPieceBoard = dirFn(newPieceBoard, 1)
			if newPieceBoard&allColorBoard == 0 {
				move := Move{OldPiecePos: pieceBoard, NewPiecePos: newPieceBoard, PieceType: PawnType}
				moves = append(moves, move)
			}
		}
	}

	// Check capture moves
	captureLeft := dirFn(moveLeft(pieceBoard, 1), 1)
	captureRight := dirFn(moveRight(pieceBoard, 1), 1)
	capturePoss := []uint64{captureLeft, captureRight}
	for _, capturePos := range capturePoss {
		if capturePos&blackMask == 0 && capturePos&board.Ctx.EnPassant == 0 {
			continue
		}

		isPromotion := isInPromotionRow(capturePos)
		move := Move{OldPiecePos: pieceBoard, NewPiecePos: capturePos, IsCapture: true, IsPromotion: isPromotion, PieceType: PawnType}
		if !isPromotion {
			moves = append(moves, move)
		} else {
			for _, pieceType := range []PieceType{QueenType, RookType, BishopType, KnightType} {
				move.NewPieceType = pieceType
				moves = append(moves, move)
			}
		}
	}

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
