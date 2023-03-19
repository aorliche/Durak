package main

import (
    "fmt"
    "reflect"
)

// Check that every item in subset is in set
// i.e. it is a subset
func IsSubset(sub []*Object, set []*Object) bool {
    start:
    for _,a := range sub {
        for _,b := range set {
            if reflect.DeepEqual(a,b) {
                continue start
            }
        }
        return false
    }
    return true
}

// All nodes must give bool results
// One arg for now
func SolveHelper(inp []*Object, arg *Object, game *Object, n *Node) []*Object {
    res := make([]*Object, 0)
    for _,a0 := range inp {
        n.BindArgs(a0, arg, game)
        if n.Eval().Val.(bool) {
            res = append(res, a0)
        }
    }
    return res
}

func Solve(inp []*Object, arg *Object, game *Object, out []*Object) []*Node {
    in := MakeNode(nil, nil, inp[0], 0)
    an := MakeNode(nil, nil, arg, 1)
    gn := MakeNode(nil, nil, game, 2)
    nodesStart := [3]*Node{in,an,gn}
    nodes := nodesStart[:]
    fns := [3]*Fn{&ExpandPropsFn, &EqualStrFn, &GreaterRankFn}
    hashes := make(map[uint32]bool)
    for i:=0; i<3; i++ {
        res := ApplyFnsNodes(fns[:], nodes[:], hashes)
        bnodes := GetBoolNodes(res)
        for _,bn := range bnodes {
            sub := SolveHelper(inp, arg, game, bn)
            if len(sub) == 0 {
                continue
            }
            if IsSubset(sub,out) {
                fmt.Println("subset found")
                fmt.Println(bn.ToString())
            }
        }
        nodes = append(nodes, res...)
    }
    return nodes
}
