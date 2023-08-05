package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "os"
    "strconv"
    "time"

    "github.com/gorilla/websocket"
)

var games = make(map[int]*Game)
var upgrader = websocket.Upgrader{}

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
        if games[key].Versus == "Human" && !games[key].joined && games[key].Recording.Winner == -1 {
            keys = append(keys, key) 
        }
    }
    jsn, _ := json.Marshal(keys)
    fmt.Fprintf(w, "%s\n", jsn)
}

/*func Join(w http.ResponseWriter, req *http.Request) {
    game := GetGame(w, req)
    if game == nil { 
        return
    }
    game.joined = true
    Info(w, req)
}

// computer = Human, Easy, Medium
func New(w http.ResponseWriter, req *http.Request) {
    comp := req.URL.Query().Get("computer")
    game := InitGame(NextGameIdx(), comp)
    games[game.Key] = game
    if comp != "Human" {
        game.StartComputer(comp)
    }
    req.URL.RawQuery = fmt.Sprintf("p=0&game=%d", game.Key)
    Info(w, req)
}

type GameInfo struct {
    Key int
    State *GameState
    Memory *Memory
    Actions [][]Action
    DeckSize int
    Winner int
}

func (game *Game) MakeGameInfo(player int) *GameInfo {
    return &GameInfo{
        Key: game.Key,
        State: game.MaskUnknownCards(player),
        Memory: game.Memory,
        Actions: [][]Action{game.State.PlayerActions(0), game.State.PlayerActions(1)},
        DeckSize: len(game.Deck),
        Winner: game.Recording.Winner,
    }
}

func Info(w http.ResponseWriter, req *http.Request) {
    game := GetGame(w, req)
    if game == nil {
        return
    }
    game.mutex.Lock()
    p, err := strconv.Atoi(req.URL.Query().Get("p"))
    if err != nil || p < 0 || p > 1 {
        fmt.Fprintf(w, "%s\n", JsonErr("Bad player"))
        game.mutex.Unlock()
        return
    }
    // Check winner, write game if done
    if game.CheckWinner() != -1 {
        jsn, _ := json.Marshal(game.Recording)
        ts := time.Now().Unix()
        err := os.WriteFile(fmt.Sprintf("games/%d.durak", ts), jsn, 0644)
        if err != nil {
            fmt.Println("Error writing game file")
        }
    }
    info := game.MakeGameInfo(p)
    jsn, _ := json.Marshal(info)
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
    // Bad action sent?
    if act.Card == Card(0) && act.Covering == Card(0) {
        game.mutex.Unlock()
        return
    }
    fmt.Printf("%s\n", act.ToStr())
    game.TakeAction(act) 
    info := game.MakeGameInfo(act.Player)
    jsn,_ := json.Marshal(info)
    fmt.Fprintf(w, "%s\n", jsn)
    game.mutex.Unlock()
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
}*/

var upgrader = websocket.Upgrader{} // Default options

type Request struct {
   Type string 
   Game int
   Computer string
   Action Action
}

        case "info": {
            info := game.MakeGameInfo(player)
            jsn, _ := json.Marshal(info)

func socket(w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Println(err)
        return
    }
    player := -1
    for {
        msgType, msg, err := conn.ReadMessage()
        if err != nil {
            log.Println(err)
            return  
        }
        // Do we ever get any other types of messages?
        if msgType != websocket.TextMessage {
            return
        }
        var req Request
        json.NewDecoder(msg).Decode(&req)
        switch req.Type {
            case "List" : {
                for key := range games {
                    if games[key].Versus == "Human" && !games[key].joined && games[key].Recording.Winner == -1 {
                        keys = append(keys, key) 
                    }
                }
                jsn, _ := json.Marshal(keys)
                err = conn.WriteMessage(websocket.TextMessage, jsn)
            }
            case "New": {
                if player != -1 {
                    log.Println("Player already joined")
                    return
                }
                player = 0
                game := InitGame(NextGameIdx(), req.Computer)
                games[game.Key] = game
                if req.Computer != "Human" {
                    game.StartComputer(comp)
                }
                Info(conn, p, game.Key)
            }
            case "Join": {
                if player != -1 {
                    log.Println("Player already joined")
                    return
                }
                player = 1
                game = games[req.Game]
                if game == nil { 
                    log.Println("No such game", req.Game)
                    return
                }
                game.joined = true
                Info(conn, p, game.Key)
            }
            case "Action": {
                for i,p := range game.State.Players {

                }
            }
        }

    }
    defer conn.Close()
}

func main() {
    log.SetFlags(0)
    serveLocalFiles()
    http.HandleFunc("/ws", socket)
    http.HandleFunc("/list", List)
    http.HandleFunc("/new", New)
    http.HandleFunc("/join", Join)
    //http.HandleFunc("/info", Headers(Info))
    http.HandleFunc("/action", TakeAction)
    log.Fatal(http.ListenAndServe("0.0.0.0:8080", nil))
}
