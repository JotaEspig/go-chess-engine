package chess

type Context struct {
	WhiteTurn              bool
	WhiteCastlingKingSide  bool
	WhiteCastlingQueenSide bool
	BlackCastlingKingSide  bool
	BlackCastlingQueenSide bool
	EnPassant              uint64
	HalfMoves              uint
	MoveNumber             uint
}
