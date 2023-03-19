package main

import (
    "fmt"
    "reflect"
)

func LearnBeats(banlist []uint32) *Pred {
    c0 := MakeCard("8", "Spades")
    c1 := MakeCard("10", "Hearts")
    c2 := MakeCard("Queen", "Diamonds")
    c3 := MakeCard("Jack", "Hearts")
    c4 := MakeCard("Jack", "Spades")
    g0 := MakeGame(c0)
    fns := []*Fn{&ExpandListFn, &ExpandPropsFn, &EqualStrFn, &GreaterRankFn}
    p := MakePred("beats", []string{"game", "card", "card"})
    ex1 := MakeExample(true, []interface{}{g0,c0,c1})
    ex2 := MakeExample(false, []interface{}{g0,c1,c0})
    ex3 := MakeExample(false, []interface{}{g0,c1,c2})
    ex4 := MakeExample(false, []interface{}{g0,c2,c1})
    ex5 := MakeExample(true, []interface{}{g0,c0,c3})
    ex6 := MakeExample(false, []interface{}{g0,c1,c3})
    ex7 := MakeExample(true, []interface{}{g0,c3,c1})
    ex8 := MakeExample(true, []interface{}{g0,c0,c1})
    ex9 := MakeExample(false, []interface{}{g0,c0,c4})
    ex10 := MakeExample(true, []interface{}{g0,c4,c2})
    exs := []*Example{ex1, ex2, ex3, ex4, ex5, ex6, ex7, ex8, ex9, ex10}
    table := MakeTable(fns, exs, 6, []int{1,2}, banlist)
    p.Nodes, p.Idx = Satisfy(3, exs, table)
    return p
}

func BanlistWrite() {
    banlist := make([]uint32, 0)
    for true {
        p := LearnBeats(banlist)
        fmt.Println(p.ToStr())
        bad := false
        for _,n := range p.Nodes {
            if n.F == &GreaterRankFn && reflect.DeepEqual(n.Args, []int{2,1})  {
                bad = true
                banlist = append(banlist, n.SavHash)
            }
        }
        if !bad {
            p.WriteFile("learned/beats.pred")
            return
        }
    }
}

func TestReadBeats() {
    c0 := MakeCard("Jack", "Spades")
    c1 := MakeCard("10", "Hearts")
    c2 := MakeCard("Queen", "Spades")
    g := MakeGame(c0)
    ex1 := MakeExample(true, []interface{}{g,c0,c1})
    ex2 := MakeExample(true, []interface{}{g,c0,c2})
    p := ReadPred("learned/beats.pred")
    p.Bind(ex1.Args)
    fmt.Println(p.Eval())
    p.Bind(ex2.Args)
    fmt.Println(p.Eval())
}

/*func ExistsNodes() {
    c0 := MakeCard("8", "Spades")
    c1 := MakeCard("10", "Hearts")
    c2 := MakeCard("Queen", "Diamonds")
    g0 := MakeGame(c0)
    fns := []*Fn{&ExpandListFn, &ExpandPropsFn, &EqualStrFn, &GreaterRankFn}
    ex1 := MakeExample(true, []interface{}{g0,c0,c1})
    ex2 := MakeExample(false, []interface{}{g0,c1,c0})
    ex3 := MakeExample(false, []interface{}{g0,c1,c2})
    ex4 := MakeExample(false, []interface{}{g0,c2,c1})
    ex5 := MakeExample(true, []interface{}{g0,c0,c3})
    ex6 := MakeExample(false, []interface{}{g0,c1,c3})
    ex7 := MakeExample(true, []interface{}{g0,c3,c1})
    ex8 := MakeExample(true, []interface{}{g0,c0,c1})
    ex9 := MakeExample(false, []interface{}{g0,c0,c4})
    ex10 := MakeExample(true, []interface{}{g0,c4,c2})
    exs := []*Example{ex1, ex2, ex3, ex4, ex5, ex6, ex7, ex8, ex9, ex10}
    table := MakeTable(fns, exs, 6, []int{1,2}, banlist)
    p.Nodes, p.Idx = Satisfy(3, exs, table)
    return p

}*/

/*func LearnCover() {
    c0 := MakeCard("Jack", "Spades")
    c1 := MakeCard("10", "Hearts")
    c2 := MakeCard("Queen", "Spades")
    g := MakeGame(c2)
    c2.SetProp("covers", "card", c0)
    c1.SetProp("covers", "card", c1)
    ex1 := MakeExample(true, []interface{}{g,c0,c1})
    ex2 := MakeExample(true, []interface{}{g,c0,c2})
}*/

func main() {
    TestReadBeats()
    /*p := learnBeats()
    p.WriteFile("preds/beats.pred")
    //p = predFromJson(p.toRec().toJson())
    p.bind(ex1.args)
    p.eval()
    fmt.Println(p.toStr())*/
    /*for _,n := range p.nodes {
        n = nodeFromJson(n.toRec().toJson())
        n.bindArgs(ex1.args)
        fmt.Println(n.eval())
        fmt.Println(n.toStr())
        //fmt.Println(nodeFromJson(n.toRec().toJson()).toStr())
    }*/
}
