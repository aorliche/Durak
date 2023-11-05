// Parametrize EvalNode search eval functions for genetic algorithms to find best parameters

package main

import (
    "encoding/json"
    "fmt"

    . "github.com/aorliche/durak"
)

func main() {
    game := InitGame(0, []string{"Human", "Human"})
    jsn, _ := json.Marshal(game)
    fmt.Println(string(jsn))
}
