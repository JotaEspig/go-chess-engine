#include <vector>

#include "board.hpp"
#include "game.hpp"
#include "move.hpp"

Game::Game() {}

Game::Game(std::string fen) { boards.push_back(Board(fen)); }

Board Game::makePseudoLegalMove(Move m) {
    Board last = boards.back();
    Board b = last.makePseudoLegalMove(m);
    boards.push_back(b);
    return b;
}