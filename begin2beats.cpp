#include "types.hpp"
#include "cards.hpp"
#include "search.hpp"

using namespace std;

// Try to teach the beats predicate

Object make_card(const string &inp) {
    string delim = " of ";
    size_t idx = inp.find(delim);
    if (idx == string::npos) {
        return nullobj;
    }
    string rank = inp.substr(0,idx);
    string suit = inp.substr(idx+4);
    try {
        return Card(rank, suit);
    } catch (Exception) {
        return nullobj;
    } 
}

bool eq_cards(const Object &o1, const Object &o2) {
    try {
        return Card(o1) == Card(o2);
    } catch (Exception) {
        return false;
    }
}

ostream &vp_fn(ostream &os, const Object &obj) {
    if (obj.is(boolean)) {
        os << (obj == yes ? "yes" : "no");
    } else if (obj.is("suit")) {
        os << Suit(obj);
    } else if (obj.is("rank")) {
        os << Rank(obj);
    }
    return os;
}

struct node_hasher {
    std::size_t operator()(const Node &n) const {
        return hash<int>{}(n.id);
    }
};

int main(void) {
    // These actions should really be called predicates
    vector<Action> actions{get_field_action, expand_list_action, higher_rank_action, same_suit_action};
    vector<string> concepts{"board", "trump", "rank", "suit", "card"};
    unordered_set<Node, node_hasher> blacklist;
    unordered_set<vector<int>, int_vector_hasher> sigs;

    Game g;
    g.set_trump(Card("6", "Clubs"));
    vector<Node> nodes{g};
    for (int i=0; i<concepts.size(); i++) {
        nodes.push_back(Concept(concepts[i]));
    }

    // Interactive teaching
    while (true) {
        string inp;
        Object c1, c2;
        cout << "Action> ";
        getline(cin, inp);
        if (inp == "expand") {
            cout << "Expanding..." << endl;
            cout << "before: " << nodes.size();
            expand(actions, nodes, 5, sigs);
            cout << " after: " << nodes.size() << endl;
            print_nodes(cout, nodes, vp_fn);
            continue;
        }
redo_c1:
        cout << "Enter card 1: ";
        getline(cin, inp);
        c1 = make_card(inp);
        cout << c1 << endl;
        if (c1 == nullobj or not insert_if_new(nodes, c1, eq_cards)) {
            goto redo_c1;
        }
redo_c2:
        cout << "Enter card 2: ";
        getline(cin, inp);
        c2 = make_card(inp);
        cout << c2 << endl;
        if (c2 == nullobj or not insert_if_new(nodes, c2, eq_cards)) {
            goto redo_c2;
        }
        nodes.push_back(c2);

//         cout << "Expanding..." << endl;
//         cout << "before: " << nodes.size();
//         expand(actions, nodes, 5, sigs);
//         cout << " after: " << nodes.size() << endl;

// redo_beats:
//         cout << "Does card 2 beat card 1? (yes/no) ";
//         getline(cin, inp);
//         if (inp != "yes" and inp != "no") {
//             cout << "Invalid" << endl;
//             goto redo_beats;
//         }
//         // Blacklist the incorrect answers
//         Object wrong = inp == "yes" ? no : yes;
//         auto bool_vec = get_concept_matches(nodes, boolean);
//         auto wrong_vec = get_object_matches(nodes, boolean, wrong);
//         blacklist.insert(wrong_vec.begin(), wrong_vec.end());
//         cout << "bools: " << bool_vec.size() << endl;
//         cout << "blacklist size: " << blacklist.size() << endl;
    }
}
