
#include <string>
#include <vector>
#include <algorithm>
#include <type_traits>
#include <memory>
#include <iostream>
#include <unordered_map>
#include <multimap>
#include <functional>

using namespace std;


// 0. objects are the base class
// 1. concepts are singleton objects
// 2. an action changes state or returns something
// 3. a relation exists between two objects

struct Concept;

struct Object {
    static int idcount;
    int id;
    Object() : id(idcount++) {}
    Object(const string &type);
    virtual bool operator==(const Object &other) const {
        return other.id == uniqid;
    }
    virtual bool operator!=(const Object &other) const {
        return !(*this == other);
    }
    // For compatibility with list and get return
    virtual Object *&operator[](int idx) const {
        return &this;
    }
    Object *get2(const Concept &c1, const Concept &c2) {
        return this->get(c1)->[0]->get(c2);
    }
    void make(const Concept &c);
    void give(const Object &o);
    Object *get(const Concept &c);
    void remove(const Object &o);
    bool is(const Concept &c) const;
};
int Object::idcount = 0;

template<>
struct hash<Object> {
    size_t operator(const Object &o) {
        return hash<int>{}(o.id);
    }
}

// A concept name map enforces single concept instance limit
// Concepts live in the concept_map
unordered_map<string,Concept> concept_map;

// Convenience "is" and "has" maps for objects->concepts and objects->objects, respectively
// Objects don't live here
multimap<Object,Concept*> is_map;
multimap<Object,Object*> has_map;

struct Concept {
    string name;
    Concept(const string &n, bool insert = true) : name(n) {
        if (insert) {
            concept_map[name] = *this;
        }
    }
    friend ostream& operator<<(ostream &os, const Concept &c);
};

ostream &operator<<(ostream &os, const Concept &c) {
    os << c.name;
    return os;
}

Concept null("null");            // null object
Concept na("not applicable");    // exception
Concept yes("yes");              // true
Concept no("no");                // false

typedef Object (*ActionFn)(vector<Object>&, Object *);

// An action takes objects as arguments and returns some value
struct Action : public Concept {
    Concept *res_type;
    vector<Concept*> arg_types;
    vector<Action*> children;
    ActionFn fn;
    Action(const string &n, const Concept *rt, const vector<Concept*> ats, ActionFn f) 
        : Concept(n), res_type(rt), args_types(at), fn(f) {}
};

// A relation between two objects
struct Relation : public Object {
    Object &from;
    Object &to;
    Relation(const string &name, Object f, Object t) : from(f), to(t) {
        Concept(name);
    }
};

vector<Relation> relations;

void Object::make(const Concept &c) {
    is_map.insert({*this,&c});
}

void Object::give(const Object &o) {
    has_map.insert({*this,&o});
}

void Object::remove(const Object &o) {
    has_map.erase(o);
}

bool Object::is(const Concept &c) const {
    return is_map[*this] == is_map[c];
}

// Convenience
struct List : public Object {
    vector<Object *> objects;
    List() {
        this->is("list");
    }
    virtual bool operator==(const List &other) const {
        if (objects.size() != other->objects.size()) {
            return false;
        }
        for (int i=0; i<objects.size() && i<other->objects.size(); i++) {
            if (objects[i] != other->objects[i]) {
                return false;
            }
        }
        return true;
    }
    virtual bool operator!=(const List &other) const {
        return !(*this == other);
    }
    virtual Object *&operator[](int idx) const {
        if (idx < 0 || idx >= objects.size()) {
            throw na;
        }
        return objects[idx];
    }
}

// Get object properties
Object *get(const Concept &c) {
    int count = has_map.count(*this);
    if (count == 0) {
        return &null;
    }
    auto range = has_map.equal_range(*this);
    auto lst = make_shared<List>();
    for (auto i = range.first; i != range.second; i++) {
        if (i->is(c)) {
            lst->objects.push_back(i);
        }
    }
    return lst;
}

Object::Object(const string &type) : id(idcount++) {
    this->make(type);
}