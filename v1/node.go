package main

import (
    "reflect"
    "sort"
)

func makeNode(f *fn, children []*node, val interface{}, bind int) *node {
    return &node{f: f, children: children, val: val, bind: bind, args: nil}
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
func (n *node) addToResult(res *[]*node, hashes map[uint32]bool) {
    h := HashFuncBind(n)
    _, ok := hashes[h]
    if ok {
        return
    }
    *res = append(*res, n)
    hashes[h] = true
}

func fNodesExpandLists(nodes []*node, hashes map[uint32]bool) []*node {
    res := make([]*node, 0)
    for _,n := range nodes {
        if reflect.TypeOf(n.val).Kind() != reflect.Slice {
            continue
        }
        switch n.val.(type) {
            case []*object: {
                for _,item := range n.val.([]*object) {
                    nn := makeNode(&expandListFn, []*node{n}, item, -1)
                    nn.addToResult(&res, hashes)
                }
            }
            default: panic("unknown type")
        }
    }
    return res
}

func fNodesExpandProps(nodes []*node, hashes map[uint32]bool) []*node {
    res := make([]*node, 0)
    for _,n := range nodes {
        switch n.val.(type) {
            case *object: {
                obj := n.val.(*object)
                for key := range obj.props {
                    kn := makeNode(nil, nil, key, -1)
                    rn := makeNode(&getPropFn, []*node{n, kn}, obj.props[key], -1)
                    rn.addToResult(&res, hashes)
                }
            }
        }
    }
    return res
}

// TODO equality compare same type of properties... (i.e. concept or type)
// Only returns new nodes
func fNodes(f *fn, nodes []*node, hashes map[uint32]bool) []*node {
    if (f == &expandPropsFn) {
        return fNodesExpandProps(nodes, hashes)
    } else if (f == &expandListFn) {
        return fNodesExpandLists(nodes, hashes)
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
    res := make([]*node, 0)
    for i := 0; i < num; i++ {
        mod := 1
        children := make([]*node, len(cargs))
        args := make([]interface{}, len(cargs))
        for j := 0; j < len(cargs); j++ {
            idx := cargs[j][(i/mod)%len(cargs[j])]
            children[j] = nodes[idx]
            args[j] = children[j].val
            mod *= len(cargs[j])
        }
        // Special case for equals: never compare same sequence of nodes
        if f == &equalStrFn {
            if HashFuncBind(children[0]) == HashFuncBind(children[1]) {
                continue
            }
        }
        // Commuting
        if f.commutes {
            // Canonicalize arguments by sorting in order of hashes
            canon := make([]*canonArg, len(cargs))
            for j := 0; j < len(cargs); j++ {
                canon[j] = &canonArg{arg: args[j], child: children[j]}
                HashFuncBind(children[j])
            }
            sort.Slice(canon, func(i, j int) bool {
                return canon[i].child.fhash < canon[j].child.fhash
            })
            for j := 0; j < len(cargs); j++ {
                args[j] = canon[j].arg
                children[j] = canon[j].child
            }
        }
        r := f.f(args)
        // Nil means not suitable
        if r == nil {
            continue
        }
        nn := makeNode(f, children, r, -1)
        nn.addToResult(&res, hashes)
    }
    return res
}

func fAllNodes(fs []*fn, nodes []*node, hashes map[uint32]bool) []*node {
    res := make([]*node, 0)
    for _,f := range fs {
        r := fNodes(f, nodes, hashes)
        res = append(res, r...)
    }
    return res
}

func fAllNodesMany(fs []*fn, nodes []*node, times int) []*node {
    hashes := make(map[uint32]bool)
    for i := 0; i < times; i++ {
        res := fAllNodes(fs, nodes, hashes)
        //res = setArgs(res)
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

/*func (n *node) hasSubnode(sn *node) bool {
    if n.children == nil {
        return false
    }
    for _,c := range n.children {
        if c == sn || c.hasSubnode(sn) {
            return true
        }
    }
    return false
}*/

func setArgs(nodes []*node) []*node {
    for _,n := range nodes {
        n.args = make([]int, 0)
        if n.bind != -1 {
            n.args = append(n.args, n.bind)
        } else {
            setArgs(n.children)
            for _,m := range n.children {
                n.args = append(n.args, m.args...)
            }
        }
        n.args = Unique(n.args)
    }
    return nodes
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

func getRequiredNodes(nodes []*node, reqArgs []int) []*node {
    nodes = setArgs(nodes)
    res := make([]*node, 0)
    for _,n := range nodes {
        for _,j := range n.args {
            if IndexOf(reqArgs, j) != -1 {
                res = append(res, n)
                break
            }
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
        n.val = n.f.f(args)
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
    for i:=0; i<len(n.children); i++ {
        if !n.children[i].equals(other.children[i]) {
            return false
        }
    }
    return true
}
