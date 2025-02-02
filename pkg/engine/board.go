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

func (b *Board) MakeMove(m Move) {
	var pb *PartialBoard
	if b.Ctx.WhiteToMove {
		pb = &b.White
	} else {
		pb = &b.Black
	}

	pb.MakeMove(m)
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
}

func (b Board) IsWhiteMated() bool {
	isKingInCheckOriginalPos := b.IsKingInCheck(b.White.King.Board, true)
	if !isKingInCheckOriginalPos {
		return false
	}

	kingMoves := b.White.King.AllPossibleMoves(b, true)
	if len(kingMoves) == 0 {
		return true
	}

	for _, move := range kingMoves {
		newBoard := b
		newBoard.MakeMove(move)
		if !newBoard.IsKingInCheck(newBoard.White.King.Board, true) {
			return false
		}
	}
	return true
}

func (b Board) IsBlackMated() bool {
	isKingInCheckOriginalPos := b.IsKingInCheck(b.Black.King.Board, false)
	if !isKingInCheckOriginalPos {
		return false
	}

	kingMoves := b.Black.King.AllPossibleMoves(b, false)
	if len(kingMoves) == 0 {
		return true
	}

	for _, move := range kingMoves {
		newBoard := b
		newBoard.MakeMove(move)
		if !newBoard.IsKingInCheck(newBoard.Black.King.Board, false) {
			return false
		}
	}
	return true
}

// IsKingInCheck returns true if the king is in check.
// kingPos is the position of the king.
// isWhite if the king is white.
func (b Board) IsKingInCheck(kingPos uint64, isWhite bool) bool {
	var pb *PartialBoard
	// Inverted
	if isWhite {
		pb = &b.Black
	} else {
		pb = &b.White
	}

	possibleCheckMoves := make([]Move, 0)
	// allWhitePossibleMovesExcludingKing = append(allWhitePossibleMovesExcludingKing, b.White.Pawns.AllPossibleMoves(b, false)...)
	possibleCheckMoves = append(possibleCheckMoves, pb.Knights.AllPossibleMoves(b, !isWhite)...)
	possibleCheckMoves = append(possibleCheckMoves, pb.Bishops.AllPossibleMoves(b, !isWhite)...)
	possibleCheckMoves = append(possibleCheckMoves, pb.Rooks.AllPossibleMoves(b, !isWhite)...)
	possibleCheckMoves = append(possibleCheckMoves, pb.Queens.AllPossibleMoves(b, !isWhite)...)

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
