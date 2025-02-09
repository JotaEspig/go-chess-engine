package chess

import "gce/pkg/utils"

// Board represents a full board with pieces of both colors on it.
type Board struct {
	White PartialBoard
	Black PartialBoard
	Ctx   Context

	MovesDone   []Move
	PreviousCtx []Context
}

func NewDefaultBoard() *Board {
	return FenToBoard(DefaultStartFen)
}

func NewBoard() *Board {
	return &Board{
		White:       NewPartialBoard(),
		Black:       NewPartialBoard(),
		MovesDone:   make([]Move, 0, 100),
		PreviousCtx: make([]Context, 0, 100),
	}
}

func (b Board) IsValidPosition() bool {
	// Check if there are no pieces in the same position
	whitePieces := b.White.AllBoardMask()
	blackPieces := b.Black.AllBoardMask()
	if whitePieces&blackPieces != 0 {
		return false
	}

	// Check if the black king is in check while it's white's turn. Meaning that an illegal move was made.
	b.Ctx.WhiteTurn = !b.Ctx.WhiteTurn
	if b.IsKingInCheck() {
		return false
	}

	return true
}

func (b Board) AllLegalMoves() []Move {
	// if b.HashWithContext == 0 {
	// 	b.HashWithContext = b.HashBoardWithContext()
	// }

	// if moves, ok := allLegalMovesHashTable[b.HashWithContext]; ok {
	// 	return moves
	// }

	moves := b.GenerateAllMoves()
	moves = utils.Filter(moves, func(m Move) bool {
		b.MakePseudoLegalMove(m)
		isValid := b.IsValidPosition()
		b.UndoMove()
		return isValid
	})
	utils.ForEach(moves, func(m *Move) {
		m.isLegal = true
	})
	// allLegalMovesHashTable[b.HashWithContext] = moves
	return moves
}

func (b Board) GenerateAllMoves() []Move {
	var pb *PartialBoard
	if b.Ctx.WhiteTurn {
		pb = &b.White
	} else {
		pb = &b.Black
	}

	// "Normal" moves (including En passant and promotion)
	moves := pb.AllPossibleMoves(b)
	moves = append(moves, pb.AllCastlingMoves(b)...)

	return moves
}

func (b *Board) IsMated() bool {
	if b.Ctx.IsMatedCacheSet {
		return b.Ctx.IsMatedCache
	}

	// Setting cached value
	b.Ctx.IsMatedCacheSet = true

	isKingInCheck := b.IsKingInCheck()
	if !isKingInCheck {
		b.Ctx.IsMatedCache = false
		return false
	}

	possibleDefensiveMoves := b.AllLegalMoves()
	b.Ctx.IsMatedCache = len(possibleDefensiveMoves) == 0
	if b.Ctx.IsMatedCache {
		if b.Ctx.WhiteTurn {
			b.Ctx.Result = BlackWin
		} else {
			b.Ctx.Result = WhiteWin
		}
	}
	return b.Ctx.IsMatedCache
}

// IsKingInCheck returns true if the king is in check.
// It generates moves for the king as if it was every piece on the board.
// Diagonal, horizontal, vertical, and knight moves are checked.
func (b *Board) IsKingInCheck() bool {
	var kingPos uint64
	var ourPb, enemyPb PartialBoard
	var dirFnToCheckPawnAttack func(uint64, int) uint64
	if b.Ctx.WhiteTurn {
		kingPos = b.White.King.Board
		ourPb = b.White
		enemyPb = b.Black
		dirFnToCheckPawnAttack = GetDirectionFunc(directionUp)
	} else {
		kingPos = b.Black.King.Board
		ourPb = b.Black
		enemyPb = b.White
		dirFnToCheckPawnAttack = GetDirectionFunc(directionDown)
	}

	ourPbMask := ourPb.AllBoardMask()
	enemyPbMask := enemyPb.AllBoardMask()
	enemySlidersAndKingMask := enemyPb.Bishops.Board | enemyPb.Rooks.Board | enemyPb.Queens.Board | enemyPb.King.Board
	directions := []int{directionUp, directionDown, directionLeft, directionRight, directionUpLeft, directionUpRight, directionDownLeft, directionDownRight}
	for _, direction := range directions {
		dirFn := GetDirectionFunc(direction)
		for i := 1; i < 8; i++ {
			newPos := dirFn(kingPos, i)
			if newPos == 0 || newPos&ourPbMask != 0 {
				break
			}

			if newPos&enemyPbMask != 0 {
				if newPos&enemySlidersAndKingMask == 0 { // if its a pawn or knight should break
					break
				}

				kingCollision := newPos&enemyPb.King.Board != 0
				// Conditional rule of propositional logic (If A then B)
				// If the piece is a king, it can only attack the king if it's one rank ahead.
				// If King then i == 1 to the statement to be true.
				// See https://plato.stanford.edu/entries/logic-conditionals/
				// A -> B is equivalent to !A || B
				if kingCollision && i == 1 {
					return true
				}

				if newPos&enemyPb.Bishops.Board != 0 {
					if direction == directionUpLeft || direction == directionUpRight || direction == directionDownLeft || direction == directionDownRight {
						return true
					}
				}
				if newPos&enemyPb.Rooks.Board != 0 {
					if direction == directionUp || direction == directionDown || direction == directionLeft || direction == directionRight {
						return true
					}
				}

				if newPos&enemyPb.Queens.Board != 0 {
					return true
				}
			}
		}
	}
	// check for pawn attack
	pawn1 := dirFnToCheckPawnAttack(moveLeft(kingPos, 1), 1)
	pawn2 := dirFnToCheckPawnAttack(moveRight(kingPos, 1), 1)
	if pawn1 != 0 && pawn1&enemyPb.Pawns.Board != 0 {
		return true
	}
	if pawn2 != 0 && pawn2&enemyPb.Pawns.Board != 0 {
		return true
	}

	knightMoves := []func(uint64) uint64{
		moveL1, moveL2, moveL3, moveL4, moveL5, moveL6, moveL7, moveL8,
	}
	for _, knightMoveFn := range knightMoves {
		newPos := knightMoveFn(kingPos)
		if newPos == 0 {
			continue
		}

		if newPos&enemyPb.Knights.Board != 0 {
			return true
		}
	}

	return false
}

func (b *Board) IsDraw() bool {
	// Threesold repetition

	// Fifty moves rule
	if b.Ctx.HalfMoves >= 100 {
		b.Ctx.Result = Draw
		return true
	}

	// Stalemate
	if !b.IsKingInCheck() {
		allLegalMoves := b.AllLegalMoves()
		if len(allLegalMoves) == 0 {
			b.Ctx.Result = Draw
			return true
		}
	}

	return false
}

func (b Board) MaterialValueBalance() int64 {
	return int64(b.White.MaterialValue()) - int64(b.Black.MaterialValue())
}
