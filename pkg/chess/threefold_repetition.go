package chess

import (
	"hash"
	"hash/fnv"
)

type BoardHash uint64

type ThreefoldRepetitionHashTable map[BoardHash]int

func (trht ThreefoldRepetitionHashTable) Copy() ThreefoldRepetitionHashTable {
	newMap := make(ThreefoldRepetitionHashTable)
	for k, v := range trht {
		newMap[k] = v
	}
	return newMap
}

func (p PiecesPosition) hash(h hash.Hash64) {
	hashUint64(h, p.Board)
	hashUint(h, uint(p.Type))
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

func hashUint64(h hash.Hash64, value uint64) {
	var buf [8]byte
	for i := uint(0); i < 8; i++ {
		buf[i] = byte(value >> (i * 8))
	}
	h.Write(buf[:])
}

func hashUint(h hash.Hash64, value uint) {
	var buf [4]byte
	for i := uint(0); i < 4; i++ {
		buf[i] = byte(value >> (i * 8))
	}
	h.Write(buf[:])
}
