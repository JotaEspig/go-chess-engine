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

func (pb PartialBoard) AllPossibleMoves(board Board) []Move {
	moves := pb.Pawns.AllPossibleMoves(board)
	moves = append(moves, pb.Knights.AllPossibleMoves(board)...)
	moves = append(moves, pb.Bishops.AllPossibleMoves(board)...)
	moves = append(moves, pb.Rooks.AllPossibleMoves(board)...)
	moves = append(moves, pb.Queens.AllPossibleMoves(board)...)
	moves = append(moves, pb.King.AllPossibleMoves(board)...)
	return moves
}

func (pb PartialBoard) AllCastlingMoves(board Board) []Move {
	anyCastleAvailable := false
	var kingSideSpaceMask uint64
	var QueenSideSpaceMask uint64
	var kingSideSafeSpot uint64
	var queenSideSafeSpot uint64
	if board.Ctx.WhiteTurn {
		anyCastleAvailable = board.Ctx.WhiteCastlingKingSide || board.Ctx.WhiteCastlingQueenSide
		kingSideSpaceMask = uint64(6)    // 6 is the bits that represents F1 and G1
		QueenSideSpaceMask = uint64(112) // 112 is the bits that represents B1, C1 and D1
		kingSideSafeSpot = uint64(2)     // G1
		queenSideSafeSpot = uint64(32)   // C1
	} else {
		anyCastleAvailable = board.Ctx.BlackCastlingKingSide || board.Ctx.BlackCastlingQueenSide
		kingSideSpaceMask = uint64(432_345_564_227_567_616)    // 432_345_564_227_567_616 is the bits that represents F8 and G8
		QueenSideSpaceMask = uint64(8_070_450_532_247_928_832) // 8_070_450_532_247_928_832 is the bits that represents B8, C8, D8
		kingSideSafeSpot = uint64(144_115_188_075_855_872)     // G8
		queenSideSafeSpot = uint64(2_305_843_009_213_693_952)  // C8
	}

	if !anyCastleAvailable {
		return []Move{}
	}

	moves := make([]Move, 0, 2)
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

func (pb *PartialBoard) MakePromotion(m Move) {
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

func (pb PartialBoard) GetPieceTypeByPos(pos uint64) PieceType {
	if pb.Pawns.Board&pos != 0 {
		return PawnType
	}
	if pb.Knights.Board&pos != 0 {
		return KnightType
	}
	if pb.Bishops.Board&pos != 0 {
		return BishopType
	}
	if pb.Rooks.Board&pos != 0 {
		return RookType
	}
	if pb.Queens.Board&pos != 0 {
		return QueenType
	}
	if pb.King.Board&pos != 0 {
		return KingType
	}
	return InvalidType
}

// MaterialValue returns the total value of all the pieces on the board.
func (pb PartialBoard) MaterialValue() uint64 {
	return pb.Pawns.Value() + pb.Knights.Value() + pb.Bishops.Value() + pb.Rooks.Value() + pb.Queens.Value()
}
