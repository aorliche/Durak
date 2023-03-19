package main

import (
    "fmt"
    "strings"
)

func MakeBoolObject(b bool) *Object {
    obj := Object{Type: "bool", Val: b}
    return &obj
}

func MakeStringObject(s string) *Object {
    obj := Object{Type: "string", Val: s}
    return &obj
}

func MakeSliceObject(s []*Object) *Object {
    obj := Object{Type: "slice", Val: s}
    return &obj
}

func MakeProdObject(a []*Object, b []*Object) *Object {
    obj := Object{Type: "prod", Val: [2][]*Object{a,b}}
    return &obj
}

/*func MakeObject(typ string) *Object {
    obj := Object{Type: typ, Slice: false, Val: nil, Props: make(map[string]*Object)}
    return &obj
}*/

func (obj *Object) SetProp(key string, val *Object) {
    obj.Props[key] = val
}

func (obj *Object) ToStr() string {
    switch obj.Type {
        case "slice": return SliceStr(obj)
        case "Card": return CardStr(obj)
        case "Player": return PlayerStr(obj)
        case "Game": return GameStr(obj)
        default: return fmt.Sprint(obj.Val)
    }
}

func SliceStr(obj *Object) string {
    sl := obj.Val.([]*Object)
    if len(sl) > 0 {
        typ := sl[0].Type
        if typ == "int" || typ == "string" || typ == "bool" || typ == "Card" {
            strs := make([]string, 0)
            for _,a := range sl {
                strs = append(strs, a.ToStr())
            }
            return fmt.Sprintf("[%s]", strings.Join(strs, ", "))
        } else {
            return fmt.Sprintf("[Slice Len: %d Type: %s]", len(sl), typ)
        }
    } else {
        return "[]"
    }
}

