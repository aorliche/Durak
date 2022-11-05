
#include <vector>
#include <string>
#include <algorithm>
#include <memory>
#include <cstdlib>
#include <iostream>

using namespace std;

vector<Object*> ranks;
vector<Object*> suits;

struct Card : public Object {
    Card(const Object &rank = null, const Object &suit = null) : Object() {
        this->make("card");
        this->give(rank);
        this->give(suit);
    }
    Card(const string &r, const string &s) : Card(Concept(r), Concept(s)) {}
    // Overload equals operator
    virtual bool operator==(const Card &other) const {
        return this->get("suit") == other->get("suit") 
            && this->get("rank") == other->get("rank");
    }
    virtual bool operator!=(const Card &other) const {
        return !(*this == other);
    }
    // To string
    friend ostream& operator<<(ostream &os, const Card &c);
};

ostream &operator<<(ostream &os, const Card &c) {
    if (c.rank == null && c.suit == null) {
        os << "NoCard";
    } else {
        os << c.rank << " of " << c.suit;
    }
    return os;
}

// Convenience object, other distinct objects may be no-cards
Card nocard;

struct Number : public Object {
    int val;
    Number() : val(0) {
        this->is("number");
    }
    bool operator==(const Number &other) const {
        return val == other.val;
    }
    bool operator!=(const Number &other) const {
        return val != other.val;
    }
    bool operator>(const Number &other) const {
        return val > other.val;
    }
    bool operator<(const Number &other) const {
        return val < other.val;
    }
};

int indexOf(const vector<Object*> &haystack, const Object *needle) {
    auto it = find(haystack.begin(), haystack.end(), *needle);
    return (it == haystack.end()) ? -1 : distance(haystack.objects.begin(), it);
}

struct Pile : public List {
    Pile(const string &type) : Object(type) {
        this->make("pile");
    }
    void add(const string &rank, const string &suit) {
        this->give(make_shared<Card>(rank, suit));
    }
    void remove(Card &card) {
        this->remove(card);
    }
};

struct Board : public Object {
    Board() {
        this->is("board");
        this->give(make_shared<Pile>("plays"));
        this->give(make_shared<Pile>("covers"));
    }
    void cover(const Object *card1, const Object *card2) {
        List *plays = this->get("plays")->[0]->get("card");
        List *covers = this->get("covers")->[0]->get("card");
        int idx = indexOf(plays->objects, card2);
        if (idx == -1) 
            throw na;
        if (*(covers->objects[idx]) != nocard)
            throw na;
        covers->objects[idx] = card1;
    }
    void play(const Object *card) {
        this->get("plays")->objects.push_back(card);
        this->get("covers")->objects.push_back(nocard);
    }
};

struct Game : public Object {
    Game() : Object("game") {
        this->give(make_shared<Board>());
        this->give(make_shared<Hand>());
    }
};

Object higherRank(const Object *card1, const Object *card2) {
    int i1 = indexOf(ranks, card1->get("rank")->[0]); 
    int i2 = indexOf(ranks, card2->get("rank")->[0]); 
    return i1 > i2 ? yes : no;
}

Object beats(vector<Object*> &args, Object *game) {
    if (args.size() != 2) 
        throw na;
    if (args[0]->get("suit") == game->get2("trump", "suit")
            && args[1]->get("suit") != game->get2("trump", "suit"))
        return yes;
    else if (getSuit(card2) == game.trump) 
        return no;
    else 
        return higherRank(card1, card2);
    // Never returned
    return null;
}

Object cover(vector<Object> &args, Game &game) {
    if (args.size() != 2)
        throw na;
    Card &card1 = dynamic_cast<Card&>(args[0]);
    Card &card2 = dynamic_cast<Card&>(args[1]);
    game.hand.remove(card1);
    game.board.cover(card1, card2);
    return null;
}

Object getItem(vector<Object> &args, Game &game) {
    if (args.size() != 2) 
        throw na;
    List &lst = dynamic_cast<List&>(args[0]);
    Number &num = dynamic_cast<Number&>(args[1]);
    return *lst.items[num.val];
}

Object getSize(vector<Object> &args, Game &game) {
    if (args.size() != 1) 
        throw na;
    List &lst = dynamic_cast<List&>(args[0]);
    return Number(lst.items.size());
}

Object randInt(vector<Object> &args, Game &game) {
    if (args.size() != 1) 
        throw na;
    Number &num = dynamic_cast<Number&>(args[0]);
    return Number(rand()%num.val);
}

Object doNothing(vector<Object> &args, Game &game) {
    if (args.size() != 0) 
        throw na;
    return null;
}

Object getTrump(vector<Object> &args, Game &game) {
    if (args.size() != 0) 
        throw na;
    return ConceptWrap(game.trump);
}

Object getBoard(vector<Object> &args, Game &game) {
    if (args.size() != 0) 
        throw na;
    return game.board;
}

Object getHand(vector<Object> &args, Game &game) {
    if (args.size() != 0) 
        throw na;
    return game.hand;
}

Action beatsAction("beats", beats);
Action coverAction("cover", cover);
Action getItemAction("getItem", getItem);
Action randIntAction("randInt", randInt);
Action getTrumpAction("getTrump", getTrump);
Action getBoardAction("getBoard", getBoard);
Action getHandAction("getHand", getHand);
Action doNothingAction("doNothing", doNothing);
