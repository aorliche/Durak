
#include "types.hpp"
#include "cards.hpp"

#include <iostream>
#include <memory>

using namespace std;

int main(void) {
    Game g;
    Card c = Card("6", "Hearts");
    Hand &h = g.get("hand");
    h.add(c);
    cout << c << endl;
    // g.hand.cards.push_back(Card());
    // List cards = g.hand.inspect();
    // for (auto p : cards.items) {
        // cout << *dynamic_pointer_cast<Card>(p) << endl;
    // }
    // g.board.play(Card("7", "Spades"));
    // List board = g.board.inspect();
    // for (auto p : board.items) {
        // for (auto card : dynamic_pointer_cast<List>(p)->items) {
            // cout << *dynamic_pointer_cast<Card>(card) << endl;
        // }
    // }
    return 0;
}
