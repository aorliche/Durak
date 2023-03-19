package main

import (
    "fmt"
)

var Card = MakeCardObject
var Sl = MakeSliceObject
var Pl = MakePlayerObject
var Str = MakeStringObject

func TestOne() {
    hand := Sl([]*Object{Card("8", "Hearts"), Card("9", "Spades")})
    player := Pl("Bob", hand)
    fmt.Println(player.ToStr())
    rs := make([]*Object, 0)
    for _,r := range ranks {
        rs = append(rs, Str(r))
    }
    fmt.Println(Sl(rs).ToStr())
    n := MakeNode(&GetPropFn, nil, MakeBoolObject(true), -1)
    fmt.Println(n.Compat("bool"))
    fmt.Println(n.Compat("bool|int"))
    fmt.Println(n.Compat("Object|int"))
    arg := Str("hello")
    n0 := MakeNode(&GetPropFn, nil, Str("hello"), -1)
    n1 := MakeNode(nil, nil, Str("hello"), -1)
    n2 := MakeNode(nil, nil, Str("goodbye"), 0)
    n3 := MakeNode(nil, nil, player, -1)
    nodes := [4]*Node{n0,n1,n2,n3}
    fns := [2]*Fn{&ExpandPropsFn, &EqualStrFn}
    hashes := make(map[uint32]bool)
    res := ApplyFnsNodes(fns[:], nodes[:], hashes)
    for _,n := range res {
        fmt.Println(n.ToStr())
        n.BindArgs(arg)
        n.Eval()
        fmt.Println(n.ToStr())
    }
}

func main() {
    c0 := Card("8", "Hearts")
    c1 := Card("9", "Spades")
    c2 := Card("10", "Diamonds")
    c3 := Card("King", "Clubs")
    c4 := Card("6", "Hearts")
    hand := []*Object{c0,c1,c2,c3,c4}
    n0 := Card("Queen", "Clubs")
    //n1 := Card("11", "Diamonds")
    //n2 := Card("Ace", "Spades")
    //n3 := Card("7", "Hearts")
    //board := []*Object{n0,n1,n2,n3}
    g0 := MakeGameObject(Card("9", "Diamonds"))
    fmt.Println(g0.Props["trump"].ToStr())
    //g1 := MakeGameObject(Card("Jack", "Clubs"))
    //g2 := MakeGameObject(Card("King", "Hearts"))
    // board-n and game-n
    out0 := []*Object{c2,c3}
    //out1 := []*Object{c3}
    //out2 := []*Object{c0,c4}
    res0 := Solve(hand, n0, g0, out0)
    for _,n := range res0 {
        fmt.Println(n.ToStr())
    }
}
