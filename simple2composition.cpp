#include "base.hpp"
#include "cards.hpp"
#include "search.hpp"

#include <cstdlib>

void add_random_card(vector<ObjectId> &board) {
    string &r = ranks[rand()%ranks.size()];
    string &s = suits[rand()%suits.size()];
    Card c(r,s);
    board.push_back(c.id);
}

int main(void) {
    Game g;
    Card c1("6", "Hearts");
    Card c2("Ace", "Spades");
    Card c3("8", "Diamonds");
    Card c4("10", "Clubs");
    vector<ObjectId> &board = get<4>(g.get("board"));
    board.push_back(c1.id);
    board.push_back(c2.id);
    board.push_back(c3.id);
    for (int i=0; i<150; i++) {
        add_random_card(board);
    }
    
    // cout << Object(get<4>(g.get("board"))[1]) << endl;
    vector<Property> props{"card", "board", g.id};
    vector<Function> fns{get_property, higher_rank, expand_vec};
    vector<Node> nodes;
    unordered_set<Node, node_hasher> sigs;
    for (auto it = props.begin(); it != props.end(); it++) {
        nodes.emplace_back(*it);
    }
    cout << nodes.size() << endl;
    expand(fns, nodes,8, sigs);
    // for (int i=0; i<nodes.size(); i++) {
    //     nodes[i].print(cout);
    // }
    cout << nodes.size() << endl;
}