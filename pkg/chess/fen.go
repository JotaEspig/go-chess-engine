package chess

import (
	"strconv"
	"strings"

	"github.com/charmbracelet/log"
)

func FenToBoard(fen string) *Board {
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
	if len(splitted) != 5 {
		log.Fatal("Invalid fen context")
	}

	ctx := Context{}
	if splitted[0] == "w" {
		ctx.WhiteTurn = true
	} else {
		ctx.WhiteTurn = false
	}

	castling := splitted[1]
	if strings.Contains(castling, "K") {
		ctx.WhiteCastlingKingSide = true
	}
	if strings.Contains(castling, "Q") {
		ctx.WhiteCastlingQueenSide = true
	}
	if strings.Contains(castling, "k") {
		ctx.BlackCastlingKingSide = true
	}
	if strings.Contains(castling, "q") {
		ctx.BlackCastlingQueenSide = true
	}

	enPassantCoord := splitted[2]
	if enPassantCoord != "-" {
		col := int(enPassantCoord[0] - 'a')
		row := int(enPassantCoord[1] - '1')
		ctx.EnPassant = PositionToUInt64(col, row)
	}

	halfMoveClock := splitted[3]
	HalfMovesInt, err := strconv.Atoi(halfMoveClock)
	if err != nil {
		log.Fatalf("Invalid half move clock: %v", halfMoveClock)
	}
	moveNumber := splitted[4]
	moveNumberInt, err := strconv.Atoi(moveNumber)
	if err != nil {
		log.Fatalf("Invalid move number: %v", moveNumber)
	}

	ctx.MoveNumber = uint(moveNumberInt)
	ctx.HalfMoves = uint(HalfMovesInt)
	return ctx
}
