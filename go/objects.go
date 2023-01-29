package main

import (
    "fmt"
    "reflect"
    "runtime"
    "strings"
    "sort"
)

// https://stackoverflow.com/questions/7052693/how-to-get-the-name-of-a-function-in-go
func GetFunctionName(i interface{}) string {
    return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

// My own integer power function
func IntPow(base int, exp int) int {
    if exp == 0 {
        return 1
    }
    res := 1
    for i := 0; i<exp; i++ {
        res *= base
    }
    return res
}

type object struct {
    typ string
    props map[string] interface{}
    propTypes map[string] string
}

func makeObject(typ string) *object {
    obj := object{typ: typ, props: make(map[string]interface{}), propTypes: make(map[string]string)}
    return &obj
}

func (obj *object) setProp(prop string, typ string, item interface{}) {
    obj.props[prop] = item
    obj.propTypes[prop] = typ
}

func makeGame(trump *object) *object {
    g := makeObject("game")
    g.setProp("trump", "*object", trump)
    return g
}

func makePlayer(name string) *object {
    p := makeObject("player")
    p.setProp("hand", "[]*object", make([]*object, 0))
    return p
}

func makeCard(rank string, suit string) *object {
    c := makeObject("card")
    c.setProp("rank", "string", rank)
    c.setProp("suit", "string", suit)
    return c
}

func getListItem(obj *object, prop string, idx int) *object {
    return obj.props[prop].([]*object)[idx]
}

func addListItem(obj *object, prop string, item *object) {
    obj.props[prop] = append(obj.props[prop].([]*object), item)
}

var ranks = []string{"6", "7", "8", "9", "10", "Jack", "Queen", "King", "Ace"}

func indexOf(slice []string, str string) int {
    for p, v := range slice {
        if (v == str) {
            return p
        }
    }
    return -1
}

func getProp(args []interface{}) interface{} {
    return args[0].(*object).props[args[1].(string)];
}

func greaterRank(args []interface{}) interface{} {
    return indexOf(ranks, args[0].(string)) > indexOf(ranks, args[1].(string));
}

func equal(args []interface{}) interface{} {
    return reflect.DeepEqual(args[0], args[1])
}

func cardStr(c *object) string {
    return fmt.Sprintf("%s of %s", c.props["rank"].(string), c.props["suit"].(string))
}

type fnSig func([]interface{}) interface{}

type fn struct {
    f fnSig
    args []string
    commutes bool
}

var getPropFn = fn{f: getProp, args: []string{"*object", "string"}}
var greaterRankFn = fn{f: greaterRank, args: []string{"string", "string"}}
var equalStrFn = fn{f: equal, args: []string{"string", "string"}, commutes: true}

type node struct {
    f *fn
    children []*node
    val interface{}
    bind int
}

func makeNode(f *fn, children []*node, val interface{}, bind int) *node {
    n := node{f: f, children: children, val: val, bind: bind}
    return &n
}

type pred struct {
    terms []*node
    typ string // conj, disj
    args []interface{}
    argTypes []string
}

func compat(f *fn, n *node, idx int) bool {
    if len(f.args) <= idx {
        return false
    }
    switch n.val.(type) {
        case int: return f.args[idx] == "int"
        case string: return f.args[idx] == "string"
        case bool: return f.args[idx] == "bool"
        case *object: return f.args[idx] == "*object"
    }
    return false
}

// TODO maybe check equal results (hash sig)
// TODO special list expand
func fNodes(f *fn, nodes []*node) []*node {
    cargs := make([]([]int), 0)
    for i := 0; i < len(f.args); i++ {
        cargs = append(cargs, make([]int, 0))
        for j, n := range nodes {
            if compat(f, n, i) {
                cargs[i] = append(cargs[i], j)
            }
        }
    }
    num := 1
    for i := 0; i < len(cargs); i++ {
        num *= len(cargs[i])
    }
    if num == 0 {
        return nil
    }
    checked := make(map[int]bool) // Commuting
    res := make([]*node, 0)
    for i := 0; i < num; i++ {
        mod := 1
        idcs := make([]int, len(cargs)) // Commuting
        children := make([]*node, len(cargs))
        args := make([]interface{}, len(cargs))
        for j := 0; j < len(cargs); j++ {
            idcs[j] = cargs[j][(i/mod)%len(cargs[j])]
            children[j] = nodes[idcs[j]]
            args[j] = children[j].val
            mod *= len(cargs[j])
        }
        // Commuting
        if f.commutes {
            sort.Ints(idcs)
            key := 0
            for k := 0; k < len(cargs); k++ {
                key += idcs[k]*IntPow(num, k)
            }
            // Commuting indices already evaluated
            _, ok := checked[key]
            if ok {
                continue;
            }
            checked[key] = true
        }
        r := f.f(args)
        if r == nil {
            continue
        }
        res = append(res, makeNode(f, children, r, -1))
    }
    return res
}

func fAllNodes(fs []*fn, nodes []*node) []*node {
    res := make([]*node, 0)
    for _,f := range fs {
        r := fNodes(f, nodes)
        res = append(res, r...)
    }
    return res
}

func objStr(obj *object) string {
    switch obj.typ {
        case "card": return cardStr(obj)
    }
    return obj.typ
}

func nodeStr(n *node, lvl int) string {
    str := strings.Repeat("  ", lvl)
    if n.f != nil {
        str += GetFunctionName(n.f.f)
        str += " (" + strings.Join(n.f.args, ",") + ")"
    }
    switch n.val.(type) {
        case *object: str += " " + objStr(n.val.(*object))
        default: str += fmt.Sprintf(" %v", n.val)
    }
    if n.bind > -1 {
        str += fmt.Sprintf(" BIND %d", n.bind)
    }
    for i := 0; i < len(n.children); i++ {
        str += "\n" + nodeStr(n.children[i], lvl+1)
    }
    return str
}

func fAllNodesMany(fs []*fn, nodes []*node, times int) []*node {
    for i := 0; i < times; i++ {
        res := fAllNodes(fs, nodes)
        uniq := make([]*node, 0)
        for _,n := range res {
            eq := false
            for _,m := range nodes {
                if reflect.DeepEqual(n,m) {
                    eq = true
                    break
                }
            }
            if !eq {
                uniq = append(uniq, n)
            }
        }
        fmt.Println(len(res))
        fmt.Println(len(uniq))
        nodes = append(nodes, uniq...)
    }
    return nodes
}

func getKeys(props map[string]interface{}) []string {
    keys := make([]string, len(props))
    i := 0
    for k := range props {
        keys[i] = k
        i++
    }
    return keys
}

// TODO use interface{}
func arrMinus(keys1 []string, keys2 []string) []string {
    set := make([]string, 0)
    for _,k1 := range keys1 {
        both := false
        for _,k2 := range keys2 {
            if k1 == k2 {
                both = true
                break
            }
        }
        if !both {
            set = append(set, k1)
        }
    }
    return set
}

func arrInt(keys1 []string, keys2 []string) []string {
    set := make([]string, 0)
    for _,k1 := range keys1 {
        for _,k2 := range keys2 {
            if k1 == k2 {
                set = append(set, k1)
                break
            }
        }
    }
    return set
}

// TODO difference between states
// Closure returns path to all different nodes
// Returns nil when there aren't any differences left
func diffObjects(obj1 *object, obj2 *object) []string {
    keys1 := getKeys(obj1.props)
    keys2 := getKeys(obj2.props)
    //diff1 := arrMinus(keys1, keys2)
    //diff2 := arrMinus(keys2, keys1)
    return diff1
}

func nodeFromPath(obj *object, path []string) *node {
    n := makeNode(nil, nil, obj, -1)
    for _,s := range path {
        sn := makeNode(nil, nil, s, -1)
        val := getProp([]interface{}{n.val, sn.val})
        n = makeNode(&getPropFn, []*node{n, sn}, val, -1)
    }
    return n
}

func nodeHasPath(n *node, path []string) bool {
    return false;
}

/*func learnPred(args []interface{}, argTypes []string, game *object, history []*pred) *pred {

}*/

func main() {
    c := makeCard("8", "Hearts")
    d := makeCard("10", "Hearts")
    g := makeGame(d)
    nrank := makeNode(nil, nil, "rank", -1)
    nsuit := makeNode(nil, nil, "suit", -1)
    nc := makeNode(nil, nil, c, 0)
    nd := makeNode(nil, nil, d, 1)
    /*fmt.Println(compat(&equalStrFn, nrank, 0))
    fmt.Println(compat(&equalStrFn, nsuit, 2))
    fmt.Println(compat(&equalStrFn, nc, 0))
    fmt.Println(cardStr(nc.val.(*object)))
    fmt.Println(cardStr(nd.val.(*object)))*/
    //nodes := fNodes(&equalStrFn, []*node{nrank, nsuit, nc, nd})
    fAllNodesMany([]*fn{&equalStrFn, &getPropFn}, []*node{nrank, nsuit, nc, nd}, 3)
    //fmt.Println(nodeStr(nodes[0], 0))
    /*for _,n := range nodes {
        fmt.Println(nodeStr(n, 0))
    }*/
    n := nodeFromPath(g, []string{"trump", "rank"})
    fmt.Println(nodeStr(n, 0))
    kc := getKeys(c.props)
    kd := getKeys(d.props)
    kg := getKeys(g.props)
    fmt.Println(arrMinus(kc, kg))
    fmt.Println(arrInt(kc, kd))
}
