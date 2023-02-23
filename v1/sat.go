package main

import "fmt"

func (idx *Index) ToStr() string {
    s := fmt.Sprintf("%d %d %d", idx.NumTerms, idx.NumNodes, idx.CurIdx)
    s += fmt.Sprint(idx.Idcs)
    return s
}

func SatStr(n int, idx int) string {
    switch n {
        case 1: return SatStr1(idx)
        case 2: return SatStr2(idx)
        case 3: return SatStr3(idx)
        case 4: return SatStr4(idx)
        default: panic("bad")
    }
}

func SatStr1(idx int) string {
    n0 := Ternary((idx & 1) == 1, "~", "")
    return n0 + "A"
}

func SatStr2(idx int) string {
    n0 := Ternary((idx & 1) == 1, "~", "")
    n1 := Ternary(((idx >> 1) & 1) == 1, "~", "")
    split := idx >> 2
    if split == 0 {
        return n0 + "A" + n1 + "B"
    } else {
        return n0 + "A+" + n1 + "B"
    }
}

func SatStr3(idx int) string {
    n0 := Ternary((idx & 1) == 1, "~", "")
    n1 := Ternary(((idx >> 1) & 1) == 1, "~", "")
    n2 := Ternary(((idx >> 2) & 1) == 1, "~", "")
    split := idx >> 3
    if split == 0 {
        return n0 + "A" + n1 + "B" + n2 + "C"
    } else if split == 1 {
        return n0 + "A+(" + n1 + "B" + n2 + "C)"
    } else {
        return n0 + "A+" + n1 + "B+" + n2 + "C"
    }
}

func SatStr4(idx int) string {
    n0 := Ternary((idx & 1) == 1, "~", "")
    n1 := Ternary(((idx >> 1) & 1) == 1, "~", "")
    n2 := Ternary(((idx >> 2) & 1) == 1, "~", "")
    n3 := Ternary(((idx >> 3) & 1) == 1, "~", "")
    split := idx >> 4
    if split == 0 {
        return n0 + "A" + n1 + "B" + n2 + "C" + n3 + "D"
    } else if split == 1 {
        return n0 + "A+(" + n1 + "B" + n2 + "C" + n3 + "D)"
    } else if split == 2 {
        return "(" + n0 + "A" + n1 + "B)+(" + n2 + "C" + n3 + "D)"
    } else if split == 3 {
        return n0 + "A+" + n1 + "B+(" + n2 + "C" + n3 + "D)"
    } else {
        return n0 + "A+" + n1 + "B+" + n2 + "C+" + n3 + "D"
    }
}

func EvalCombo(idx int, nodes []*Node) bool {
   switch len(nodes) {
       case 1: return EvalCombo1(idx, nodes)
       case 2: return EvalCombo2(idx, nodes)
       case 3: return EvalCombo3(idx, nodes)
       case 4: return EvalCombo4(idx, nodes)
       default: return false
   }
}

func EvalCombo1(idx int, nodes []*Node) bool {
    return ((idx & 1) == 1) != nodes[0].Val.(bool)
}

func EvalCombo2(idx int, nodes []*Node) bool {
    n0 := (idx & 1) == 1
    n1 := ((idx >> 1) & 1) == 1
    r0 := nodes[0].Val.(bool)
    r1 := nodes[1].Val.(bool)
    split := idx >> 2
    if split == 0 {
        return n0 != r0 && n1 != r1
    } else {
        return n0 != r0 || n1 != r1
    }
}

func EvalCombo3(idx int, nodes []*Node) bool {
    n0 := (idx & 1) == 1
    n1 := ((idx >> 1) & 1) == 1
    n2 := ((idx >> 2) & 1) == 1
    r0 := nodes[0].Val.(bool)
    r1 := nodes[1].Val.(bool)
    r2 := nodes[2].Val.(bool)
    split := idx >> 3
    if split == 0 {
        return n0 != r0 && n1 != r1 && n2 != r2
    } else if split == 1 {
        return n0 != r0 || (n1 != r1 && n2 != r2)
    } else {
        return n0 != r0 || n1 != r1 || n2 != r2
    }
}
func EvalCombo4(idx int, nodes []*Node) bool {
    n0 := (idx & 1) == 1
    n1 := ((idx >> 1) & 1) == 1
    n2 := ((idx >> 2) & 1) == 1
    n3 := ((idx >> 3) & 1) == 1
    r0 := nodes[0].Val.(bool)
    r1 := nodes[1].Val.(bool)
    r2 := nodes[2].Val.(bool)
    r3 := nodes[3].Val.(bool)
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

func HasNil(n int, nodes []*Node) bool {
    for i := 0; i<n; i++ {
        if nodes[i] == nil {
            return true
        }
    }
    return false
}

func HasMissingArgs(nargs int, nodes []*Node) bool {
    got := make([]bool, nargs)
    for _,n := range nodes {
        for _,i := range n.Args {
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

func CreateIndex(n int, nNodes int) *Index {
    return &Index{NumTerms: n, NumNodes: nNodes, CurIdx: 0, Idcs: [4]int{0,0,0,0}}
}

func (idx *Index) Done() bool {
    return idx.CurIdx == idx.NumTerms
}

func (idx *Index) Inc() {
    for idx.CurIdx < idx.NumTerms && idx.Idcs[idx.CurIdx] == idx.NumNodes-1 {
        idx.Idcs[idx.CurIdx] = 0
        idx.CurIdx++
    }
    if idx.CurIdx < idx.NumTerms {
        idx.Idcs[idx.CurIdx]++
        idx.CurIdx = 0
    }
}

func (idx *Index) GetCombo(row []*Node) []*Node {
    nodes := make([]*Node, 0)
    for i := 0; i < idx.NumTerms; i++ {
        nodes = append(nodes, row[idx.Idcs[i]])
    }
    return nodes
}

func SubIdxFromTerms(nTerms int) int {
    switch nTerms {
        case 1: return 2
        case 2: return 8
        case 3: return 24
        case 4: return 80
        default: return -1
    }
}

func Satisfy(nargs int, exs []*Example, table [][]*Node) ([]*Node, int) {
    for n := 1; n <= 4; n++ {
        subIdx := SubIdxFromTerms(n)
        for idx := CreateIndex(n, len(table[0])); !idx.Done(); idx.Inc() {
            for i := 0; i < subIdx; i++ {
                succ := true
                var nodes []*Node
                for j,row := range table {
                    nodes = idx.GetCombo(row)
                    ex := exs[j]
                    if HasNil(n, nodes) || HasMissingArgs(nargs, nodes) {
                        succ = false
                        break
                    }
                    if EvalCombo(i, nodes) != ex.Val {
                        succ = false
                        break
                    }
                }
                if succ {
                    return nodes, i
                }
            }
        }
    }
    return nil, -1
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
