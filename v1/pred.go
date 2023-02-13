package main

//import "fmt"

func makePred(name string, argTypes []string, fns []*fn) *pred {
    return &pred{
        nodes: make([]*node, 0), 
        idx: -1,
        name: name, 
        argTypes: argTypes, 
        exs: make([]*example, 0), 
        hist: make([]*history, 0), 
        fns: fns}
}

func makeExample(val bool, args []interface{}) *example {
    return &example{val: val, args: args}
}

// Get all bound nodes
func getBound(t interface{}) []*node {
    bnodes := make([]*node, 0)
    switch t.(type) {
        case *node: {
            n := t.(*node)
            if n.children != nil {
                for _,m := range n.children {
                    bnodes = append(bnodes, getBound(m)...)
                }
            } else if n.bind != -1 {
                bnodes = append(bnodes, n)
            }
        }
        case *pred: {
            for _,n := range t.(*pred).nodes {
                bnodes = append(bnodes, getBound(n)...)
            }
        }
    }
    return bnodes
}

// For pred or history
func eval(idx int, nodes []*node, ex *example) bool {
    for _,n := range nodes {
        for _,bn := range getBound(n) {
            i := bn.bind
            bn.val = ex.args[i]
        }
    }
    return evalCombo(idx, nodes)
}

func expandNodes(fns []*fn, ex *example, times int, reqArgs []int) []*node {
    nargs := make([]*node, 0)
    for i,arg := range ex.args {
        n := makeNode(nil, nil, arg, i)
        nargs = append(nargs, n)
    }
    nodes := fAllNodesMany(fns, nargs, times)
    nodes = getRequiredNodes(getBoolNodes(nodes), reqArgs)
    return nodes
}

// TODO add history
// Uses all args is checked in sat.go
func makeTable(fns []*fn, exs []*example, times int, reqArgs []int) [][]*node {
    ncount := 0
    hash2idx := make(map[uint32]int)
    idx2node := make(map[[2]int]*node)
    for i,ex := range exs {
        nodes := expandNodes(fns, ex, times, reqArgs)
        for _,n := range nodes {
            h := HashFuncBind(n)
            //v := Hash(n)
            //fmt.Println(v,h)
            j,ok := hash2idx[h]
            if !ok {
                hash2idx[h] = ncount
                j = ncount
                ncount++
            } 
            idx2node[[2]int{i,j}] = n
        }
        //fmt.Println("...")
    }
    table := make([][]*node, len(exs))
    for i:=0; i<len(exs); i++ {
        table[i] = make([]*node, ncount)
    }
    for pair,n := range idx2node {
        table[pair[0]][pair[1]] = n
    }
    return table
}
