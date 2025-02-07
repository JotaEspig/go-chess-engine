package chess

import (
	"gce/pkg/utils"
	"strconv"
	"strings"

	"github.com/charmbracelet/log"
)

// Board represents a full board with pieces of both colors on it.
type Board struct {
	White PartialBoard
	Black PartialBoard
	Ctx   Context

	// These 2 fields are used to keep track of the previous board state (i.e. every move made on board)
	MoveDone  Move   // MoveDone is the move that was done in the current board
	PrevBoard *Board // PrevBoard is the previous board state
}

func NewDefaultBoard() *Board {
	return FenToBoard(DefaultStartFen)
}

func NewBoard() *Board {
	return &Board{
		White: NewPartialBoard(),
		Black: NewPartialBoard(),
	}
}

func FenToBoard(fen string) *Board {
	b := NewBoard()
	splitted := strings.Split(fen, " ")
	pos := splitted[0]
	rows := strings.Split(pos, "/")
	row := 7
	for _, rowString := range rows {
		col := 0
		for _, char := range rowString {
			if char >= '1' && char <= '8' {
				col += int(char - '0')
				continue
			}

			isWhite := char >= 'A' && char <= 'Z'
			pieceType := PieceTypeFromChar(char)
			var pb *PartialBoard
			if isWhite {
				pb = &b.White
			} else {
				pb = &b.Black
			}

			var pp *PiecesPosition
			switch pieceType {
			case PawnType:
				pp = &pb.Pawns
			case KnightType:
				pp = &pb.Knights
			case BishopType:
				pp = &pb.Bishops
			case RookType:
				pp = &pb.Rooks
			case QueenType:
				pp = &pb.Queens
			case KingType:
				pp = &pb.King
			default:
				log.Fatalf("Invalid piece type: %v. Got from char: %c", pieceType, char)
			}

			pp.SetPieceAt(col, row)
			col++
		}

		row--
	}

	if len(splitted) > 1 {
		b.Ctx = FenToContext(splitted[1:])
	}

	// Setup threefold repetition hash
	boardHash := BoardHash(b.Hash())
	b.Ctx.ThreesoldRepetition[boardHash] = 1
	return b
}

func FenToContext(splitted []string) Context {
	if len(splitted) != 5 {
		log.Fatal("Invalid fen context")
	}

	ctx := Context{}
	if splitted[0] == "w" {
		ctx.WhiteTurn = true
	} else {
		ctx.WhiteTurn = false
	}

	castling := splitted[1]
	if strings.Contains(castling, "K") {
		ctx.WhiteCastlingKingSide = true
	}
	if strings.Contains(castling, "Q") {
		ctx.WhiteCastlingQueenSide = true
	}
	if strings.Contains(castling, "k") {
		ctx.BlackCastlingKingSide = true
	}
	if strings.Contains(castling, "q") {
		ctx.BlackCastlingQueenSide = true
	}

	enPassantCoord := splitted[2]
	if enPassantCoord != "-" {
		col := int(enPassantCoord[0] - 'a')
		row := int(enPassantCoord[1] - '1')
		ctx.EnPassant = PositionToUInt64(col, row)
	}

	halfMoveClock := splitted[3]
	HalfMovesInt, err := strconv.Atoi(halfMoveClock)
	if err != nil {
		log.Fatalf("Invalid half move clock: %v", halfMoveClock)
	}
	moveNumber := splitted[4]
	moveNumberInt, err := strconv.Atoi(moveNumber)
	if err != nil {
		log.Fatalf("Invalid move number: %v", moveNumber)
	}

	ctx.MoveNumber = uint(moveNumberInt)
	ctx.HalfMoves = uint(HalfMovesInt)
	ctx.ThreesoldRepetition = make(ThreefoldRepetitionHashTable)
	return ctx
}

func (b Board) IsValidPosition() bool {
	// Check if there are no pieces in the same position
	whitePieces := b.White.AllBoardMask()
	blackPieces := b.Black.AllBoardMask()
	if whitePieces&blackPieces != 0 {
		return false
	}

	// originalWhiteToMove := b.Ctx.WhiteToMove
	if b.Ctx.WhiteTurn {
		// Check if the black king is in check while it's white's turn. Meaning that an illegal move was made.
		b.Ctx.WhiteTurn = false
		if b.IsKingInCheck() {
			return false
		}
	} else {
		// Check if the white king is in check while it's black's turn. Meaning that an illegal move was made.
		b.Ctx.WhiteTurn = true
		if b.IsKingInCheck() {
			return false
		}
	}

	return true
}

func (b Board) AllLegalMoves() []Move {
	hashForMoves := b.HashForAllMoves()
	if cachedMoves, ok := allLegalBoardMovesHashTable[hashForMoves]; ok {
		return cachedMoves
	}

	moves := b.AllPossibleMoves()
	// Filter out moves that are not legal
	moves = utils.Filter(moves, func(m Move) bool {
		newBoard := b.Copy()
		return newBoard.MakeMove(m)
	})
	// Set IsCheck, IsCheckFieldSet and CapturedPieceType
	utils.ForEach(moves, func(m *Move) {
		var enemy PartialBoard
		if b.Ctx.WhiteTurn {
			enemy = b.Black
		} else {
			enemy = b.White
		}

		// Set Capture Piece type
		capturedPieceType := InvalidType
		if m.IsCapture {
			if m.NewPiecePos&enemy.Pawns.Board != 0 || m.NewPiecePos&b.Ctx.EnPassant != 0 {
				capturedPieceType = PawnType
			} else if m.NewPiecePos&enemy.Knights.Board != 0 {
				capturedPieceType = KnightType
			} else if m.NewPiecePos&enemy.Bishops.Board != 0 {
				capturedPieceType = BishopType
			} else if m.NewPiecePos&enemy.Rooks.Board != 0 {
				capturedPieceType = RookType
			} else if m.NewPiecePos&enemy.Queens.Board != 0 {
				capturedPieceType = QueenType
			}
			m.CapturedPieceType = capturedPieceType
		}

		// Set IsCheck
		b.MakeMove(*m)
		if b.IsKingInCheck() {
			m.IsCheck = true
		}
		b = *b.PrevBoard
	})

	allLegalBoardMovesHashTable[hashForMoves] = moves
	return moves
}

func (b Board) AllPossibleMoves() []Move {
	var pb *PartialBoard
	if b.Ctx.WhiteTurn {
		pb = &b.White
	} else {
		pb = &b.Black
	}

	// "Normal" moves (including En passant and promotion)
	moves := pb.AllPossibleMoves(b)
	// Castling
	moves = append(moves, b.AllCastlingMoves()...)

	return moves
}

func (b Board) AllCastlingMoves() []Move {
	var pb PartialBoard
	anyCastleAvailable := false
	var kingSideSpaceMask uint64
	var QueenSideSpaceMask uint64
	var kingSideSafeSpot uint64
	var queenSideSafeSpot uint64
	if b.Ctx.WhiteTurn {
		pb = b.White
		anyCastleAvailable = b.Ctx.WhiteCastlingKingSide || b.Ctx.WhiteCastlingQueenSide
		kingSideSpaceMask = ^uint64(6)    // 6 is the bits that represents F1 and G1
		QueenSideSpaceMask = ^uint64(112) // 112 is the bits that represents B1, C1 and D1
		kingSideSafeSpot = uint64(2)      // G1
		queenSideSafeSpot = uint64(32)    // C1
	} else {
		pb = b.Black
		anyCastleAvailable = b.Ctx.BlackCastlingKingSide || b.Ctx.BlackCastlingQueenSide
		kingSideSpaceMask = ^uint64(432_345_564_227_567_616)    // 432_345_564_227_567_616 is the bits that represents F8 and G8
		QueenSideSpaceMask = ^uint64(8_070_450_532_247_928_832) // 8_070_450_532_247_928_832 is the bits that represents B8, C8, D8
		kingSideSafeSpot = uint64(144_115_188_075_855_872)      // G8
		queenSideSafeSpot = uint64(2_305_843_009_213_693_952)   // C8
	}

	if !anyCastleAvailable {
		return []Move{}
	}

	moves := make([]Move, 0)
	allBoardMask := pb.AllBoardMask()
	// king side is empty, can castle
	if kingSideSpaceMask&allBoardMask == 0 {
		move := Move{OldPiecePos: pb.King.Board, NewPiecePos: kingSideSafeSpot, IsCastling: true, PieceType: KingType}
		moves = append(moves, move)
	}
	// queen side is empty, can castle
	if QueenSideSpaceMask&allBoardMask == 0 {
		move := Move{OldPiecePos: pb.King.Board, NewPiecePos: queenSideSafeSpot, IsCastling: true, PieceType: KingType}
		moves = append(moves, move)
	}
	return moves
}

// MakeMove makes a move on the board.
// This function does not treat draw positions by threesold repetition, You should check if before hand.
// Returns true if it's a valid move, false otherwise.
func (b *Board) MakeMove(m Move) bool {
	prevBoard := &Board{}
	*prevBoard = b.Copy()

	var pb *PartialBoard
	if b.Ctx.WhiteTurn {
		pb = &b.White
	} else {
		pb = &b.Black
	}

	var mask uint64
	switch m.PieceType {
	case PawnType:
		mask = pb.Pawns.Board
	case KnightType:
		mask = pb.Knights.Board
	case BishopType:
		mask = pb.Bishops.Board
	case RookType:
		mask = pb.Rooks.Board
	case QueenType:
		mask = pb.Queens.Board
	case KingType:
		mask = pb.King.Board
	default:
		log.Fatalf("Invalid piece type: %v", m.PieceType)
	}

	// Check if the piece is at the position
	if mask&m.OldPiecePos == 0 {
		return false
	}

	if m.IsPromotion {
		pb.MakePromotion(m)
	} else if m.IsCastling { // Castling verifications
		// Cannot castle if the king is in check
		if b.IsKingInCheck() {
			// Restore to the previous board
			*b = *prevBoard
			return false
		}

		isKingSide := m.NewPiecePos < m.OldPiecePos
		copyBoard := *b
		copyMove := m
		copyMove.IsCastling = false
		// King side
		if isKingSide {
			copyMove.NewPiecePos <<= 1
		} else { // Queen side
			copyMove.NewPiecePos >>= 1
		}

		// Means that king moves in a square attacked by enemy piece
		if !copyBoard.MakeMove(copyMove) {
			// Restore to the previous board
			*b = *prevBoard
			return false
		}

		// Check if it is allowed castling giving the context of the board
		if b.Ctx.WhiteTurn {
			if isKingSide {
				if !b.Ctx.WhiteCastlingKingSide {
					// Restore to the previous board
					*b = *prevBoard
					return false
				}
			} else {
				if !b.Ctx.WhiteCastlingQueenSide {
					// Restore to the previous board
					*b = *prevBoard
					return false
				}
			}
		} else {
			if isKingSide {
				if !b.Ctx.BlackCastlingKingSide {
					// Restore to the previous board
					*b = *prevBoard
					return false
				}
			} else {
				if !b.Ctx.BlackCastlingQueenSide {
					// Restore to the previous board
					*b = *prevBoard
					return false
				}
			}
		}

		// Move rook
		if isKingSide {
			if b.Ctx.WhiteTurn {
				b.White.Rooks.Board &= ^uint64(1) // H1
				b.White.Rooks.Board |= uint64(4)  // F1
			} else {
				b.Black.Rooks.Board &= ^uint64(72_057_594_037_927_936) // H8
				b.Black.Rooks.Board |= uint64(288_230_376_151_711_744) // F8
			}
		} else {
			if b.Ctx.WhiteTurn {
				b.White.Rooks.Board &= ^uint64(128) // A1
				b.White.Rooks.Board |= uint64(16)   // D1
			} else {
				b.Black.Rooks.Board &= ^uint64(9_223_372_036_854_775_808) // A8
				b.Black.Rooks.Board |= uint64(1_152_921_504_606_846_976)  // D8
			}
		}
		// Move king
		pb.MakeMove(m)
	} else {
		// Then, normal move
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

		// Clear the piece at the position at every piece type board
		// This is necessary because the piece could be at any of the piece type boards
		// and we don't know which one it is.
		pb.Pawns.Board &= ^m.NewPiecePos
		pb.Knights.Board &= ^m.NewPiecePos
		pb.Bishops.Board &= ^m.NewPiecePos
		pb.Rooks.Board &= ^m.NewPiecePos
		pb.Queens.Board &= ^m.NewPiecePos
		pb.King.Board &= ^m.NewPiecePos
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

	// Add new position to Threefold repetition hash table
	boardHash := b.Hash()
	if n, ok := b.Ctx.ThreesoldRepetition[boardHash]; ok {
		b.Ctx.ThreesoldRepetition[boardHash] = n + 1
	} else {
		b.Ctx.ThreesoldRepetition[boardHash] = 1
	}

	// Reset cached values to false, because it's a new position
	b.Ctx.IsKingInCheckCacheSet = false
	b.Ctx.IsMatedCacheSet = false
	b.Ctx.IsDrawCacheSet = false

	// Check if it's a valid position
	// Should be the last thing to do (besides setting the previous board)
	if !b.IsValidPosition() {
		// Restore to the previous board
		*b = *prevBoard
		return false
	}
	prevBoard.MoveDone = m
	b.PrevBoard = prevBoard
	return true
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
func (b *Board) IsKingInCheck() bool {
	if b.Ctx.IsKingInCheckCacheSet {
		return b.Ctx.IsKingInCheckCache
	}

	// Setting cached value
	b.Ctx.IsKingInCheckCacheSet = true

	var kingPos uint64
	isWhite := b.Ctx.WhiteTurn
	if isWhite {
		kingPos = b.White.King.Board
	} else {
		kingPos = b.Black.King.Board
	}

	// Invert the color to get the enemy moves and see if it's possible to "capture" the king
	b.Ctx.WhiteTurn = !b.Ctx.WhiteTurn
	kingCaptureMoves := b.AllPossibleMoves()
	kingCaptureMoves = utils.Filter(kingCaptureMoves, func(m Move) bool {
		return m.PieceType != KingType && m.NewPiecePos&kingPos != 0
	})
	b.Ctx.WhiteTurn = !b.Ctx.WhiteTurn // Restore to the original value

	if len(kingCaptureMoves) > 0 {
		b.Ctx.IsKingInCheckCache = true
		return true
	} else {
		b.Ctx.IsKingInCheckCache = false
		return false
	}
}

func (b *Board) IsDraw() bool {
	if b.Ctx.IsDrawCacheSet {
		return b.Ctx.IsDrawCache
	}

	// Setting cached value
	b.Ctx.IsDrawCacheSet = true

	// 50 moves rule
	if b.Ctx.HalfMoves >= 100 {
		b.Ctx.IsDrawCache = true
		b.Ctx.Result = Draw
		return true
	}

	// Stalemate
	if !b.IsKingInCheck() {
		possibleMoves := b.AllLegalMoves()
		if len(possibleMoves) == 0 {
			b.Ctx.IsDrawCache = true
			b.Ctx.Result = Draw
			return true
		}
	}

	// Threefold repetition
	boardHash := BoardHash(b.Hash())
	if n, ok := b.Ctx.ThreesoldRepetition[boardHash]; ok {
		if n >= 3 {
			b.Ctx.IsDrawCache = true
			b.Ctx.Result = Draw
			return true
		}
	}
	b.Ctx.IsDrawCache = false
	return false
}

func (b Board) MaterialValueBalance() int64 {
	return int64(b.White.MaterialValue()) - int64(b.Black.MaterialValue())
}

func (b Board) Copy() Board {
	return Board{
		White:     b.White,
		Black:     b.Black,
		Ctx:       b.Ctx.Copy(),
		MoveDone:  b.MoveDone,
		PrevBoard: b.PrevBoard,
	}
}
