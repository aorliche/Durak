package main

import (
    "encoding/json"
    "fmt"
    "net/http"
    //"net/url"
    "strconv"
)

var games = make(map[int]*Game)

func NextGameIdx() int {
    max := -1
    for key := range games {
        if key > max {
            max = key
        }
    }
    return max
}

func JsonErr(s string) string {
    jsn, _ := json.Marshal(s)
    return string(jsn)
}

func GetGame(w http.ResponseWriter, req *http.Request) *Game {
    key, err := strconv.Atoi(req.URL.Query().Get("game"))
    fmt.Println(key)
    if err != nil {
        fmt.Fprintf(w, "%s\n", JsonErr("No such game A"))
        return nil
    }
    game, ok := games[key]
    fmt.Println(game)
    if !ok {
        fmt.Fprintf(w, "%s\n", JsonErr("No such game B"))
        return nil
    }
    return game
}

func List(w http.ResponseWriter, req *http.Request) {
    keys := make([]int, len(games))
    for key := range games {
        keys = append(keys, key) 
    }
    jsn, _ := json.Marshal(keys)
    fmt.Fprintf(w, "%s\n", jsn)
}

func Join(w http.ResponseWriter, req *http.Request) {
    game := GetGame(w, req)
    if game == nil { 
        return
    }
    Info(w, req)
}

func New(w http.ResponseWriter, req *http.Request) {
    game := InitGame(NextGameIdx())
    fmt.Println(game.Key)
    games[game.Key] = game
    req.URL.RawQuery = fmt.Sprintf("p=0&game=%d", game.Key)
    Info(w, req)
}

func Info(w http.ResponseWriter, req *http.Request) {
    game := GetGame(w, req)
    if game == nil {
        return
    }
    game.mutex.Lock()
    actions := make([][]*Action, 0)
    for _,p := range game.Players {
        actions = append(actions, game.PlayerActions(p))
    }
    p, err := strconv.Atoi(req.URL.Query().Get("p"))
    if err != nil || p < 0 || p > 1 {
        fmt.Fprintf(w, "%s\n", JsonErr("Bad player"))
        game.mutex.Unlock()
        return
    }
    mp := 1-p
    upd := GameUpdate{Board: game.Board, Deck: len(game.Deck), Trump: game.Trump, Players: game.MaskedPlayers(mp), Actions: actions, Winner: game.CheckWinner()}
    game.RecordUpdate(&upd, true)
    jsn,_ := json.Marshal(upd)
    fmt.Fprintf(w, "%s\n", jsn)
    game.mutex.Unlock()
}

func TakeAction(w http.ResponseWriter, req *http.Request) {
    game := GetGame(w, req)
    if game == nil {
        return
    }
    game.mutex.Lock()
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
        game.RecordAction(&act)
        game.RecordUpdate(gameUpd, false)
    } else {
        actions := game.PlayerActions(game.Players[act.PlayerIdx])
        game.RecordAction(&act)
        game.RecordPossibleActions(act.PlayerIdx, actions)
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
    http.HandleFunc("/list", Headers(List))
    http.HandleFunc("/new", Headers(New))
    http.HandleFunc("/join", Headers(Join))
    http.HandleFunc("/info", Headers(Info))
    http.HandleFunc("/action", Headers(TakeAction))
    http.ListenAndServe("0.0.0.0:8080", nil)
}
