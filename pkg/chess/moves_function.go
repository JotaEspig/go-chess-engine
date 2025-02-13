package chess

import (
	"github.com/charmbracelet/log"
)

// MovesFunction is a function that returns all possible new positions for a piece position in the complete Board.
type MovesFunction func(Board, uint64) []Move

func normalMoves(board Board, pieceBoard uint64, directions []int, pieceType PieceType) []Move {
	moves := make([]Move, 0, len(directions)*7)
	var ourPb PartialBoard
	var enemyPb PartialBoard
	if board.Ctx.WhiteTurn {
		ourPb = board.White
		enemyPb = board.Black
	} else {
		ourPb = board.Black
		enemyPb = board.White
	}
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
			allColorBoard := ourPb.AllBoardMask() & ^pieceBoard // Removes the piece from the board
			if newPieceBoard&allColorBoard != 0 {
				break
			}

			isCapture := newPieceBoard&enemyPb.AllBoardMask() != 0
			capturedPiece := InvalidType
			if isCapture {
				capturedPiece = enemyPb.GetPieceTypeByPos(newPieceBoard)
			}

			move := Move{
				OldPiecePos:       pieceBoard,
				NewPiecePos:       newPieceBoard,
				IsCapture:         isCapture,
				CapturedPieceType: capturedPiece,
				PieceType:         pieceType,
			}
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
	moves := make([]Move, 0, 1)
	var ourPb PartialBoard
	var enemyPb PartialBoard
	if board.Ctx.WhiteTurn {
		ourPb = board.White
		enemyPb = board.Black
	} else {
		ourPb = board.Black
		enemyPb = board.White
	}

	newPieceBoard := fn(pieceBoard)
	if newPieceBoard != 0 {
		allColorBoard := ourPb.AllBoardMask() & ^pieceBoard // Removes the piece from the board
		if newPieceBoard&allColorBoard == 0 {
			isCapture := newPieceBoard&enemyPb.AllBoardMask() != 0
			capturedPiece := InvalidType
			if isCapture {
				capturedPiece = enemyPb.GetPieceTypeByPos(newPieceBoard)
			}
			move := Move{
				OldPiecePos:       pieceBoard,
				NewPiecePos:       newPieceBoard,
				IsCapture:         isCapture,
				CapturedPieceType: capturedPiece,
				PieceType:         KnightType,
			}
			moves = append(moves, move)
		}
	}
	return moves
}

// Includes En Passant and promotion
func PawnMoves(board Board, pieceBoard uint64) []Move {
	moves := make([]Move, 0, 8)

	// Color configs
	var dirFn func(uint64, int) uint64
	var isInPromotionRow func(uint64) bool
	var isInInitialRow func(uint64) bool
	if board.Ctx.WhiteTurn {
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
	var enemyMask uint64
	var enemyPb PartialBoard
	if board.Ctx.WhiteTurn {
		enemyMask = blackMask
		enemyPb = board.Black
	} else {
		enemyMask = whiteMask
		enemyPb = board.White
	}
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

		// Double pawn move
		// Check if pawn is at 2 row
		if isInInitialRow(pieceBoard) {
			// Check if row in front is not blocked
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
		if capturePos&enemyMask == 0 && capturePos&board.Ctx.EnPassant == 0 {
			continue
		}

		IsEnPassant := capturePos&board.Ctx.EnPassant != 0
		isPromotion := isInPromotionRow(capturePos)
		capturedPiece := enemyPb.GetPieceTypeByPos(capturePos)
		if IsEnPassant {
			capturedPiece = PawnType
		}
		move := Move{
			OldPiecePos:       pieceBoard,
			NewPiecePos:       capturePos,
			IsCapture:         true,
			CapturedPieceType: capturedPiece,
			IsPromotion:       isPromotion,
			IsEnPassant:       IsEnPassant,
			PieceType:         PawnType,
		}
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
	moves := make([]Move, 0, 8)
	moves = append(moves, knightMove(board, pieceBoard, moveL1)...)
	moves = append(moves, knightMove(board, pieceBoard, moveL2)...)
	moves = append(moves, knightMove(board, pieceBoard, moveL3)...)
	moves = append(moves, knightMove(board, pieceBoard, moveL4)...)
	moves = append(moves, knightMove(board, pieceBoard, moveL5)...)
	moves = append(moves, knightMove(board, pieceBoard, moveL6)...)
	moves = append(moves, knightMove(board, pieceBoard, moveL7)...)
	moves = append(moves, knightMove(board, pieceBoard, moveL8)...)

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
	moves := make([]Move, 0, 8)

	var ourPb PartialBoard
	var enemyPb PartialBoard
	if board.Ctx.WhiteTurn {
		ourPb = board.White
		enemyPb = board.Black
	} else {
		ourPb = board.Black
		enemyPb = board.White
	}
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
		allColorBoard := ourPb.AllBoardMask() & ^pieceBoard // Removes the piece from the board
		if newPieceBoard&allColorBoard != 0 {
			continue
		}

		isCapture := newPieceBoard&enemyPb.AllBoardMask() != 0
		capturedPiece := InvalidType
		if isCapture {
			capturedPiece = enemyPb.GetPieceTypeByPos(newPieceBoard)
		}

		move := Move{
			OldPiecePos:       pieceBoard,
			NewPiecePos:       newPieceBoard,
			IsCapture:         isCapture,
			CapturedPieceType: capturedPiece,
			PieceType:         KingType,
		}
		moves = append(moves, move)
	}

	return moves
}
