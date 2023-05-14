package main

import (
    "fmt"
    "testing"
)

/*func TestSearch(t *testing.T) {
    game := InitGame(0, false)
    state := InitGameState(game)
    act,val := state.Move(0, 0)
    fmt.Println(state.Trump)
    fmt.Println(ChainString(act),val)
    fmt.Println(HandString(state.Hands[0]))
    fmt.Println(HandString(state.Hands[1]))
}*/

/*func TestMystery(t *testing.T) {
    game := InitGame(0, false)
    state := InitGameState(game)
    for i,_ := range state.Hands[1] {
        c := state.Hands[1][i] 
        c.Rank = "?"
        c.Suit = "?"
    }
    act,val := state.Move(0, 0)
    fmt.Println(state.Trump)
    fmt.Println(ChainString(act),val)
    fmt.Println(HandString(state.Hands[0]))
}*/

func TestPartialMystery(t *testing.T) {
    game := InitGame(0, false)
    state := InitGameState(game)
    for i,_ := range state.Hands[1][:4] {
        c := state.Hands[1][i] 
        c.Rank = "?"
        c.Suit = "?"
    }
    act,val := state.Move(0, 0)
    fmt.Println(state.Trump)
    fmt.Println(ChainString(act),val)
    fmt.Println(HandString(state.Hands[0]))
    fmt.Println(HandString(state.Hands[1]))
}
