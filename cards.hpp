
#include <vector>
#include <string>
#include <algorithm>
#include <memory>
#include <cstdlib>
#include <iostream>

using namespace std;

struct Rank : public Object {
    string rank;
    Rank(const string &r) : Object("rank"), rank(r) {}
    bool operator==(const Rank &other) const {
        return rank == other.rank;
    }
    bool operator!=(const Rank &other) const {
        return !(*this == other);
    }
    bool operator>(const Rank &other) const;
    bool operator<(const Rank &other) const;
}

struct Suit : public Object {
    string suit;
    Rank(const string &s) : Object("suit"), suit(s) {}
    bool operator==(const Suit &other) const {
        return suit == other.suit;
    }
    bool operator!=(const Suit &other) const {
        return !(*this == other);
    }
}

vector<Rank> ranks = {
    Rank("6"),
    Rank("7"),
    Rank("8"),
    Rank("9"),
    Rank("10"),
    Rank("Jack"),
    Rank("Queen"),
    Rank("King"),
    Rank("Ace")
};

bool Rank::operator>(const Rank &other) {
    for (auto it = ranks.begin(); it != ranks.end(); it++) {
        if (*it == *this && *it != other) {
            return false;
        } else if (*it == *this && *it == other) {
            return false;
        } else if (*it == other) {
            return true;
        }
    }
}

bool Rank::operator<(const Rank &other) {
    return !(*this > other) && this != other;
}

vector<Suit> suits = {
    Suit("Hearts"),
    Suit("Diamonds"),
    Suit("Spades"),
    Suit("Clubs")
};

struct Card : public Object {
    static Object no_rank("rank");
    static Object no_suit("suit");
    static Object no_card("card");
    Card(const Object &rank = no_rank, const Object &suit = no_suit) : Object("card") {
        this->give(rank);
        this->give(suit);
    }
    Card(const string &r, const string &s) : Card(Rank(r), Suit(s)) {}
    // Overload equals operator
    virtual bool operator==(const Card &other) const {
        return this.get("suit") == other.get("suit") 
            && this.get("rank") == other.get("rank");
    }
    virtual bool operator!=(const Card &other) const {
        return !(*this == other);
    }
    // To string
    friend ostream& operator<<(ostream &os, const Card &c);
};

ostream &operator<<(ostream &os, const Card &c) {
    if (c.get("rank") == no_rank && c.get("suit) == no_suit) {
        os << "NoCard";
    } else {
        os << c.get("rank") << " of " << c.get("suit");
    }
    return os;
}

struct Pile : public List {
    Pile(const string &type) : Object(type) {
        this->make("pile");
    }
    void add(const string &rank, const string &suit) {
        this->give(make_shared<Card>(rank, suit));
    }
};

struct Board : public Object {
    Board() {
        this->is("board");
        this->give(make_shared<Pile>("plays"));
        this->give(make_shared<Pile>("covers"));
    }
    void cover(const Object &card1, const Object &card2) {
        List &plays = get("plays");
        List &covers = get("covers");
        int idx = plays.index_of(card2);
        if (idx == -1) 
            throw na;
        if (covers[idx] != Card::no_card)
            throw na;
        covers[idx] = card1;
    }
    void play(const Object &card) {
        get("plays").add(&card)
        get("covers").add(Card::no_card);
    }
};

struct Game : public Object {
    Game() : Object("game") {
        give(make_shared<Board>());
        give(make_shared<Hand>());
    }
};

# define MAKEACTION(a) Action a##_act(#a, a)

// Composition functions
Object *higher_rank(List &args, Game &game) {
    Rank &r1 = dynamic_cast<Rank&>(args[0]->get("rank"));
    Rank &r2 = dynamic_cast<Rank&>(args[1]->get("rank"));
    if (r1 > r2) {
        return &yes;
    } 
    return &no;
}
MAKEACTION(higher_rank);

Object *beats(List &args, Object &game) {
    Suit &s0 = args[0]->get("suit");
    Suit &s1 = args[1]->get("suit");
    Suit &ts = game.get("trump").get("suit");
    if (s0 == ts && s1 != ts) {
        return &yes;
    } else if (s1 == ts) {
        return &no;
    }
    return higher_rank(args, game);
}
MAKEACTION(beats);

Object *cover(List &args, Game &game) {
    Card &c1 = dynamic_cast<Card&>(args[0]);
    Card &c2 = dynamic_cast<Card&>(args[1]);
    Hand &h = dynamic_cast<Hand&>(game.get("hand"));
    Board &b = dynamic_cast<Board&>(board.get("board"));
    h.remove(c1);
    board.cover(c1, c2);
    return &null;
}
MAKEACTION(cover);

Object *get_item(List &args, Game &game) {
    List &lst = dynamic_cast<List&>(args[0]);
    Number &num = dynamic_cast<Number&>(args[1]);
    return &lst[num.val]
}
MAKEACTION(get_item);

Object *get_size(List &args, Game &game) {
    List &lst = dynamic_cast<List&>(args[0]);
    return make_shared<Number>(lst.size());
}
MAKEACTION(get_size);

Object *randint(List &args, Game &game) {
    Number &num = dynamic_cast<Number&>(args[0]);
    return make_shared<Number>(rand()%num.val);
}
MAKEACTION(randint);

Object *do_nothing(List &args, Game &game) {
    return &null;
}
MAKEACTION(do_nothing);

Object *get_trump(List &args, Game &game) {
    return &game.get("trump");
}
MAKEACTION(get_trump);

Object *get_suit(List &args, Game &game) {
    return &args[0].get("suit");
}
MAKEACTION(get_suit);

Object *get_board(vector<Object> &args, Game &game) {
    return &game.get("board");
}

Object get_hand(vector<Object> &args, Game &game) {
    return &game.get("hand");
}
MAKEACTION(get_hand);
