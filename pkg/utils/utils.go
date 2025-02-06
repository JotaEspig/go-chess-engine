package utils

import "hash"

func Filter[T any](s []T, f func(T) bool) []T {
	var r []T
	for _, v := range s {
		if f(v) {
			r = append(r, v)
		}
	}
	return r
}

func HashUint64(h hash.Hash64, value uint64) {
	var buf [8]byte
	for i := uint(0); i < 8; i++ {
		buf[i] = byte(value >> (i * 8))
	}
	h.Write(buf[:])
}

func HashUint(h hash.Hash64, value uint) {
	var buf [4]byte
	for i := uint(0); i < 4; i++ {
		buf[i] = byte(value >> (i * 8))
	}
	h.Write(buf[:])
}

func HashBool(h hash.Hash64, value bool) {
	if value {
		h.Write([]byte{1})
	} else {
		h.Write([]byte{0})
	}
}
