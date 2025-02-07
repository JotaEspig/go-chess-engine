package chess

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

func (pb PartialBoard) AllPossibleMoves(board Board) []*Move {
	moves := pb.Pawns.AllPossibleMoves(board)
	moves = append(moves, pb.Knights.AllPossibleMoves(board)...)
	moves = append(moves, pb.Bishops.AllPossibleMoves(board)...)
	moves = append(moves, pb.Rooks.AllPossibleMoves(board)...)
	moves = append(moves, pb.Queens.AllPossibleMoves(board)...)
	moves = append(moves, pb.King.AllPossibleMoves(board)...)
	return moves
}

func (pb *PartialBoard) MakeMove(m *Move) {
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

func (pb *PartialBoard) MakePromotion(m *Move) {
	var pp, pp2 *PiecesPosition
	switch m.PieceType {
	case PawnType:
		pp = &pb.Pawns
	default:
		log.Fatalf("Invalid piece type for promotion: %v", m.PieceType)
	}
	switch m.NewPieceType {
	case QueenType:
		pp2 = &pb.Queens
	case RookType:
		pp2 = &pb.Rooks
	case BishopType:
		pp2 = &pb.Bishops
	case KnightType:
		pp2 = &pb.Knights
	default:
		log.Fatalf("Invalid piece type for promotion target: %v", m.NewPieceType)
	}

	pp.Board &= ^m.OldPiecePos
	pp2.Board |= m.NewPiecePos
}

func (pb PartialBoard) AllBoardMask() uint64 {
	return pb.Pawns.Board | pb.Knights.Board | pb.Bishops.Board | pb.Rooks.Board | pb.Queens.Board | pb.King.Board
}

// MaterialValue returns the total value of all the pieces on the board.
func (pb PartialBoard) MaterialValue() uint64 {
	return pb.Pawns.Value() + pb.Knights.Value() + pb.Bishops.Value() + pb.Rooks.Value() + pb.Queens.Value()
}
