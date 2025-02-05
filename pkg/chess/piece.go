package chess

// Enum for the different types of pieces.
const (
	InvalidType PieceType = iota
	PawnType
	KnightType
	BishopType
	RookType
	QueenType
	KingType
) // Enum for the different values for each pieces.
// The int values represent the value of the piece
const (
	InvalidValue PieceValue = 0
	PawnValue    PieceValue = 1
	KnightValue  PieceValue = 3
	BishopValue  PieceValue = 3
	RookValue    PieceValue = 5
	QueenValue   PieceValue = 9
	KingValue    PieceValue = 9999999
)

type PieceType int
type PieceValue int

func runeToLower(char rune) rune {
	if char >= 'A' && char <= 'Z' {
		return char + 32
	}
	return char
}

func PieceTypeFromChar(char rune) PieceType {
	char = runeToLower(char)
	switch char {
	case 'p':
		return PawnType
	case 'n':
		return KnightType
	case 'b':
		return BishopType
	case 'r':
		return RookType
	case 'q':
		return QueenType
	case 'k':
		return KingType
	default:
		return InvalidType
	}
}

func (p PieceType) String() string {
	switch p {
	case PawnType:
		return "p"
	case KnightType:
		return "n"
	case BishopType:
		return "b"
	case RookType:
		return "r"
	case QueenType:
		return "q"
	case KingType:
		return "k"
	default:
		return " "
	}
}

func (p PieceType) Value() uint64 {
	switch p {
	case PawnType:
		return uint64(PawnValue)
	case KnightType:
		return uint64(KnightValue)
	case BishopType:
		return uint64(BishopValue)
	case RookType:
		return uint64(RookValue)
	case QueenType:
		return uint64(QueenValue)
	case KingType:
		return uint64(KingValue)
	default:
		return 0
	}
}

func GetMovesFunction(pieceType PieceType) MovesFunction {
	switch pieceType {
	case PawnType:
		return PawnMoves
	case KnightType:
		return KnightMoves
	case BishopType:
		return BishopMoves
	case RookType:
		return RookMoves
	case QueenType:
		return QueenMoves
	case KingType:
		return KingMoves
	default:
		return nil
	}
}
