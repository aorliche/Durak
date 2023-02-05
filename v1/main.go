package main

import (
    "fmt"
)

func main() {
    c := makeCard("8", "Spades")
    d := makeCard("10", "Hearts")
    g := makeGame(d)
    nodes := makePred("beats", []string{"card", "card"}, []interface{}{c,d}, g, []*fn{&expandListFn, &expandPropsFn, &equalStrFn, &greaterRankFn}, 6) 
    for _,n := range nodes {
        fmt.Println(nodeStr(n.n, 0))
        fmt.Println(Hash(n.n))
    }
    fmt.Println(len(nodes))
}
