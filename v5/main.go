package main

import (
    "encoding/json"
    "fmt"
    "math/rand"
)

func main() {
    game := InitGame(0, "Computer")
    game.Deck = make([]Card, 0)
    for i := 0; i < 20; i++ {
        fmt.Println(game.State.Trump)
        c0, r0 := game.State.EvalNode(nil, 0, 0, 0, true)
        c1, r1 := game.State.EvalNode(nil, 1, 0, 0, true)
        var act Action
        var r int
        fmt.Println(len(c0), len(c1))
        if len(c0) == 0 {
           act = c1[len(c1)-1]
           r = r1
        } else if len(c1) == 0 {
           act = c0[len(c0)-1]
           r = r0
        } else {
            if rand.Intn(2) == 0 {
                act = c0[len(c0)-1]
                r = r0
            } else {
                act = c1[len(c1)-1]
                r = r1
            }
        }
        fmt.Println(r, act.ToStr())
        game.TakeAction(act)
        /*act := game.State.RandomAction()
        fmt.Println(act.ToStr())
        game.TakeAction(act)*/
    }
    jsn, _ := json.Marshal(game)
    fmt.Println(string(jsn))
    /*c, _ := game.State.EvalNode(nil, 0, 0, 0, false)
    for _, act := range c {
        fmt.Println(act.ToStr())
    }*/
}
