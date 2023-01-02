#include "base.hpp"
#include "cards.hpp"
#include "search.hpp"

int main(void) {
    Game g;
    vector<Function> fns{get_property, expand_vec, higher_rank, same_suit};
    vector<Node*> nodes;
    unordered_set<size_t> sigs;
    nodes.push_back(new Node(g.id));
    nodes.push_back(new Node("cards"));
    nodes.push_back(new Node("players"));
    nodes.push_back(new Node("trump"));

    cout << nodes.size() << endl;
    expand(fns, nodes, 10, sigs);
    // for (int i=0; i<nodes.size(); i++) {
    //     nodes[i].print(cout);
    // }
    cout << nodes.size() << endl;
}