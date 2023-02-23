package main

import (
    "encoding/binary"
    "encoding/json"
    "fmt"
    "hash/crc32"
    "reflect"
    "sort"
    "strings"
)

func AppendNode(b []byte, n *Node) []byte {
    if n.F != nil {
        b = append(b, []byte(GetFunctionName(n.F.F))...)
    }
    if n.Bind != -1 {
        return binary.LittleEndian.AppendUint32(b, uint32(n.Bind))
    }
    // For props
    if n.Children == nil && n.Bind == -1 {
        switch n.Val.(type) {
            case string: b = append(b, []byte(n.Val.(string))...)
        }
    }
    for _,m := range n.Children {
        b = AppendNode(b, m)
    }
    return b
}

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

func (n *Node) ToStr2(lvl int) string {
    str := strings.Repeat("  ", lvl)
    if n.F != nil {
        str += GetFunctionName(n.F.F)
        str += " (" + strings.Join(n.F.Args, ",") + ") "
    }
    if n.Args != nil {
        str += fmt.Sprint(n.Args) + " "
    }
    switch n.Val.(type) {
        case *Object: str += n.Val.(*Object).ToStr()
        case []*Object: {
            strSlice := make([]string, 0)
            for _,obj := range n.Val.([]*Object) {
                strSlice = append(strSlice, obj.ToStr())
            }
            str += "[" + strings.Join(strSlice, ", ") + "]"
        }
        default: str += fmt.Sprintf("%v", n.Val)
    }
    if n.Bind > -1 {
        str += fmt.Sprintf(" BIND %d", n.Bind)
    }
    for i := 0; i < len(n.Children); i++ {
        str += "\n" + n.Children[i].ToStr2(lvl+1)
    }
    return str
}

func (n *Node) ToStr() string {
    return n.ToStr2(0)
}

func MakeNode(f *Fn, children []*Node, val interface{}, bind int) *Node {
    return &Node{F: f, Children: children, Val: val, Bind: bind, Args: nil}
}

// Special code for expanding Slices means no need to get element type of Slice
// See fNodesExpandLists 
func Compat(f *Fn, n *Node, idx int) bool {
    if len(f.Args) <= idx {
        return false
    }
    switch n.Val.(type) {
        case int: return f.Args[idx] == "int"
        case string: return f.Args[idx] == "string"
        case bool: return f.Args[idx] == "bool"
        case *Object: return f.Args[idx] == n.Val.(*Object).Type
    }
    return false
}

// Do not re-evaluate with same arguments, even if reached by different path
// New evidence will cause suppression of hash disallow
func (n *Node) AddToResult(res *[]*Node, hashes map[uint32]bool) {
    h := n.Hash()
    _, ok := hashes[h]
    if ok {
        return
    }
    *res = append(*res, n)
    hashes[h] = true
}

func FNodesExpandLists(nodes []*Node, hashes map[uint32]bool) []*Node {
    res := make([]*Node, 0)
    for _,n := range nodes {
        if reflect.TypeOf(n.Val).Kind() != reflect.Slice {
            continue
        }
        switch n.Val.(type) {
            case []*Object: {
                for _,item := range n.Val.([]*Object) {
                    nn := MakeNode(&ExpandListFn, []*Node{n}, item, -1)
                    nn.AddToResult(&res, hashes)
                }
            }
            default: panic("unknown type")
        }
    }
    return res
}

func FNodesExpandProps(nodes []*Node, hashes map[uint32]bool) []*Node {
    res := make([]*Node, 0)
    for _,n := range nodes {
        switch n.Val.(type) {
            case *Object: {
                obj := n.Val.(*Object)
                for key := range obj.Props {
                    kn := MakeNode(nil, nil, key, -1)
                    rn := MakeNode(&GetPropFn, []*Node{n, kn}, obj.Props[key], -1)
                    rn.AddToResult(&res, hashes)
                }
            }
        }
    }
    return res
}

// Only returns new nodes
func FNodes(f *Fn, nodes []*Node, hashes map[uint32]bool) []*Node {
    if (f == &ExpandPropsFn) {
        return FNodesExpandProps(nodes, hashes)
    } else if (f == &ExpandListFn) {
        return FNodesExpandLists(nodes, hashes)
    }
    cargs := make([]([]int), 0)
    for i := 0; i < len(f.Args); i++ {
        cargs = append(cargs, make([]int, 0))
        for j, n := range nodes {
            if Compat(f, n, i) {
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
    res := make([]*Node, 0)
    for i := 0; i < num; i++ {
        mod := 1
        children := make([]*Node, len(cargs))
        args := make([]interface{}, len(cargs))
        for j := 0; j < len(cargs); j++ {
            idx := cargs[j][(i/mod)%len(cargs[j])]
            children[j] = nodes[idx]
            args[j] = children[j].Val
            mod *= len(cargs[j])
        }
        // Special case for equals: never compare same sequence of nodes
        // This makes us learn the wrong "beats" Pred so we have to ban by hash
        if f == &EqualStrFn {
            if children[0].Hash() == children[1].Hash() {
                continue
            }
        }
        // Commuting
        if f.Commutes {
            // Canonicalize arguments by sorting in order of hashes
            canon := make([]*CanonArg, len(cargs))
            for j := 0; j < len(cargs); j++ {
                canon[j] = &CanonArg{Arg: args[j], Child: children[j]}
                children[j].Hash()
            }
            sort.Slice(canon, func(i, j int) bool {
                return canon[i].Child.SavHash < canon[j].Child.SavHash
            })
            for j := 0; j < len(cargs); j++ {
                args[j] = canon[j].Arg
                children[j] = canon[j].Child
            }
        }
        // Check if this is a learned function
        r := f.F(args)
        // Nil means not suitable
        if r == nil {
            continue
        }
        nn := MakeNode(f, children, r, -1)
        nn.AddToResult(&res, hashes)
    }
    return res
}

func FAllNodes(fs []*Fn, nodes []*Node, hashes map[uint32]bool) []*Node {
    res := make([]*Node, 0)
    for _,f := range fs {
        r := FNodes(f, nodes, hashes)
        res = append(res, r...)
    }
    return res
}

func FAllNodesMany(fs []*Fn, nodes []*Node, times int, banlist []uint32) []*Node {
    hashes := make(map[uint32]bool)
    // Banned from satisfiability
    for _,ban := range banlist {
        hashes[ban] = true
    }
    for i := 0; i < times; i++ {
        res := FAllNodes(fs, nodes, hashes)
        nodes = append(nodes, res...)
    }
    return nodes
}

/*func NodeFromPath(obj *object, path []string) *node {
    n := makeNode(nil, nil, obj, -1)
    for _,s := range path {
        sn := makeNode(nil, nil, s, -1)
        val := getProp([]interface{}{n.val, sn.val})
        n = makeNode(&getPropFn, []*node{n, sn}, val, -1)
    }
    return n
}*/

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

// Find the list of bind numbers associated with node and children
func SetArgs(nodes []*Node) []*Node {
    for _,n := range nodes {
        n.Args = make([]int, 0)
        if n.Bind != -1 {
            n.Args = append(n.Args, n.Bind)
        } else {
            SetArgs(n.Children)
            for _,m := range n.Children {
                n.Args = append(n.Args, m.Args...)
            }
        }
        n.Args = Unique(n.Args)
    }
    return nodes
}

func GetBoolNodes(nodes []*Node) []*Node {
    res := make([]*Node, 0)
    for _,n := range nodes {
        if reflect.TypeOf(n.Val).Kind() == reflect.Bool {
            res = append(res, n)
        }
    }
    return res
}

func GetRequiredNodes(nodes []*Node, reqArgs []int) []*Node {
    nodes = SetArgs(nodes)
    res := make([]*Node, 0)
    for _,n := range nodes {
        for _,j := range n.Args {
            if IndexOf(reqArgs, j) != -1 {
                res = append(res, n)
                break
            }
        }
    }
    return res
}

func (n *Node) BindArgs(args []interface{}) {
    if n.Bind != -1 {
        n.Val = args[n.Bind]
    }
    for _,m := range n.Children {
        m.BindArgs(args)
    }
}

func (n *Node) Eval() interface{} {
    if len(n.Children) > 0 {
        args := make([]interface{}, 0)
        for _,c := range n.Children {
            args = append(args, c.Eval())
        }
        n.Val = n.F.F(args)
        return n.Val
    } else {
        return n.Val
    }
}

/*func (n *Node) equals(other *node) bool {
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
}*/

// For serialization
func (n *Node) ToRec() *NodeRec {
    r := &NodeRec{}
    r.Children = make([]*NodeRec, 0)
    for _,m := range n.Children {
        r.Children = append(r.Children, m.ToRec())
    }
    if n.F != nil {
        r.Name = GetFunctionName(n.F.F)
    }
    r.Bind = n.Bind
    if len(n.Children) == 0 && n.Bind == -1 {
        r.Val = n.Val.(string)
    }
    return r
}

// For use in Pred.Eval()
func (r *NodeRec) ToNode() *Node {
    n := &Node{}
    n.Children = make([]*Node, 0)
    for _,m := range r.Children {
        n.Children = append(n.Children, m.ToNode())
    }
    if r.Name != "" {
        n.F = GetFnFromName(r.Name)
    }
    n.Bind = r.Bind
    if r.Val != "" {
        n.Val = r.Val
    }
    return n
}

func (r *NodeRec) ToJson() []byte {
    res,_ := json.Marshal(*r)
    return res
}

func NodeFromJson(jsn []byte) *Node {
    var r NodeRec
    json.Unmarshal(jsn, &r)
    return r.ToNode()
}
