package chess

import (
	"errors"
	"fmt"
	"gce/pkg/utils"
	"strings"

	"github.com/charmbracelet/log"
)

func (b Board) ParseMove(notation string) (Move, error) {
	originalNotation := notation
	notation = strings.ReplaceAll(notation, "+", "")
	notation = strings.ReplaceAll(notation, "#", "")

	if notation == "O-O" {
		move := Move{IsCastling: true, PieceType: KingType}
		if b.Ctx.WhiteTurn {
			move.OldPiecePos = b.White.King.Board
		} else {
			move.OldPiecePos = b.Black.King.Board
		}
		move.NewPiecePos = moveRight(move.OldPiecePos, 2)
		return move, nil
	} else if notation == "O-O-O" {
		move := Move{IsCastling: true, PieceType: KingType}
		if b.Ctx.WhiteTurn {
			move.OldPiecePos = b.White.King.Board
		} else {
			move.OldPiecePos = b.Black.King.Board
		}
		move.NewPiecePos = moveLeft(move.OldPiecePos, 2)
		return move, nil
	}

	var pieceType PieceType
	var newPieceType PieceType
	// Check if it's a promotion
	if strings.Contains(notation, "=") {
		parts := strings.Split(notation, "=")
		switch parts[1] {
		case "Q":
			newPieceType = QueenType
		case "R":
			newPieceType = RookType
		case "B":
			newPieceType = BishopType
		case "N":
			newPieceType = KnightType
		default:
			return Move{}, errors.New("Invalid promotion piece type")
		}

		notation = parts[0]
		pieceType = PawnType
	} else {
		// Get type of piece
		switch notation[0] {
		case 'N':
			pieceType = KnightType
		case 'B':
			pieceType = BishopType
		case 'R':
			pieceType = RookType
		case 'Q':
			pieceType = QueenType
		case 'K':
			pieceType = KingType
		default:
			pieceType = PawnType
		}
	}

	var piecePossibleMoves []Move
	var pb PartialBoard
	if b.Ctx.WhiteTurn {
		pb = b.White
	} else {
		pb = b.Black
	}

	switch pieceType {
	case PawnType:
		piecePossibleMoves = pb.Pawns.AllPossibleMoves(b)
	case KnightType:
		piecePossibleMoves = pb.Knights.AllPossibleMoves(b)
	case BishopType:
		piecePossibleMoves = pb.Bishops.AllPossibleMoves(b)
	case RookType:
		piecePossibleMoves = pb.Rooks.AllPossibleMoves(b)
	case QueenType:
		piecePossibleMoves = pb.Queens.AllPossibleMoves(b)
	case KingType:
		piecePossibleMoves = pb.King.AllPossibleMoves(b)
	default:
		return Move{}, errors.New(fmt.Sprintf("Invalid piece type: %v", pieceType))
	}

	// Check if it's a capture
	isCapture := strings.Contains(notation, "x")
	// Get the destination position
	var destination string
	if isCapture {
		splitted := strings.Split(notation, "x")
		destination = splitted[1]
	} else {
		destination = notation[len(notation)-2:]
	}

	// Get the destination position
	col := int(destination[0] - 'a')
	row := int(destination[1] - '1')
	destinationPos := PositionToUInt64(col, row)

	// Filter out moves that are not the destination position
	piecePossibleMoves = utils.Filter(piecePossibleMoves, func(m Move) bool {
		return m.NewPiecePos == destinationPos && m.NewPieceType == newPieceType
	})

	var move Move
	if len(piecePossibleMoves) == 0 {
		return Move{}, errors.New(fmt.Sprintf("Invalid move: %v", originalNotation))
	} else if len(piecePossibleMoves) == 1 {
		move = piecePossibleMoves[0]
	} else {
		var source string
		if pieceType == PawnType {
			source = string(notation[0])
		} else {
			source = string(notation[1])
		}

		// if source is a column
		if source >= "a" && source <= "h" {
			col := int(source[0] - 'a')
			// Filter out moves that are not the source column
			piecePossibleMoves = utils.Filter(piecePossibleMoves, func(m Move) bool {
				return Int64toPositions(m.OldPiecePos)[0][0] == col
			})
		} else {
			row := int(source[0] - '1')
			// Filter out moves that are not the source row
			piecePossibleMoves = utils.Filter(piecePossibleMoves, func(m Move) bool {
				return Int64toPositions(m.OldPiecePos)[0][1] == row
			})
		}

		if len(piecePossibleMoves) == 0 {
			return Move{}, errors.New(fmt.Sprintf("Invalid move: %v", originalNotation))
		}

		if len(piecePossibleMoves) > 1 {
			remaningAmbiguityRemoval := notation[2]
			if remaningAmbiguityRemoval >= 'a' && remaningAmbiguityRemoval <= 'h' {
				col := int(remaningAmbiguityRemoval - 'a')
				// Filter out moves that are not the source column
				piecePossibleMoves = utils.Filter(piecePossibleMoves, func(m Move) bool {
					return Int64toPositions(m.OldPiecePos)[0][0] == col
				})
			} else {
				row := int(remaningAmbiguityRemoval - '1')
				// Filter out moves that are not the source row
				piecePossibleMoves = utils.Filter(piecePossibleMoves, func(m Move) bool {
					return Int64toPositions(m.OldPiecePos)[0][1] == row
				})
			}
		}

		if len(piecePossibleMoves) != 1 {
			return Move{}, errors.New(fmt.Sprintf("Invalid move: %v", originalNotation))
		}

		move = piecePossibleMoves[0]
	}

	return move, nil
}

func (b Board) MoveToNotation(move Move) string {
	if move.IsCastling {
		castle := "O-O"
		if move.NewPiecePos > move.OldPiecePos {
			castle = "O-O-O"
		}
		if move.IsCheck {
			castle += "+"
		}
		return castle
	}

	notation := ""
	switch move.PieceType {
	case KnightType:
		notation += "N"
	case BishopType:
		notation += "B"
	case RookType:
		notation += "R"
	case QueenType:
		notation += "Q"
	case KingType:
		notation += "K"
	default:
		sourceCol := Int64toPositions(move.OldPiecePos)[0][0]
		notation += string(rune('a' + sourceCol))
	}

	// Check for ambiguity
	var pb PartialBoard
	var possiblePieceMoves []Move
	if b.Ctx.WhiteTurn {
		pb = b.White
	} else {
		pb = b.Black
	}
	switch move.PieceType {
	case PawnType:
		possiblePieceMoves = pb.Pawns.AllPossibleMoves(b)
	case KnightType:
		possiblePieceMoves = pb.Knights.AllPossibleMoves(b)
	case BishopType:
		possiblePieceMoves = pb.Bishops.AllPossibleMoves(b)
	case RookType:
		possiblePieceMoves = pb.Rooks.AllPossibleMoves(b)
	case QueenType:
		possiblePieceMoves = pb.Queens.AllPossibleMoves(b)
	case KingType:
		possiblePieceMoves = pb.King.AllPossibleMoves(b)
	default:
		log.Fatalf("Invalid piece type: %v", move.PieceType)
	}

	// Filter out moves that are not the destination position
	possiblePieceMoves = utils.Filter(possiblePieceMoves, func(m Move) bool {
		return m.NewPiecePos == move.NewPiecePos && m.NewPieceType == move.NewPieceType
	})
	length := len(possiblePieceMoves)
	if length > 1 {
		// Check for ambiguity
		// Check for column ambiguity
		possiblePieceMoves = utils.Filter(possiblePieceMoves, func(m Move) bool {
			return Int64toPositions(m.OldPiecePos)[0][0] == Int64toPositions(move.OldPiecePos)[0][0]
		})
		if length != len(possiblePieceMoves) {
			notation += string(rune('a' + Int64toPositions(move.OldPiecePos)[0][0]))
		}
		if len(possiblePieceMoves) > 1 {
			// Check for row ambiguity
			possiblePieceMoves = utils.Filter(possiblePieceMoves, func(m Move) bool {
				return Int64toPositions(m.OldPiecePos)[0][1] == Int64toPositions(move.OldPiecePos)[0][1]
			})
		}
		if len(possiblePieceMoves) > 1 {
			log.Fatalf("Invalid move: %v", move)
		}
		if len(possiblePieceMoves) == 1 {
			notation += string(rune('1' + Int64toPositions(move.OldPiecePos)[0][1]))
		}
	}

	if move.IsCapture {
		notation += "x"
	}
	dest := Int64toPositions(move.NewPiecePos)[0]
	destCol, destRow := dest[0], dest[1]
	// if it's a pawn move and not a capture, then the column was already added
	if move.PieceType != PawnType || move.IsCapture { // De morgan baby
		notation += string(rune('a' + destCol))
	}
	notation += string(rune('1' + destRow))

	if move.IsPromotion {
		notation := "="
		switch move.NewPieceType {
		case QueenType:
			notation += "Q"
		case RookType:
			notation += "R"
		case BishopType:
			notation += "B"
		case KnightType:
			notation += "N"
		default:
			log.Fatalf("Invalid piece type: %v", move.NewPieceType)
		}
	}
	if move.IsCheck {
		notation += "+"
	}

	if strings.TrimSpace(notation) == "" {
		log.Fatalf("Invalid move: %v", move)
	}

	return notation
}

func (b Board) getMoveListInNotation() string {
	s := ""
	for i, move := range b.MovesDone {
		if i%2 == 0 {
			s += fmt.Sprintf("%d. ", i/2+1)
		}
		s += b.MoveToNotation(move) + " "
	}
	if b.IsMated() {
		s = strings.TrimSpace(s)
		s += "#"
	}
	return s
}

func (b Board) GetMoveListInNotation() string {
	moveList := b.getMoveListInNotation()
	moveList = strings.TrimSpace(moveList)
	switch b.Ctx.Result {
	case WhiteWin:
		moveList += " 1-0"
	case BlackWin:
		moveList += " 0-1"
	case Draw:
		moveList += " 1/2-1/2"
	}
	return moveList
}
