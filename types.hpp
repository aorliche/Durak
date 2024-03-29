#ifndef DURAK_TYPES_H
#define DURAK_TYPES_H

#include <string>
#include <vector>
#include <algorithm>
#include <type_traits>
#include <memory>
#include <iostream>
#include <unordered_map>
#include <map>
#include <set>
#include <sstream>
#include <iterator>
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
    // Constructors
    Object(int _id) : id(_id) {}
    Object(const Object &o) : id(o.id) {}
    Object();
    Object(const string &type);
    Object(const Concept &c);
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
    bool is(const string &str) const;
    void give(const Object obj) const;
    void make(const Concept c) const;
    void make(const string &str) const;
    void unmake(const Concept &c) const;
    void unmake(const string &str) const;
    Object get(const Concept c) const;
    Object get(const string &str) const;
    void remove(const Concept c) const;
    void remove(const string &str) const;
    // Search
    size_t get_sig() const;
    // Printing and debug
    friend ostream& operator<<(ostream &os, const Object &o);
    void dump(ostream &os) const;
};

// Should go in source file
int Object::idcount = 0;

struct Concept : public Object {
    Concept(const string &n);
    Concept(int _id) : Object(_id) {}
    Concept(const Object &o) : Object(o.id) {}
    Concept(const Concept &c) : Object(c.id) {}
    virtual bool operator==(const Concept &other) const {
        return id == other.id;
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
unordered_map<int,vector<int>> list_map;

// "is" and "has" maps 
// for Object to Concepts and Object to Objects, respectively
multimap<int,int> is_map;
multimap<int,int> has_map;

// Hash operator for vector<int>
// https://stackoverflow.com/questions/20511347/a-good-hash-function-for-a-vector
struct int_vector_hasher {
    std::size_t operator()(std::vector<int> const& vec) const {
        std::size_t seed = vec.size();
        for(auto x : vec) {
            x = ((x >> 16) ^ x) * 0x45d9f3b;
            x = ((x >> 16) ^ x) * 0x45d9f3b;
            x = (x >> 16) ^ x;
            seed ^= x + 0x9e3779b9 + (seed << 6) + (seed >> 2);
        }
        return seed;
    }
};

// Get list of concepts for the object
size_t Object::get_sig() const {
    auto range = is_map.equal_range(id);
    vector<int> ids;
    for (auto it = range.first; it != range.second; it++) {
        ids.push_back(it->second);
    }
    sort(ids.begin(), ids.end());
    return int_vector_hasher{}(ids);
}

void Object::dump(ostream &os) const {
    os << *this << endl;
    auto range = has_map.equal_range(id);
    int i=0; 
    for (auto it = range.first; it != range.second; it++) {
        Object(it->second).dump(os);
        cout << endl;
    }
}

ostream &operator<<(ostream &os, const Concept &c) {
    os << name_map[c.id];
    return os;
}

// Only for the static initializer
struct ConceptInit {
    ConceptInit() {
        name_map[CONCEPT_ID] = "concept";
        rev_name_map["concept"] = CONCEPT_ID;
        concepts.insert("concept");
    }
} concept_init;

// History stack
vector<tuple<int, int, int, int>> history;

// History Operations
enum MapType : int {
    HasMap, IsMap, ListMap
};

enum MapOp : int {
    Add, Remove
};

// Functional interface to has, is, and list maps
// While searching for correct action, keep track of state history
void map_add(multimap<int,int> &map, MapType type, int a, int b, bool record) {
    map.insert({a,b});
    if (record)
        history.push_back({type, MapOp::Add, a, b});
}

void map_remove(multimap<int,int> &map, MapType type, int a, int b, bool record) {
    auto range = map.equal_range(a);
    for (auto it = range.first; it != range.second; it++) {
        if (it->second == b) {
            if (record)            
                history.push_back({type, MapOp::Remove, a, b});
            map.erase(it);
            return;
        }
    }
}

void has_add(int a, int b, bool record = true) {
    map_add(has_map, MapType::HasMap, a, b, record);
}

void is_add(int a, int b, bool record = true) {
    map_add(is_map, MapType::IsMap, a, b, record);
}

void has_remove(int a, int b, bool record = true) {
    map_remove(has_map, MapType::HasMap, a, b, record);
}

void is_remove(int a, int b, bool record = true) {
    map_remove(is_map, MapType::IsMap, a, b, record);
}

bool list_contains(int a, int b) {
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
        history.push_back({MapType::ListMap, MapOp::Add, a, b});
}

void list_remove(int a, int b, bool record = true) {
    if (!list_contains(a, b)) return;
    auto lst = list_map[a];
    for (auto it = lst.begin(); it != lst.end(); it++) {
        if (*it == b) {
            lst.erase(it);
            if (record)
                history.push_back({MapType::ListMap, MapOp::Remove, a, b});
            return;
        }
    }
}

class Exception {};

// Special objects and concepts
Concept null("null");           // null object
Concept boolean("bool");
Object nullobj("null");         // null
Object yes("bool");             // true
Object no("bool");              // false

ostream &operator<<(ostream &os, const Object &o) {
    if (o.is("concept")) {
        os << Concept(o);
        return os;
    }
    if (o == nullobj) {
        os << "nullobj";
        return os;
    }
    auto range = is_map.equal_range(o.id);
    int i=0;
    for (auto it = range.first; it != range.second; it++) {
        if (i++ > 0) os << ',';
        os << name_map[it->second];
    }
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
    is_add(id, CONCEPT_ID);
}

// For concept objects
Object::Object() : id(idcount++) {}

// Objects of concrete type
Object::Object(const string &type) : id(idcount++) {
    Concept c(type);
    Concept o("object");
    is_add(id, c.id);
    is_add(id, o.id);
}

// Concrete type again
Object::Object(const Concept &c) : id(idcount++) {
    Concept o("object");
    is_add(id, c.id);
    is_add(id, o.id);
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

bool Object::is(const string &str) const {
    return is(Concept(str));
}

void Object::give(const Object o) const {
    has_add(id, o.id);
}

void Object::make(const string &str) const {
    make(Concept(str));
}

void Object::make(const Concept c) const {
    is_add(id, c.id);
}

void Object::unmake(const Concept &c) const {
    is_remove(id, c.id);
}

void Object::unmake(const string &str) const {
    unmake(Concept(str));
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

Object Object::get(const string &str) const {
    return get(Concept(str));
}

// Remove object property
void Object::remove(const Concept c) const {
    has_remove(id, c.id);
    /*auto range = has_map.equal_range(id);
    for (auto it = range.first; it != range.second; it++) {
        auto obj = Object(it->second);
        if (obj.is(c)) {
            has_remove(obj.id, c.id);
            return;
        }
    }*/
}

void Object::remove(const string &str) const {
    remove(Concept(str));
}

vector<int> objvec2intvec(const vector<Object> &objs) {
    vector<int> res;
    for (size_t i=0; i<objs.size(); i++) {
        res.push_back(objs[i].id);
    }
    return res;
}

vector<Object> intvec2objvec(const vector<int> &ints) {
    vector<Object> res;
    for (size_t i=0; i<ints.size(); i++) {
        res.push_back(Object(ints[i]));
    }
    return res;
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
        vector<Object> objects = intvec2objvec(list_map[id]);
        if (objects.size() > idx) {
            return objects.at(idx);
        }
        throw Exception();
    }
    vector<Object> get_objects() {
        return intvec2objvec(list_map[id]);
    }
    int size() const {
        return list_map[id].size();
    }
    void add(Object obj) {
        list_add(id, obj.id);
    }
    void remove(Object obj) {
        list_remove(id, obj.id);
    }
    int index_of(Object obj) {
        vector<Object> objects = intvec2objvec(list_map[id]);
        for (size_t i=0; i<objects.size(); i++) {
            if (objects[i].id == obj.id) {
                return i;
            }
        }
        return -1;
    }
};

// Action function type
typedef Object (*ActionFn)(const vector<Object> &);

// Action data structures
unordered_map<int,int> act_res_map;
unordered_map<int,ActionFn> act_map;
unordered_map<int,vector<int>> act_args_map;
//unordered_map<int,vector<int>> act_inst_map;

// An action takes objects as arguments and returns some value
struct Action : public Object {
    static Action no_action;
    Action() : Object("action") {}
    Action(int _id) : Object(_id) {}
    Action(const Action &a) : Object(a.id) {}
    Action(const string &name, const string &rtype, const vector<string> &atypes, ActionFn fn) : Object(name) {
        make("action");
        act_res_map[id] = Concept(rtype).id;
        auto &act_args = act_args_map[id];
        for (size_t i=0; i<atypes.size(); i++) {
            act_args.push_back(Concept(atypes[i]).id);
        }
        act_map[id] = fn;
    }
    Action(Concept name, Concept rtype, const vector<int> &atypes, ActionFn fn) : Object(name) {
        make("action");
        act_res_map[id] = rtype.id;
        act_args_map[id] = atypes;
        act_map[id] = fn;
    }
    Object eval(vector<Object> &args) {
        return act_map[id](args);
    }
    Concept get_concept() {
        auto range = is_map.equal_range(id);
        for (auto it = range.first; it != range.second; it++) {
            int cid = it->second;
            if (name_map[cid] != "object" and name_map[cid] != "action") {
                return Concept(cid);
            }
        }
        return null;
    }
    string get_name() {
        return name_map[get_concept().id];
    }
    size_t get_args_size() const {
        return act_args_map[id].size();
    }
    /*Action instantiate(const vector<Object> &args) {
        if (!is_compatible(args)) 
            throw na;
        Action a(get_concept(), Concept(act_res_map[id]), act_args_map[id], act_map[id]);
        inst_map[a.id] = objvec2intvec(args);
    }*/
    bool is_compatible_arg(int idx, Object arg) {
        return arg.is(Concept(act_args_map[id][idx]));
    }
    bool is_compatible(const vector<Object> &args) {
        auto &act_args = act_args_map[id];
        if (args.size() != act_args.size()) {
            return false;
        }
        for (size_t i=0; i<args.size(); i++) {
            Concept type(act_args[i]);
            Object arg(args[i]);
            if (!arg.is(type)) {
                return false;
            }
        }
        return true;
    }
};

Action Action::no_action;

// TODO not used
// typedef tuple<int,int,int,int> History;

unordered_map<int,int> node_res_map;
unordered_map<int,vector<int>> node_parent_map;
unordered_map<int,int> node_act_map;
unordered_map<int,vector<int>> node_sig_map;

typedef const function<ostream &(ostream &, const Object&)> value_print;

struct Node {
    static int idcount;
    int id;
    Node(int i) : id(i) {}
    Node(const Node &n) : id(n.id) {}
    Node(const Object &r) : id(idcount++) {
        node_res_map[id] = r.id;
        node_act_map[id] = Action::no_action.id;
        sign();
    }
    Node(const Action &a) : id(idcount++) {
        node_res_map[id] = nullobj.id;
        node_act_map[id] = a.id;
        sign();
    }
    Node(const Object &r, const Action &a, const vector<int> &p) 
            : id(idcount++) {
        node_res_map[id] = r.id;
        node_act_map[id] = a.id;
        node_parent_map[id] = p;
        sign();
    }
    Object get_res() const {
        return node_res_map[id];
    }
    Action get_act() const {
        return node_act_map[id];
    }
    vector<int> &get_sig() const {
        return node_sig_map[id];
    }
    vector<int> get_shallow() const {
        vector<int> shallow{node_act_map[id]};
        vector<int> &parents = node_parent_map[id];
        for (size_t i=0; i<parents.size(); i++) {
            int resid = node_res_map[parents[i]];
            shallow.push_back(resid);
        }
        return shallow;
    }
    void sign() {
        vector<int> sig;
        vector<int> &parents = node_parent_map[id];
        int resid = node_res_map[id];
        int actid = node_act_map[id];
        sig.push_back(Object(resid).get_sig());
        sig.push_back(actid);
        for (size_t i=0; i<parents.size(); i++) {
            sig.insert(
                sig.end(), 
                node_sig_map[parents[i]].begin(), 
                node_sig_map[parents[i]].end());
        }
        node_sig_map[id] = sig;
    }
    string sig_str() const {
        stringstream str;
        vector<int> sig = node_sig_map[id];
        copy(sig.begin(), sig.end(), ostream_iterator<int>(str, " "));
        return str.str();
    }
    Object eval() {
        vector<Object> args;
        vector<int> parents = node_parent_map[id];
        Action act(node_act_map[id]);
        for (size_t i=0; i<parents.size(); i++) {
            Object res = node_res_map[parents[i]];
            args.push_back(res);
        }
        return act.eval(args);
    }
    void remove() {
        node_act_map.erase(id);
        node_res_map.erase(id);
        node_parent_map.erase(id);
        node_sig_map.erase(id);
    }
    void print(ostream &os, size_t lvl = 0, value_print &vp = {}) {
        Action act(node_act_map[id]);
        Object res(node_res_map[id]);
        vector<int> parents = node_parent_map[id];
        for (size_t i=0; i<lvl; i++) {
            os << '\t';
        }
        os << act.get_name() << " (" << res << ") ";
        if (vp) vp(os, res);
        os << endl;
        for (size_t i=0; i<parents.size(); i++) {
            Node(parents[i]).print(os, lvl+1, vp);
        }
    }
    bool operator==(const Node &other) const {
        return id == other.id;
    }
};

int Node::idcount = 0;

void print_nodes(ostream &os, vector<Node> &nodes, value_print &vp = {}) {
    for (size_t i=0; i<nodes.size(); i++) {
        nodes[i].print(os, 0, vp);
    }
}

void print_objs(ostream &os, vector<int> &ids) {
    for (size_t i=0; i<ids.size(); i++) {
        os << Object(ids[i]) << ' ';
    }
    os << endl;
}

# define MAKEACTION(a,ret,args) Action a##_action(#a,ret,args,a)

// Most general functions
Object get_field(const vector<Object> &args) {
    return args[0].get(Concept(args[1].id));
}

MAKEACTION(get_field, "get-field", (vector<string>{"object", "concept"}));

// Expand list has special behavior in search code
Object expand_list(const vector<Object> &args) {
    return args[0];
}

MAKEACTION(expand_list, "list", (vector<string>{"list"}));

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

#endif