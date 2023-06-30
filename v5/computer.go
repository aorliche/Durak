package main

import (
    "fmt"
    "math/rand"
    "time"
)

func (game *Game) MakeEasyPlay() {
    game.mutex.Lock()
    if game.CheckWinner() != -1 {
        game.mutex.Unlock()
        return
    }
    var act Action
    acts := game.State.PlayerActions(1)
    rand.Shuffle(len(acts), func(i, j int) {
        acts[i], acts[j] = acts[j], acts[i]
    })
    found := false
    for _, a := range acts {
        switch a.Verb {
            case PlayVerb, CoverVerb, ReverseVerb: {
                found = true
                act = a
                break
            }
        }
    }
    if !found && len(acts) > 0 {
        found = true
        act = acts[0]
    }
    if found && act.Verb != DeferVerb {
        fmt.Println(act.ToStr())
        game.TakeAction(act)
    }
    game.mutex.Unlock()
}

func (game *Game) MakeMediumPlay() {
    game.mutex.Lock()
    if game.CheckWinner() != -1 {
        game.mutex.Unlock()
        return
    }
    var state *GameState
    if len(game.Deck) > 1 {
        state = game.MaskUnknownCards(1)
    } else {
        state = game.State
    }
    c, r := state.EvalNode(state, 1, 0, 0, len(game.Deck) == 0)
    if len(c) > 0 {
        act := c[len(c)-1]
        fmt.Println(r, act.ToStr())
        game.TakeAction(act)
    }
    game.mutex.Unlock()
}

// comp is "Easy" or "Medium"
func (game *Game) StartComputer(comp string) {
    go game.RandomLoop(comp)
}

func (game *Game) RandomLoop(comp string) {
    for {
        if game.CheckWinner() != -1 {
            break
        }
        time.Sleep(100 * time.Millisecond)
        if comp == "Easy" {
            game.MakeEasyPlay()
        } else {
            game.MakeMediumPlay()
        }
    }
}
