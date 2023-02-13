package main

import (
    "fmt"
    "strings"
)

func objStr(obj *object) string {
    switch obj.typ {
        case "card": return cardStr(obj)
    }
    return obj.typ
}

func nodeStr(n *node, lvl int) string {
    str := strings.Repeat("  ", lvl)
    if n.f != nil {
        str += GetFunctionName(n.f.f)
        str += " (" + strings.Join(n.f.args, ",") + ") "
    }
    if n.args != nil {
        str += fmt.Sprint(n.args) + " "
    }
    switch n.val.(type) {
        case *object: str += objStr(n.val.(*object))
        case []*object: {
            strSlice := make([]string, 0)
            for _,obj := range n.val.([]*object) {
                strSlice = append(strSlice, objStr(obj))
            }
            str += "[" + strings.Join(strSlice, ", ") + "]"
        }
        default: str += fmt.Sprintf("%v", n.val)
    }
    if n.bind > -1 {
        str += fmt.Sprintf(" BIND %d", n.bind)
    }
    for i := 0; i < len(n.children); i++ {
        str += "\n" + nodeStr(n.children[i], lvl+1)
    }
    return str
}

func indexStr(idx *index) string {
    s := fmt.Sprintf("%d %d %d", idx.nTerms, idx.nNodes, idx.curIdx)
    s += fmt.Sprint(idx.idcs)
    return s
}

func satStr(n int, idx int) string {
    switch n {
        case 1: return satStr1(idx)
        case 2: return satStr2(idx)
        default: panic("bad")
    }
}

func satStr1(idx int) string {
    n0 := Ternary((idx & 1) == 1, "~", "")
    return n0 + "A"
}

func satStr2(idx int) string {
    n0 := Ternary((idx & 1) == 1, "~", "")
    n1 := Ternary(((idx >> 1) & 1) == 1, "~", "")
    split := idx >> 2
    if split == 0 {
        return n0 + "A" + n1 + "B"
    } else {
        return n0 + "A+" + n1 + "B"
    }
}
