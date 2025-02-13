#include <string>

#include "board.hpp"
#include "consts.hpp"
#include "move.hpp"
#include "utils.hpp"

Board::Board() : bitboards{{0}}, ctx{false, 0, 0ULL, 0, 0} {}

Board::Board(std::string fen) : bitboards{{0}}, ctx{false, 0, 0ULL, 0, 0} {
    auto splitted = split(fen, ' ');
    if (splitted.size() != 6) {
        return;
    }

    std::memset(bitboards[0], 0ULL, 6);
    std::memset(bitboards[1], 0ULL, 6);
    int rank = 7;
    int file = 0;
    for (char c : splitted[0]) {
        if (c == '/') {
            --rank;
            file = 0;
            continue;
        }
        if (c >= '1' && c <= '8') {
            file += c - '0'; // skip empty squares
        } else {
            int boardIdx = 0;
            if (c >= 'A' && c <= 'Z') {
                boardIdx = 0;
            } else if (c >= 'a' && c <= 'z') {
                boardIdx = 1;
            }

            int pieceIdx = getPieceIdx(c);
            bitboards[boardIdx][pieceIdx] |=
                1ULL << (rank * 8 + (7 - file)); // 7 - x because of MSB and
                                                 // LSB. See bitboards comment
            file++;
        }
    }

    // Set context
    // White turn
    if (splitted[1] == "w") {
        ctx.whiteTurn = true;
    } else {
        ctx.whiteTurn = false;
    }
    // Castling rights
    if (splitted[2].find('K') != std::string::npos) {
        ctx.castlingRights |= 1;
    }
    if (splitted[2].find('Q') != std::string::npos) {
        ctx.castlingRights |= 1 << 1;
    }
    if (splitted[2].find('k') != std::string::npos) {
        ctx.castlingRights |= 1 << 2;
    }
    if (splitted[2].find('q') != std::string::npos) {
        ctx.castlingRights |= 1 << 3;
    }
    // En passant square
    if (splitted[3] != "-") {
        ctx.enPassantSquare =
            1ULL
            << (8 * (splitted[3][1] - '1') +
                (7 - splitted[3][0] -
                 'a')); // 7 - x because of MSB and LSB. See bitboards comment
    }
    // Half move clock
    ctx.halfMoveClock = std::stoi(splitted[4]);
    // Full move number
    ctx.fullMoveNumber = std::stoi(splitted[5]);

    _isValidFlag = true;
}

Board Board::makePseudoLegalMove(Move m) {
    int color = -1;
    int capturedColor = -1;
    int pieceType = -1;
    int capturedPieceType = -1;
    Board newBoard = *this;

    // Update bitboards
    for (pieceType = PAWN; pieceType < KING; ++pieceType) {
        if (bitboards[WHITE][pieceType] & m.from) {
            newBoard.bitboards[WHITE][pieceType] &= ~m.from;
            newBoard.bitboards[WHITE][pieceType] |= m.to;
            color = WHITE;
        } else if (bitboards[BLACK][pieceType] & m.from) {
            newBoard.bitboards[BLACK][pieceType] &= ~m.from;
            newBoard.bitboards[BLACK][pieceType] |= m.to;
            color = BLACK;
        }
        // Check for captures
        if (bitboards[WHITE][pieceType] & m.to) {
            capturedPieceType = pieceType;
            capturedColor = WHITE;
        } else if (bitboards[BLACK][pieceType] & m.to) {
            capturedPieceType = pieceType;
            capturedColor = BLACK;
        }
    }

    if (color == -1 || pieceType == -1 || color == capturedColor) {
        newBoard._isValidFlag = false;
        return newBoard;
    }

    // Update context
    newBoard.ctx.whiteTurn = !ctx.whiteTurn;
    // TODO Implement update castling rights
    // newBoard.ctx.castlingRights = ctx.castlingRights;
    // TODO implement enPassant Square
    newBoard.ctx.enPassantSquare = 0ULL;
    if (capturedColor == -1 || pieceType != PAWN) {
        newBoard.ctx.halfMoveClock = ctx.halfMoveClock + 1;
    } else {
        newBoard.ctx.halfMoveClock = 0;
    }
    newBoard.ctx.fullMoveNumber += newBoard.ctx.whiteTurn;
    return newBoard;
}

bool Board::isValid() const { return _isValidFlag; }

std::string Board::stringify() const {
    char visualBoard[64] = {0};
    for (int i = 0; i < 6; ++i) {
        for (int j = 0; j < 64; ++j) {
            if (bitboards[0][i] & (1ULL << j)) {
                visualBoard[j] = getPieceChar(i, true);
            } else if (bitboards[1][i] & (1ULL << j)) {
                visualBoard[j] = getPieceChar(i, false);
            }
        }
    }

    std::string result = "-";
    for (int i = 0; i < 8; ++i) {
        result += "----";
    }
    result += "\n";
    // Reverse because the bits are stored from a8 to h1
    for (int i = 63; i >= 0; --i) {
        result += "| ";
        if (visualBoard[i]) {
            result += visualBoard[i];
        } else {
            result += ' '; // two spaces for empty squares
        }
        result += ' ';
        if (i % 8 == 0) {
            result += "|\n-";
            for (int j = 0; j < 8; ++j) {
                result += "----";
            }
            result += '\n';
        }
    }

    return result;
}
