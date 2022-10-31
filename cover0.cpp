
#include "types.hpp"
#include "cards.hpp"

#include <iostream>
#include <memory>

using namespace std;

int main(void) {
    Game g;
    g.hand.add("6", "Hearts");
    g.hand.cards.push_back(Card());
    List cards = g.hand.inspect();
    for (auto p : cards.items) {
        cout << *dynamic_pointer_cast<Card>(p) << endl;
    }
    g.board.play(Card("7", "Spades"));
    List board = g.board.inspect();
    for (auto p : board.items) {
        for (auto card : dynamic_pointer_cast<List>(p)->items) {
            cout << *dynamic_pointer_cast<Card>(card) << endl;
        }
    }
    return 0;
}
