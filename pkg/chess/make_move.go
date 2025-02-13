package chess

import "github.com/charmbracelet/log"

func (b *Board) MakeLegalMove(m Move) {
	if !m.isLegal {
		log.Fatalf("Trying to make an illegal move: %v", m)
	}
	b.MakePseudoLegalMove(m)
}

// MakeMove makes a move on the board.
func (b *Board) MakePseudoLegalMove(m Move) {
	var pb *PartialBoard
	if b.Ctx.WhiteTurn {
		pb = &b.White
	} else {
		pb = &b.Black
	}

	// Update MovesDone and PreviousCtx
	b.MovesDone = append(b.MovesDone, m)
	b.PreviousCtx = append(b.PreviousCtx, b.Ctx)

	if m.IsCastling { // Castling verifications
		isKingSide := m.NewPiecePos < m.OldPiecePos
		// Move rook, king is moved on normal MakeMove
		if isKingSide {
			if b.Ctx.WhiteTurn {
				b.White.Rooks.Board &= ^H1
				b.White.Rooks.Board |= F1
				b.Ctx.WhiteCastlingKingSide = false
				b.Ctx.WhiteCastlingQueenSide = false
			} else {
				b.Black.Rooks.Board &= ^H8
				b.Black.Rooks.Board |= F8
				b.Ctx.BlackCastlingKingSide = false
				b.Ctx.BlackCastlingQueenSide = false
			}
		} else {
			if b.Ctx.WhiteTurn {
				b.White.Rooks.Board &= ^A1
				b.White.Rooks.Board |= D1
				b.Ctx.WhiteCastlingKingSide = false
				b.Ctx.WhiteCastlingQueenSide = false
			} else {
				b.Black.Rooks.Board &= ^A8
				b.Black.Rooks.Board |= D8
				b.Ctx.BlackCastlingKingSide = false
				b.Ctx.BlackCastlingQueenSide = false
			}
		}
	}
	// Removes castling rights if the piece is a rook or king
	if m.PieceType == RookType || m.PieceType == KingType {
		if b.Ctx.WhiteTurn {
			if m.OldPiecePos == A1 {
				b.Ctx.WhiteCastlingQueenSide = false
			} else if m.OldPiecePos == H1 {
				b.Ctx.WhiteCastlingKingSide = false
			} else {
				b.Ctx.WhiteCastlingKingSide = false
				b.Ctx.WhiteCastlingQueenSide = false
			}
		} else {
			if m.OldPiecePos == A8 {
				b.Ctx.BlackCastlingQueenSide = false
			} else if m.OldPiecePos == H8 {
				b.Ctx.BlackCastlingKingSide = false
			} else {
				b.Ctx.BlackCastlingKingSide = false
				b.Ctx.BlackCastlingQueenSide = false
			}
		}
	}
	if m.IsPromotion {
		pb.MakePromotion(m)
	} else {
		pb.MakeMove(m)
	}
	// Removes enemy piece if it's a capture
	if m.IsCapture {
		// Inverted color to erase the piece from the board
		if b.Ctx.WhiteTurn {
			pb = &b.Black
		} else {
			pb = &b.White
		}

		colrow := Int64toPositions(m.NewPiecePos)
		if len(colrow) != 1 {
			log.Fatal("Invalid NewPiecePos: %v", m.NewPiecePos)
		}

		switch m.CapturedPieceType {
		case PawnType:
			enemyPawnPos := m.NewPiecePos
			if m.IsEnPassant {
				if b.Ctx.WhiteTurn {
					enemyPawnPos = moveDown(enemyPawnPos, 1) // Black pawn
				} else {
					enemyPawnPos = moveUp(enemyPawnPos, 1) // White pawn
				}
			}
			pb.Pawns.Board &= ^enemyPawnPos
		case KnightType:
			pb.Knights.Board &= ^m.NewPiecePos
		case BishopType:
			pb.Bishops.Board &= ^m.NewPiecePos
		case RookType:
			pb.Rooks.Board &= ^m.NewPiecePos
		case QueenType:
			pb.Queens.Board &= ^m.NewPiecePos
		default:
			log.Fatalf("Invalid piece type: %v", m.CapturedPieceType)
		}
	}

	// Check for next move En passant
	// Default value is 0
	var enPassantPos uint64 = 0
	// If is 2 square pawn move, set the en passant position for the one row behind the pawn
	if m.Is2SquarePawnMove() {
		isWhite := m.NewPiecePos > m.OldPiecePos // Assumes it's a pawn move
		if isWhite {
			enPassantPos = moveUp(m.OldPiecePos, 1)
		} else {
			enPassantPos = moveDown(m.OldPiecePos, 1)
		}
	}

	// means black completed its turn
	if !b.Ctx.WhiteTurn {
		b.Ctx.MoveNumber++
	}
	// Check for half moves
	if !m.IsCapture && m.PieceType != PawnType {
		b.Ctx.HalfMoves++
	} else {
		b.Ctx.HalfMoves = 0
	}
	b.Ctx.WhiteTurn = !b.Ctx.WhiteTurn
	b.Ctx.EnPassant = enPassantPos

	// Reset cached values to false, because it's a new position
	b.Ctx.IsKingInCheckCacheSet = false
	b.Ctx.IsMatedCacheSet = false
	b.Ctx.IsDrawCacheSet = false
}

func (b *Board) UndoMove() {
	// Get the last move and context
	lastMove := b.MovesDone[len(b.MovesDone)-1]
	b.MovesDone = b.MovesDone[:len(b.MovesDone)-1]       // Pop
	b.Ctx = b.PreviousCtx[len(b.PreviousCtx)-1]          // Restores the context
	b.PreviousCtx = b.PreviousCtx[:len(b.PreviousCtx)-1] // Pop

	var ourPb, enemyPb *PartialBoard
	if b.Ctx.WhiteTurn {
		ourPb = &b.White
		enemyPb = &b.Black
	} else {
		ourPb = &b.Black
		enemyPb = &b.White
	}

	// Restore the piece to the previous position
	switch lastMove.PieceType {
	case PawnType:
		ourPb.Pawns.Board &= ^lastMove.NewPiecePos
		ourPb.Pawns.Board |= lastMove.OldPiecePos
	case KnightType:
		ourPb.Knights.Board &= ^lastMove.NewPiecePos
		ourPb.Knights.Board |= lastMove.OldPiecePos
	case BishopType:
		ourPb.Bishops.Board &= ^lastMove.NewPiecePos
		ourPb.Bishops.Board |= lastMove.OldPiecePos
	case RookType:
		ourPb.Rooks.Board &= ^lastMove.NewPiecePos
		ourPb.Rooks.Board |= lastMove.OldPiecePos
	case QueenType:
		ourPb.Queens.Board &= ^lastMove.NewPiecePos
		ourPb.Queens.Board |= lastMove.OldPiecePos
	case KingType:
		ourPb.King.Board &= ^lastMove.NewPiecePos
		ourPb.King.Board |= lastMove.OldPiecePos
	default:
		log.Fatalf("Invalid piece type: %v", lastMove.PieceType)
	}

	// Restore the captured piece if it's a capture
	if lastMove.IsCapture {
		switch lastMove.CapturedPieceType {
		case PawnType:
			pos := lastMove.NewPiecePos
			if lastMove.IsEnPassant {
				if b.Ctx.WhiteTurn {
					pos = moveUp(pos, 1) // White pawn
				} else {
					pos = moveDown(pos, 1) // Black pawn
				}
			}
			enemyPb.Pawns.Board |= pos
		case KnightType:
			enemyPb.Knights.Board |= lastMove.NewPiecePos
		case BishopType:
			enemyPb.Bishops.Board |= lastMove.NewPiecePos
		case RookType:
			enemyPb.Rooks.Board |= lastMove.NewPiecePos
		case QueenType:
			enemyPb.Queens.Board |= lastMove.NewPiecePos
		default:
			log.Fatalf("Invalid piece type: %v", lastMove.CapturedPieceType)
		}
	}

	// Restore the castling if it's a castling
	if lastMove.IsCastling {
		if lastMove.NewPiecePos < lastMove.OldPiecePos { // King side
			if b.Ctx.WhiteTurn {
				b.White.Rooks.Board &= ^uint64(4) // F1
				b.White.Rooks.Board |= uint64(1)  // H1
			} else {
				b.Black.Rooks.Board &= ^uint64(288_230_376_151_711_744) // F8
				b.Black.Rooks.Board |= uint64(72_057_594_037_927_936)   // H8
			}
		} else { // Queen side
			if b.Ctx.WhiteTurn {
				b.White.Rooks.Board &= ^uint64(16) // D1
				b.White.Rooks.Board |= uint64(128) // A1
			} else {
				b.Black.Rooks.Board &= ^uint64(1_152_921_504_606_846_976) // D8
				b.Black.Rooks.Board |= uint64(9_223_372_036_854_775_808)  // A8
			}
		}
	}
}
