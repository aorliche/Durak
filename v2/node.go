package main

import (
    "encoding/binary"
    //"encoding/json"
    "fmt"
    "hash/crc32"
    //"reflect"
    "sort"
    "strings"
    "unicode"
    "unicode/utf8"
)

// Hash of Node and children
func AppendNode(b []byte, n *Node) []byte {
    if n.F != nil {
        b = append(b, []byte(GetName(n.F))...)
    }
    if n.Bind != -1 {
        return binary.LittleEndian.AppendUint32(b, uint32(n.Bind))
    }
    // For Val Nodes
    if n.Children == nil {
        switch n.Val.Type {
            case "string": b = append(b, []byte(n.Val.Val.(string))...)
        }
    }
    // For intermediate Nodes
    for _,m := range n.Children {
        b = AppendNode(b, m)
    }
    return b
}

// Get (possibly cached) hash value of node
func (n *Node) Hash() uint32 {
    if n.SavHash != 0 {
        return n.SavHash
    }
    c := crc32.NewIEEE()
    b := make([]byte, 0)
    b = AppendNode(b, n)
    c.Write(b)
    n.SavHash = c.Sum32()
    return n.SavHash
}

// Create Node
func MakeNode(f *Fn, children []*Node, val *Object, bind int) *Node {
    return &Node{F: f, Children: children, Val: val, Bind: bind, Type: val.Type}
}

// ToStr
func (n *Node) ToStr() string {
    return n.ToStr2(0)
}

// ToStr helper
func (n *Node) ToStr2(lvl int) string {
    str := strings.Repeat("  ", lvl)
    if n.F != nil {
        str += GetName(n.F)
        str += " (" + strings.Join(n.F.Args, ",") + ") "
    }
    if n.BindTree != nil {
        str += fmt.Sprint(n.BindTree) + " "
    }
    str += n.Val.ToStr() + " "
    if n.Bind > -1 {
        str += fmt.Sprintf("BIND %d", n.Bind)
    }
    for i := 0; i < len(n.Children); i++ {
        str += "\n" + n.Children[i].ToStr2(lvl+1)
    }
    return str
}

// Check if type begins with uppercase letter
func (n *Node) IsObjectType() bool {
    r, _ := utf8.DecodeRuneInString(n.Type)
    return unicode.IsUpper(r)
}

// Node compatible with arg type
func (n *Node) Compat(arg string) bool {
    if arg == n.Type {
        return true
    }
    if arg == "Object" || strings.Contains(arg, "Object") {
        return n.IsObjectType()
    }
    if strings.Contains(arg, "|") {
        for _,allow := range strings.Split(arg, "|") {
            if allow == n.Val.Type {
                return true
            }
        }
    }
    return false
}

// Adding to result
func (n *Node) AddToResult(res *[]*Node, hashes map[uint32]bool) {
    h := n.Hash()
    _, ok := hashes[h]
    if ok {
        return
    }
    *res = append(*res, n)
    hashes[h] = true
}

// Bind args
func (n *Node) BindArgs(args ...*Object) {
    if n.Bind != -1 {
        n.Val = args[n.Bind]
    }
    for _,m := range n.Children {
        m.BindArgs(args...)
    }
}

// Eval
func (n *Node) Eval() *Object {
    if len(n.Children) > 0 {
        args := make([]*Object, 0)
        for _,c := range n.Children {
            args = append(args, c.Eval())
        }
        // Only Atoms for now
        n.Val = n.F.A(args)
        n.Type = n.Val.Type
        return n.Val
    } else {
        return n.Val
    }
}

// Apply function to Nodes, skipping duplicates
// Only return new Nodes
func ApplyFnNodes(f *Fn, nodes []*Node, hashes map[uint32]bool) []*Node {
    // Special behavior
    if (f == &ExpandPropsFn) {
        return ApplyExpandProps(nodes, hashes)
    } 
    // Get compatible args for building powerset later
    cargs := make([][]int, 0)
    pSetSize := 1
    for _,a := range f.Args {
        c := make([]int, 0)
        for i,n := range nodes {
            if n.Compat(a) {
                c = append(c, i)
            }
        }
        pSetSize *= len(c)
        if pSetSize == 0 {
            return nil
        }
        cargs = append(cargs, c)
    }
    // Powerset
    res := make([]*Node, 0)
    for i := 0; i < pSetSize; i++ {
        mod := 1
        children := make([]*Node, len(cargs))
        // Tricky code
        for j := 0; j < len(cargs); j++ {
            idx := cargs[j][(i/mod)%len(cargs[j])]
            children[j] = nodes[idx]
            mod *= len(cargs[j])
        }
        // Special case for equals and greater rank: never compare same sequence of nodes
        if f == &EqualStrFn || f == &GreaterRankFn {
            if children[0].Hash() == children[1].Hash() {
                continue
            }
        }
        // Commuting
        if f.Commutes {
            // Canonicalize arguments by sorting in order of hashes
            sort.Slice(children, func(i, j int) bool {
                return children[i].Hash() < children[j].Hash()
            })
        }
        // Evaluate Fn
        // Only Atoms for now
        args := make([]*Object, 0)
        for _,n := range children {
            args = append(args, n.Val)
        }
        r := f.A(args)
        // Nil means not suitable 
        // TODO: does anything ever return nil?
        if r == nil {
            continue
        }
        n := MakeNode(f, children, r, -1)
        n.AddToResult(&res, hashes)
        // Bool nodes also get their negation
        if n.Type == "bool" {
            m := MakeNode(&NotFn, []*Node{n}, Not([]*Object{r}), -1)
            m.AddToResult(&res, hashes)
        }
    }
    return res
}

// Special behavior
func ApplyExpandProps(nodes []*Node, hashes map[uint32]bool) []*Node {
    res := make([]*Node, 0)
    for _,n := range nodes {
        if !n.IsObjectType() {
           continue 
        }
        for key := range n.Val.Props {
            kn := MakeNode(nil, nil, MakeStringObject(key), -1)
            rn := MakeNode(&GetPropFn, []*Node{n, kn}, n.Val.Props[key], -1)
            rn.AddToResult(&res, hashes)
        }
    }
    return res
}

// Apply many functions
func ApplyFnsNodes(fs []*Fn, nodes []*Node, hashes map[uint32]bool) []*Node {
    res := make([]*Node, 0)
    for _,f := range fs {
        r := ApplyFnNodes(f, nodes, hashes)
        res = append(res, r...)
    }
    return res
}

// Get bool nodes
func GetBoolNodes(nodes []*Node) []*Node {
    res := make([]*Node, 0)
    for _,n := range nodes {
        if n.Type == "bool" {
            res = append(res, n)
        }
    }
    return res
}
