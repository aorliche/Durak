package main

import (
    "log"
    "math/rand"
    "time"
)

func (game *Game) MakeEasyPlay(player int) Action {
    game.mutex.Lock()
    var act Action
    if game.CheckGameOver() {
        return act
    }
    acts := game.State.PlayerActions(player)
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
    }
    game.mutex.Unlock()
    return act
}

func (game *Game) MakeMediumPlay(player int) Action {
    var act Action
    game.mutex.Lock()
    if game.CheckGameOver() {
        game.mutex.Unlock()
        return act
    }
    var state *GameState
    if len(game.Deck) > 1 || (len(game.State.Hands) - len(game.Recording.Winners)) > 2 {
        state = game.MaskUnknownCards(player)
    } else {
        state = game.State.Clone()
    }
    game.mutex.Unlock()
    c, r := state.EvalNode(state, player, 0, 0, nil)
    if len(c) > 0 {
        act = c[len(c)-1]
        if act.Verb != DeferVerb {
            log.Println(r, act.ToStr())
        }
        game.mutex.Lock()
        game.TakeAction(act)
        game.mutex.Unlock()
    }
    return act
}

// comp is "Easy" or "Medium"
func (game *Game) StartComputer(comp string, player int) {
    go func() {
        for {
            if game.CheckGameOver() {
                break
            }
            time.Sleep(100 * time.Millisecond)
            var act Action
            if comp == "Easy" {
                act = game.MakeEasyPlay(player)
            } else if comp == "Medium" {
                act = game.MakeMediumPlay(player)
            }
            game.CheckGameOver()
            // Send info to human players
            if !act.IsNull() && act.Verb != DeferVerb {
                game.SendInfoHumans()
            }
            if game.CheckGameOver() {
                game.WriteGame()
                break
            }
        }
    }()
}
