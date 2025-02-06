package chess

import (
	"gce/pkg/utils"
	"hash/fnv"
)

var allLegalBoardMovesHashTable = make(BoardMovesHashTable)

type AllLegalBoardMovesHash uint64

type BoardMovesHashTable map[AllLegalBoardMovesHash][]Move

func (b Board) HashForAllLegalMoves() AllLegalBoardMovesHash {
	h := fnv.New64()
	utils.HashUint64(h, uint64(b.Hash()))
	utils.HashBool(h, b.Ctx.WhiteTurn)
	return AllLegalBoardMovesHash(h.Sum64())
}
