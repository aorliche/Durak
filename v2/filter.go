package main

import (
    "fmt"
    "sort"
    "strings"
)

// Check that every item in subset is in set
// i.e. it is a subset
func IsSubset(sub []*Object, set []*Object) bool {
    for _,a := range sub {
        if IndexOf(set, a) == -1 {
            return false
        }
    }
    return true
}

func Intersection(a []*Object, b []*Object) []*Object {
    res := make([]*Object, 0)
    for _,aa := range a {
        if IndexOf(b, aa) != -1 {
            res = append(res, aa)
        }
    }
    return res
}

func Union(a []*Object, b []*Object) []*Object {
    res := a[:]
    for _,bb := range b {
        if IndexOf(res, bb) == -1 {
            res = append(res, bb)
        }
    }
    return res
}

func ExpandNodes(nodes []*Node, n int) []*Node {
    fns := []*Fn{&ExpandPropsFn, &EqualStrFn, &GreaterRankFn}
    hashes := make(map[uint32]bool)
    for i:=0; i<n; i++ {
        res := ApplyFnsNodes(fns, nodes, hashes)
        nodes = append(nodes, res...)
    }
    return nodes
}

// Return
// 1. All perfect Filters or best matches
// 2. Percent score for matches
// 3. Target for perfect match
func Solve(inp [][]*Object, args []*Object, games []*Object, out [][]*Object) ([]*Filter, float64, int) {
    in := MakeNode(nil, nil, inp[0][0], 0)
    an := MakeNode(nil, nil, args[0], 1)
    gn := MakeNode(nil, nil, games[0], 2)
    nodes := []*Node{in,an,gn}
    bnodes := GetBoolNodes(ExpandNodes(nodes, 3))
    filts := make([]*Filter, 0)
    for _,n := range bnodes {
        fi := Filter{N: n, Out: make([][]*Object, len(inp))}
        good := false
        for i:=0; i<len(inp); i++ {
            sub := SolveNode(inp[i], args[i], games[i], n)
            fi.Out[i] = sub
            // Hueristic that we require some negative and positive selection 
            // from a filter atom across all examples
            if len(sub) != 0 && len(sub) != len(inp[i]) {
                good = true
            }
        }
        if !good {
            continue
        }
        filts = append(filts, &fi)
    }
    // Conjunctions and disjunctions
    singles := ConjunctionSingle(filts, out)
    doubles := Conjunction(filts, out)
    both := Cat(singles, doubles)
    dis := Disjunction(both, out)
    // Ranking
    res, _, tgt := RankFilters(Cat(dis, both), out)
    // Get best results
    best := 0
    if len(res) > 0 {
        best = res[0].NumObjects
    }
    topRes := make([]*Filter, 0)
    for _,fi := range res {
        if fi.NumObjects < best {
            break
        }
        topRes = append(topRes, fi)
    }
    return topRes, float64(best)/float64(tgt), tgt
}

// All nodes must give bool results
// One arg for now
func SolveNode(inp []*Object, arg *Object, game *Object, n *Node) []*Object {
    res := make([]*Object, 0)
    for _,a0 := range inp {
        n.BindArgs(a0, arg, game)
        if n.Eval().Val.(bool) {
            res = append(res, a0)
        }
    }
    return res
}

// Check if an atom is a subset all by itself
func ConjunctionSingle(filts []*Filter, out [][]*Object) []*Filter {
    res := make([]*Filter, 0)
    start:
    for _,fi := range filts {
        for i:=0; i<len(out); i++ {
            if !IsSubset(fi.Out[i], out[i]) {
                continue start
            }
        }
        res = append(res, fi)
    }
    return res
}

// Maximum 2 way intersections for now
func Conjunction(filts []*Filter, out [][]*Object) []*Filter {
    n := len(filts)
    m := len(out)
    res := make([]*Filter, 0)
    for i:=0; i<n; i++ {
        start:
        for j:=i+1; j<n; j++ {
            fi := Filter{C: []*Filter{filts[i], filts[j]}, Out: make([][]*Object, len(out))}
            allEmpty := true
            triviallySame := true
            for k:=0; k<m; k++ {
                sub := Intersection(filts[i].Out[k], filts[j].Out[k])
                fi.Out[k] = sub
                // Same logic as before about no empty results for all examples
                if len(sub) > 0 {
                   allEmpty  = false
                }
                // Except now check that we also don't have "!f(A,B) AND f(B,A)"
                // TODO this does not catch GreaterRank(A,B) and !GreaterRank(B,A)
                if len(sub) != len(filts[i].Out[k]) || len(sub) != len(filts[j].Out[k]) {
                    triviallySame = false
                }
                // Now also check not being a subset of out
                if !IsSubset(sub, out[k]) {
                    continue start
                }
            }
            if allEmpty || triviallySame {
                continue
            }
            res = append(res, &fi)
        }
    }
    return res
}

// Find unions that get as close as possible to out
// Only 2 for now
// Only take unions that add to result set
func Disjunction(filts []*Filter, out [][]*Object) []*Filter {
    n := len(filts)
    m := len(out)
    res := make([]*Filter, 0)
    for i:=0; i<n; i++ {
        for j:=i+1; j<n; j++ {
            fi := Filter{D: []*Filter{filts[i], filts[j]}, Out: make([][]*Object, len(out))}
            increase := false
            for k:=0; k<m; k++ {
                sub := Union(filts[i].Out[k], filts[j].Out[k])
                fi.Out[k] = sub
                if len(sub) > len(filts[i].Out[k]) && len(sub) > len(filts[j].Out[k]) {
                    increase = true
                }
            }
            if !increase {
                continue
            }
            res = append(res, &fi)
        }
    }
    return res
}

// Out arg for future penalty? due to going beyond subsets
// Returns sorted filts argument
func RankFilters(filts []*Filter, out [][]*Object) ([]*Filter, int, int) {
    nPerf := 0
    // Get target num
    tgt := 0
    for _,sub := range out {
        tgt += len(sub)
    }
    // Need to init int field
    for _,fi := range filts {
        sum := 0
        for _,sub := range fi.Out {
            sum += len(sub)
        }
        fi.NumObjects = sum
        if sum == tgt {
            nPerf++
        }
    }
    sort.Slice(filts, func(i, j int) bool {
        return filts[i].NumObjects > filts[j].NumObjects
    })
    return filts, nPerf, tgt
}

// String rep of any of the SolverAtom types
// Node, Conjunction, Disjunction of Conjunctions
func (fi *Filter) ToStr() string {
    strs := make([]string,0)
    if fi.N != nil {
        strs = append(strs, fi.N.ToStr())
    }
    // Always 2 for now
    if fi.C != nil {
        strs = append(strs, "Conjunction")
        strs = append(strs, fi.C[0].N.ToStr())
        strs = append(strs, fi.C[1].N.ToStr())
    }
    if fi.D != nil {
        strs = append(strs, "Disjunction of Conjunctions")
        for i,d := range fi.D {
            strs = append(strs, fmt.Sprintf("%d. %s", i, d.ToStr()))
        }
    }
    for _,out := range fi.Out {
        strs = append(strs, MakeSliceObject(out).ToStr())
    }
    return strings.Join(strs, "\n")
}
