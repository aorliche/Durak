#include "base.hpp"

// Conjunctive
struct Predicate {
    string name;
    vector<Node> parents;
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

void make_predicate(vector<vector<Node>> &train, vector<bool> &res, int maxargs = 3) {
    
}