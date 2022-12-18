// Search over actions

#include <unordered_set>

#include "types.hpp"

using namespace std;

void search(
        vector<Action> &actions, 
        vector<Node> &nodes,
        size_t depth=1) {
    // Init signatures
    unordered_set<string> sigs;
    for (size_t i=0; i<nodes.size(); i++) {
        sigs.insert(nodes[i].sig_str());
    }
    for (size_t iter=0; iter<depth; iter++) {
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
                // cout << 'a' << psetsize << endl;
                if (psetsize == 0) break;
            }
            // cout << 'b' << psetsize << endl;
            if (psetsize == 0) 
                continue;
            // Evaluate compatible nodes and store results
            for (size_t j=0; j<psetsize; j++) {
                Node n(actions[i]);
                vector<Node*> pps;
                size_t jj = j;
                for (size_t k=0; k<nargs; k++) {
                    size_t sz = compat_nodes[k].size();
                    size_t kk = jj%sz;
                    jj /= sz;
                    pps.push_back(compat_nodes[k][kk]);
                }
                // Special code for expand-list
                if (actions[i] == expand_list_action) {
                    List lst(pps[0]->res);
                    vector<Object> objs = lst.get_objects();
                    for (size_t k=0; k<objs.size(); k++) {
                        Node n(objs[k], actions[i], vector<Node>{*pps[0]});
                        string sig = n.sig_str();
                        if (sigs.count(sig) == 0) {
                            sigs.insert(sig);
                            newnodes.push_back(n);
                        }
                    }
                    continue;
                }
                // TODO: don't eval duplicates
                Object res = n.eval(pps);
                if (res == nullobj)
                    continue;
                n.res = res;
                for (auto pp = pps.begin(); pp != pps.end(); pp++) {
                    n.parents.push_back(**pp);
                }
                n.sign();
                string sig = n.sig_str();
                if (sigs.count(sig) == 0) {
                    sigs.insert(sig);
                    newnodes.push_back(n);
                }
            }
        }
        for (size_t n=0; n<newnodes.size(); n++) {
            nodes.push_back(newnodes[n]);
        }
    }
    for (size_t n=0; n<nodes.size(); n++)
        nodes[n].print(cout);
}