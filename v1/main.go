package main

import (
    "fmt"
)

func main() {
    c := makeCard("8", "Spades")
    d := makeCard("10", "Hearts")
    e := makeCard("Jack", "Hearts")
    f := makeCard("Jack", "Spades")
    g := makeGame(d)
    h := makeGame(c)
    p := makePred("beats", []string{"game", "card", "card"}, []*fn{&expandListFn, &expandPropsFn, &equalStrFn, &greaterRankFn})
    ex1 := makeExample(false, []interface{}{g,c,d})
    ex2 := makeExample(true, []interface{}{g,d,c})
    ex3 := makeExample(false, []interface{}{g,d,e})
    ex4 := makeExample(true, []interface{}{g,d,f})
    ex5 := makeExample(false, []interface{}{h,d,e})
    ex6 := makeExample(true, []interface{}{h,e,d})
    exs := []*example{ex1, ex2, ex3, ex4, ex5, ex6}
    /*for _,ex := range []*example{ex1, ex2, ex3, ex4} {
        nodes := expandNodes(p.fns, ex, 6, []int{0,1,2})
        for _,n := range nodes {
            if n == nil {
                fmt.Println("nil")
                continue
            }
            fmt.Println(nodeStr(n, 0))
            fmt.Println(HashFuncBind(n))
        }
        fmt.Println(len(nodes))
        fmt.Println("-----")
    }*/
    table := makeTable(p.fns, exs, 6, []int{1,2})
    for _,nodes := range table {
        for i,n := range nodes {
            if n == nil {
                fmt.Println("nil")
                continue
            }
            fmt.Print(i," ")
            //fmt.Println(nodeStr(n, 0))
            fmt.Println(HashFuncBind(n))
        }
        fmt.Println(len(nodes))
        fmt.Println("---")
    }
    ps := satisfy(3, exs, table) 
    for _,n := range ps.nodes {
        fmt.Println(nodeStr(n, 0))
    }
    fmt.Println(satStr(len(ps.nodes), ps.idx))
}
