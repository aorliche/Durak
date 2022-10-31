
#include <vector>
#include <string>
#include <algorithm>
#include <memory>
#include <cstdlib>
#include <iostream>

using namespace std;

Concept card("Card");
Concept suit("Suit");
Concept rank("Rank");
vector<Concept> suits{
    Concept("Diamonds"), 
    Concept("Hearts"), 
    Concept("Clubs"), 
    Concept("Spades")};
vector<Concept> ranks{
    Concept("6"), 
    Concept("7"), 
    Concept("8"), 
    Concept("9"), 
    Concept("10"), 
    Concept("Jack"), 
    Concept("Queen"), 
    Concept("King"), 
    Concept("Ace")};
Concept number("Number");
Concept board("Board");
Concept hand("Hand");
Concept game("Game");
Concept player("Player");

struct Card : public Object {
    Concept suit;
    Concept rank;
    Card(const Concept &r = null, const Concept &s = null) : Object("Card"), rank(r), suit(s) {
        relations.push_back(Relation("is", *this, card));
    }
    Card(const string &r, const string &s) : Card(Concept(r), Concept(s)) {}
    Card(const Card &card) : Card(card.rank, card.suit) {}
    virtual bool operator==(const Card &other) const {
        return suit == other.suit && rank == other.rank;
    }
    virtual bool operator!=(const Card &other) const {
        return !(*this == other);
    }
    Card &operator=(const Card &other) {
        suit = other.suit;
        rank = other.rank;
        return *this;
    }
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

Card NoCard;

struct Number : public Object {
    int val;
    Number(int v) : Object("Number"), val(v) {
        Relation r = Relation("is", *this, number);
        relations.push_back(r);
    }
    bool operator==(const Number &other) const {
        return val == other.val;
    }
    bool operator>(const Number &other) const {
        return val > other.val;
    }
    bool operator<(const Number &other) const {
        return val < other.val;
    }
};

Concept getSuit(string &s) {
    for (int i=0; i<suits.size(); i++) {
        if (suits[i].name == s) return suits[i];
    }
    throw na;
}

Concept getSuit(Card &card) {
    return card.suit;
}

Concept getRank(string &r) {
    for (int i=0; i<ranks.size(); i++) {
        if (ranks[i].name == r) return ranks[i];
    }
    throw na;
}

Concept getRank(Card &card) {
    return card.rank;
}

template <typename T>
int indexOf(const vector<T> &haystack, const T &needle) {
    auto it = find(haystack.begin(), haystack.end(), needle);
    return (it == haystack.end()) ? -1 : distance(haystack.begin(), it);
}

struct Hand : public Object {
    vector<Card> cards;
    Hand(const vector<Card> &c = vector<Card>()) : Object("Hand"), cards(c) {
        Relation r = Relation("is", *this, hand);
        relations.push_back(r);
    }
    void add(const string &rank, const string &suit) {
        cards.push_back(Card(rank, suit));
    }
    virtual List inspect() {
        /*vector<Object*> list;
        for (auto &card : cards) {
            list.push_back(new Card(card));
        }
        return List(list);*/
        //return List(listify(cards));
        return List(cards);
    }
    void remove(Card &card) {
        int idx = indexOf(cards, card);
        if (idx == -1) 
            throw na;
        cards.erase(cards.begin()+idx);
    }
};

struct Board : public Object {
    vector<Card> plays;
    vector<Card> covers;
    Board() : Object("Board") {
        Relation r = Relation("is", *this, board);
        relations.push_back(r);
    }
    void cover(const Object &card1, const Object &card2) {
        const Card &c2 = dynamic_cast<const Card&>(card2);
        int idx = indexOf(plays, c2);
        if (idx == -1) 
            throw na;
        if (covers[idx] != NoCard)
            throw na;
        covers[idx] = c2;
    }
    virtual List inspect() {
        //return List(vector{new List(listify(plays)), new List(listify(covers))});
        return List(vector{make_shared<List>(plays), make_shared<List>(covers)});
    }
    void play(const Object &card) {
        plays.push_back(dynamic_cast<const Card&>(card));
        covers.resize(plays.size(), NoCard);
    }
};

struct Game : public Object {
    Concept trump;
    Hand hand;
    Board board;
    Game() : Object("Game") {}
};

Object higherRank(Object &card1, Object &card2) {
    int i1 = indexOf(ranks, getRank(dynamic_cast<Card&>(card1))); 
    int i2 = indexOf(ranks, getRank(dynamic_cast<Card&>(card2))); 
    return i1 > i2 ? yes : no;
}

Object beats(vector<Object> &args, Game &game) {
    if (args.size() != 2) 
        throw na;
    Card &card1 = dynamic_cast<Card&>(args[0]);
    Card &card2 = dynamic_cast<Card&>(args[1]);
    if (getSuit(card1) == game.trump && getSuit(card1) != game.trump)
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
