
#include <string>
#include <vector>
#include <algorithm>
#include <type_traits>
#include <memory>
#include <iostream>

using namespace std;

// A concept can be 
// 0. an abstract idea
// 1. an object
// 2. an action
// 3. a relation between objects
struct Concept {
    string name;
    Concept(const string &n) : name(n) {}
    Concept();
    virtual bool operator==(const Concept &other) const {
        return other.name == name;
    }
    virtual bool operator!=(const Concept &other) const {
        return !(*this == other);
    }
    friend ostream& operator<<(ostream &os, const Concept &c);
};

ostream &operator<<(ostream &os, const Concept &c) {
    os << c.name;
    return os;
}

struct List;

struct Object : public Concept {
    static int idcounter;
    int uniqid;
    Object(const string &n) : Concept(n), uniqid(idcounter++) {}
    virtual bool operator==(const Object &other) const {
        return other.uniqid == uniqid;
    }
    virtual bool operator!=(const Object &other) const {
        return !(*this == other);
    }
    // Get object properties
    [[noreturn]] virtual List inspect();
};
int Object::idcounter = 0;

struct ConceptWrap : public Object {
    ConceptWrap(const Concept &c) : Object(c.name) {}
};

Object null("Null");
Object na("Not Applicable");
Object yes("Yes");
Object no("No");

Concept::Concept() : name(null.name) {}

struct Game;
typedef Object (*ActionFn)(vector<Object>&, Game &);

// An action takes objects as arguments and returns some value
struct Action : public Object {
    vector<Object*> args;
    ActionFn fn;
    Action(const string &n, ActionFn f) : Object(n), fn(f) {}
};

// A relation between two objects
struct Relation : public Object {
    Concept from;
    Concept to;
    Relation(const string &n, Concept f, Concept t) : Object(n), from(f), to(t) {}
};

vector<Relation> relations;

// List workhorse
Concept list("List");

struct List : public Object {
    vector<shared_ptr<Object>> items;
    List(const vector<shared_ptr<Object>> &objs = vector<shared_ptr<Object>>()) : Object("List"), items(objs) {
        relations.push_back(Relation("is", *this, list));
    }
    template<typename T, typename enable_if<is_base_of<Object, T>::value>::type* = nullptr>
    List(const vector<T> &vec) : Object("List") {
        for (const T &t : vec) {
            items.push_back(make_shared<T>(t));
        }
    }
    template<typename T, typename enable_if<is_base_of<Object, T>::value>::type* = nullptr>
    List(const vector<shared_ptr<T>> &vec) : Object("List") {
        for (const shared_ptr<T> &t : vec) {
            items.push_back(t);
        }
    }
};

[[noreturn]] List Object::inspect() {
    throw na;
}
