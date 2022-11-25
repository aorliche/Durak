
#include <string>
#include <vector>
#include <algorithm>
#include <type_traits>
#include <memory>
#include <iostream>
#include <unordered_map>
#include <map>
#include <set>

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
    // Constructors
    Object(int _id) : id(_id) {}
    Object();
    Object(const string &type = "");
    Object(const Concept c);
    // Operators
    virtual bool operator==(const Object &other) const {
        return id == other.id;
    }
    virtual bool operator!=(const Object &other) const {
        return !(*this == other);
    }
    virtual bool operator<(const Object &other) const {
        return id < other.id;
    }
    bool is(const Concept c) const;
    void give(const Object obj) const;
    void make(const Concept c) const;
    Object get(const Concept c) const;
    void remove(const Concept c) const;
};

// Should go in source file
int Object::idcount = 0;

struct Concept : public Object {
    Concept(const string &n);
    Concept(int _id) : Object(_id) {}
    virtual bool operator==(const Concept &other) const {
        return id = other.id;
    }
    virtual bool operator!=(const Concept &other) const {
        return !(*this == other);
    }
    friend ostream& operator<<(ostream &os, const Concept &c);
};

int History::idcount = 0; 

// The set of all concepts
set<string> concepts;

// "name" map for concepts
// no storing extra data in "budget polymorphic" Object subclasses Concept and List
unordered_map<int,string> name_map;
unordered_map<string,int> rev_name_map;

// Don't use -1 because that's a default in History
int CONCEPT_OJBECT_ID = -2;

// "list" map for lists
unordered_map<int,vector<Object>> list_map;

// "is" and "has" maps 
// for Object to Concepts and Object to Objects, respectively
multimap<int,int> is_map;
multimap<int,int> has_map;

struct ConceptInit {
    ConceptInit() {
        name_map[CONCEPT_ID] = "concept";
        rev_name_map["concept"] = CONCEPT_ID;
    }
} dummy_concept_init;

// While searching for correct action, keep track of state history
// for easy manipulation without copying
struct History {
    static int idcount;
    int id;
    int iskey, haskey, listkey, isval, hasval;
    vector<Object> listval;
    History(int isk = -1, int hask = -1, int listk = -1, 
        int isv = -1, int hasv = -1, const vector<Object> &listv = vector<Object>()) 
        : 
        id(idcount++), iskey(isk), haskey(hask), listkey(listk),
        isval(isv), hasval(hasv), listval(listv) {}
    static void mm_rollback(multimap<int,int> &mm, int key, int val);
    void rollback() {
        if (iskey != -1) {
            mm_rollback(is_map, iskey, isval);
        }
        if (haskey != -1) {
            mm_rollback(has_map, haskey, hasval);
        }
        if (listkey != -1) {
            list_map.erase(listkey);
        }
    }
    void putback() {
        if (iskey != -1) {
            is_map.insert({iskey, isval});
        }
        if (haskey != -1) {
            has_map.insert({haskey, hasval});
        }
        if (listkey != -1) {
            list_map[listkey] = listval;
        }
    }
};

static void mm_rollback(multimap<int,int> &mm, int key, int val) {
    auto pair = mm.find(key);
    for (; pair.first != pair.second; pair.first++) {
        if (*pair.first == val) {
            mm.erase(pair.first);
            return;
        }
    }
    cout << "Failed to roll back " << key << endl; // Maybe more object info
}

ostream &operator<<(ostream &os, const Concept &c) {
    os << name_map[c.id];
    return os;
}

// Concepts are singletons
Concept::Concept(const string &n) : Object() {
    if (concepts.count(n) != 0) {
        // Keep object equality
        id = rev_name_map[n];
    } else {
        name_map[id] = n;
        rev_name_map[n] = id;
        concepts.insert(n);
    }
}

Concept null("null");           // null object
Concept na("not applicable");   // exception
Concept boolean("boolean");
Object nullobj("null");         // null
Object yes("boolean");          // true
Object no("boolean");           // false

// For concept objects
Object::Object() : id(idcount++) {
    is_map.insert({id,CONCEPT_ID});
}

// Objects of concrete type
Object::Object(const string &type) : id(idcount++) {
    Concept c(type);
    is_map.insert({id,c.id});
}

// Concrete type again
Object::Object(const Concept c) : id(idcount++) {
    is_map.insert({id,c.id});
}

bool Object::is(const Concept c) const {
    auto range = is_map.equal_range(id);
    for (auto it = range.first; it != range.second; it++) {
        if (it->second == c.id) {
            return true;
        }
    }
    return false;
}

void Object::give(const Object o) const {
    has_map.insert({id,o.id});
}

void Object::make(const Concept c) const {
    is_map.insert({id,c.id});
}

// Get object property
Object Object::get(const Concept c) const {
    auto range = has_map.equal_range(id);
    for (auto it = range.first; it != range.second; it++) {
        auto obj = Object(it->second);
        if (obj.is(c)) {
            return obj;
        }
    }
    return nullobj;
}

// Remove object property
void Object::remove(const Concept c) const {
    auto range = has_map.equal_range(id);
    for (auto it = range.first; it != range.second; it++) {
        auto obj = Object(it->second);
        if (obj.is(c)) {
            has_map.erase(it);
            return;
        }
    }
}

struct List : public Object {
    List(int _id) : Object(_id) {}
    List(const Object &o) : Object(o.id) {}
    List(const string &field_name = "") : Object("list") {
        list_map[id];
        if (field_name.length() != 0) {
            Object(id).make(Concept(field_name));
        }
    }
    vector<Object> &List::get_objects() const {
        return list_map.at(id);
    }
    virtual bool operator==(const List &other) const {
        vector<Object> &obj = get_objects();
        vector<Object> &oth = other.get_objects();
        return obj == oth;
    }
    virtual bool operator!=(const List &other) const {
        return !(*this == other);
    }
    Object operator[](int idx) const {
        vector<Object> &objects = get_objects();
        if (objects.size() > idx) {
            return objects.at(idx);
        }
        return na;
    }
    int size() const {
        vector<Object> &objects = get_objects();
        return objects.size();
    }
    void add(Object obj) {
        vector<Object> &objects = get_objects();
        objects.push_back(obj);
    }
    void remove(Object obj) {
        vector<Object> &objects = get_objects();
        for (auto it = objects.begin(); it != objects.end(); it++) {
            if (it->id == obj.id) {
                objects.erase(it);
                return;
            }
        }
    }
    int index_of(Object obj) {
        vector<Object> &objects = get_objects();
        for (size_t i=0; i<objects.size(); i++) {
            if (objects[i].id == obj.id) {
                return i;
            }
        }
        return -1;
    }
};

/*

typedef Object (*ActionFn)(vector<Object> &, Object &);

// An action takes objects as arguments and returns some value
struct Action : public Object {
    const Concept res_type;
    vector<Concept> arg_types;
    ActionFn fn;
    Action(const string &name, const string &rtype, const vector<string> atypes, ActionFn f) 
        : Object(name), res_type(rtype), fn(f) {
            for (auto atype : atypes) {
                arg_types.push_back(Concept(atype));
            }
        }
    bool is_compatible(const vector<Object> &args) {
        for (size_t i=0; i<args.size(); i++) {
            Concept type = action.arg_types[i];
            Object arg = args[i];
            if (!arg.is(type)) {
                return false;
            }
        }
        return true;
    }
    Object eval(vector<Object> &args, Object &ctx) {
        return fn(args, ctx);
    }
};


// A relation between two objects
struct Relation : public Object {
    Object &from;
    Object &to;
    Relation(const string &typ, Object f, Object t) : Object(typ), from(f), to(t) {}
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
};*/