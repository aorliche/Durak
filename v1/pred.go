package main

func makePred(name string, argTypes []string, fns []*fn) *pred {
    return &pred{terms: make([]*disj, 0), name: name, artTypes: argTypes, hist: make([]*predHist, 0), fns: fns}
}

func getBindChildren[T pred|history|disj|conj|node](t *T) []*node {
    bnodes := make([]*node, 0)
    switch t.(type) {
        case *node: {
            if t.children != nil {
                for _,n := range t.children {
                    bnodes = append(bnodes, getBindChildren(n)...)
                }
            } else if t.bind != -1 {
                bnodes = append(bnodes, t)
            }
        }
        default: {
            for _,tt := range t.terms {
                bnodes = append(bnodes, getBindChildren(tt)...)
            }
        }
    }
    return bnodes
}

func eval[T pred|history|disj|conj](t *T) interface{} {
    switch t.(type) {
        case *conj: {
            for i,n := t.terms {
                res := t.neg[i] != n.eval()
                if !res {
                    return false
                }
            }
            return true
        }
        case *disj: {
            for _,tt := t.terms {
                if eval(tt) {
                    return true
                }
            }
            return false
        }
        default: {
            return eval(t.dis)
        }
    }
}

func evalPred(p *pred, args []interface{}) bool {
    bnodes := getBindChildren(p)
    for i,n := range bnodes {
        n.val = args[i]
    }
    return eval(p)
}

// TODO add history
// TODO must use args, each node at least 1 arg, pred uses all args
// Return value is true if predicate changed, false otherwise
func makeTable(exs []*example, times int) [][]*nodeArgs {
    ncount := 0
    hash2idx := make(map[uint32]int)
    idx2node := make(map[[2]int]*node)
    for i,ex := range exs {
        disreg := make(map[uint32]bool)
        nodes := expandNodes(ex, times, disreg)
        for _,n := range nodes {
            h := Hash(n.n)
            j,ok := hash2idx[h]
            if !ok {
                j = ncount++
                hash2idx[h] = j
            } 
            idx2node[[2]int{i,j}] = n
        }
    }
    table := make([][]*node, len(ex))
    for i=0; i<len(exs); i++ {
        table[i] = make([]*node, ncount)
    }
    for pair,n := range idx2node {
        table[pair[0]][pair[1]] = n
    }
    return table
}

func makeExample(val bool, args []interface{}) *example {
    return &example{val: val, args: args}
}

func expandNodes(ex *example, times int, disreg map[uint32]bool) []*nodeArgs {
    nargs := make([]*node, 0)
    for i,arg := range ex.args {
        n := makeNode(nil, nil, arg, i)
        nargs = append(nargs, n)
    }
    nodes := fAllNodesMany(fs, nargs, times, disreg)
    nodes = getBoolNodes(nodes)
    return getNodesWithArgs(nodes, nargs)
}
