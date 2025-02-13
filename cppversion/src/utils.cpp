#include <cctype>
#include <iosfwd>
#include <sstream>
#include <string>
#include <vector>

std::vector<std::string> split(const std::string &s, char delim) {
    std::vector<std::string> result;
    std::stringstream ss(s);
    std::string item;
    while (std::getline(ss, item, delim)) {
        result.push_back(item);
    }
    return result;
}

int getPieceIdx(char c) {
    c = std::tolower(c);
    int result = -1;
    switch (c) {
    case 'p':
        result = 0;
        break;
    case 'n':
        result = 1;
        break;
    case 'b':
        result = 2;
        break;
    case 'r':
        result = 3;
        break;
    case 'q':
        result = 4;
        break;
    case 'k':
        result = 5;
        break;
    default:
        break;
    }
    return result;
}

char getPieceChar(int i, bool isWhite) {
    char result = ' ';
    switch (i) {
    case 0:
        result = 'p';
        break;
    case 1:
        result = 'n';
        break;
    case 2:
        result = 'b';
        break;
    case 3:
        result = 'r';
        break;
    case 4:
        result = 'q';
        break;
    case 5:
        result = 'k';
        break;
    default:
        break;
    }
    if (isWhite) {
        return std::toupper(result);
    }
    return result;
}
