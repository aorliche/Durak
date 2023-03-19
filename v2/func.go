package main

import (
    //"fmt"
    //"reflect"
)

func GetProp(args []*Object) *Object {
    return args[0].Props[args[1].Val.(string)]
}

func Empty(args []*Object) *Object {
    return MakeBoolObject(len(args[0].Val.([]*Object)) == 0)
}

// Cartesian product
// Only 2D for now
func Prod(args []*Object) *Object {
    a := args[0].Val.([]*Object)
    b := args[1].Val.([]*Object)
    return MakeProdObject(a, b)
}

// Pairs, Prods, and Slices
func First(args []*Object) *Object {
    return args[0].Val.([]*Object)[0]
}

// Pairs, Prods, and Slices
func Second(args []*Object) *Object {
    return args[0].Val.([]*Object)[1]
}

func GreaterRank(args []*Object) *Object {
    i1 := IndexOf(ranks, args[0].Val.(string))
    i2 := IndexOf(ranks, args[1].Val.(string))
    if i1 == -1 || i2 == -1 {
        return nil
    }
    return MakeBoolObject(i1 > i2)
}

func EqualStr(args []*Object) *Object {
    return MakeBoolObject(args[0].Val.(string) == args[1].Val.(string))
}

var GetPropFn = Fn{A: GetProp, Args: []string{"Object", "string"}}
var ExpandPropsFn = Fn{A: GetProp, Args: []string{"Object"}}
var EmptyFn = Fn{A: Empty, Args: []string{"slice"}}
var ProdFn = Fn{A: Prod, Args: []string{"slice", "slice"}}
var FirstFn = Fn{A: First, Args: []string{"slice|pair"}}
var SecondFn = Fn{A: Second, Args: []string{"slice|pair"}}
var GreaterRankFn = Fn{A: GreaterRank, Args: []string{"string", "string"}}
var EqualStrFn = Fn{A: EqualStr, Args: []string{"string", "string"}, Commutes: true}

func GetName(f *Fn) string {
    if f.A != nil {
        return GetFunctionName(f.A)
    } 
    if f.N != nil {
        return GetName(f.N.F)
    }
    return GetName(f.F.Pred)
}

// For serialization
/*func GetFnFromName(name string) *Fn {
    list := []*Fn{&GetPropFn, &EmptySliceFn, &, &EqualStrFn, &ExpandPropsFn, &ExpandListFn}
    for _,f := range list {
        if GetFunctionName(f.F) == name {
            return f
        }
    }
    return nil
}*/
