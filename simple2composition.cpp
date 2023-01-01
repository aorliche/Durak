#include "base.hpp"

struct Game : Object {
    Game() : Object() {
        make("game");
        set("board", vector<ObjectId>());
    };
};

int main(void) {
    Game g;
    Card c1("6", "Hearts");
    Card c2("Ace", "Spades");
    vector<ObjectId> &board = get<4>(g.get("board"));
    board.push_back(c1.id);
    board.push_back(c2.id);
    cout << Object(get<4>(g.get("board"))[1]) << endl;
    vector<Property> props{"card", "board", g.id};
    vector<Function> fns{get_property, higher_rank};
    vector<NodeId> nodes;
    unordered_set<NodeId, node_hasher> sigs;
    for (auto it = props.begin(); it != props.end(); it++) {
        Node n(*it);
        node_map[n.id] = n;
        nodes.push_back(n.id);
        // cout << 'b' << Node::get(nodes[0]).res.index() << endl;
    }
    cout << nodes.size() << endl;
    // for (int i=0; i<nodes.size(); i++) {
        // cout << 'a' << Node::get(nodes[i]).res.index() << endl;
    // }
    expand(fns, nodes, 5, sigs);
    cout << nodes.size() << endl;
}