package chess

import (
	"gce/pkg/utils"
	"hash/fnv"
)

var allLegalBoardMovesHashTable = make(BoardMovesHashTable)
var allPossibleBoardMovesHashTable = make(BoardMovesHashTable)

type AllBoardMovesHash uint64

type BoardMovesHashTable map[AllBoardMovesHash][]*Move

func (b Board) HashForAllMoves() AllBoardMovesHash {
	h := fnv.New64()
	utils.HashUint64(h, uint64(b.Hash()))
	utils.HashBool(h, b.Ctx.WhiteTurn)
	return AllBoardMovesHash(h.Sum64())
}
