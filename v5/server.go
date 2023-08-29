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
   Players []string
   Action *Action
}

type GameInfo struct {
    Key int
    State *GameState
    Memory *Memory
    Actions [][]Action
    DeckSize int
    Winners []int
}

func (game *Game) MakeGameInfo(player int) *GameInfo {
    acts := make([][]Action, len(game.State.Hands)) 
    acts[player] = game.State.PlayerActions(player)
    return &GameInfo{
        Key: game.Key,
        State: game.MaskUnknownCards(player),
        Memory: game.Memory,
        Actions: acts,
        DeckSize: len(game.Deck),
        Winners: game.Recording.Winners,
    }
}

func SendInfo(player int, game *Game) {
    game.mutex.Lock()
    conn := game.conns[player]
    info := game.MakeGameInfo(player)
    //log.Println(info.State.ToStr())
    //log.Println(game.Recording.Winners)
    jsn, _ := json.Marshal(info)
    conn.WriteMessage(websocket.TextMessage, jsn)   
    game.mutex.Unlock()
}

func (game *Game) WriteGame() {
    jsn, _ := json.Marshal(game.Recording)
    ts := time.Now().Unix()
    err := os.WriteFile(fmt.Sprintf("games/%d.durak", ts), jsn, 0644)
    if err != nil {
        log.Println("Error writing game file")
    }
}

func (game *Game) SendInfoHumans() {
    for i,p := range game.Players {
        // Check nil for automated test
        if p == "Human" && game.conns[i] != nil {
            SendInfo(i, game)
        }
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
                    // Check if has not been won and has open slots
                    game := games[key]
                    if game.CheckGameOver() {
                        continue
                    }
                    for i := 0; i < len(game.joined); i++ {
                        if !game.joined[i] {
                            keys = append(keys, key) 
                            break
                        }
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
                game := InitGame(NextGameIdx(), req.Players)
                games[game.Key] = game
                for i,typ := range req.Players {
                    if typ != "Human" {
                        game.StartComputer(typ, i)
                        game.joined[i] = true
                    }
                }
                game.conns[0] = conn
                game.joined[0] = true
                SendInfo(player, game)
            }
            case "Join": {
                if player != -1 {
                    log.Println("Player already joined")
                    continue
                }
                game := games[req.Game]
                if game == nil { 
                    log.Println("No such game", req.Game)
                    continue
                }
                for i := 1; i < len(game.joined); i++ {
                    if !game.joined[i] {
                        player = i
                    }
                }
                if player == -1 {
                    log.Println("No open slots")
                    continue
                }
                game.joined[player] = true
                game.conns[player] = conn
                SendInfo(player, game)
            }
            case "Action": {
                game := games[req.Game]
                if game == nil { 
                    log.Println("No such game", req.Game)
                    continue
                }
                if game.Players[req.Action.Player] == "Human" && !game.joined[req.Action.Player] {
                    log.Println("Player not joined")
                    continue
                }
                if game.CheckGameOver() {
                    log.Println("Game already won")
                    continue
                }
                // Bad action sent?
                if req.Action.Card == Card(0) && req.Action.Covering == Card(0) {
                    log.Println("Bad action")
                    continue
                }
                log.Printf("%s\n", req.Action.ToStr())
                game.mutex.Lock()
                game.TakeAction(*req.Action) 
                game.mutex.Unlock()
                // Check winner, write game if done
                // Also in computer.go checks for computer games
                if game.CheckGameOver() {
                    game.WriteGame()
                }
                // Send update to all humans
                game.SendInfoHumans()
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
