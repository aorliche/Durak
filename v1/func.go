package main

import (
    "reflect"
)

func getProp(args []interface{}) interface{} {
    return args[0].(*object).props[args[1].(string)];
}

func greaterRank(args []interface{}) interface{} {
    i1 := IndexOf(ranks, args[0].(string)) 
    i2 := IndexOf(ranks, args[1].(string))
    if i1 == -1 || i2 == -1 {
        return nil
    }
    return i1 > i2
}

func equal(args []interface{}) interface{} {
    return reflect.DeepEqual(args[0], args[1])
}

// Never called
// Special action in fNodes
func expandList(args []interface{}) interface{} {
    return nil
}

var getPropFn = fn{f: getProp, args: []string{"*object", "string"}}
var greaterRankFn = fn{f: greaterRank, args: []string{"string", "string"}}
var equalStrFn = fn{f: equal, args: []string{"string", "string"}, commutes: true}
var expandPropsFn = fn{f: getProp, args: []string{"*object"}}
var expandListFn = fn{f: expandList, args: []string{"[]any"}}
