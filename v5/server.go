package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "os"
    "time"

    "github.com/gorilla/websocket"
)

var games = make(map[int]*Game)
var upgrader = websocket.Upgrader{} // Default options

func NextGameIdx() int {
    max := -1
    for key := range games {
        if key > max {
            max = key
        }
    }
    return max+1
}

type Request struct {
   Type string 
   Game int
   Computer string
   Action *Action
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

func SendInfo(player int, game *Game) {
    game.mutex.Lock()
    conn := game.conns[player]
    info := game.MakeGameInfo(player)
    game.mutex.Unlock()
    log.Println(game.Recording.Winner)
    jsn, _ := json.Marshal(info)
    conn.WriteMessage(websocket.TextMessage, jsn)   
}

func (game *Game) WriteGame() {
    jsn, _ := json.Marshal(game.Recording)
    ts := time.Now().Unix()
    err := os.WriteFile(fmt.Sprintf("games/%d.durak", ts), jsn, 0644)
    if err != nil {
        fmt.Println("Error writing game file")
    }
}

func Socket(w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Println(err)
        return
    }
    defer conn.Close()
    player := -1
    for {
        msgType, msg, err := conn.ReadMessage()
        if err != nil {
            log.Println(err)
            return  
        }
        // Do we ever get any other types of messages?
        if msgType != websocket.TextMessage {
            log.Println("Not a text message")
            return
        }
        var req Request
        json.NewDecoder(bytes.NewBuffer(msg)).Decode(&req)
        //log.Println(string(msg))
        switch req.Type {
            case "List" : {
                keys := make([]int, 0)
                for key := range games {
                    if games[key].Versus == "Human" && !games[key].joined && games[key].Recording.Winner == -1 {
                        keys = append(keys, key) 
                    }
                }
                jsn, _ := json.Marshal(keys)
                err = conn.WriteMessage(websocket.TextMessage, jsn)
                if err != nil {
                    log.Println(err)
                    continue
                }
            }
            case "New": {
                if player != -1 {
                    log.Println("Player already joined")
                    continue
                }
                player = 0
                game := InitGame(NextGameIdx(), req.Computer)
                games[game.Key] = game
                if req.Computer != "Human" {
                    game.StartComputer(req.Computer)
                }
                game.conns[0] = conn
                SendInfo(player, game)
            }
            case "Join": {
                if player != -1 {
                    log.Println("Player already joined")
                    continue
                }
                player = 1
                game := games[req.Game]
                if game == nil { 
                    log.Println("No such game", req.Game)
                    continue
                }
                game.joined = true
                game.conns[1] = conn
                SendInfo(player, game)
            }
            case "Action": {
                game := games[req.Game]
                if game == nil { 
                    log.Println("No such game", req.Game)
                    continue
                }
                if game.Versus == "Human" && !game.joined {
                    log.Println("Player not joined")
                    continue
                }
                if game.Recording.Winner != -1 {
                    log.Println("Game already won")
                    continue
                }
                // TODO multiple computer players multiple threads
                game.mutex.Lock()
                // Bad action sent?
                if req.Action.Card == Card(0) && req.Action.Covering == Card(0) {
                    log.Println("Bad action")
                    game.mutex.Unlock()
                    continue
                }
                fmt.Printf("%s\n", req.Action.ToStr())
                game.TakeAction(*req.Action) 
                game.CheckWinner()
                // Check winner, write game if done
                // Also in computer.go checks for computer games
                if game.Versus == "Human" && game.Recording.Winner != -1 {
                    game.WriteGame()
                }
                game.mutex.Unlock()
                // Send update to all
                if game.Versus == "Human" { 
                    SendInfo(0, game)
                    SendInfo(1, game)
                } else {
                    SendInfo(0, game)
                }
            }
        }
    }
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
func ServeStatic(w http.ResponseWriter, req *http.Request, file string) {
    w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
    http.ServeFile(w, req, file)
}

func ServeLocalFiles(dirs []string) {
    for _, dirName := range dirs {
        fsDir := "static/" + dirName
        dir, err := os.Open(fsDir)
        if err != nil {
            log.Fatal(err)
        }
        files, err := dir.Readdir(0)
        if err != nil {
            log.Fatal(err)
        }
        for _, v := range files {
            //fmt.Println(v.Name(), v.IsDir())
            if v.IsDir() {
                continue
            }
            reqFile := dirName + "/" + v.Name()
            file := fsDir + "/" + v.Name()
            http.HandleFunc(reqFile, Headers(func (w http.ResponseWriter, req *http.Request) {ServeStatic(w, req, file)}))
        }
    }
}

func main() {
    log.SetFlags(0)
    ServeLocalFiles([]string{"", "/cards/backs", "/cards/fronts"})
    http.HandleFunc("/ws", Socket)
    log.Fatal(http.ListenAndServe(":8000", nil))
}
