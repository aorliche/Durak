
#include <string>
#include <vector>
#include <algorithm>

using namespace std;

struct List;

// A concept can be 
// 0. an abstract idea
// 1. an object
// 2. an action
// 3. a relation between objects
struct Concept {
    string name;
    Concept(const string &n) : name(n) {}
    virtual bool operator==(Concept &other) {
        return &other->name == this->name;
    }
    virtual List inspect();
}

struct Object : public Concept {
    static int idcounter;
    int uniqid;
    Object(const string &n) : Concept(n), uniqid(idcounter++) {}
    virtual bool operator==(Object &other) {
        return &other->uniqid == this->uniqid;
    }
}
Object::idcounter = 0;

typedef Object (*)(vector<Object>&) ActionFn;

// An action takes objects as arguments and returns some value
struct Action : public Object {
    vector<Object> args;
    ActionFn fn;
    Action(const string &n, ActionFn f) : Object(n), fn(f) {}
}

// A relation between two objects
struct Relation : public Object {
    Concept from;
    Concept to;
    Relation(const string &n, Concept f, Concept t) : Object(name), from(f), to(t) {}
}

