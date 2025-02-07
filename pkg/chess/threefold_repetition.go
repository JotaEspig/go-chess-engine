package chess

import (
	"gce/pkg/utils"
	"hash"
	"hash/fnv"
)

type BoardHash uint64

type ThreefoldRepetitionHashTable map[BoardHash]int

func (trht ThreefoldRepetitionHashTable) Copy() ThreefoldRepetitionHashTable {
	newMap := make(ThreefoldRepetitionHashTable, len(trht))
	for k, v := range trht {
		newMap[k] = v
	}
	return newMap
}

func (p PiecesPosition) hash(h hash.Hash64) {
	utils.HashUint64(h, p.Board)
	utils.HashUint(h, uint(p.Type))
}

func (pb PartialBoard) hash(h hash.Hash64) {
	pb.Pawns.hash(h)
	pb.Knights.hash(h)
	pb.Bishops.hash(h)
	pb.Rooks.hash(h)
	pb.Queens.hash(h)
	pb.King.hash(h)
}

func (b Board) Hash() BoardHash {
	h := fnv.New64()
	b.White.hash(h)
	b.Black.hash(h)
	return BoardHash(h.Sum64())
}
