package engine

import "github.com/charmbracelet/log"

// PartialBoard represents a board with only pieces of the same color on it.
type PartialBoard struct {
	Pawns   PiecesPosition
	Knights PiecesPosition
	Bishops PiecesPosition
	Rooks   PiecesPosition
	Queens  PiecesPosition
	King    PiecesPosition
}

func NewPartialBoard() PartialBoard {
	return PartialBoard{
		Pawns:   PiecesPosition{Type: PawnType},
		Knights: PiecesPosition{Type: KnightType},
		Bishops: PiecesPosition{Type: BishopType},
		Rooks:   PiecesPosition{Type: RookType},
		Queens:  PiecesPosition{Type: QueenType},
		King:    PiecesPosition{Type: KingType},
	}
}

func (pb *PartialBoard) MakeMove(m Move) {
	var pp *PiecesPosition
	switch m.PieceType {
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
		log.Fatalf("Invalid piece type: %v", m.PieceType)
	}

	pp.Board &= ^m.OldPiecePos
	pp.Board |= m.NewPiecePos
}

func (pb PartialBoard) AllBoardMask() uint64 {
	return pb.Pawns.Board | pb.Knights.Board | pb.Bishops.Board | pb.Rooks.Board | pb.Queens.Board | pb.King.Board
}

// MaterialValue returns the total value of all the pieces on the board.
func (pb PartialBoard) MaterialValue() uint64 {
	return pb.Pawns.Value() + pb.Knights.Value() + pb.Bishops.Value() + pb.Rooks.Value() + pb.Queens.Value()
}
