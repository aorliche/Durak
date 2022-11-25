
#include <vector>
#include <string>
#include <algorithm>
#include <memory>
#include <cstdlib>
#include <iostream>
#include <map>

using namespace std;

int NO_RANK = -1;

struct Rank : public Object {
    static vector<string> ranks;
    static unordered_map<int,int> rank_map;
    Rank(const string &r = "") : Object("rank") {
        if (r.length() == 0) {
            rank_map[id] = NO_RANK;
            return;
        }
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

int NO_SUIT = -1;

struct Suit : public Object {
    static vector<string> suits;
    static unordered_map<int,int> suit_map;
    Suit(const string &s = "") : Object("suit") {
        if (s.length() == 0) {
            suit_map[id] = NO_RANK;
            return;
        }
        for (size_t i=0; i<suits.size(); i++) {
            if (suits[i] == r) {
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
    static Card no_card;
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

Object Card::no_card;

ostream &operator<<(ostream &os, const Card &c) {
    if (c == Card::no_card) {
        os << "No Card";
    } else {
        os << Rank(c.get("rank")) << " of " << Suit(c.get("suit"));
    }
    return os;
}

struct Pile : public List {
    Pile(const string &type) : List() {
        make("pile");
        make(type);
    }
    Pile(const Object &o) : Object(o.id) {}
    void add(const string &rank, const string &suit) {
        add(Card(rank, suit));
    }
    void add(Card c) {
        List::add(cp);
    }
    Concept get_type() const {
        auto pair = is_map.equal_range(id);
        for (; pair.first != pair.second; pair.first++) {
            if (!pair.first->is("list") && !pair.first->is("pile")) {
                return *pair.first;
            }
        }
        return nullobj;
    }
};

ostream &operator<<(ostream &os, const Pile &p) {
    auto lst = list_map[p.id];
    cout << p.get_type() << ":" << endl;
    for (size_t i=0; i<lst.size(); i++) {
        cout << Card(lst[i]) << endl;
    }
}

struct Board : public Object {
    Board(const Object &o) : Object(o.id) {}
    Board() {
        make("board");
        give(Pile("plays"));
        give(Pile("covers"));
    }
    void cover(Card c1, Card c2) {
        List plays = List(get("plays"));
        List covers = List(get("covers"));
        int idx = plays.index_of(c2);
        if (idx == -1) 
            throw na;
        if (Card(covers[idx]) != Card::no_card)
            throw na;
        covers[idx] = c1;
    }
    void play(shared_ptr<Card> card) {
        List &plays = dynamic_cast<List&>(*get("plays"));
        List &covers = dynamic_cast<List&>(*get("covers"));
        plays.add(card);
        covers.add(Card::no_card);
    }
};

struct Game : public Object {
    Game() : Object("game") {
        give(make_shared<Board>());
        give(make_shared<Pile>("hand"));
    }
};

# define MAKEACTION(a,ret,args) Action a##_action(#a,ret,args,a)

// Composition functions
shared_ptr<Object> higher_rank(List &args, Object &game) {
    Rank &r1 = dynamic_cast<Rank&>(*args[0]->get("rank"));
    Rank &r2 = dynamic_cast<Rank&>(*args[1]->get("rank"));
    if (r1 > r2) {
        return yes;
    } 
    return no;
}
MAKEACTION(higher_rank, "boolean", (vector<string>{"rank", "rank"}));

shared_ptr<Object> beats(List &args, Object &game) {
    Suit &s0 = dynamic_cast<Suit&>(*args[0]->get("suit"));
    Suit &s1 = dynamic_cast<Suit&>(*args[1]->get("suit"));
    Suit &ts = dynamic_cast<Suit&>(*game.get("trump"));
    if (s0 == ts && s1 != ts) {
        return yes;
    } else if (s1 == ts) {
        return no;
    }
    return higher_rank(args, game);
}
MAKEACTION(beats, "boolean", (vector<string>{"card", "card"}));

shared_ptr<Object> cover(List &args, Object &game) {
    Pile &hand = dynamic_cast<Pile&>(*game.get("hand"));
    Board &board = dynamic_cast<Board&>(*game.get("board"));
    hand.remove(args[1]);
    board.cover(args[0], args[1]);
    return nullobj;
}
MAKEACTION(cover, "null", (vector<string>{"card", "card"}));

shared_ptr<Object> get_item(List &args, Object &game) {
    List &lst = dynamic_cast<List&>(*args[0]);
    Number &num = dynamic_cast<Number&>(*args[1]);
    return lst[num.val];
}
MAKEACTION(get_item, "object", (vector<string>{"list", "number"}));

shared_ptr<Object> get_size(List &args, Object &game) {
    List &lst = dynamic_cast<List&>(*args[0]);
    return make_shared<Number>(lst.size());
}
MAKEACTION(get_size, "number", (vector<string>{"list"}));

shared_ptr<Object> randint(List &args, Object &game) {
    Number &num = dynamic_cast<Number&>(*args[0]);
    return make_shared<Number>(rand()%num.val);
}
MAKEACTION(randint, "number", (vector<string>{"number"}));

shared_ptr<Object> do_nothing(List &args, Object &game) {
    return nullobj;
}
MAKEACTION(do_nothing, "null", (vector<string>{}));

shared_ptr<Object> get_trump(List &args, Object &game) {
    return game.get("trump");
}
MAKEACTION(get_trump, "suit", (vector<string>{}));

shared_ptr<Object> get_suit(List &args, Object &game) {
    return args[0]->get("suit");
}
MAKEACTION(get_suit, "suit", (vector<string>{"card"}));

shared_ptr<Object> get_board(List &args, Object &game) {
    return game.get("board");
}
MAKEACTION(get_board, "board", (vector<string>{}));

shared_ptr<Object> get_hand(List &args, Object &game) {
    return game.get("hand");
}
MAKEACTION(get_hand, "hand", (vector<string>{}));
