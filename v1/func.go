package main

import (
    "reflect"
)

func GetProp(args []interface{}) interface{} {
    return args[0].(*Object).Props[args[1].(string)];
}

// Workaround for HasProp
func Exists(args []interface{}) interface{} {
    return true
}

func GreaterRank(args []interface{}) interface{} {
    i1 := IndexOf(ranks, args[0].(string))
    i2 := IndexOf(ranks, args[1].(string))
    if i1 == -1 || i2 == -1 {
        return nil
    }
    return i1 > i2
}

func Equal(args []interface{}) interface{} {
    return reflect.DeepEqual(args[0], args[1])
}

// Never called
// Still used for its name in hashing and printouts
// Special action in fNodes
// Also special action in pred.eval()
// This is due to multiple return values
func ExpandList(args []interface{}) interface{} {
    return nil
}

var GetPropFn = Fn{F: GetProp, Args: []string{"*Object", "string"}}
var GreaterRankFn = Fn{F: GreaterRank, Args: []string{"string", "string"}}
var CardExistsFn = Fn{F: Exists, Args: []string{"card"}}
var EqualStrFn = Fn{F: Equal, Args: []string{"string", "string"}, Commutes: true}
var ExpandPropsFn = Fn{F: GetProp, Args: []string{"*Object"}}
var ExpandListFn = Fn{F: ExpandList, Args: []string{"[]any"}}

// For serialization
func GetFnFromName(name string) *Fn {
    list := []*Fn{&GetPropFn, &GreaterRankFn, &CardExistsFn, &EqualStrFn, &ExpandPropsFn, &ExpandListFn}
    for _,f := range list {
        if GetFunctionName(f.F) == name {
            return f
        }
    }
    return nil
}
