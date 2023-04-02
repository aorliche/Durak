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

func GetPlayerIdx(w http.ResponseWriter, req *http.Request) int {
    idx, err := strconv.Atoi(req.URL.Query().Get("p"))
    if err != nil || idx < 0 || idx >= len(game.Players) {
        fmt.Fprintf(w, "No such player %v\n", idx)
        return -1
    }
    return idx
}

func SendGameState(w http.ResponseWriter, req *http.Request) {
    state := GameState{Board: game.Board, Deck: len(game.Deck), Trump: game.Trump}
    jsn,_ := json.Marshal(state)
    fmt.Fprintf(w, "%s\n", jsn)
}

func NewGame(w http.ResponseWriter, req *http.Request) {
    game = InitGame()
    SendGameState(w, req)
}

func GetHand(w http.ResponseWriter, req *http.Request) {
    if BadGame(w, req) {
        return
    }
    idx := GetPlayerIdx(w, req)
    if idx == -1 {
        return 
    }
    jsn,_ := json.Marshal(game.Players[idx].Hand)
    fmt.Fprintf(w, "%s\n", jsn)
}

func GetActions(w http.ResponseWriter, req *http.Request) {
    if BadGame(w, req) {
        return
    }
    idx := GetPlayerIdx(w, req)
    if idx == -1 {
        return 
    }
    a := game.PlayerActions(game.Players[idx])
    jsn,_ := json.Marshal(a)
    fmt.Fprintf(w, "%s\n", jsn)
}

func TakeAction(w http.ResponseWriter, req *http.Request) {
    if BadGame(w, req) {
        return
    }
    idx := GetPlayerIdx(w, req)
    if idx == -1 {
        return 
    }
    var act Action
    json.NewDecoder(req.Body).Decode(&act)
    jsn,_ := json.Marshal(act)
    fmt.Println(jsn)
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
    http.HandleFunc("/hand", Headers(GetHand))
    http.HandleFunc("/actions", Headers(GetActions))
    http.HandleFunc("/action", Headers(TakeAction))
    http.ListenAndServe("0.0.0.0:8080", nil)
}
