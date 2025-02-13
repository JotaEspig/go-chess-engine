#include <iostream>

#include "board.hpp"
#include "consts.hpp"

int main() {
    std::cout << "Hello, World!" << std::endl;

    Board b{DEFAULT_FEN};
    std::string s = b.stringify();
    std::cout << s << std::endl;

    Move m{2048ULL, 134217728ULL};
    Board nb = b.makePseudoLegalMove(m);
    std::string ns = nb.stringify();
    std::cout << ns << std::endl;

    return 0;
}