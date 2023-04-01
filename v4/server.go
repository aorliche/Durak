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
    idx, err := strconv.Atoi(req.URL.Query().Get("p"))
    if err != nil || idx < 0 || idx >= len(game.Players) {
        fmt.Fprintf(w, "No such player %v\n", idx)
        return 0,true
    }
    return idx,false
}

func NewGame(w http.ResponseWriter, req *http.Request) {
    game = InitGame()
    fmt.Fprintf(w, fmt.Sprint(game.PlayerNames()))
}

func GetActions(w http.ResponseWriter, req *http.Request) {
    if BadGameOrPlayer(w, req) {
        return
    }
    a := game.PlayerActions(game.Players[idx])
    jsn,_ := json.Marshal(a)
    fmt.Fprintf(w, "%s\n", jsn)
}

func TakeAction(w http.ResponseWriter, req *http.Request) {
    if BadGameOrPlayer(w, req) {
        return
    }
    var act Action
    err := json.NewDecoder(req.Body).Decode(&act)
    jsn,_ := json.Marshal(act)
    fmt.Println(jsn)
}

func main() {
    http.HandleFunc("/game", NewGame)
    http.HandleFunc("/actions", GetActions)
    http.HandleFunc("/action", TakeAction)
    http.ListenAndServe(":8080", nil)
}
