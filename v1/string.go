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
