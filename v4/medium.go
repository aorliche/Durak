package main

import (
    "fmt"
    "time" 
)

func MakeBestPlay(game *Game) {
    game.mutex.Lock()
    if game.CheckWinner() != -1 {
        game.mutex.Unlock()
        return
    }
    state := InitGameState(game)
    for i,_ := range state.Hands[0] {
        state.Hands[0][i] = 36
        //c := state.Hands[0][i] 
        //c = 36
        /*c.Rank = "?"
        c.Suit = "?"*/
    }
    game.memory.SetKnownCards(state, 1, 0)
    chain,val := state.Move(1, 0, -1)
    if len(chain) == 0 {
        game.mutex.Unlock()
        return
    }
    act := chain[len(chain)-1]
    fmt.Println(act.ToAction().String(),val)
    if act.Verb == DeferV {
        game.mutex.Unlock()
        return
    }
    upd := game.TakeAction(act.ToAction())
    if upd != nil {
        game.Recording = append(game.Recording, &Record{Action: act.ToAction()})
        game.Recording = append(game.Recording, &Record{Update: upd})
    } else {
        fmt.Println(act)
        panic("Error in Best AI B")
    }
    game.mutex.Unlock()
}

func MediumLoop(game *Game) {
    for {
        if game.CheckWinner() != -1 {
            break
        }
        time.Sleep(time.Second)
        MakeBestPlay(game)  
    }
}
