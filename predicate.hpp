#include "base.hpp"

// Conjunctive
struct Predicate {
    string name;
    vector<Node*> parents;
    vector<bool> weights;
    int nargs;
    Predicate(const string &n) : name(n) {}
    bool eval(const vector<Node*> ps) {
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
    }
};

void match(const vector<Node*> &nodes, bool res, int maxargs = 3) {
    // Get bool-valued nodes
    vector<Node*> bnodes;
    for (int i=0; i<nodes.size(); i++) {
        if (nodes[i]->res.index() == 0) {
            bnodes.push_back(nodes[i]);
        }
    }
}