package chess

import "strings"

type VisualPiece struct {
	IsWhite bool
	Type    PieceType
}

// Represents a visual board with pieces of both colors on it.
type VisualBoard struct {
	Board [8][8]VisualPiece
}

func (b Board) VisualBoard() VisualBoard {
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

func (vb VisualBoard) String() string {
	var sb strings.Builder
	sb.WriteRune('-')
	sb.WriteString(strings.Repeat("----", 8) + "\n")
	for i := 7; i >= 0; i-- { // Row are inverted to print the board correctly.
		for j := 0; j < 8; j++ {
			str := vb.Board[i][j].Type.String()
			if vb.Board[i][j].IsWhite {
				str = strings.ToUpper(str)
			}
			sb.WriteString("| " + str + " ")
		}
		sb.WriteString("|\n-")
		sb.WriteString(strings.Repeat("----", 8) + "\n")
	}
	return sb.String()
}
