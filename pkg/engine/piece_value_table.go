package engine

import "gce/pkg/chess"

var whitePawnTable = [64]float64{
	0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0,
	0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5,
	0.1, 0.1, 0.2, 0.3, 0.3, 0.2, 0.1, 0.1,
	0.05, 0.05, 0.1, 0.25, 0.25, 0.1, 0.05, 0.05,
	0.0, 0.0, 0.0, 0.2, 0.2, 0.0, 0.0, 0.0,
	0.05, -0.05, -0.1, 0.0, 0.0, -0.1, -0.05, 0.05,
	0.05, 0.1, 0.1, -0.2, -0.2, 0.1, 0.1, 0.05,
	0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0,
}

var whiteKnightTable = [64]float64{
	-0.5, -0.4, -0.3, -0.3, -0.3, -0.3, -0.4, -0.5,
	-0.4, -0.2, 0.0, 0.0, 0.0, 0.0, -0.2, -0.4,
	-0.3, 0.0, 0.1, 0.15, 0.15, 0.1, 0.0, -0.3,
	-0.3, 0.05, 0.15, 0.2, 0.2, 0.15, 0.05, -0.3,
	-0.3, 0.0, 0.15, 0.2, 0.2, 0.15, 0.0, -0.3,
	-0.3, 0.05, 0.1, 0.15, 0.15, 0.1, 0.05, -0.3,
	-0.4, -0.2, 0.0, 0.05, 0.05, 0.0, -0.2, -0.4,
	-0.5, -0.4, -0.3, -0.3, -0.3, -0.3, -0.4, -0.5,
}

var whiteBishopTable = [64]float64{
	-0.2, -0.1, -0.1, -0.1, -0.1, -0.1, -0.1, -0.2,
	-0.1, 0.05, 0.0, 0.0, 0.0, 0.0, 0.05, -0.1,
	-0.1, 0.1, 0.1, 0.1, 0.1, 0.1, 0.1, -0.1,
	-0.1, 0.0, 0.1, 0.1, 0.1, 0.1, 0.0, -0.1,
	-0.1, 0.05, 0.05, 0.1, 0.1, 0.05, 0.05, -0.1,
	-0.1, 0.0, 0.05, 0.1, 0.1, 0.05, 0.0, -0.1,
	-0.1, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, -0.1,
	-0.2, -0.1, -0.1, -0.1, -0.1, -0.1, -0.1, -0.2,
}

var whiteRookTable = [64]float64{
	0.25, 0.3, 0.3, 0.3, 0.3, 0.3, 0.3, 0.25,
	0.25, 0.3, 0.3, 0.3, 0.3, 0.3, 0.3, 0.25,
	0.0, 0.0, 0.0, 0.05, 0.05, 0.0, 0.0, 0.0,
	-0.05, 0.0, 0.0, 0.05, 0.05, 0.0, 0.0, -0.05,
	-0.05, 0.0, 0.0, 0.05, 0.05, 0.0, 0.0, -0.05,
	0.0, 0.0, 0.0, 0.05, 0.05, 0.0, 0.0, 0.0,
	0.25, 0.3, 0.3, 0.3, 0.3, 0.3, 0.3, 0.25,
	0.25, 0.3, 0.3, 0.3, 0.3, 0.3, 0.3, 0.25,
}

var whiteQueenTable = [64]float64{
	-0.2, -0.1, -0.1, -0.05, -0.05, -0.1, -0.1, -0.2,
	-0.1, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, -0.1,
	-0.1, 0.0, 0.05, 0.05, 0.05, 0.05, 0.0, -0.1,
	-0.05, 0.0, 0.05, 0.05, 0.05, 0.05, 0.0, -0.05,
	0.0, 0.0, 0.05, 0.05, 0.05, 0.05, 0.0, 0.0,
	-0.1, 0.05, 0.05, 0.05, 0.05, 0.05, 0.05, -0.1,
	-0.1, 0.0, 0.05, 0.0, 0.0, 0.0, 0.0, -0.1,
	-0.2, -0.1, -0.1, -0.05, -0.05, -0.1, -0.1, -0.2,
}

var blackPawnTable = mirrorPST(whitePawnTable)
var blackKnightTable = mirrorPST(whiteKnightTable)
var blackBishopTable = mirrorPST(whiteBishopTable)
var blackRookTable = mirrorPST(whiteRookTable)
var blackQueenTable = mirrorPST(whiteQueenTable)

func mirrorPST(pst [64]float64) [64]float64 {
	var mirrored [64]float64
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			mirrored[i*8+j] = pst[(7-i)*8+j] // Flip along horizontal axis
		}
	}
	return mirrored
}

func pieceTypeTableValue(piecesPosition uint64, pieceType chess.PieceType, isWhite bool) float64 {
	var table *[64]float64
	switch pieceType {
	case chess.PawnType:
		if isWhite {
			table = &whitePawnTable
		} else {
			table = &blackPawnTable
		}
	case chess.KnightType:
		if isWhite {
			table = &whiteKnightTable
		} else {
			table = &blackKnightTable
		}
	case chess.BishopType:
		if isWhite {
			table = &whiteBishopTable
		} else {
			table = &blackBishopTable
		}
	case chess.RookType:
		if isWhite {
			table = &whiteRookTable
		} else {
			table = &blackRookTable
		}
	case chess.QueenType:
		if isWhite {
			table = &whiteQueenTable
		} else {
			table = &blackQueenTable
		}
	}

	var value float64
	for i := 0; i < 64; i++ {
		if (piecesPosition>>uint(i))&1 == 1 {
			value += table[i]
		}
	}
	if !isWhite {
		value = -value
	}
	return value
}

func BoardEvaluationByPieceSquareTable(board chess.Board) float64 {
	evaluation := pieceTypeTableValue(board.White.Pawns.Board, chess.PawnType, true)
	evaluation += pieceTypeTableValue(board.White.Knights.Board, chess.KnightType, true)
	evaluation += pieceTypeTableValue(board.White.Bishops.Board, chess.BishopType, true)
	evaluation += pieceTypeTableValue(board.White.Rooks.Board, chess.RookType, true)
	evaluation += pieceTypeTableValue(board.White.Queens.Board, chess.QueenType, true)
	evaluation += pieceTypeTableValue(board.Black.Pawns.Board, chess.PawnType, false)
	evaluation += pieceTypeTableValue(board.Black.Knights.Board, chess.KnightType, false)
	evaluation += pieceTypeTableValue(board.Black.Bishops.Board, chess.BishopType, false)
	evaluation += pieceTypeTableValue(board.Black.Rooks.Board, chess.RookType, false)
	evaluation += pieceTypeTableValue(board.Black.Queens.Board, chess.QueenType, false)
	return evaluation
}
