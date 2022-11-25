
#include "types.hpp"
#include "cards.hpp"

#include <iostream>
#include <memory>

using namespace std;

int main(void) {
    Game g;
    shared_ptr<Card> c = make_shared<Card>("6", "Hearts");
    Pile &hand = dynamic_cast<Pile&>(*g.get("hand"));
    hand.add(c);
    hand.add(make_shared<Card>("Ace", "Spades"));
    cout << *c << endl;
    for (auto p : hand.get_objects()) {
        cout << dynamic_cast<Card&>(*p) << endl;
    }
    Board &board = dynamic_cast<Board&>(*g.get("board"));
    board.play(make_shared<Card>("7", "Spades"));
    for (auto cov : dynamic_cast<Pile&>(*board.get("plays")).get_objects()) {
        cout << dynamic_cast<Card&>(*cov) << endl;
    }
    try{
    for (auto cov : dynamic_cast<Pile&>(*board.get("covers")).get_objects()) {
        cout << dynamic_cast<Card&>(*cov) << endl;
    }
    }catch(bad_cast b) {
        cout << "error occurred" << endl;
    }
    return 0;
}
