package main

//import "fmt"

func evalCombo(idx int, nodes []*node) bool {
   switch len(nodes) {
       case 1: return evalCombo1(idx, nodes)
       case 2: return evalCombo2(idx, nodes)
       case 3: return evalCombo3(idx, nodes)
       case 4: return evalCombo4(idx, nodes)
       default: return false
   }
}

func evalCombo1(idx int, nodes []*node) bool {
    return ((idx & 1) == 1) != nodes[0].val.(bool)
}

func evalCombo2(idx int, nodes []*node) bool {
    n0 := (idx & 1) == 1
    n1 := ((idx >> 1) & 1) == 1
    r0 := nodes[0].val.(bool)
    r1 := nodes[1].val.(bool)
    split := idx >> 2
    if split == 0 {
        return n0 != r0 && n1 != r1
    } else {
        return n0 != r0 || n1 != r1
    }
}

func evalCombo3(idx int, nodes []*node) bool {
    n0 := (idx & 1) == 1
    n1 := ((idx >> 1) & 1) == 1
    n2 := ((idx >> 2) & 1) == 1
    r0 := nodes[0].val.(bool)
    r1 := nodes[1].val.(bool)
    r2 := nodes[2].val.(bool)
    split := idx >> 3
    if split == 0 {
        return n0 != r0 && n1 != r1 && n2 != r2
    } else if split == 1 {
        return n0 != r0 || (n1 != r1 && n2 != r2)
    } else {
        return n0 != r0 || n1 != r1 || n2 != r2
    }
}
func evalCombo4(idx int, nodes []*node) bool {
    n0 := (idx & 1) == 1
    n1 := ((idx >> 1) & 1) == 1
    n2 := ((idx >> 2) & 1) == 1
    n3 := ((idx >> 3) & 1) == 1
    r0 := nodes[0].val.(bool)
    r1 := nodes[1].val.(bool)
    r2 := nodes[2].val.(bool)
    r3 := nodes[3].val.(bool)
    split := idx >> 4
    if split == 0 {
        return n0 != r0 && n1 != r1 && n2 != r2 && n3 != r3
    } else if split == 1 {
        return n0 != r0 || (n1 != r1 && n2 != r2 && n3 != r3)
    } else if split == 2 {
        return (n0 != r0 && n1 != r1) || (n2 != r2 && n3 != r3)
    } else if split == 3 {
        return (n0 != r0) || (n1 != r1) || (n2 != r2 && n3 != r3) 
    } else {
        return n0 != r0 || n1 != r1 || n2 != r2 || n3 != r3
    }
}

func hasNil(n int, nodes []*node) bool {
    for i := 0; i<n; i++ {
        if nodes[i] == nil {
            return true
        }
    }
    return false
}

func hasMissingArgs(nargs int, nodes []*node) bool {
    got := make([]bool, nargs)
    for _,n := range nodes {
        for _,i := range n.args {
            got[i] = true
        }
    }
    for i := 0; i<nargs; i++ {
        if !got[i] {
            return true
        }
    }
    return false
}

func createIndex(n int, nNodes int) *index {
    return &index{nTerms: n, nNodes: nNodes, curIdx: 0, idcs: [4]int{0,0,0,0}}
}

func (idx *index) done() bool {
    return idx.curIdx == idx.nTerms
}

func (idx *index) inc() {
    for idx.curIdx < idx.nTerms && idx.idcs[idx.curIdx] == idx.nNodes-1 {
        idx.idcs[idx.curIdx] = 0
        idx.curIdx++
    }
    if idx.curIdx < idx.nTerms {
        idx.idcs[idx.curIdx]++
        idx.curIdx = 0
    }
}

func (idx *index) getCombo(row []*node) []*node {
    nodes := make([]*node, 0)
    for i := 0; i < idx.nTerms; i++ {
        nodes = append(nodes, row[idx.idcs[i]])
    }
    return nodes
}

func subIdxFromTerms(nTerms int) int {
    switch nTerms {
        case 1: return 2
        case 2: return 8
        case 3: return 24
        case 4: return 80
        default: return -1
    }
}

// Name, argTypes, history to be filled in later
func satisfy(nargs int, exs []*example, table [][]*node) *pred {
    for n := 1; n <= 4; n++ {
        subIdx := subIdxFromTerms(n)
        for idx := createIndex(n, len(table[0])); !idx.done(); idx.inc() {
            //fmt.Println(indexStr(idx))
            for i := 0; i < subIdx; i++ {
                    /*if j > 2 {
                        fmt.Println(i,subIdx)
                    }*/
                succ := true
                var nodes []*node
                for j,row := range table {
                    nodes = idx.getCombo(row)
                    ex := exs[j]
                    if hasNil(n, nodes) || hasMissingArgs(nargs, nodes) {
                        succ = false
                        break
                    }
                    if evalCombo(i, nodes) != ex.val {
                        succ = false
                        break
                    }
                }
                if succ {
                    return &pred{nodes: nodes, idx: i, name: "TEMP", argTypes: nil, exs: exs, hist: nil, fns: nil}
                }
            }
        }
    }
    return nil
}

// Combinations of nodes where each argument appears at least once
// (e.g. node 1: a1, node 2: a1,a2, node 3: a1)
// yields 1-2, 2-2 (i.e. just 2), 2-3
// Probably not incorporated into final sat algorithm
// For one, what if you need argument to appear in in multiple nodes (e.g. n1a1,n2a1,n3a1,n4a2)
/*func argComb(int nargs, row []*node) [][]int {
    sets := make([][]int, nargs)
    for i:=0; i<nargs; i++ {
        sets[i] = make([]int, 0)
    }
    for i,n := range row {
        for _,j := range n.args {
            sets[j] = append(sets[j], i)
        }
    }
    num := 1
    for i:=0; i<nargs; i++ {
        num *= len(sets[i])
    }
    res := make([][]int, num)
    for i:=0; i<num; i++ {
        mod := 1
        res[i] = make([]int, nargs) 
        for j := 0; j < nargs; j++ {
            idx = sets[j][(i/mod)%len(sets[j])]
            res[i] = append(res[i], idx)
            mod *= len(cargs[j])
        }
    }
    return res
}*/
