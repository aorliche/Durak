package main

import (
    "encoding/json"
    "fmt"
    "net/http"
    //"net/url"
    "strconv"
)

var game *Game

func BadGame(w http.ResponseWriter, req *http.Request) bool {
    if game == nil {
        fmt.Fprintf(w, "No active game\n")
        return true
    }
    return false
}

/*func GetPlayerIdx(w http.ResponseWriter, req *http.Request) int {
    idx, err := strconv.Atoi(req.URL.Query().Get("p"))
    if err != nil || idx < 0 || idx >= len(game.Players) {
        fmt.Fprintf(w, "No such player %v\n", idx)
        return -1
    }
    return idx
}*/

func Join(w http.ResponseWriter, req *http.Request) {
    req.URL.RawQuery = "p=1"
    Update(w, req)
}

func NewGame(w http.ResponseWriter, req *http.Request) {
    game = InitGame()
    Update(w, req)
}

func Update(w http.ResponseWriter, req *http.Request) {
    if BadGame(w, req) {
        return
    }
    game.mutex.Lock()
    actions := make([][]*Action, 0)
    for _,p := range game.Players {
        actions = append(actions, game.PlayerActions(p))
    }
    p := req.URL.Query().Get("p")
    mp := 1
    if p != "" {
       val, err := strconv.Atoi(p)  
        if err != nil {
            fmt.Println("Bad player query", err)
            game.mutex.Unlock()
            return
        }
       mp = 1-val
    }
    upd := GameUpdate{Board: game.Board, Deck: len(game.Deck), Trump: game.Trump, Players: game.MaskedPlayers(mp), Actions: actions, Winner: game.CheckWinner()}
    game.RecordUpdate(&upd, true)
    jsn,_ := json.Marshal(upd)
    fmt.Fprintf(w, "%s\n", jsn)
    game.mutex.Unlock()
}

/*func GetHand(w http.ResponseWriter, req *http.Request) {
    if BadGame(w, req) {
        return
    }
    idx := GetPlayerIdx(w, req)
    if idx == -1 {
        return 
    }
    jsn,_ := json.Marshal(game.Players[idx].Hand)
    fmt.Fprintf(w, "%s\n", jsn)
}*/

/*func GetActions(w http.ResponseWriter, req *http.Request) {
    game.mutex.Lock()
    if BadGame(w, req) {
        game.mutex.Unlock()
        return
    }
    idx := GetPlayerIdx(w, req)
    if idx == -1 {
        game.mutex.Unlock()
        return 
    }
    a := game.PlayerActions(game.Players[idx])
    game.RecordPossibleActions(idx, a)
    jsn,_ := json.Marshal(a)
    fmt.Fprintf(w, "%s\n", jsn)
    game.mutex.Unlock()
}*/

func TakeAction(w http.ResponseWriter, req *http.Request) {
    game.mutex.Lock()
    if BadGame(w, req) {
        game.mutex.Unlock()
        return
    }
    var act Action
    json.NewDecoder(req.Body).Decode(&act)
    if act.Verb == "" {
        game.mutex.Unlock()
        return
    }
    jsn,_ := json.Marshal(act)
    fmt.Printf("%s\n", jsn)
    gameUpd,err := game.TakeAction(&act) 
    if gameUpd != nil {
        //jsn,_ = json.Marshal(gameUpd)
        game.RecordAction(&act)
        game.RecordUpdate(gameUpd, false)
        //fmt.Fprintf(w, "%s\n", jsn)
        //game.mutex.Unlock()
        //return
    } else {
        actions := game.PlayerActions(game.Players[act.PlayerIdx])
        /*update := ActionResponse{
            Success: err == nil, 
            Actions: actions,
        }*/
        game.RecordAction(&act)
        game.RecordPossibleActions(act.PlayerIdx, actions)
        /*jsn,_ = json.Marshal(update)
        fmt.Fprintf(w, "%s\n", jsn)*/
    }
    jsn,_ = json.Marshal(err == nil)
    fmt.Fprintf(w, "%s\n", jsn)
    game.mutex.Unlock()
}

type HFunc func (http.ResponseWriter, *http.Request)

func Headers(fn HFunc) HFunc {
    return func (w http.ResponseWriter, req *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
        w.Header().Set("Access-Control-Allow-Headers",
            "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
        fn(w, req)
    }
}

func main() {
    http.HandleFunc("/game", Headers(NewGame))
    http.HandleFunc("/join", Headers(Join))
    http.HandleFunc("/update", Headers(Update))
    //http.HandleFunc("/actions", Headers(GetActions))
    http.HandleFunc("/action", Headers(TakeAction))
    http.ListenAndServe("0.0.0.0:8080", nil)
}
