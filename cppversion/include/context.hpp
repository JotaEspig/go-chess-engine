#pragma once

#include <cstdint>

struct Context {
    bool whiteTurn;
    // a nibble representing the castling rights of both
    // players, with the following mapping: 0b0000 -> no
    // rights, 0b0001 -> white king side, 0b0010 ->
    // white queen side, 0b0011 -> white both, 0b0100 ->
    // black king side, 0b1000 -> black queen side,
    // 0b1100 -> black both
    uint8_t castlingRights;
    uint64_t enPassantSquare;
    uint16_t halfMoveClock;
    uint16_t fullMoveNumber;
};