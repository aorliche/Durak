#include "base.hpp"
#include "cards.hpp"
#include "search.hpp"

int main(void) {
    Game g;
    vector<Function> fns{get_property, expand_vec, higher_rank, same_suit};
    vector<Node> nodes;
    unordered_set<Node, node_hasher> sigs;
    nodes.emplace_back(g.id);
    nodes.emplace_back("cards");
    nodes.emplace_back("players");
    nodes.emplace_back("trump");

    cout << nodes.size() << endl;
    expand(fns, nodes, 4, sigs);
    // for (int i=0; i<nodes.size(); i++) {
    //     nodes[i].print(cout);
    // }
    cout << nodes.size() << endl;
}