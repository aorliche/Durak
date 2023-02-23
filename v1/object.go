package main

import (
    //"reflect"
)

func MakeObject(typ string) *Object {
    obj := Object{Type: typ, Props: make(map[string]interface{}), PropTypes: make(map[string]string)}
    return &obj
}

func (obj *Object) SetProp(prop string, typ string, item interface{}) {
    obj.Props[prop] = item
    obj.PropTypes[prop] = typ
}

func (obj *Object) ToStr() string {
    switch obj.Type {
        case "card": return CardStr(obj)
    }
    return obj.Type
}

// TODO difference between states
// TODO Closure returns path to all different nodes
// TODO slices as property types
/*func diffObjects(obj1 *object, obj2 *object, path []string) []([]string) {
    keys1 := GetKeys(obj1.props)
    keys2 := GetKeys(obj2.props)
    diff := StrArrMinus(keys2, keys1)
    common := StrArrInt(keys1, keys2)
    res := make([]([]string), 0)
    for _,key := range diff {
        npath := append(path, key)
        res = append(res, npath)
    }
    tObj := reflect.TypeOf(obj1)
    for _,key := range common {
        t1 := reflect.TypeOf(obj1.props[key])
        t2 := reflect.TypeOf(obj2.props[key])
        if t1 == t2 && t1 == tObj {
            npath := append(path, key)
            subres := diffObjects(obj1.props[key].(*object), obj2.props[key].(*object), npath)
            res = append(res, subres...)
        } else if t1 == t2 && obj1.props[key] != obj2.props[key] {
            npath := append(path, key)
            res = append(res, npath)
        }
    }
    return res
}*/
