package chess

import (
	"gce/pkg/utils"
	"hash/fnv"
)

var allLegalBoardMovesHashTable = make(BoardMovesHashTable)
var allPossibleBoardMovesHashTable = make(BoardMovesHashTable)

type BoardWithContextHash uint64

type BoardMovesHashTable map[BoardWithContextHash][]*Move

func (b Board) HashBoardWithContext() BoardWithContextHash {
	h := fnv.New64()
	utils.HashUint64(h, uint64(b.Hash()))
	utils.HashBool(h, b.Ctx.WhiteTurn)
	utils.HashBool(h, b.Ctx.WhiteCastlingKingSide)
	utils.HashBool(h, b.Ctx.WhiteCastlingQueenSide)
	utils.HashBool(h, b.Ctx.BlackCastlingKingSide)
	utils.HashBool(h, b.Ctx.BlackCastlingQueenSide)
	utils.HashUint64(h, b.Ctx.EnPassant)
	return BoardWithContextHash(h.Sum64())
}
