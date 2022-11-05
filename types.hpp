
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
struct List;

struct Object {
    static int idcount;
    int id;
    string type;
    Object(const string &t = "none") : id(idcount++), type(t) {}
    virtual bool operator==(const Object &other) const {
        return other.id == uniqid;
    }
    virtual bool operator!=(const Object &other) const {
        return !(*this == other);
    }
    // For compatibility with list
    virtual Object &operator[](int idx) const {
        return dynamic_cast<List>(this)[idx];
    }
    void make(const Concept &c);
    bool is(const Concept &c) const;
    void give(const Object &o);
    Object &get(const Concept &c);
    void remove(const Object &o);
};
int Object::idcount = 0;

template<>
struct hash<Object> {
    size_t operator(const Object &o) {
        return hash<int>{}(o.id);
    }
}

// A concept map enforces the single concept instance limit
// Concepts live in the concept_map
unordered_map<string,Concept> concept_map;

// "is" and "has" maps for objects->concepts and objects->objects, respectively
// Objects don't live here
// Unless they are a shared_ptr
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

void Object::make(const Concept &c) {
    is_map.insert({*this,&c});
}

bool Object::is(const Concept &c) const {
    return *is_map[*this] == c;
}

void Object::give(const Object &o) {
    has_map.insert({*this,&o});
}

// Get object property
Object &get(const Concept &c) {
    auto range = has_map.equal_range(*this);
    for (auto it = range.first; it != range.second; it++) {
        if (it->is(c)) {
            return *it;
        }
    }
    return null;
}

// Remove object property
void Object::remove(const Object &o) {
    auto range = has_map.equal_range(*this);
    for (auto it = range.first; it != range.second; it++) {
        if (it->is(c)) {
            has_map.erase(it);
            return;
        }
    }
}

// Object constructor
// Make object of specific type
Object::Object(const string &type) : id(idcount++) {
    this->make(type);
}

// Convenience
struct List : public Object {
    vector<Object *> objects;
    List() : Object("list") {}
    virtual bool operator==(const List &other) const {
        if (objects.size() != other.objects.size()) {
            return false;
        }
        for (int i=0; i<objects.size() && i<other.objects.size(); i++) {
            if (*objects[i] != *other.objects[i]) {
                return false;
            }
        }
        return true;
    }
    virtual bool operator!=(const List &other) const {
        return !(*this == other);
    }
    virtual Object &operator[](int idx) const {
        if (idx < 0 || idx >= objects.size()) {
            throw na;
        }
        return *objects[idx];
    }
    int size() const {
        return objects.size();
    }
    virtual void add(const Object &obj) {
        objects.push_back(&obj);
    }
    virtual void remove(const Object &obj) {
        for (auto it = objects.begin(); it != objects.end(); it++) {
            if (*it == obj) {
                objects.erase(it);
                return;
            }
        }
    }
    int index_of(const Object &obj) {
        int i = 0;
        for (auto it = objects.begin(); it != objects.end(); it++, i++) {
            if (*it == obj) {
                return i;
            }
        }
        return -1;
    }
}

typedef Object *(*ActionFn)(List &, Object &);

// An action takes objects as arguments and returns some value
struct Action : public Concept {
    Concept &res_type;
    vector<Concept*> arg_types;
    vector<Action*> children;
    ActionFn fn;
    Action(const string &n, const Concept &rt, const vector<Concept*> ats, ActionFn f) 
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

// Convenience
struct Number : public Object {
    int val;
    Number(int v = 0) : Object("number"), val(v) {}
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