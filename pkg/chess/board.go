package chess

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/charmbracelet/log"
)

// Board represents a full board with pieces of both colors on it.
type Board struct {
	White PartialBoard
	Black PartialBoard
	Ctx   Context

	// These 2 fields are used to keep track of the previous board state (i.e. every move made on board)
	MoveDone  Move   // MoveDone is the move that was done in the current board
	PrevBoard *Board // PrevBoard is the previous board state
}

func NewDefaultBoard() *Board {
	return FenToBoard(DefaultStartFen)
}

func NewBoard() *Board {
	return &Board{
		White: NewPartialBoard(),
		Black: NewPartialBoard(),
	}
}

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
		ctx.WhiteToMove = true
	} else {
		ctx.WhiteToMove = false
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

func (b Board) IsValidPosition() bool {
	// Check if there are no pieces in the same position
	whitePieces := b.White.AllBoardMask()
	blackPieces := b.Black.AllBoardMask()
	if whitePieces&blackPieces != 0 {
		return false
	}

	// originalWhiteToMove := b.Ctx.WhiteToMove
	if b.Ctx.WhiteToMove {
		// Check if the black king is in check while it's white's turn. Meaning that an illegal move was made.
		b.Ctx.WhiteToMove = false
		if b.IsKingInCheck() {
			return false
		}
	} else {
		// Check if the white king is in check while it's black's turn. Meaning that an illegal move was made.
		b.Ctx.WhiteToMove = true
		if b.IsKingInCheck() {
			return false
		}
	}

	return true
}

func (b Board) AllLegalMoves() []Move {
	moves := b.AllPossibleMoves()
	// Filter out moves that are not legal
	moves = Filter(moves, func(m Move) bool {
		newBoard := b
		return newBoard.MakeMove(m)
	})

	return moves
}

func (b Board) AllPossibleMoves() []Move {
	moves := make([]Move, 0)
	// "Normal" moves (including En passant and promotion)
	if b.Ctx.WhiteToMove {
		moves = append(moves, b.White.AllPossibleMoves(b)...)
	} else {
		moves = append(moves, b.Black.AllPossibleMoves(b)...)
	}

	// Castling
	moves = append(moves, b.AllCastlingMoves()...)

	return moves
}

func (b Board) allMovesToDefendCheck() []Move {
	if !b.IsKingInCheck() {
		log.Fatal("King is not in check")
	}

	return b.AllLegalMoves()
}

func (b Board) AllCastlingMoves() []Move {
	var pb PartialBoard
	anyCastleAvailable := false
	var kingSideSpaceMask uint64
	var QueenSideSpaceMask uint64
	var kingSideSafeSpot uint64
	var queenSideSafeSpot uint64
	if b.Ctx.WhiteToMove {
		pb = b.White
		anyCastleAvailable = b.Ctx.WhiteCastlingKingSide || b.Ctx.WhiteCastlingQueenSide
		kingSideSpaceMask = ^uint64(6)    // 6 is the bits that represents F1 and G1
		QueenSideSpaceMask = ^uint64(112) // 112 is the bits that represents B1, C1 and D1
		kingSideSafeSpot = uint64(2)      // G1
		queenSideSafeSpot = uint64(32)    // C1
	} else {
		pb = b.Black
		anyCastleAvailable = b.Ctx.BlackCastlingKingSide || b.Ctx.BlackCastlingQueenSide
		kingSideSpaceMask = ^uint64(432_345_564_227_567_616)    // 432_345_564_227_567_616 is the bits that represents F8 and G8
		QueenSideSpaceMask = ^uint64(8_070_450_532_247_928_832) // 8_070_450_532_247_928_832 is the bits that represents B8, C8, D8
		kingSideSafeSpot = uint64(144_115_188_075_855_872)      // G8
		queenSideSafeSpot = uint64(2_305_843_009_213_693_952)   // C8
	}

	if !anyCastleAvailable {
		return []Move{}
	}

	moves := make([]Move, 0)
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

// MakeMove makes a move on the board.
// Returns true if it's a valid move, false otherwise.
func (b *Board) MakeMove(m Move) bool {
	prevBoard := &Board{}
	*prevBoard = b.Copy()

	var pb *PartialBoard
	if b.Ctx.WhiteToMove {
		pb = &b.White
	} else {
		pb = &b.Black
	}

	var mask uint64
	switch m.PieceType {
	case PawnType:
		mask = pb.Pawns.Board
	case KnightType:
		mask = pb.Knights.Board
	case BishopType:
		mask = pb.Bishops.Board
	case RookType:
		mask = pb.Rooks.Board
	case QueenType:
		mask = pb.Queens.Board
	case KingType:
		mask = pb.King.Board
	default:
		log.Fatalf("Invalid piece type: %v", m.PieceType)
	}

	// Check if the piece is at the position
	if mask&m.OldPiecePos == 0 {
		return false
	}

	// Castling verifications
	if m.IsCastling {
		if b.IsKingInCheck() {
			return false
		}

		isKingSide := m.NewPiecePos < m.OldPiecePos
		copyBoard := *b
		copyMove := m
		copyMove.IsCastling = false
		// King side
		if isKingSide {
			copyMove.NewPiecePos <<= 1
		} else { // Queen side
			copyMove.NewPiecePos >>= 1
		}

		// Means that king moves in a square attacked by enemy piece
		if !copyBoard.MakeMove(copyMove) {
			return false
		}

		// Check if it is allowed castling giving the context of the board
		if b.Ctx.WhiteToMove {
			if isKingSide {
				if !b.Ctx.WhiteCastlingKingSide {
					return false
				}
			} else {
				if !b.Ctx.WhiteCastlingQueenSide {
					return false
				}
			}
		} else {
			if isKingSide {
				if !b.Ctx.BlackCastlingKingSide {
					return false
				}
			} else {
				if !b.Ctx.BlackCastlingQueenSide {
					return false
				}
			}
		}
	}

	if m.IsPromotion {
		pb.MakePromotion(m)
	} else {
		pb.MakeMove(m)
	}

	// Check for next move En passant
	// Default value is 0
	var enPassantPos uint64 = 0
	// If is 2 square pawn move, set the en passant position for the one row behind the pawn
	if m.Is2SquarePawnMove() {
		isWhite := m.NewPiecePos > m.OldPiecePos // Assumes it's a pawn move
		if isWhite {
			enPassantPos = moveUp(m.OldPiecePos, 1)
		} else {
			enPassantPos = moveDown(m.OldPiecePos, 1)
		}
	}
	if m.IsCapture {
		// Inverted color to erase the piece from the board
		if b.Ctx.WhiteToMove {
			pb = &b.Black
		} else {
			pb = &b.White
		}

		colrow := Int64toPositions(m.NewPiecePos)
		if len(colrow) != 1 {
			log.Fatal("Invalid NewPiecePos: %v", m.NewPiecePos)
		}

		// Clear the piece at the position at every piece type board
		// This is necessary because the piece could be at any of the piece type boards
		// and we don't know which one it is.
		pb.Pawns.Board &= ^m.NewPiecePos
		pb.Knights.Board &= ^m.NewPiecePos
		pb.Bishops.Board &= ^m.NewPiecePos
		pb.Rooks.Board &= ^m.NewPiecePos
		pb.Queens.Board &= ^m.NewPiecePos
		pb.King.Board &= ^m.NewPiecePos
	}

	// means black completed its turn
	if !b.Ctx.WhiteToMove {
		b.Ctx.MoveNumber++
	}
	// Check for half moves
	if !m.IsCapture && m.PieceType != PawnType {
		b.Ctx.HalfMoves++
	} else {
		b.Ctx.HalfMoves = 0
	}
	b.Ctx.WhiteToMove = !b.Ctx.WhiteToMove
	b.Ctx.EnPassant = enPassantPos

	if !b.IsValidPosition() {
		return false
	}
	b.PrevBoard = prevBoard
	b.PrevBoard.MoveDone = m
	return true
}

func (b Board) IsMated() bool {
	isKingInCheckOriginalPos := b.IsKingInCheck()
	if !isKingInCheckOriginalPos {
		return false
	}

	possibleDefensiveMoves := b.allMovesToDefendCheck()
	return len(possibleDefensiveMoves) == 0
}

// IsKingInCheck returns true if the king is in check.
func (b Board) IsKingInCheck() bool {
	var kingPos uint64
	var enemyPb *PartialBoard

	isWhite := b.Ctx.WhiteToMove
	if isWhite {
		kingPos = b.White.King.Board
		enemyPb = &b.Black
	} else {
		kingPos = b.Black.King.Board
		enemyPb = &b.White
	}

	// Invert the color to get the enemy moves and see if it's possible to "capture" the king
	b.Ctx.WhiteToMove = !b.Ctx.WhiteToMove
	possibleCheckMoves := make([]Move, 0)
	possibleCheckMoves = append(possibleCheckMoves, enemyPb.Pawns.AllPossibleMoves(b)...)
	possibleCheckMoves = append(possibleCheckMoves, enemyPb.Knights.AllPossibleMoves(b)...)
	possibleCheckMoves = append(possibleCheckMoves, enemyPb.Bishops.AllPossibleMoves(b)...)
	possibleCheckMoves = append(possibleCheckMoves, enemyPb.Rooks.AllPossibleMoves(b)...)
	possibleCheckMoves = append(possibleCheckMoves, enemyPb.Queens.AllPossibleMoves(b)...)

	for _, move := range possibleCheckMoves {
		if move.NewPiecePos == kingPos {
			return true
		}
	}
	return false
}

func (b Board) MaterialValueBalance() int64 {
	return int64(b.White.MaterialValue()) - int64(b.Black.MaterialValue())
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

func (b Board) ParseMove(notation string) (Move, error) {
	originalNotation := notation
	notation = strings.ReplaceAll(notation, "+", "")
	notation = strings.ReplaceAll(notation, "#", "")

	if notation == "O-O" {
		move := Move{IsCastling: true, PieceType: KingType}
		if b.Ctx.WhiteToMove {
			move.OldPiecePos = b.White.King.Board
		} else {
			move.OldPiecePos = b.Black.King.Board
		}
		move.NewPiecePos = moveRight(move.OldPiecePos, 2)
		return move, nil
	} else if notation == "O-O-O" {
		move := Move{IsCastling: true, PieceType: KingType}
		if b.Ctx.WhiteToMove {
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
	if b.Ctx.WhiteToMove {
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
	piecePossibleMoves = Filter(piecePossibleMoves, func(m Move) bool {
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
			piecePossibleMoves = Filter(piecePossibleMoves, func(m Move) bool {
				return Int64toPositions(m.OldPiecePos)[0][0] == col
			})
		} else {
			row := int(source[0] - '1')
			// Filter out moves that are not the source row
			piecePossibleMoves = Filter(piecePossibleMoves, func(m Move) bool {
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
				piecePossibleMoves = Filter(piecePossibleMoves, func(m Move) bool {
					return Int64toPositions(m.OldPiecePos)[0][0] == col
				})
			} else {
				row := int(remaningAmbiguityRemoval - '1')
				// Filter out moves that are not the source row
				piecePossibleMoves = Filter(piecePossibleMoves, func(m Move) bool {
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
		if move.NewPiecePos < move.OldPiecePos {
			return "O-O-O"
		}
		return "O-O"
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
	if b.Ctx.WhiteToMove {
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
	possiblePieceMoves = Filter(possiblePieceMoves, func(m Move) bool {
		return m.NewPiecePos == move.NewPiecePos && m.NewPieceType == move.NewPieceType
	})
	length := len(possiblePieceMoves)
	if length > 1 {
		// Check for ambiguity
		// Check for column ambiguity
		possiblePieceMoves = Filter(possiblePieceMoves, func(m Move) bool {
			return Int64toPositions(m.OldPiecePos)[0][0] == Int64toPositions(move.OldPiecePos)[0][0]
		})
		if length != len(possiblePieceMoves) {
			notation += string(rune('a' + Int64toPositions(move.OldPiecePos)[0][0]))
		}
		if len(possiblePieceMoves) > 1 {
			// Check for row ambiguity
			possiblePieceMoves = Filter(possiblePieceMoves, func(m Move) bool {
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
	notation += string(rune('a' + destCol))
	notation += string(rune('1' + destRow))

	return ""
}

func (b Board) getMoveListInNotation() string {
	if b.MoveDone != (Move{}) {
		return b.PrevBoard.getMoveListInNotation()
	}
	moveNotation := b.MoveToNotation(b.MoveDone)
	if b.PrevBoard == nil {
		return "1. " + moveNotation
	}
	moveNumberIfNeeded := ""
	if b.Ctx.WhiteToMove {
		moveNumberInt := b.Ctx.MoveNumber
		moveNumberIfNeeded = strconv.Itoa(int(moveNumberInt)) + ". "
	}
	return b.PrevBoard.getMoveListInNotation() + " " + moveNumberIfNeeded + moveNotation
}

func (b Board) GetMoveListInNotation() string {
	moveList := b.getMoveListInNotation()
	moveList = strings.TrimSpace(moveList)
	return moveList
}

func (b Board) Copy() Board {
	return Board{
		White:     b.White,
		Black:     b.Black,
		Ctx:       b.Ctx,
		PrevBoard: b.PrevBoard,
	}
}
