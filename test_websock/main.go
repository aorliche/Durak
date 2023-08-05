package main

import (
    "fmt"
    "log"
    "net/http"

    "github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{} // use default options

func home(w http.ResponseWriter, r *http.Request) {
    http.ServeFile(w, r, "index.html")
}

func wsEndpoint(w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Println(err)
        return
    }
    for {
        msgType, msg, err := conn.ReadMessage()
        if err != nil {
            log.Println(err)
            return
        }
        fmt.Println(msgType, string(msg))
        err = conn.WriteMessage(msgType, msg)
        if err != nil {
            log.Println(err)
            return
        }
    }
}

func main() {
    fmt.Println("Hello World")
    http.HandleFunc("/", home)
    http.HandleFunc("/ws", wsEndpoint)    
    log.Fatal(http.ListenAndServe(":8080", nil))
}
