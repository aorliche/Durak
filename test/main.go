package main

import (
    "time"
    "os"
    "encoding/json"
    "fmt"
)

type Cake struct {
    Name string
}

func main() {
    glados := Cake{Name: "OpenAI"}
    jsn,_ := json.Marshal(glados)
    ts := time.Now().Unix()
    os.WriteFile(fmt.Sprintf("games/%d.durak", ts), jsn, 0644)
}
