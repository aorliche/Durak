#include "types.hpp"
#include "cards.hpp"
#include "search.hpp"

using namespace std;

// Init function
vector<Node> init_nodes() {
    vector<string> concepts{"board", "hand", "played", "trump", 
        "cover", "rank", "suit", "card"};
    Game g;
    // cout << g.get("trump").get("rank") << endl;
    g.set_trump(Card("Ace", "Diamonds"));
    g.give(Pile("hand"));
    Board b = g.get("board");
    Pile h = g.get("hand");
    h.give(Card("10", "Spades"));
    b.play(Card("6", "Hearts"));
    b.play(Card("7", "Spades"));
    b.play(Card("8", "Clubs"));
    b.play(Card("King", "Diamonds"));
    cout << Object(Card("9", "Clubs").id) << endl;
    vector<Node> nodes{g};
    for (int i=0; i<concepts.size(); i++) {
        nodes.push_back(Concept(concepts[i]));
        //cout << c.id << " " << nodes.back().res.id << endl;
    }
    return nodes;
}

int main(void) {
    // Actions
    vector<Action> actions{get_field_action, expand_list_action, higher_rank_action, beats_action};

    // Working set
    vector<Node> nodes = init_nodes();
    vector<Node> graveyard;
    unordered_set<vector<int>, int_vector_hasher> sigs;

    expand(actions, nodes, 2, sigs);
    expand(actions, nodes, 3, sigs);
    expand(actions, nodes, 2, sigs);

    vector<Node> bnodes = get_concept_matches(nodes, boolean);
    // print_nodes(cout, nodes);
    cout << nodes.size() << endl;
    cout << bnodes.size() << endl;

    Card needle("6", "Hearts");
    vector<Node> cnodes = get_matches(nodes, [&] (const Node &n) {
        try {
            Card c(n.get_res().id);
            return c == needle;
        } catch (Exception) {
            return false;
        }
    });
    cout << cnodes.size() << endl;

    cout << get_object_matches(nodes, Concept("card"), Card("6", "Hearts")).size() << endl;
    cout << get_object_matches(nodes, boolean, yes).size() << endl;
    cout << get_object_matches(nodes, boolean, no).size() << endl;
}
