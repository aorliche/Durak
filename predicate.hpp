#include <algorithm>

#include "base.hpp"

using namespace std;

// Conjunctive
struct Predicate {
    string name;
    vector<Function> fns;
    vector<vector<int>> arg_binds;
    vector<int> set_idcs;
    vector<bool> neg;
    Predicate(const string &n) : name(n) {}
    /*bool eval(const vector<Node*> ps) {
        for (int i=0; i<ps.size(); i++) {
            bool val = get<0>(ps[i]->res);
            if (weights[i] ^ val) {
                return false;
            }
        }
    }
    bool allow(const vector<Node*> ps) {
        for (int i=0; i<ps.size(); i++) {
            if (ps[i]->res.index() != 0) {
                return false;
            }
        }
        return true;
    }
    void update_parents(vector<Node*> ps) {
        parents.clear();
        for (int i=0; i<ps.size(); i++) {
            parents.push_back(*ps[i]);
        }
    }*/
    bool eval(const vector<Property> &args) {
        vector<bool> res;
        for (int i=0; i<fns.size(); i++) {
            vector<Property> bound_args;
            for (int j=0; j<arg_binds[i].size(); i++) {
                int k = arg_binds[i][j];
                bound_args.push_back(args[k]);
            }
            bool bonly = get<bool>(fns[i](bound_args));
            if (neg[i]) {
                bonly = !bonly;
            }
            res.push_back(bonly);
        }
        int midx = *max_element(set_idcs.begin(), set_idcs.end());
        bool *conjs = new bool[midx+1];
        for (int i=0; i<set_idcs.size(); i++) {
            conjs[set_idcs[i]] = res[i] and conjs[set_idcs[i]];
        }
        bool disj = false;
        for (int i=0; i<midx; i++) {
            if (conjs[i]) {
                disj = true;
                break;
            }
        }
        delete conjs;
        return disj;
    }
};

Predicate *fit(bool res, const vector<Property> &args, const vector<Node *> &lib, Predicate &prev) {
    bool prev_res = prev.eval(args);
    if (res == prev_res) {
        return &prev;
    }
}

// Get bool nodes consistent with existing predicate
void get_bool_nodes(const vector<Node*> &nodes, const Property &pred) {
    vector<Node*> bnodes;
    for (int i=0; i<nodes.size(); i++) {
        if (nodes[i]->res.index() == 0) {
            bnodes.push_back(nodes[i]);
        }
    }
}

void match_bool(bool res, const vector<Node*> &nodes, const Property &pred) {
    // Get bool-valued nodes
    

}