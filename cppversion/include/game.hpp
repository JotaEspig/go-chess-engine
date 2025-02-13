#pragma once

#include <vector>

#include "board.hpp"
#include "move.hpp"

class Game {
  public:
    std::vector<Board> boards;

    Game();
    Game(std::string fen);

    Board makePseudoLegalMove(Move m);
    Board undoMove();
};