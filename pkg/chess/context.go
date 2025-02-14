package chess

const (
	NoResult uint = iota
	WhiteWin
	BlackWin
	Draw
)

type ContextCache struct {
	IsDrawCache           bool
	IsDrawCacheSet        bool
	IsMatedCache          bool
	IsMatedCacheSet       bool
	IsKingInCheckCache    bool
	IsKingInCheckCacheSet bool
}

type Context struct {
	ContextCache
	WhiteTurn              bool
	WhiteCastlingKingSide  bool
	WhiteCastlingQueenSide bool
	BlackCastlingKingSide  bool
	BlackCastlingQueenSide bool
	EnPassant              uint64
	HalfMoves              uint
	MoveNumber             uint
	Result                 uint
}
