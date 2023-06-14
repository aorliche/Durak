package main

import (
    "encoding/json"
    "fmt"
    "os"
    "net/http"
    "strconv"
    "strings"
    "time"
)

var games = make(map[int]*Game)

func NextGameIdx() int {
    max := -1
    for key := range games {
        if key > max {
            max = key
        }
    }
    return max+1
}

func JsonErr(s string) string {
    jsn, _ := json.Marshal(s)
    return string(jsn)
}

func GetGame(w http.ResponseWriter, req *http.Request) *Game {
    key, err := strconv.Atoi(req.URL.Query().Get("game"))
    if err != nil {
        fmt.Fprintf(w, "%s\n", JsonErr("No such game A"))
        return nil
    }
    game, ok := games[key]
    if !ok {
        fmt.Fprintf(w, "%s\n", JsonErr("No such game B"))
        return nil
    }
    return game
}

func List(w http.ResponseWriter, req *http.Request) {
    keys := make([]int, 0)
    for key := range games {
        if games[key].Versus == "Human" && !games[key].joined && games[key].CheckWinner() == -1 {
            keys = append(keys, key) 
        }
    }
    jsn, _ := json.Marshal(keys)
    fmt.Fprintf(w, "%s\n", jsn)
}

func Join(w http.ResponseWriter, req *http.Request) {
    game := GetGame(w, req)
    if game == nil { 
        return
    }
    game.joined = true
    Info(w, req)
}

func New(w http.ResponseWriter, req *http.Request) {
    comp := req.URL.Query().Get("computer")
    game := InitGame(NextGameIdx(), comp)
    games[game.Key] = game
    req.URL.RawQuery = fmt.Sprintf("p=0&game=%d", game.Key)
    Info(w, req)
}

func WriteGameIfWinner(w http.ResponseWriter, req *http.Request, game *Game) {
    if game.CheckWinner() != -1 {
        str := fmt.Sprintf("[\n\t%s\n]", strings.Join(game.Recording, ",\n\t"))
         ts := time.Now().Unix()
         err := os.WriteFile(fmt.Sprintf("games/%d.durak", ts), []byte(str), 0644)
         if err != nil {
             fmt.Println("Error writing game file")
         }
    }
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
    upd := Update{
        Key: game.Key, 
        Board: game.Board, 
        Deck: len(game.Deck), 
        Trump: game.Trump, 
        Players: game.MaskedPlayers(mp), 
        Actions: actions, 
        Winner: game.CheckWinner()}
    jsn,_ := json.Marshal(upd)
    fmt.Fprintf(w, "%s\n", jsn)
    if len(game.Recording) == 0 {
        // Unmask player cards
        upd.Players = game.Players
        jsn,_ = json.Marshal(upd)
        game.Recording = append(game.Recording, fmt.Sprintf("\"%s\"", game.Versus))
        game.Recording = append(game.Recording, string(jsn))
	    WriteGameIfWinner(w, req, game)
    }
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
    actJsn,_ := json.Marshal(act)
    //fmt.Printf("%s\n", jsn)
    upd := game.TakeAction(&act) 
    jsn,_ := json.Marshal(upd != nil)
    fmt.Fprintf(w, "%s\n", jsn)
    updJsn,_ := json.Marshal(upd)
    if upd != nil {
        game.Recording = append(game.Recording, string(actJsn))
        game.Recording = append(game.Recording, string(updJsn))
	    WriteGameIfWinner(w, req, game)
    }
    game.mutex.Unlock()
}

func Knowledge(w http.ResponseWriter, req *http.Request) {
    game := GetGame(w, req)
    if game == nil {
        return
    }
    jsn,_ := json.Marshal(game.memory)
    //fmt.Printf("%s\n", jsn)
    fmt.Fprintf(w, "%s\n", jsn)
}

type HFunc func (http.ResponseWriter, *http.Request)

func Headers(fn HFunc) HFunc {
    return func (w http.ResponseWriter, req *http.Request) {
        //fmt.Println(req.Method)
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
    http.HandleFunc("/memory", Headers(Knowledge))
    http.ListenAndServe("0.0.0.0:8080", nil)
}
