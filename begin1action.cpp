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
        //cout << c.id << " " << nodes.back().res.id << endl;
    }
    return nodes;
}
vector<Action> actions{get_field_action, expand_list_action, higher_rank_action, beats_action};
vector<Node> nodes = init_nodes();

int main(void) {
    vector<Node> newnodes;
    for (size_t i=0; i<actions.size(); i++) {
        size_t nargs = actions[i].get_args_size();
        // Find compatible nodes for jth arg
        vector<vector<Node*>> compat_nodes(nargs);
        size_t psetsize = 1;
        for (size_t j=0; j<nargs; j++) {
            for (size_t k=0; k<nodes.size(); k++) {
                if (actions[i].is_compatible_arg(j, nodes[k].res)) {
                    compat_nodes[j].push_back(&nodes[k]);
                }
            }
            psetsize *= compat_nodes[j].size();
            cout << 'a' << psetsize << endl;
            if (psetsize == 0) break;
        }
        cout << 'b' << psetsize << endl;
        if (psetsize == 0) 
            continue;
        for (size_t j=0; j<psetsize; j++) {
            Node n(actions[i]);
            int jj = j;
            for (size_t k=0; k<nargs; k++) {
                size_t sz = compat_nodes[k].size();
                size_t kk = jj%sz;
                jj /= sz;
                n.parents.push_back(compat_nodes[k][kk]);
            }
            // TODO: don't eval duplicates
            n.eval();
            newnodes.push_back(n);
        }
    }
    for (auto node : newnodes) {
        node.print(cout);
    }
}