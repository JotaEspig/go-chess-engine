package engine

import (
	"strings"

	"github.com/charmbracelet/log"
)

// Board represents a full board with pieces of both colors on it.
type Board struct {
	White PartialBoard
	Black PartialBoard
	Ctx   Context
}

func NewDefaultBoard() Board {
	return FenToBoard(DefaultStartFen)
}

func NewBoard() Board {
	return Board{
		White: NewPartialBoard(),
		Black: NewPartialBoard(),
	}
}

func FenToBoard(fen string) Board {
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

	return b
}

func FenToContext(splitted []string) Context {
	ctx := Context{}
	if len(splitted) > 0 {
		var whiteToMove bool
		if splitted[0] == "w" {
			whiteToMove = true
		} else {
			whiteToMove = false
		}
		ctx.WhiteToMove = whiteToMove
	} else {
	}
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
	if b.Ctx.WhiteToMove {
		// Check if the black king is in check while it's white's turn. Meaning that an illegal move was made.
		b.Ctx.WhiteToMove = false
		if b.IsKingInCheck() {
			return false
		}
	} else {
		// Check if the white king is in check while it's black's turn. Meaning that an illegal move was made.
		b.Ctx.WhiteToMove = true
		if b.IsKingInCheck() {
			return false
		}
	}

	return true
}

func (b Board) AllPossibleMoves() []Move {
	moves := make([]Move, 0)
	// "Normal" moves (including En passant and promotion)
	if b.Ctx.WhiteToMove {
		moves = append(moves, b.White.AllPossibleMoves(b)...)
	} else {
		moves = append(moves, b.Black.AllPossibleMoves(b)...)
	}

	// Castling

	return moves
}

func (b Board) AllMovesToDefendCheck() []Move {
	if !b.IsKingInCheck() {
		log.Fatal("King is not in check")
	}

	moves := make([]Move, 0)
	if b.Ctx.WhiteToMove {
		moves = append(moves, b.White.AllPossibleMoves(b)...)
	} else {
		moves = append(moves, b.Black.AllPossibleMoves(b)...)
	}

	// Filter out moves that don't defend the king
	moves = Filter(moves, func(m Move) bool {
		newBoard := b
		newBoard.MakeMove(m)
		return newBoard.IsValidPosition()
	})

	return moves
}

func (b *Board) MakeMove(m Move) {
	var pb *PartialBoard
	if b.Ctx.WhiteToMove {
		pb = &b.White
	} else {
		pb = &b.Black
	}

	if m.IsPromotion {
		pb.MakePromotion(m)
	} else {
		pb.MakeMove(m)
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
	if m.IsCapture {
		// Inverted color to erase the piece from the board
		if b.Ctx.WhiteToMove {
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

	b.Ctx.WhiteToMove = !b.Ctx.WhiteToMove
	b.Ctx.MoveNumber++
	b.Ctx.EnPassant = enPassantPos

	if !b.IsValidPosition() {
		log.Fatalf("Invalid position after move: %v", m)
	}
}

func (b Board) IsMated() bool {
	isKingInCheckOriginalPos := b.IsKingInCheck()
	if !isKingInCheckOriginalPos {
		return false
	}

	possibleDefensiveMoves := b.AllMovesToDefendCheck()
	return len(possibleDefensiveMoves) == 0
}

// IsKingInCheck returns true if the king is in check.
func (b Board) IsKingInCheck() bool {
	var kingPos uint64
	var enemyPb *PartialBoard

	isWhite := b.Ctx.WhiteToMove
	if isWhite {
		kingPos = b.White.King.Board
		enemyPb = &b.Black
	} else {
		kingPos = b.Black.King.Board
		enemyPb = &b.White
	}

	// Invert the color to get the enemy moves and see if it's possible to "capture" the king
	b.Ctx.WhiteToMove = !b.Ctx.WhiteToMove
	possibleCheckMoves := make([]Move, 0)
	possibleCheckMoves = append(possibleCheckMoves, enemyPb.Pawns.AllPossibleMoves(b)...)
	possibleCheckMoves = append(possibleCheckMoves, enemyPb.Knights.AllPossibleMoves(b)...)
	possibleCheckMoves = append(possibleCheckMoves, enemyPb.Bishops.AllPossibleMoves(b)...)
	possibleCheckMoves = append(possibleCheckMoves, enemyPb.Rooks.AllPossibleMoves(b)...)
	possibleCheckMoves = append(possibleCheckMoves, enemyPb.Queens.AllPossibleMoves(b)...)

	for _, move := range possibleCheckMoves {
		if move.NewPiecePos == kingPos {
			return true
		}
	}
	return false
}

func (b Board) MaterialValueBalance() int64 {
	return int64(b.White.MaterialValue()) - int64(b.Black.MaterialValue())
}

func (b Board) ToVisualBoard() VisualBoard {
	vb := VisualBoard{}

	// White pieces
	for _, pos := range Int64toPositions(b.White.Pawns.Board) {
		col, row := pos[0], pos[1]
		vb.Board[row][col] = VisualPiece{IsWhite: true, Type: PawnType}
	}
	for _, pos := range Int64toPositions(b.White.Knights.Board) {
		col, row := pos[0], pos[1]
		vb.Board[row][col] = VisualPiece{IsWhite: true, Type: KnightType}
	}
	for _, pos := range Int64toPositions(b.White.Bishops.Board) {
		col, row := pos[0], pos[1]
		vb.Board[row][col] = VisualPiece{IsWhite: true, Type: BishopType}
	}
	for _, pos := range Int64toPositions(b.White.Rooks.Board) {
		col, row := pos[0], pos[1]
		vb.Board[row][col] = VisualPiece{IsWhite: true, Type: RookType}
	}
	for _, pos := range Int64toPositions(b.White.Queens.Board) {
		col, row := pos[0], pos[1]
		vb.Board[row][col] = VisualPiece{IsWhite: true, Type: QueenType}
	}
	for _, pos := range Int64toPositions(b.White.King.Board) {
		col, row := pos[0], pos[1]
		vb.Board[row][col] = VisualPiece{IsWhite: true, Type: KingType}
	}

	// Black pieces
	for _, pos := range Int64toPositions(b.Black.Pawns.Board) {
		col, row := pos[0], pos[1]
		vb.Board[row][col] = VisualPiece{IsWhite: false, Type: PawnType}
	}
	for _, pos := range Int64toPositions(b.Black.Knights.Board) {
		col, row := pos[0], pos[1]
		vb.Board[row][col] = VisualPiece{IsWhite: false, Type: KnightType}
	}
	for _, pos := range Int64toPositions(b.Black.Bishops.Board) {
		col, row := pos[0], pos[1]
		vb.Board[row][col] = VisualPiece{IsWhite: false, Type: BishopType}
	}
	for _, pos := range Int64toPositions(b.Black.Rooks.Board) {
		col, row := pos[0], pos[1]
		vb.Board[row][col] = VisualPiece{IsWhite: false, Type: RookType}
	}
	for _, pos := range Int64toPositions(b.Black.Queens.Board) {
		col, row := pos[0], pos[1]
		vb.Board[row][col] = VisualPiece{IsWhite: false, Type: QueenType}
	}
	for _, pos := range Int64toPositions(b.Black.King.Board) {
		col, row := pos[0], pos[1]
		vb.Board[row][col] = VisualPiece{IsWhite: false, Type: KingType}
	}

	return vb
}
