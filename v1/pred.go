package main

import (
    "encoding/json"
    //"fmt"
    "os"
)

func MakePred(name string, args []string) *Pred {
    return &Pred{
        Nodes: make([]*Node, 0),
        Idx: -1,
        Name: name,
        Args: args}
}

func MakeExample(val bool, args []interface{}) *Example {
    return &Example{Val: val, Args: args}
}

// Get all bound nodes
func GetBound(t interface{}) []*Node {
    bnodes := make([]*Node, 0)
    switch t.(type) {
        case *Node: {
            n := t.(*Node)
            if n.Children != nil {
                for _,m := range n.Children {
                    bnodes = append(bnodes, GetBound(m)...)
                }
            } else if n.Bind != -1 {
                bnodes = append(bnodes, n)
            }
        }
        case *Pred: {
            for _,n := range t.(*Pred).Nodes {
                bnodes = append(bnodes, GetBound(n)...)
            }
        }
    }
    return bnodes
}

// While doing satistfiability
func Eval(idx int, nodes []*Node, ex *Example) bool {
    for _,n := range nodes {
        for _,bn := range GetBound(n) {
            bn.Val = ex.Args[bn.Bind]
        }
    }
    return EvalCombo(idx, nodes)
}

func ExpandNodes(fns []*Fn, ex *Example, times int, reqArgs []int, banlist []uint32) []*Node {
    nargs := make([]*Node, 0)
    for i,arg := range ex.Args {
        n := MakeNode(nil, nil, arg, i)
        nargs = append(nargs, n)
    }
    nodes := FAllNodesMany(fns, nargs, times, banlist)
    nodes = GetRequiredNodes(GetBoolNodes(nodes), reqArgs)
    return nodes
}

// TODO add history
// Uses all args is checked in sat.go
func MakeTable(fns []*Fn, exs []*Example, times int, reqArgs []int, banlist []uint32) [][]*Node {
    ncount := 0
    hash2idx := make(map[uint32]int)
    idx2node := make(map[[2]int]*Node)
    for i,ex := range exs {
        nodes := ExpandNodes(fns, ex, times, reqArgs, banlist)
        for _,n := range nodes {
            h := n.Hash()
            j,ok := hash2idx[h]
            if !ok {
                hash2idx[h] = ncount
                j = ncount
                ncount++
            }
            idx2node[[2]int{i,j}] = n
        }
    }
    table := make([][]*Node, len(exs))
    for i:=0; i<len(exs); i++ {
        table[i] = make([]*Node, ncount)
    }
    for pair,n := range idx2node {
        table[pair[0]][pair[1]] = n
    }
    return table
}

func (p *Pred) Bind(args []interface{}) {
    // Bind this pred's args
    for _,n := range p.Nodes {
        n.BindArgs(args)
    }
}

func (p *Pred) Eval() bool {
    for _,n := range p.Nodes {
        n.Eval()
    }
    return EvalCombo(p.Idx, p.Nodes)
}

func (p *Pred) ToStr() string {
    str := p.Name + "\n"
    for _,n := range p.Nodes {
        str += n.ToStr() + "\n"
    }
    str += SatStr(len(p.Nodes), p.Idx)
    return str
}

func (p *Pred) ToRec() *PredRec {
    r := &PredRec{Name: p.Name, Idx: p.Idx, Args: p.Args, Nodes: make([]*NodeRec, 0)}
    for _,n := range p.Nodes {
        r.Nodes = append(r.Nodes, n.ToRec())
    }
    return r
}

func (r *PredRec) ToPred() *Pred {
    p := &Pred{Name: r.Name, Idx: r.Idx, Args: r.Args, Nodes: make([]*Node, 0)}
    for _,n := range r.Nodes {
        p.Nodes = append(p.Nodes, n.ToNode())
    }
    return p
}

func (r *PredRec) ToJson() []byte {
    res,_ := json.Marshal(*r)
    return res
}

func PredFromJson(jsn []byte) *Pred {
    var r PredRec
    json.Unmarshal(jsn, &r)
    return r.ToPred()
}

func (p *Pred) WriteFile(path string) {
    os.WriteFile(path, p.ToRec().ToJson(), 0644)
}

func ReadPred(path string) *Pred {
    dat, err := os.ReadFile(path)
    if err != nil {
        return nil
    }
    return PredFromJson(dat)
}
