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

// Data types
type object struct {
    typ string
    props map[string] interface{}
    propTypes map[string] string
}

type fnSig func([]interface{}) interface{}

type fn struct {
    f fnSig
    args []string
    commutes bool
}

type node struct {
    f *fn
    children []*node
    val interface{}
    bind int
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
    i1 := indexOf(ranks, args[0].(string)) 
    i2 := indexOf(ranks, args[1].(string))
    if i1 == -1 || i2 == -1 {
        return nil
    }
    return i1 > i2
}

func equal(args []interface{}) interface{} {
    return reflect.DeepEqual(args[0], args[1])
}

func cardStr(c *object) string {
    return fmt.Sprintf("%s of %s", c.props["rank"].(string), c.props["suit"].(string))
}

// Never called
func expandList(args[] interface{}) interface{} {
    return nil
}

var getPropFn = fn{f: getProp, args: []string{"*object", "string"}}
var greaterRankFn = fn{f: greaterRank, args: []string{"string", "string"}}
var equalStrFn = fn{f: equal, args: []string{"string", "string"}, commutes: true}
var expandPropsFn = fn{f: getProp, args: []string{"*object"}}
var expandListFn = fn{f: expandList, args: []string{"[]any"}}

func makeNode(f *fn, children []*node, val interface{}, bind int) *node {
    return &node{f: f, children: children, val: val, bind: bind}
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

func fNodesExpandLists(nodes []*node) []*node {
    res := make([]*node, 0)
    for _,n := range nodes {
        if reflect.TypeOf(n.val).Kind() != reflect.Slice {
            continue
        }
        switch n.val.(type) {
            case []*object: {
                for _,item := range n.val.([]*object) {
                    res = append(res, makeNode(&expandListFn, []*node{n}, item, -1))
                }
            }
            default: panic("unknown type")
        }
    }
    return res
}

func fNodesExpandProps(nodes []*node) []*node {
    res := make([]*node, 0)
    for _,n := range nodes {
        switch n.val.(type) {
            case *object: {
                obj := n.val.(*object)
                for key := range obj.props {
                    kn := makeNode(nil, nil, key, -1)
                    rn := makeNode(&getPropFn, []*node{n, kn}, obj.props[key], -1)
                    res = append(res, rn)
                }
            }
        }
    }
    return res
}

// TODO maybe check equal results (hash sig)
// TODO special list expand
// Only returns new nodes
func fNodes(f *fn, nodes []*node) []*node {
    if (f == &expandPropsFn) {
        return fNodesExpandProps(nodes)
    } else if (f == &expandListFn) {
        return fNodesExpandLists(nodes)
    }
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

// TODO hash equals instead of deep equals
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
        //fmt.Println(len(res))
        //fmt.Println(len(uniq))
        nodes = append(nodes, uniq...)
    }
    return nodes
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

func getKeys(props map[string]interface{}) []string {
    keys := make([]string, len(props))
    i := 0
    for k := range props {
        keys[i] = k
        i++
    }
    sort.Strings(keys)
    return keys
}

// Assumes sorted keys
func strArrMinus(keys1 []string, keys2 []string) []string {
    set := make([]string, 0)
    for i,j := 0,0; i < len(keys1); i++ {
        for j < len(keys2) && keys2[j] < keys1[i] {
            j++
        }
        if j < len(keys2) && keys1[i] == keys2[j] {
            j++
            continue
        }
        set = append(set, keys1[i])
    }
    return set
}

// Assumes sorted keys
func strArrInt(keys1 []string, keys2 []string) []string {
    set := make([]string, 0)
    for i,j := 0,0; i < len(keys1); i++ {
        for j < len(keys2) && keys2[j] < keys1[i] {
            j++
        }
        if j == len(keys2) {
            break
        }
        if keys1[i] == keys2[j] {
            set = append(set, keys1[i])
            j++
        }
    }
    return set
}

// TODO difference between states
// TODO Closure returns path to all different nodes
// TODO slices as property types
func diffObjects(obj1 *object, obj2 *object, path []string) []([]string) {
    keys1 := getKeys(obj1.props)
    keys2 := getKeys(obj2.props)
    diff := strArrMinus(keys2, keys1)
    common := strArrInt(keys1, keys2)
    res := make([]([]string), 0)
    for _,key := range diff {
        npath := append(path, key)
        res = append(res, npath)
    }
    tObj := reflect.TypeOf(obj1)
    for _,key := range common {
        t1 := reflect.TypeOf(obj1.props[key])
        t2 := reflect.TypeOf(obj2.props[key])
        if t1 == t2 && t1 == tObj {
            npath := append(path, key)
            subres := diffObjects(obj1.props[key].(*object), obj2.props[key].(*object), npath)
            res = append(res, subres...)
        } else if t1 == t2 && obj1.props[key] != obj2.props[key] {
            npath := append(path, key)
            res = append(res, npath)
        }
    }
    return res
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

/*func nodeHasPath(n *node, path []string) bool {
    return false;
}*/

/*func learnPred(args []interface{}, argTypes []string, game *object, history []*pred) *pred {

}*/

type conj struct {
    terms []*node
    neg []bool
}

type disj struct {
    terms []*conj
}

type predHist struct {
    terms []*disj
    val bool
    game *object
    args []interface{}
}

type pred struct {
    terms []*disj
    name string
    argTypes []string
    hist []*predHist
}

// -1: game
type nodeWithArg struct {
    n *node
    arg int
}

func (n *node) hasSubnode(sn *node) bool {
    if n.children == nil {
        return false
    }
    for _,c := range n.children {
        if c == sn || c.hasSubnode(sn) {
            return true
        }
    }
    return false
}

func getNodesWithArgs(nodes []*node, args []*node) []*nodeWithArg {
    res := make([]*nodeWithArg, 0)
    for _,n := range nodes {
        for i,sn := range args {
            if n.hasSubnode(sn) {
                res = append(res, &nodeWithArg{n: n, arg: i})
            }
        }
    }
    return res
}

func getBoolNodes(nodes []*node) []*node {
    res := make([]*node, 0)
    for _,n := range nodes {
        if reflect.TypeOf(n.val).Kind() == reflect.Bool {
            res = append(res, n)
        }
    }
    return res
}

// TODO must use args, each node at least 1 arg or game, pred uses all args
func makePred(name string, argTypes []string, args []interface{}, game *object, fs []*fn, times int) []*node {
    //p := &pred{name: name, argTypes: argTypes, hist: make([]*predHist, 0), terms: make([]*disj, 0)}
    n := makeNode(nil, nil, game, -1)
    nargs := make([]*node, 0)
    nargs = append(nargs, n)
    for i,arg := range args {
        n = makeNode(nil, nil, arg, i)
        nargs = append(nargs, n)
    }
    nodes := fAllNodesMany(fs, nargs, times)
    nodes = getBoolNodes(nodes)
    //nodesArgs := getNodesWithArgs(nodes, nargs) 
    return nodes //nodesArgs
}

func main() {
    c := makeCard("8", "Spades")
    d := makeCard("10", "Hearts")
    g := makeGame(d)
    //g := makeObject("x factor")
    g.setProp("x", "[]*object", []*object{c, d})
    //x := makeNode(nil, nil, []bool{true, false}, -1)
    //h := makeGame(c)
    //nrank := makeNode(nil, nil, "rank", -1)
    //nsuit := makeNode(nil, nil, "suit", -1)
    //nc := makeNode(nil, nil, c, 0)
    //nd := makeNode(nil, nil, d, 1)
    /*fmt.Println(compat(&equalStrFn, nrank, 0))
    fmt.Println(compat(&equalStrFn, nsuit, 2))
    fmt.Println(compat(&equalStrFn, nc, 0))
    fmt.Println(cardStr(nc.val.(*object)))
    fmt.Println(cardStr(nd.val.(*object)))*/
    //nodes := fNodes(&equalStrFn, []*node{nrank, nsuit, nc, nd})
    //fAllNodesMany([]*fn{&equalStrFn, &getPropFn}, []*node{nrank, nsuit, nc, nd}, 3)
    //fmt.Println(nodeStr(nodes[0], 0))
    /*for _,n := range nodes {
        fmt.Println(nodeStr(n, 0))
    }*/
    /*n := nodeFromPath(g, []string{"trump", "rank"})
    fmt.Println(nodeStr(n, 0))
    kc := getKeys(c.props)
    kd := getKeys(d.props)
    kg := getKeys(g.props)
    fmt.Println(strArrMinus(kc, kg))
    fmt.Println(strArrMinus(kg, kc))
    fmt.Println(strArrMinus(kc, kd))
    fmt.Println(strArrInt(kc, kd))
    fmt.Println(strArrInt(kg, kc))
    fmt.Println(kg)
    do := diffObjects(g,h,make([]string,0))
    fmt.Println(do)
    fmt.Println(nodeStr(nodeFromPath(h, do[0]), 0))*/
    nodes := makePred("beats", []string{"card", "card"}, []interface{}{c,d}, g, []*fn{&expandListFn, &expandPropsFn, &equalStrFn, &greaterRankFn}, 6) 
    for _,n := range nodes {
        fmt.Println(nodeStr(n, 0))
    }
}
