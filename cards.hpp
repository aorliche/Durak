
#include <vector>
#include <string>
#include <algorithm>
#include <memory>
#include <cstdlib>
#include <iostream>
#include <map>

using namespace std;

struct Rank : public Object {
    static vector<string> ranks;
    static unordered_map<int,int> rank_map;
    Rank(const string &r = "") : Object("rank") {
        for (size_t i=0; i<ranks.size(); i++) {
            if (ranks[i] == r) {
                rank_map[id] = i;
                return;
            }
        }
    }
    Rank(const Object &o) : Object(o.id) {}
    bool operator==(const Rank &other) const {
        return rank_map[id] == rank_map[other.id];
    }
    bool operator!=(const Rank &other) const {
        return !(*this == other);
    }
    bool operator>(const Rank &other) const {
        return rank_map[id] > rank_map[other.id];
    }
    bool operator<(const Rank &other) const {
        return rank_map[id] < rank_map[other.id];
    }
};

vector<string> Rank::ranks = {"6", "7", "8", "9", "10", "Jack", "Queen", "King", "Ace"};
unordered_map<int,int> Rank::rank_map;

ostream & operator<<(ostream &os, const Rank &r) {
    return os << Rank::ranks[Rank::rank_map[r.id]];
}

struct Suit : public Object {
    static vector<string> suits;
    static unordered_map<int,int> suit_map;
    Suit(const string &s) : Object("suit") {
        for (size_t i=0; i<suits.size(); i++) {
            if (suits[i] == s) {
                suit_map[id] = i;
                return;
            }
        }
    }
    Suit(const Object &o) : Object(o.id) {}
    bool operator==(const Suit &other) const {
        return suit_map[id] == suit_map[other.id];
    }
    bool operator!=(const Suit &other) const {
        return !(*this == other);
    }
};

vector<string> Suit::suits = {"Hearts", "Diamonds", "Spades", "Clubs"};
unordered_map<int,int> Suit::suit_map;

ostream & operator<<(ostream &os, const Suit &s) {
    return os << Suit::suits[Suit::suit_map[s.id]];
}

struct Card : public Object {
    Card(const string &rank = "", const string &suit = "") : Object("card") {
        give(Rank(rank));
        give(Suit(suit));
    }
    Card(const Object &o) : Object(o.id) {}
    // Overload equals operator
    virtual bool operator==(const Card &other) const {
        return get("suit") == other.get("suit") && get("rank") == other.get("rank");
    }
    virtual bool operator!=(const Card &other) const {
        return !(*this == other);
    }
    // To string
    friend ostream& operator<<(ostream &os, shared_ptr<Card> c);
};

struct Cover : Card {
    Cover(const Object &o) : Card(o) {
        make("cover");
    }
};

ostream &operator<<(ostream &os, const Card &c) {
    os << Rank(c.get("rank")) << " of " << Suit(c.get("suit"));
    if (c.get("cover") != nullobj) {
        Cover cov(c.get("cover"));
        os << " covered by: " << cov;
    }
    return os;
}

struct Pile : public List {
    Pile(const string &type) : List() {
        make("pile");
        make(type);
    }
    Pile(const Object &o) : List(o.id) {}
    void add(const string &rank, const string &suit) {
        add(Card(rank, suit));
    }
    void add(Card c) {
        List::add(c);
    }
    string get_type() const {
        auto range = is_map.equal_range(id);
        for (auto it = range.first; it != range.second; it++) {
            string &name = name_map[it->second];
            if (name != "object" and name != "pile" and name != "list") {
                return name;
            }
        }
        return "ERROR";
    }
};

ostream &operator<<(ostream &os, const Pile &p) {
    auto lst = list_map[p.id];
    os << p.get_type() << ":";
    for (size_t i=0; i<lst.size(); i++) {
        os << endl << Card(lst[i]);
    }
    return os;
}

struct Game : public Object {
    Game();
    Game(const Object &obj) : Object(obj.id) {}
};

Object beats(const vector<Object> &args);

struct Board : public Pile {
    Board(const Object &o) : Pile(o) {}
    Board() : Pile("board") {}
    void cover(Card c1, Card c2, Game g) {
        int idx = index_of(c1);
        if (idx == -1) 
            throw na;
        if (c1.get("cover") != nullobj)
            throw na;
        if (beats(vector<Object>{c1, c2, g}) != yes) 
            throw na;
        c1.give(Cover(c2));
    }
    void play(Card c) {
        add(c);
    }
};

Game::Game() : Object("game") {
    give(Board());
    give(Pile("hand"));
    give(Card("Ace", "Spades"));
    get("card").make("trump");
    get("trump").unmake("card");
}

// Composition functions
Object higher_rank(const vector<Object> &args) {
    Rank r1(args[0].get("rank"));
    Rank r2(args[1].get("rank"));
    if (r1 > r2) {
        return yes;
    } 
    return no;
}
MAKEACTION(higher_rank, "bool", (vector<string>{"rank", "rank"}));

Object beats(const vector<Object> &args) {
    Suit s0(args[0].get("suit"));
    Suit s1(args[1].get("suit"));
    Suit ts(args[2].get("trump"));
    if (s0 == ts && s1 != ts) {
        return yes;
    } else if (s1 == ts) {
        return no;
    }
    return higher_rank(args);
}
MAKEACTION(beats, "bool", (vector<string>{"card", "card", "game"}));

// Context-changing functions
/*Object cover(vector<Object> &args, Object &game) {
    Pile hand(game.get("hand"));
    Board board(game.get("board"));
    hand.remove(args[1]);
    board.cover(args[0], args[1]);
    return nullobj;
}
MAKEACTION(cover, "null", (vector<string>{"card", "card"}));*/

/*
shared_ptr<Object> get_trump(List &args, Object &game) {
    return game.get("trump");
}
MAKEACTION(get_trump, "suit", (vector<string>{}));

shared_ptr<Object> get_suit(List &args, Object &game) {
    return args[0]->get("suit");
}
MAKEACTION(get_suit, "suit", (vector<string>{"card"}));
*/
