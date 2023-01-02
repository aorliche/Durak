#include "base.hpp"

void expand(
    vector<Function> &fns, 
    vector<Node> &nodes, 
    int depth,
    unordered_set<Node, node_hasher> &sigs) 
{
    for (int iter=0; iter<depth; iter++) {
        // NOTE! Must have a separate newnodes for each iteration
        // Otherwise an infinite loop might happen?
        vector<Node> newnodes;
        for (size_t i=0; i<fns.size(); i++) {
            int nargs = fns[i].nargs;
            // Find compatible nodes for jth arg
            // Use pointers in inner loop only
            vector<vector<Node *>> compat_node_pps(nargs);
            size_t psetsize = 1;
            for (size_t j=0; j<nargs; j++) {
                for (size_t k=0; k<nodes.size(); k++) {
                    if (fns[i].compatible(nodes[k].res, j)) {
                        compat_node_pps[j].push_back(&nodes[k]);
                    }
                }
                psetsize *= compat_node_pps[j].size();
                // cout << 'a' << psetsize << endl;
                if (psetsize == 0) break;
            }
            // cout << 'b' << psetsize << endl;
            if (psetsize == 0) 
                continue;
            // Evaluate compatible nodes and store results
            for (size_t j=0; j<psetsize; j++) {
                vector<Node*> node_pps;
                size_t jj = j;
                for (size_t k=0; k<nargs; k++) {
                    size_t sz = compat_node_pps[k].size();
                    size_t kk = jj%sz;
                    jj /= sz;
                    node_pps.push_back(compat_node_pps[k][kk]);
                }
                Node n(fns[i]);
                Property res = n.eval(node_pps);
                // Skip null result
                if (res == none) 
                    continue;
                // Bake node
                n.res = res;
                n.update_parents(node_pps);
                // Skip shallow result
                if (sigs.count(n) > 0) {
                    continue;
                } else {
                    sigs.insert(n);
                }
                // Special code for get_property(vector<ObjectId>)
                if (n.fn.name == "expand_vec") {
                    vector<ObjectId> &objs = get<4>(res);
                    for (int i=0; i<objs.size(); i++) {
                        Node nn(n.fn);
                        nn.res = Property(objs[i]);
                        nn.update_parents(node_pps);
                        newnodes.push_back(nn);
                    }
                } else {
                    newnodes.push_back(n);
                }
            }
        }
        // cout << "Iteration " << iter << " added " << newnodes.size() << endl;
        cout << "Sigsize: " << sigs.size() << endl;
        // cout << "nodes size: " << nodes.size() << endl;
        for (size_t n=0; n<newnodes.size(); n++) {
            // newnodes[n].print(cout);
            nodes.push_back(newnodes[n]);
        }
    }
    // for (size_t n=0; n<nodes.size(); n++)
    //     nodes[n].print(cout);
}