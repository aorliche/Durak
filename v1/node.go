package main

import (
    "reflect"
    "sort"
)

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

// Do not re-evaluate with same arguments, even if reached by different path
// New evidence will cause suppression of hash disallow
func (n *node) addToRes(res *[]*node, hashes map[uint32]bool, disreg map[uint32]bool) {
    h := Hash(n)
    _, ok := hashes[h]
    _, dis := disreg[h]
    if ok && !dis {
        return
    }
    *res = append(*res, n)
    if !dis {
        hashes[h] = true
    }
}

func fNodesExpandLists(nodes []*node, hashes map[uint32]bool, disreg map[uint32]bool) []*node {
    res := make([]*node, 0)
    for _,n := range nodes {
        if reflect.TypeOf(n.val).Kind() != reflect.Slice {
            continue
        }
        switch n.val.(type) {
            case []*object: {
                for _,item := range n.val.([]*object) {
                    nn := makeNode(&expandListFn, []*node{n}, item, -1)
                    nn.addToRes(&res, hashes, disreg)
                }
            }
            default: panic("unknown type")
        }
    }
    return res
}

func fNodesExpandProps(nodes []*node, hashes map[uint32]bool, disreg map[uint32]bool) []*node {
    res := make([]*node, 0)
    for _,n := range nodes {
        switch n.val.(type) {
            case *object: {
                obj := n.val.(*object)
                for key := range obj.props {
                    kn := makeNode(nil, nil, key, -1)
                    rn := makeNode(&getPropFn, []*node{n, kn}, obj.props[key], -1)
                    rn.addToRes(&res, hashes, disreg)
                }
            }
        }
    }
    return res
}

// TODO equality compare same type of properties... (i.e. concept or type)
// Only returns new nodes
func fNodes(f *fn, nodes []*node, hashes map[uint32]bool, disreg map[uint32]bool) []*node {
    if (f == &expandPropsFn) {
        return fNodesExpandProps(nodes, hashes, disreg)
    } else if (f == &expandListFn) {
        return fNodesExpandLists(nodes, hashes, disreg)
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
        // Special case for equals: never compare same sequence of nodes
        if f == &equalStrFn {
            if children[0].val == children[1].val {
                continue
            }
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
                continue
            }
            checked[key] = true
        }
        r := f.f(args)
        // Nil means not suitable somehow
        if r == nil {
            continue
        }
        nn := makeNode(f, children, r, -1)
        nn.addToRes(&res, hashes, disreg)
    }
    return res
}

func fAllNodes(fs []*fn, nodes []*node, hashes map[uint32]bool, disreg map[uint32]bool) []*node {
    res := make([]*node, 0)
    for _,f := range fs {
        r := fNodes(f, nodes, hashes, disreg)
        res = append(res, r...)
    }
    return res
}

// TODO back to deep equals or hash equals?
func fAllNodesMany(fs []*fn, nodes []*node, times int, disreg map[uint32]bool) []*node {
    hashes := make(map[uint32]bool)
    for i := 0; i < times; i++ {
        res := fAllNodes(fs, nodes, hashes, disreg)
        nodes = append(nodes, res...)
    }
    return nodes
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

func getNodeArgs(nodes []*node, args []*node) []*nodeArgs {
    res := make([]*nodeArgs, 0)
    for _,n := range nodes {
        created := false
        for i,sn := range args {
            if n.hasSubnode(sn) {
                if !created {
                    res = append(res, &nodeArg{n :n, args: make([]int,0)})
                }
                n.args = append(n.args, i)
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

func (n *node) eval() interface{} {
    if n.children != nil {
        args := make([]interface{}, 0)
        for _,c := range n.children {
            args = append(args, c.eval())
        }
        n.val = f.f(args)
        return n.val
    } else {
        return n.val
    }
}

func (n *node) equals(other *node) bool {
    if n.f != other.f {
        return false
    }
    if n.children == nil {
        if n.bind > -1 {
            return n.bind == other.bind
        }
        // Get property
        return n.val == other.val
    } 
    for i=0; i<len(n.children); i++ {
        if !n.children[i].equals(other.children[i]) {
            return false
        }
    }
    return true
}
