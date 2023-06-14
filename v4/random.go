package main

import (
    "encoding/json"
    "fmt"
    "time"
) 

func MakeRandomPlay(game *Game) {
    game.mutex.Lock()
    if game.CheckWinner() != -1 {
        game.mutex.Unlock()
        return
    }
    // Make non-pass, non-pickup play if possible
    actions := game.PlayerActions(game.Players[1]) 
    found := false
    var act *Action
    var upd *Update
    for _,a := range Shuffle(actions) {
        if a.Verb == "Attack" || a.Verb == "Defend" || a.Verb == "Reverse" {
            found = true
            upd = game.TakeAction(a)
            if upd != nil {
                act = a
                break
            } 
        }
    }
    if !found {
        for _,a := range actions {
            if a.Verb == "Pass" || a.Verb == "Pickup" {
                found = true
                upd = game.TakeAction(a)
                act = a
            }
        }
    }
    if upd != nil {
        actJsn,_ := json.Marshal(act)
        updJsn,_ := json.Marshal(upd)
        game.Recording = append(game.Recording, string(actJsn))
        game.Recording = append(game.Recording, string(updJsn))
    } else if found {
        fmt.Println("Error in Random AI B")
    }
    game.mutex.Unlock()
}

func RandomLoop(game *Game) {
    for {
        if game.CheckWinner() != -1 {
            break
        }
        time.Sleep(time.Second)
        MakeRandomPlay(game)  
    }
}
