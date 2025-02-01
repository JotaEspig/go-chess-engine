package engine

import (
	"strings"

	"github.com/charmbracelet/log"
)

// Board represents a full board with pieces of both colors on it.
type Board struct {
	White PartialBoard
	Black PartialBoard
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
	rows := strings.Split(fen, "/")
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

	return b
}

func (b Board) MaterialValueBalance() int64 {
	return b.White.MaterialValue() - b.Black.MaterialValue()
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
