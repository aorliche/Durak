#include "base.hpp"

struct Card : Object {
    Card(const ObjectId &obj) : Object(obj) {}
    Card(const string &rank, const string &suit) {
        make("card");
        set("rank", rank);
        set("suit", suit);
    }
};

vector<string> ranks{"6","7","8","9","10","Jack","Queen","King","Ace"};
vector<string> suits{"Hearts", "Clubs", "Diamonds", "Spades"};

vector<ObjectId> get_all_cards() {
    vector<ObjectId> vec;
    for (int i=0; i<ranks.size(); i++) {
        for (int j=0; j<suits.size(); j++) {
            Card c(ranks[i], suits[j]);
            vec.push_back(c.id);
        }
    }
    return vec;
}

struct Player : Object {
    Player() : Object() {
        make("player");
        set("hand", vector<ObjectId>());
    }
};

vector<ObjectId> make_players(int n) {
    vector<ObjectId> players;
    for (int i=0; i<n; i++) {
        Player p;
        players.push_back(p.id);
    }
    return players;
}

struct Game : Object {
    Game() : Object() {
        make("game");
        set("board", vector<ObjectId>());
        set("cards", get_all_cards());
        set("players", make_players(2));
        set("trump", std::get<4>(get("cards"))[rand()%36]);
    };
};

bool higher_rank_eval(const vector<Property> &objs) {
    return index_of(ranks, get<string>(get_helper(objs[0], "rank")))
        > index_of(ranks, get<string>(get_helper(objs[1], "rank")));
}

bool same_suit_eval(const vector<Property> &objs)  {
    return get_helper(objs[0], "suit") == get_helper(objs[1], "suit");
}

bool two_cards_allow(const Property &p, int n) {
    return n < 2 and p.index() == 3 and Object(get<ObjectId>(p)).is("card");
}

bool always_allow(const Property &p, int n) {
    return true;
}

Function higher_rank("higher_rank", higher_rank_eval, two_cards_allow, 2);
Function same_suit("same_suit", same_suit_eval, two_cards_allow, 2);