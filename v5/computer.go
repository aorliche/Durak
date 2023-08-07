package main

import (
    "log"
    "math/rand"
    "time"
)

func (game *Game) MakeEasyPlay() Action {
    var act Action
    game.mutex.Lock()
    if game.CheckWinner() != -1 {
        game.mutex.Unlock()
        return act
    }
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
        log.Println(act.ToStr())
        game.TakeAction(act)
        SendInfo(0, game)
    }
    game.mutex.Unlock()
    return act
}

func (game *Game) MakeMediumPlay() Action {
    var act Action
    game.mutex.Lock()
    if game.CheckWinner() != -1 {
        game.mutex.Unlock()
        return act
    }
    var state *GameState
    if len(game.Deck) > 1 {
        state = game.MaskUnknownCards(1)
    } else {
        state = game.State
    }
    c, r := state.EvalNode(state, 1, 0, 0, len(game.Deck))
    if len(c) > 0 {
        act = c[len(c)-1]
        if act.Verb != DeferVerb {
            log.Println(r, act.ToStr())
        }
        game.TakeAction(act)
    }
    game.mutex.Unlock()
    return act
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
        var act Action
        if comp == "Easy" {
            act = game.MakeEasyPlay()
        } else {
            act = game.MakeMediumPlay()
        }
        game.CheckWinner()
        if !act.IsNull() && act.Verb != DeferVerb {
            SendInfo(0, game)
        }
        if game.Recording.Winner != -1 {
            game.WriteGame()
            break
        }
    }
}
