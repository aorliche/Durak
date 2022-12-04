#include "types.hpp"
#include "cards.hpp"

#include <iostream>

using namespace std;

// Init function
vector<Node> init_nodes() {
    vector<string> concepts{"board", "hand", "played", "trump", "cover"};
    vector<Node> nodes{Game()};
    for (int i=0; i<concepts.size(); i++) {
        nodes.push_back(Concept(concepts[i]));
    }
    return nodes;
}
vector<Action> actions{get_field_action, expand_list_action, higher_rank_action, beats_action};
vector<Node> nodes = init_nodes();

int main(void) {
    Game g;
    Card c0("6", "Hearts");
    Card c1("King", "Spades");
    cout << c0 << endl;
    Pile hand(g.get("hand"));
    hand.add(c0);
    hand.add(c1);
    cout << hand << endl;
    Board board(g.get("board"));
    board.play(c0);
    board.cover(c0, c1, g);
    cout << board << endl;
    try {
        board.cover(c1, c0, g);
    } catch (Object ex) {
        cout << "Got exception " << (ex == na) << endl;
    }
    return 0;
}
