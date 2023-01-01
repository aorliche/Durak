#include <variant>
#include <string>
#include <vector>
#include <unordered_map>
#include <unordered_set>
#include <algorithm>
#include <iostream>
#include <functional>

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

typedef variant<bool,int,string,ObjectId,vector<ObjectId>> Property;

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

void Object::set(const string &key, const Property &val) const {
    auto &keys = keys_map[id];
    auto &vals = props_map[id];
    if (find(keys.begin(), keys.end(), key) == keys.end()) {
        keys.push_back(key);
    }
    vals[key] = val;
}

Property &Object::get(const string &key) const {
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
    for (auto it = types.begin(); it != types.end(); it++) {
        os << *it << ',';
    }
    os << " [props: ";
    for (auto it = keys.begin(); it != keys.end(); it++) {
        os << *it;
        auto &val = vals[*it];
        switch (val.index()) {
            case 0: os << ':' << get<bool>(val); break;
            case 1: os << ':' << get<int>(val); break;
            case 2: os << ':' << get<string>(val); break;
        }
        os << ", ";
    }
    os << ']';
    return os;
}

struct Card : Object {
    Card(const ObjectId &obj) : Object(obj) {}
    Card(const string &rank, const string &suit) {
        make("card");
        set("rank", rank);
        set("suit", suit);
    }
};

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
    bool compatible(const Property &p, int n) const {
        if (name == "null") return false;
        return allow(p, n);
    }
    Property operator()(const vector<Property> &args) const {
        return eval(args);
    }
};

// Missing ObjectId and vector<ObjectId>
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

void default_pfn(ostream &os, const Property &p) {
    switch(p.index()) {
        case 0: os << get<bool>(p); break;
        case 1: os << get<int>(p); break;
        case 2: os << get<string>(p); break;
        case 3: os << Object(get<ObjectId>(p)); break;
        case 4: os << "Vector"; break;
    }
}

typedef int NodeId;

struct Node {
    NodeId id;
    static NodeId idcount;
    Property res;
    Function fn;
    vector<NodeId> parents;
    Node() : id(idcount++) {}
    Node(Property p) : id(idcount++), res(p) {}
    Node(Function f) : id(idcount++), fn(f) {}
    Node(const Node &n) : id(n.id), res(n.res), fn(n.fn), parents(n.parents) {} 
    static Node &get(NodeId id);
    Property operator()() {
        vector<Property> objs;
        for (auto it = parents.begin(); it != parents.end(); it++) {
            objs.push_back(Node::get(*it).res);
        }
        res = fn(objs);
        return res;
    }
    bool operator==(const Node &other) const {
        if (!(res == other.res 
            and fn == other.fn 
            and parents.size() == other.parents.size())) 
            return false;
        for (int i=0; i<parents.size(); i++) {
            if (Node::get(parents[i]).res != Node::get(other.parents[i]).res) 
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
        os << fn.name;
        pfn(os, res);
        os << endl;
        for (auto it = parents.begin(); it != parents.end(); it++) {
            Node::get(*it).print(os, lvl+1, pfn);
        }
    }
};

int Node::idcount = 0;
unordered_map<NodeId,Node> node_map;

Node &Node::get(NodeId id) {
    return node_map[id];
}

struct node_hasher {
    size_t operator()(const NodeId &nid) const {
        auto &n = Node::get(nid);
        size_t h = hash<string>{}(n.fn.name) + hash_prop(n.res);
        for (auto it = n.parents.begin(); it != n.parents.end(); it++) {
            auto &par = Node::get(*it);
            auto &res = par.res;
            auto &name = par.fn.name;
            h += hash<string>{}(name) + hash_prop(res);
            h >>= 8;
        }
        return h;
    }
};

template <typename T>
size_t index_of(const vector<T> &vec, const T &item) {
    for (size_t i=0; i<vec.size(); i++) {
        if (vec[i] == item) {
            return i;
        }
    }
    return -1;
}

vector<string> ranks{"6","7","8","9","10","Jack","Queen","King","Ace"};

Property &get_helper(const Property &p, const string &s) {
    return Object(get<ObjectId>(p)).get(s);
}

// Functions
bool higher_rank_eval(const vector<Property> &objs) {
    return index_of(ranks, get<string>(get_helper(objs[0], "rank")))
        > index_of(ranks, get<string>(get_helper(objs[1], "rank")));
}

bool same_suit_eval(const vector<Property> &objs)  {
    return get_helper(objs[0], "suit") == get_helper(objs[1], "suit");
}

Property get_property_eval(const vector<Property> &objs) {
    return get_helper(objs[0], get<string>(objs[1]));
}

bool two_cards(const Property &p, int n) {
    return n < 2 and p.index() == 3 and Object(get<ObjectId>(p)).is("card");
}

bool get_property_allow(const Property &p, int n) {
    return (n == 0 and p.index() == 3 ) or (n == 1 and p.index() == 2);
}

bool expand_vec_allow(const Property &p, int n) {
    return (n == 0 and p.index() == 4) or (n == 1 and p.index() == 2);
}

Function higher_rank("higher_rank", higher_rank_eval, two_cards, 2);
Function same_suit("same_suit", same_suit_eval, two_cards, 2);
Function get_property("get_property", get_property_eval, get_property_allow, 2);
Function expand_vec("expand_vec", get_property_eval, expand_vec_allow, 2);


void expand(
    vector<Function> &fns, 
    vector<NodeId> &nodes, 
    int depth,
    unordered_set<NodeId, node_hasher> &sigs) 
{
    for (int iter=0; iter<depth; iter++) {
        // NOTE! Must have a separate newnodes for each iteration
        vector<NodeId> newnodes;
        for (size_t i=0; i<fns.size(); i++) {
            int nargs = fns[i].nargs;
            // Find compatible nodes for jth arg
            vector<vector<NodeId>> compat_nodes(nargs);
            size_t psetsize = 1;
            for (size_t j=0; j<nargs; j++) {
                for (size_t k=0; k<nodes.size(); k++) {
                    if (fns[i].compatible(Node::get(nodes[k]).res, j)) {
                        compat_nodes[j].push_back(nodes[k]);
                    }
                }
                psetsize *= compat_nodes[j].size();
                // cout << 'a' << psetsize << endl;
                if (psetsize == 0) break;
            }
            // cout << 'b' << psetsize << endl;
            if (psetsize == 0) 
                continue;
            // Evaluate compatible nodes and store results
            for (size_t j=0; j<psetsize; j++) {
                Node n(fns[i]);
                node_map[n.id] = n;
                size_t jj = j;
                for (size_t k=0; k<nargs; k++) {
                    size_t sz = compat_nodes[k].size();
                    size_t kk = jj%sz;
                    jj /= sz;
                    n.parents.push_back(compat_nodes[k][kk]);
                }
                // Skip shallow result
                if (sigs.count(n.id) > 0) {
                    node_map.erase(n.id);
                    continue;
                } else {
                    sigs.insert(n.id);
                }
                Property res = n();
                // Special code for get_property(vector<ObjectId>)
                if (res.index() == 4) {
                    vector<ObjectId> &objs = get<4>(res);
                    for (auto it = objs.begin(); it != objs.end(); it++) {
                        Node nn(n.fn);
                        nn.res = Property(*it);
                        nn.parents = n.parents;
                        node_map[nn.id] = nn;
                        newnodes.push_back(nn.id);
                    }
                }
            }
        }
        cout << "Iteration " << iter << " added " << newnodes.size() << endl;
        cout << "Sigsize: " << sigs.size() << endl;
        for (size_t n=0; n<newnodes.size(); n++) {
            Node::get(newnodes[n]).print(cout);
            nodes.push_back(newnodes[n]);
        }
    }
    // for (size_t n=0; n<nodes.size(); n++)
    //     nodes[n].print(cout);
}