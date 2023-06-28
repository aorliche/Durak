package main

import (
    //"fmt"
    "testing"
)

func TestBeats(t *testing.T) {
    if Card(10).Beats(Card(11), Card(20)) {
        t.Errorf("Card(10).Beats(Card(11), Card(20))")
    }
    if !Card(4).Beats(Card(17), Card(1)) {
        t.Errorf("Card(4).Beats(Card(17), Card(1))")
    }
}

func TestInitGame(t *testing.T) {
    game := InitGame(0, "Computer")
    if game == nil {
        t.Errorf("InitGame failed")
    }
}

func TestGetActions(t *testing.T) {
    game := InitGame(0, "Computer")
    acts0 := game.State.PlayerActions(0)
    acts1 := game.State.PlayerActions(1)
    if len(acts0) == 0 {
        t.Errorf("No actions for player 0")
    }
    if len(acts1) != 0 {
        t.Errorf("Actions for player 1")
    }
    /*for _, act := range acts0 {
        fmt.Println(act.ToStr())
    }*/
}

func TestTakeAction(t *testing.T) {
    game := InitGame(0, "Computer")
    game.TakeAction(game.State.RandomAction())
    acts1 := game.State.PlayerActions(1)
    if len(acts1) == 0 {
        t.Errorf("No actions for player 1")
    }
}

func TestSearchStart(t *testing.T) {
    game := InitGame(0, "Computer")
    c, _ := game.State.EvalNode(nil, 0, 0, 0, false)
    if len(c) == 0 {
        t.Errorf("No action chain for search")
    }
}

func TestSearchEnd(t *testing.T) {
    game := InitGame(0, "Computer")
    game.Deck = make([]Card, 0)
    c, _ := game.State.EvalNode(nil, 0, 0, 0, true)
    if len(c) == 0 {
        t.Errorf("No action chain for search")
    }
}
