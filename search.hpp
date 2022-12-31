// Search over actions

#include <unordered_set>
#include <functional>

#include "types.hpp"

using namespace std;

// Boolean results are never used in further calculations
// Except with an aggregation predicate

void expand(
        vector<Action> &actions, 
        vector<Node> &nodes,
        size_t depth,
        unordered_set<vector<int>, int_vector_hasher> &sigs) {
    for (size_t iter=0; iter<depth; iter++) {
        // NOTE! You must have a separate newnodes for each iteration
        vector<Node> newnodes;
        for (size_t i=0; i<actions.size(); i++) {
            size_t nargs = actions[i].get_args_size();
            // Find compatible nodes for jth arg
            vector<vector<int>> compat_nodes(nargs);
            size_t psetsize = 1;
            for (size_t j=0; j<nargs; j++) {
                for (size_t k=0; k<nodes.size(); k++) {
                    if (actions[i].is_compatible_arg(j, 
                        nodes[k].get_res())) {
                        compat_nodes[j].push_back(nodes[k].id);
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
                vector<int> &parents = node_parent_map[n.id];
                size_t jj = j;
                for (size_t k=0; k<nargs; k++) {
                    size_t sz = compat_nodes[k].size();
                    size_t kk = jj%sz;
                    jj /= sz;
                    parents.push_back(compat_nodes[k][kk]);
                }
                // Skip shallow result
                vector<int> shallow = n.get_shallow();
                if (sigs.count(shallow) > 0) {
                    n.remove();
                    continue;
                } else {
                    sigs.insert(shallow);
                }
                // Special code for expand-list
                if (actions[i] == expand_list_action) {
                    List lst(Node(parents[0]).get_res());
                    vector<Object> objs = lst.get_objects();
                    for (size_t k=0; k<objs.size(); k++) {
                        Node n(objs[k], actions[i], parents);
                        newnodes.push_back(n);
                    }
                    continue;
                }
                Object res = n.eval();
                // Delete null results
                if (res == nullobj) {
                    n.remove();
                    continue;
                }
                // No duplicates
                node_res_map[n.id] = res.id;
                n.sign();
                newnodes.push_back(n);
            }
        }
        // cout << "Iteration " << iter << " added " << newnodes.size() << endl;
        // cout << "Sigsize: " << sigs.size() << endl;
        for (size_t n=0; n<newnodes.size(); n++) {
            // newnodes[n].print(cout);
            nodes.push_back(newnodes[n]);
        }
    }
    // for (size_t n=0; n<nodes.size(); n++)
    //     nodes[n].print(cout);
}

// Returns true if successfully inserted
// false if object exists
bool insert_if_new(vector<Node> &nodes, const Object &obj, const function<bool(const Object&, const Object&)>& eqpred) {
    for (size_t i=0; i<nodes.size(); i++) {
        if (eqpred(nodes[i].get_res(), obj)) {
            return false;
        }
    }
    nodes.push_back(obj);
    return true;
}

vector<Node> get_matches(vector<Node> &nodes, const function<bool(const Node&)>& pred) {
    vector<Node> res;
    for (size_t i=0; i<nodes.size(); i++) {
        if (pred(nodes[i])) {
            res.push_back(nodes[i]);
        }
    }
    return res;
}

vector<Node> get_concept_matches(vector<Node> &nodes, Concept c) {
    return get_matches(nodes, [&] (const Node &n) {
        return n.get_res().is(c);
    });
}

template <typename T>
vector<Node> get_object_matches(vector<Node> &nodes, Concept c, T t) {
    return get_matches(nodes, [&] (const Node &n) {
        Object obj = n.get_res();
        if (!obj.is(c)) {
            return false;
        }
        return T(obj) == t;
    });
}