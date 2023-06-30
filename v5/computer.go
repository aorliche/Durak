package main

import (
    "fmt"
    "math/rand"
    "time"
)

func (game *Game) MakeEasyPlay() {
    game.mutex.Lock()
    if game.Recording.Winner != -1 {
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

// comp is "Easy" or "Medium"
func (game *Game) StartComputer(comp string) {
    go game.RandomLoop(comp)
}

func (game *Game) RandomLoop(comp string) {
    for {
        if game.Recording.Winner != -1 {
            break
        }
        time.Sleep(100 * time.Millisecond)
        if comp == "Easy" {
            game.MakeEasyPlay()
        } /*else {
            game.MakeMediumPlay()
        }*/
    }
}
