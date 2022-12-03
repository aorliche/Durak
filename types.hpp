
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

// The set of all concepts
set<string> concepts;

// "name" map for concepts
// no storing extra data in "budget polymorphic" Object subclasses Concept and List
unordered_map<int,string> name_map;
unordered_map<string,int> rev_name_map;

// Concept of concepts
int CONCEPT_ID = -1;

// "list" map for lists
unordered_map<int,vector<Object>> list_map;

// "is" and "has" maps 
// for Object to Concepts and Object to Objects, respectively
multimap<int,int> is_map;
multimap<int,int> has_map;

ostream &operator<<(ostream &os, const Concept &c) {
    os << name_map[c.id];
    return os;
}

// Only for the static initializer
struct ConceptInit {
    ConceptInit() {
        name_map[CONCEPT_ID] = "concept";
        rev_name_map["concept"] = CONCEPT_ID;
    }
} concept_init;

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

// Special objects and concepts
Concept null("null");           // null object
Concept na("not applicable");   // exception
Concept boolean("boolean");
Object nullobj("null");         // null
Object yes("boolean");          // true
Object no("boolean");           // false

// History stack
vector<tuple<int, int, int, int>> history;

// History Operations
enum MapType {
    HasMap, IsMap, ListMap
};

enum MapOp {
    Add, Remove
};

// Functional interface to has, is, and list maps
// While searching for correct action, keep track of state history
void map_add(multimap<int,int> &map, MapType type, int a, int b, bool record) {
    map.insert({a,b});
    if (record)
        history.push_back({type, MapOp.Add, a, b});
}

void map_remove(multimap<int,int> &map, MapType type, int a, int b, bool record) {
    auto range = has_map.equal_range(a);
    for (auto it = range.first; it != range.second; it++) {
        if (*it == b) {
            if (record)            
                history.push_back({type, MapOp.Remove, a, b});
            map.erase(it);
            return;
        }
    }
}

void has_add(int a, int b, bool record = true) {
    map_add(has_map, MapType.HasMap, a, b, record);
}

void is_add(int a, int b, bool record = true) {
    map_add(is_map, MapType.IsMap, a, b, record);
}

void has_remove(int a, int b, bool record = true) {
    map_remove(has_map, MapType.HasMap, a, b, record);
}

void is_remove(int a, int b, bool record = true) {
    map_remove(is_map, MapType.IsMap, a, b, record);
}

void list_contains(int a, int b) {
    auto lst = list_map[a];
    for (int i=0; i<lst.size(); i++) {
        if (lst[i] == b) return true;
    }
    return false;
}

void list_add(int a, int b, bool record = true) {
    if (list_contains(a, b)) return;
    list_map[a].push_back(b);
    if (record)
        history.push_back({MapType.ListMap, MapOp.Add, a, b});
}

void list_remove(int a, int b, bool record = true) {
    if (!list_contains(a, b)) return;
    auto lst = list_map[a];
    for (auto it = lst.begin(); it != lst.end(); it++) {
        if (*it == b) {
            lst.erase(it);
            if (record)
                history.push_back({MapType.ListMap, MapOp.Remove, a, b});
            return;
        }
    }
}

// For concept objects
Object::Object() : id(idcount++) {
    is_add(id, CONCEPT_ID);
}

// Objects of concrete type
Object::Object(const string &type) : id(idcount++) {
    Concept c(type);
    is_add(id, c.id);
}

// Concrete type again
Object::Object(const Concept c) : id(idcount++) {
    is_add(id, c.id);
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
    has_add(id, o.id);
}

void Object::make(const Concept c) const {
    is_add(id, c.id);
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

// Useful because it's indexed
struct List : public Object {
    List(int _id) : Object(_id) {}
    List(const Object &o) : Object(o.id) {}
    List(const string &field_name = "") : Object("list") {
        if (field_name.length() != 0) {
            make(Concept(field_name));
        }
    }
    virtual bool operator==(const List &other) const {
        return list_map[id] == list_map[other.id];
    }
    virtual bool operator!=(const List &other) const {
        return !(*this == other);
    }
    Object operator[](int idx) const {
        vector<Object> &objects = list_map[id];
        if (objects.size() > idx) {
            return objects.at(idx);
        }
        return na;
    }
    int size() const {
        list_map[id].size();
    }
    void add(Object obj) {
        list_add(id, obj.id);
    }
    void remove(Object obj) {
        list_remove(id, obj.id);
    }
    int index_of(Object obj) {
        vector<Object> &objects = list_map[id];
        for (size_t i=0; i<objects.size(); i++) {
            if (objects[i].id == obj.id) {
                return i;
            }
        }
        return -1;
    }
};

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

/*
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