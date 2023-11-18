package durak

import (
    "log"
    "math/rand"
    "time"
)

func (game *Game) MakeEasyPlay(player int) Action {
    game.Mutex.Lock()
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
    game.Mutex.Unlock()
    return act
}

func (game *Game) MakeMediumPlay(player int, params *EvalParams) Action {
    var act Action
    game.Mutex.Lock()
    if game.CheckGameOver() {
        game.Mutex.Unlock()
        return act
    }
    var state *GameState
    if len(game.Deck) > 1 || (len(game.State.Hands) - len(game.Recording.Winners)) > 2 {
        state = game.MaskUnknownCards(player)
    } else {
        state = game.State.Clone()
    }
    game.Mutex.Unlock()
    c, _ := state.EvalNode(state, player, 0, 0, params)
    if len(c) > 0 {
        act = c[len(c)-1]
        /*if act.Verb != DeferVerb {
            log.Println(r, act.ToStr())
        }*/
        game.Mutex.Lock()
        game.TakeAction(act)
        game.Mutex.Unlock()
    }
    return act
}

// comp is "Easy" or "Medium"
func (game *Game) StartComputer(comp string, player int, params *EvalParams, actionCb func (*Game), gameOverCb func (*Game)) *bool {
    if params == nil {
        params = &DefaultEvalParams
    }
    kill := false
    go func() {
        for {
            if kill {
                break
            }
            if game.CheckGameOver() {
                break
            }
            time.Sleep(100 * time.Millisecond)
            var act Action
            if comp == "Easy" {
                act = game.MakeEasyPlay(player)
            } else if comp == "Medium" {
                act = game.MakeMediumPlay(player, params)
            }
            game.CheckGameOver()
            // Send info to human players
            if !act.IsNull() && act.Verb != DeferVerb && actionCb != nil {
                //game.SendInfoHumans()
                actionCb(game)
            }
            if game.CheckGameOver() && gameOverCb != nil {
                //game.WriteGame()
                gameOverCb(game)
                break
            }
        }
    }()
    return &kill
}
