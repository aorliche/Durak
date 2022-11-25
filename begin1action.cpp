#include "types.hpp"
#include "cards.hpp"

#include <iostream>
#include <memory>
#include <set>

using namespace std;

vector<Actions> base{
    beats_action,
    get_item_action,
    get_size_action,
    randint_action,
    get_trump_action
    get_hand_action
};

struct ActionTreeNode {
    vector<ActionTreeNode *> children;
    vector<shared_ptr<Object
    ActionTreeNode() {
        
    }
};    

template<>
struct less<shared_ptr<Object>> {
    constexpr bool operator()(const shared_ptr<Object> &lhs, const shared_ptr<Object> &rhs) {
        return lhs->id < rhs->id;
    }
}

set<Object> eval(set<shared_ptr<Object>> input) {
    set<shared_ptr<Object>> output;
    for (auto action : base) {
        for (auto arg : action.arg_types) {
            for (auto p : input) {
                if (p->is(arg)) {
                    
                }
            }
    }
}

int main(void) {

}