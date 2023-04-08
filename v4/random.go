package main

import (
    "fmt"
    "time"
) 

var running bool

func MakeRandomPlay() {
    if game == nil {
        return
    }
    game.mutex.Lock()
    if game.CheckWinner() != -1 {
        game.mutex.Unlock()
        return
    }
    // Make non-pass, non-pickup play if possible
    actions := game.PlayerActions(game.Players[1]) 
    found := false
    for _,a := range Shuffle(actions) {
        if a.Verb == "Attack" || a.Verb == "Defend" || a.Verb == "Reverse" {
            found = true
            _,err := game.TakeAction(a)
            if err != nil {
                fmt.Println(err)
            } else {
                fmt.Println("Broke")
                break
            }
        }
    }
    if !found {
        for _,a := range actions {
            if a.Verb == "Pass" || a.Verb == "Pickup" {
                _,err := game.TakeAction(a)
                if err != nil {
                    fmt.Println(err)
                }
            }
        }
    }
    game.mutex.Unlock()
}

func StartRandom() {
    if running {
        return
    }
    running = true
    go RandomLoop()
}

func RandomLoop() {
    for {
        if !running {
            break
        }
        time.Sleep(time.Second)
        MakeRandomPlay()  
    }
}

func StopRandom() {
    running = false
}
