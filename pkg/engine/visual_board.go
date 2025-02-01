package engine

import "strings"

type VisualPiece struct {
	IsWhite bool
	Type    PieceType
}

// Represents a visual board with pieces of both colors on it.
type VisualBoard struct {
	Board [8][8]VisualPiece
}

func (vb VisualBoard) String() string {
	var sb strings.Builder
	sb.WriteRune('-')
	sb.WriteString(strings.Repeat("----", 8) + "\n")
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			str := vb.Board[i][j].Type.String()
			if !vb.Board[i][j].IsWhite {
				str = strings.ToUpper(str)
			}
			sb.WriteString("| " + str + " ")
		}
		sb.WriteString("|\n-")
		sb.WriteString(strings.Repeat("----", 8) + "\n")
	}
	return sb.String()
}
