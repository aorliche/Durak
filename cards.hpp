
Concept null("Null");
Concept card("Card");
Concept suit("Suit");
Concept rank("Rank");
vector<Concept> suits{"Diamonds", "Hearts", "Clubs", "Spades"};
vector<Concept> ranks{"6","7","8","9","10","Jack","Queen","King","Ace"};
Concept list("List");
Concept board("Board");
Concept hand("Hand");
Concept game("Game");
Concept player("Player");

Object NA("Not Applicable");
Object True("True");
Object False("False");

vector<Relation> relations;

struct Card : public Object {
    Concept suit;
    Concept rank;
    Card() {
        relations.push_back(Relation("is", *this, card));
    }
    Card(string s) : suit(null), rank(null) {}
    Card(Concept &s, Concept &r) : Object("Card"), suit(&s), rank(&r) {}
}

Card NoCard("No Card");

struct Number : public Object {
    int val;
    Number(int nn) : n(nn) {
        Relation r = Relation("is", *this, number);
        relations.push_back(r);
    }
    bool operator==(Number &other) {
        return val == other.val;
    }
    bool operator>(Number &other) {
        return val > other.val;
    }
    bool operator<(Number &other) {
        return val < other.val;
    }
}

struct List : public Object {
    vector<Object> items;
    List(vector<Object> &it) : items(it) {
        relations.push_back(Relation("is", *this, list));
    }
}

virtual List Concept::inspect() {
    throw NA;
}

Concept getSuit(string &s) {
    for (int i=0; i<suits.size(); i++) {
        if (suits[i].name == s) return suits[i];
    }
    throw NA;
}

Concept getSuit(Card &card) {
    return card.suit;
}

Concept getRank(string &r) {
    for (int i=0; i<ranks.size(); i++) {
        if (ranks[i].name == r) return ranks[i];
    }
    throw NA;
}

Concept getRank(Card &card) {
    return card.rank;
}

int indexOf(vector<Concept> &haystack, Concept &needle) {
    auto it = find(haystack.begin(), haystack.end(), needle);
    return (it == haystack.end()) ? -1 : distance(haystack.begin(), it);
}

struct Hand : public Object {
    vector<Card> cards;
    Hand(const string &n) : name(n) {
        Relation r = Relation("is", *this, hand);
        relations.push_back(r);
    }
    void add(string &suit, string &rank) {
        cards.push_back(Card(suit, rank));
    }
    virtual List inspect() {
        return List(cards);
    }
    void remove(Card &card) {
        int idx = indexOf(cards, card);
        if (idx == -1) 
            throw NA;
        cards.erase(cards.begin()+idx);
    }
}

struct Board : public Object {
    vector<Card> plays;
    vector<Card> covers;
    Board(const string &n) : name(n) {
        Relation r = Relation("is", *this, board);
        relations.push_back(r);
    }
    void cover(Object &card1, Object &card2) {
        int idx = indexOf(plays, card2);
        if (idx == -1) 
            throw NA;
        if (covers[idx] != NoCard)
            throw NA;
        covers[idx] = card2;
    }
    virtual List inspect() {
        return List(vector{List(plays), List(covers)});
    }
    void play(Object &card) {
        plays.push_back(card);
        covers.resize(plays.size(), NoCard);
    }
}

Object higherRank(Object &card1, Object &card2) {
    int i1 = indexOf(ranks, getRank(card1)); 
    int i2 = indexOf(ranks, getRank(card2)); 
    return i1 > i2;
}

Object beats(vector<Object> &args, Object &ctx) {
    if (args.size() != 2) 
        throw NA;
    Card &card1 = dynamic_cast<Card&>(args[0]);
    Card &card2 = dynamic_cast<Card&>(args[1]);
    Suit &trump = dynamic_cast<Suit&>(ctx.trump);
    if (eq(getSuit(card1), trump) && ne(getSuit(card1), trump)) 
        return True;
    else if (eq(getSuit(card2), trump)) 
        return False;
    else 
        return higherRank(card1, card2);
}

Object cover(vector<Object> &args, Object &ctx) {
    if (args.size() != 2)
        throw NA;
    Card &card1 = dynamic_cast<Card&>(args[0]);
    Card &card2 = dynamic_cast<Card&>(args[1]);
    ctx.hand.remove(card1);
    ctx.board.cover(card1, card2);
}

Object getItem(vector<Object> &args, Object &ctx) {
    if (args.size() != 2) 
        throw NA;
    List &lst = dynamic_cast<List&>(args[0]);
    Number &num = dynamic_cast<Number&>(args[1]);
    return lst.items[num.val];
}

Object getSize(vector<Object> &args, Object &ctx) {
    if (args.size() != 1) 
        throw NA;
    List &lst = dynamic_cast<List&>(args[0]);
    return Number(lst.items.size());
}

Object randInt(vector<Object> &args, Object &ctx) {
    if (args.size() != 1) 
        throw NA;
    Number &num = dynamic_cast<Number&>args[0];
    return Number(rand()%num.val);
}

Object doNothing(vector<Object> &args, Object &ctx) {
    if (args.size() != 0) 
        throw NA;
}

Object getTrump(vector<Object> &args, Object &ctx) {
    if (args.size() != 0) 
        throw NA;
    return ctx.trump;
}

Object getBoard(vector<Object> &args, Object &ctx) {
    if (args.size() != 0) 
        throw NA;
    return ctx.board;
}

Object getHand(vector<Object> &args, Object &ctx) {
    if (args.size() != 0) 
        throw NA;
    return ctx.hand;
}

Action beatsAction("beats", beats);
Action coverAction("cover", cover);
Action getItem("getItem", getItem);
Action randInt("randInt", randInt);
Action getTrump("getTrump", getTrump);
Action getBoard("getBoard", getBoard);
Action getHand("getHand", getHand);
Action doNothing("doNothing", doNothing);
