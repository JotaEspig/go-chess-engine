package engine

type PieceType int
type PieceValue int

// Enum for the different types of pieces.
const (
	InvalidType PieceType = iota
	PawnType
	KnightType
	BishopType
	RookType
	QueenType
	KingType
)

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

func (p PieceType) Value() int64 {
	switch p {
	case PawnType:
		return int64(PawnValue)
	case KnightType:
		return int64(KnightValue)
	case BishopType:
		return int64(BishopValue)
	case RookType:
		return int64(RookValue)
	case QueenType:
		return int64(QueenValue)
	case KingType:
		return int64(KingValue)
	default:
		return 0
	}
}

// Enum for the different values for each pieces.
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
