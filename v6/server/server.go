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

    . "github.com/aorliche/durak"
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
   Name string
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
    Names []string
}

func MakeGameInfo(game *Game, player int) *GameInfo {
    acts := make([][]Action, len(game.State.Hands)) 
    acts[player] = game.State.PlayerActions(player)
    return &GameInfo{
        Key: game.Key,
        State: game.MaskUnknownCards(player),
        Memory: game.Memory,
        Actions: acts,
        DeckSize: len(game.Deck),
        Winners: game.Recording.Winners,
        Names: game.Names,
    }
}

func SendInfo(player int, game *Game) {
    game.Mutex.Lock()
    conn := game.Conns[player]
    info := MakeGameInfo(game, player)
    //log.Println(info.State.ToStr())
    //log.Println(game.Recording.Winners)
    jsn, _ := json.Marshal(info)
    conn.WriteMessage(websocket.TextMessage, jsn)   
    game.Mutex.Unlock()
}

func WriteGame(game *Game) {
    jsn, _ := json.Marshal(game.Recording)
    ts := time.Now().Unix()
    err := os.WriteFile(fmt.Sprintf("games/%d.durak", ts), jsn, 0644)
    if err != nil {
        log.Println("Error writing game file")
    }
}

func SendInfoHumans(game *Game) {
    for i,p := range game.Players {
        // Check nil for automated test
        if p == "Human" && game.Conns[i] != nil {
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
                    for i := 0; i < len(game.Joined); i++ {
                        if !game.Joined[i] {
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
                game.Names[0] = req.Name
                games[game.Key] = game
                j := 0
                for i,typ := range req.Players {
                    if typ != "Human" {
                        j += 1
                        game.StartComputer(typ, i, func (game *Game) {
                            SendInfoHumans(game)
                        },
                        func (game *Game) {
                            WriteGame(game)
                        })
                        game.Joined[i] = true
                        game.Names[i] = fmt.Sprintf("%s%d", typ, j)
                    }
                }
                game.Conns[0] = conn
                game.Joined[0] = true
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
                for i := 1; i < len(game.Joined); i++ {
                    if !game.Joined[i] {
                        player = i
                        game.Names[i] = req.Name
                        break
                    }
                }
                if player == -1 {
                    log.Println("No open slots")
                    continue
                }
                game.Joined[player] = true
                game.Conns[player] = conn
                SendInfoHumans(game)
            }
            case "Action": {
                game := games[req.Game]
                if game == nil { 
                    log.Println("No such game", req.Game)
                    continue
                }
                if player == -1 || player > len(game.Joined) {
                    log.Println("Invalid player")
                    continue
                }
                if game.Players[player] == "Human" && !game.Joined[player] {
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
                game.Mutex.Lock()
                game.TakeAction(*req.Action) 
                game.Mutex.Unlock()
                // Check winner, write game if done
                // Also in computer.go checks for computer games
                if game.CheckGameOver() {
                    WriteGame(game)
                }
                // Send update to all humans
                SendInfoHumans(game)
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
    ServeLocalFiles([]string{"", "/cards/backs", "/cards/fronts", "/images"})
    http.HandleFunc("/ws", Socket)
    log.Fatal(http.ListenAndServe(":8000", nil))
}
