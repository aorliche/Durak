#include "types.hpp"
#include "cards.hpp"
#include "search.hpp"

using namespace std;

// Init function
vector<Node> init_nodes() {
    vector<string> concepts{"board", "hand", "played", "trump", "cover", "rank", "suit", "card"};
    Game g;
    cout << g.get("trump").get("rank") << endl;
    Board(g.get("board")).play(Card("6", "hearts"));
    vector<Node> nodes{g};
    for (int i=0; i<concepts.size(); i++) {
        nodes.push_back(Concept(concepts[i]));
        //cout << c.id << " " << nodes.back().res.id << endl;
    }
    return nodes;
}

// Actions
vector<Action> actions{get_field_action, expand_list_action, higher_rank_action, beats_action};

// Working set
vector<Node> nodes = init_nodes();
vector<Node> invalid_nodes;

int main(void) {
    search(actions, nodes, 3);
}