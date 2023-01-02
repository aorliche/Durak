#ifndef _base_hpp
#define _base_hpp

#include <variant>
#include <string>
#include <vector>
#include <unordered_map>
#include <unordered_set>
#include <algorithm>
#include <iostream>
#include <functional>
#include <bit>

using namespace std;

struct ObjectId {
    int id;
    ObjectId(int i) : id(i) {}
    bool operator==(const ObjectId &other) const {
        return id == other.id;
    }
    bool operator!=(const ObjectId &other) const {
        return !(*this == other);
    }
};

struct object_hasher {
    size_t operator()(const ObjectId &obj) const {
        return hash<int>{}(obj.id);
    }
};

typedef variant<bool,int,string,ObjectId,vector<ObjectId>,nullptr_t> Property;
Property none(nullptr);

struct Object {
    ObjectId id;
    bool fungible = false;
    static int idcount;
    Object() : id(idcount++) {}
    Object(const Object &obj) : id(obj.id) {}
    Object(const ObjectId &objid) : id(objid) {}
    void make(const string &type);
    void set(const string &key, const Property &val) const;
    bool is(const string &type) const;
    bool has(const string &key) const;
    Property &get(const string &key) const;
    bool operator==(const Object &other) const;
    friend ostream &operator<<(ostream &, const Object &);
};

unordered_map<ObjectId,vector<string>,object_hasher> types_map;
unordered_map<ObjectId,vector<string>,object_hasher> keys_map;
unordered_map<ObjectId,unordered_map<string,Property>,object_hasher> props_map;

int Object::idcount = 0;

void Object::make(const string &type) {
    types_map[id].push_back(type);
}

bool Object::is(const string &type) const {
    auto &types = types_map[id];
    for (auto it = types.begin(); it != types.end(); it++) {
        if (*it == type) return true;
    }
    return false;
}

bool Object::has(const string &key) const {
    auto &keys = keys_map[id];
    return find(keys.begin(), keys.end(), key) != keys.end();
}

void Object::set(const string &key, const Property &val) const {
    auto &keys = keys_map[id];
    auto &vals = props_map[id];
    if (find(keys.begin(), keys.end(), key) == keys.end()) {
        keys.push_back(key);
    }
    vals[key] = val;
}

Property &Object::get(const string &key) const {
    if (!has(key)) return none;
    auto &vals = props_map[id];
    return vals[key];
}

bool Object::operator==(const Object &other) const {
    if (id == other.id) return true;
    else if (fungible) return false;
    auto &keys = keys_map[id];
    auto &other_keys = keys_map[other.id];
    if (keys != other_keys) return false;
    auto &vals = props_map[id];
    auto &other_vals = props_map[other.id];
    for (int i=0; i<keys.size(); i++) {
        auto &val = vals[keys[i]];
        auto &other_val = other_vals[keys[i]];
        if (val != other_val) return false;
    }
    return true;
}

ostream &operator<<(ostream &os, const Object &obj) {
    auto &types = types_map[obj.id];
    auto &keys = keys_map[obj.id];
    auto &vals = props_map[obj.id];
    for (int i=0; i<types.size(); i++) {
        os << types[i];
        if (i < types.size()-1) {
            os << ',';
        }
    }
    if (keys.size() > 0) {
        os << " [";
    }
    for (int i=0; i<keys.size(); i++) {
        os << keys[i];
        auto &val = vals[keys[i]];
        switch (val.index()) {
            case 0: os << ": " << get<bool>(val); break;
            case 1: os << ": " << get<int>(val); break;
            case 2: os << ": " << get<string>(val); break;
            case 3: os << ": " << "Object"; break;
            case 4: os << ": " << "Vector<Object>"; break;
            case 5: os << ": " << "_none_"; break;
        }
        if (i < keys.size()-1) {
            os << ", ";
        }
    }
    if (keys.size() > 0) {
        os << ']';
    }
    return os;
}

typedef function<Property(const vector<Property> &)> FunctionEval;
typedef function<bool(const Property &, int)> FunctionAllow;

struct Function {
    string name;
    FunctionEval eval;
    FunctionAllow allow;
    int nargs;
    Function() : name("null") {}
    Function(const string &n, FunctionEval e, FunctionAllow a, int na) 
        : name(n), eval(e), allow(a), nargs(na) {}
    bool operator==(const Function &other) const {
        return name == other.name;
    }
    bool operator!=(const Function &other) const {
        return !(*this == other);
    }
    bool compatible(const Property &p, int n) const {
        if (name == "null") return false;
        return allow(p, n);
    }
    Property operator()(const vector<Property> &args) const {
        return eval(args);
    }
};

void default_pfn(ostream &os, const Property &p) {
    switch(p.index()) {
        case 0: os << get<bool>(p); break;
        case 1: os << get<int>(p); break;
        case 2: os << get<string>(p); break;
        case 3: os << Object(get<ObjectId>(p)); break;
        case 4: os << "Vector<Object>"; break;
        case 5: os << "_none_"; break;
    }
}

struct Node {
    Property res;
    Function fn;
    vector<Node> parents;
    Node() : res("null") {};
    Node(const Property &p) : res(p) {}
    Node(const Function &f) : fn(f) {}
    Node(const Node &n) : res(n.res), fn(n.fn), parents(n.parents) {}
    Property eval(vector<Node*> ps) const {
        vector<Property> args;
        for (int i=0; i<ps.size(); i++) {
            args.push_back(ps[i]->res);
        }
        return fn(args);
    }
    void update_parents(vector<Node*> ps) {
        parents.clear();
        for (int i=0; i<ps.size(); i++) {
            parents.push_back(*ps[i]);
        }
    }
    bool operator==(const Node &other) const {
        if (res != other.res 
            or fn != other.fn 
            or parents.size() != other.parents.size()) 
            return false;
        for (int i=0; i<parents.size(); i++) {
            if (parents[i].res != other.parents[i].res) 
                return false;
        }
        return true;
    }
    void print(ostream &os, 
        int lvl = 0, 
        function<void(ostream &, const Property&)> pfn = default_pfn) {
        for (int i=0; i<lvl; i++) {
            os << '\t';
        }
        os << "(" << fn.name << ") ";
        pfn(os, res);
        os << endl;
        for (int i=0; i<parents.size(); i++) {
            parents[i].print(os, lvl+1, pfn);
        }
    }
};

// Hasher
// Missing vector<ObjectId>
// Maybe okay because of get_property
size_t hash_prop(const Property &res) {
    size_t h = 0;
    switch (res.index()) {
        case 0: h = hash<bool>{}(get<bool>(res)); break;
        case 1: h = hash<int>{}(get<int>(res)); break;
        case 2: h = hash<string>{}(get<string>(res)); break;
        case 3: h = hash<int>{}(get<ObjectId>(res).id); break;
    }
    return h;
}

struct node_hasher {
    size_t operator()(const Node &n) const {
        size_t h = hash<string>{}(n.fn.name) + hash_prop(n.res);
        for (auto it = n.parents.begin(); it != n.parents.end(); it++) {
            auto &res = it->res;
            auto &name = it->fn.name;
            h += hash<string>{}(name) + hash_prop(res);
        }
        // cout << "hash" << h << endl;
        return h;
    }
};

// Helper functions
template <typename T>
int index_of(const vector<T> &vec, const T &item) {
    for (size_t i=0; i<vec.size(); i++) {
        if (vec[i] == item) {
            return i;
        }
    }
    return -1;
}

Property &get_helper(const Property &p, const string &s) {
    return Object(get<ObjectId>(p)).get(s);
}

// Generic functions
Property get_property_eval(const vector<Property> &objs) {
    return get_helper(objs[0], get<string>(objs[1]));
}

vector<ObjectId> expand_vec_eval(const vector<Property> &objs) {
    return get<vector<ObjectId>>(objs[0]);
}

bool get_property_allow(const Property &p, int n) {
    return (n == 0 and p.index() == 3) or (n == 1 and p.index() == 2);
}

bool expand_vec_allow(const Property &p, int n) {
    return n == 0 and p.index() == 4;
}

Function get_property("get_property", get_property_eval, get_property_allow, 2);
Function expand_vec("expand_vec", expand_vec_eval, expand_vec_allow, 1);

#endif // _base_hpp