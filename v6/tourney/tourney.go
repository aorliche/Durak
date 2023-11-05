// Genetic algorithm to find best eval function parameters for Durak AI

package main

import (
    //"encoding/json"
    "fmt"
    "math/rand"
    "time"

    . "github.com/aorliche/durak"
)

var popSize = 10
var numGen = 100
var nGames = 1 // Each player plays each other player this many times
var nSimulGames = 2 // We have 12 cores (2 players each game)

// Inclusive
func IntFrom(a int, b int) int {
    return a + rand.Intn(b-a+1)
}

// Inclusive
func FloatFrom(a float64, b float64) float64 {
    return a + rand.Float64()*(b-a)
}

// Based on default parameters
func RandomEvalParams() *EvalParams {
    params := EvalParams{
        IntFrom(3, 9),
        IntFrom(2, 6),
        FloatFrom(1.0, 2.0),
        IntFrom(1, 5),
        FloatFrom(0.5, 3.0),
        IntFrom(0, 5),
        IntFrom(10, 50),
        FloatFrom(0.5, 2.0),
        FloatFrom(1.0, 4.0),
    }
    return &params
}

func change() bool {
    return rand.Float64() < 0.2 
}

// 20% Chance of mutating each particular parameter
// Sometimes a successful mutation mutates ints by zero
func Mutate(params *EvalParams) {
    if change() {
        params.CardValueTrumpBonus += IntFrom(-1, 1)
    }
    if change() {
        params.CardValueCardOffset += IntFrom(-1, 1)
    }
    if change() {
        params.HandSizePickingUpMult += FloatFrom(-1.0, 1.0)
    }
    if change() {
        params.HandSizeSmallDeckLimit += IntFrom(-1, 1)
    }
    if change() {
        params.HandSizeSmallDeckMult += FloatFrom(-1.0, 1.0)
    }
    if change() {
        params.SmallDeckLimit += IntFrom(-1, 1)
    }
    if change() {
        params.NotLastWinnerValue += IntFrom(-10, 10)
    }
    if change() {
        params.HandMult += FloatFrom(-0.5, 0.5)
    }
    if change() {
        params.KnownMult += FloatFrom(-1.0, 1.0)
    }
}

func InitPops() []*EvalParams {
    pops := make([]*EvalParams, popSize)
    for i := 0; i < popSize; i++ {
        pops[i] = RandomEvalParams()
    }
    return pops
}

// Blocks
func PlayGames(finished chan bool, games []chan bool) {
    for i := 0; i < len(games); i++ {
        if i < nSimulGames {
            games[i] <- true
            continue
        }
        <-finished
        games[i] <- true
    }
}

func PlayGame(pops []*EvalParams, results [][]int, i int, j int, start chan bool, finished chan bool) {
    game := InitGame(i, []string{"Medium", "Medium"})
    <-start
    game.StartComputer("Medium", 0, pops[i], nil, nil)
    game.StartComputer("Medium", 1, pops[j], nil, nil)
    n := 0
    for !game.CheckGameOver() {
        time.Sleep(1000 * time.Millisecond)
        n++
        if n > 500 {
            break
        }
    }
    fmt.Println(game.Recording.Winners)
    if len(game.Recording.Winners) > 1 {
        w := game.Recording.Winners[0]
        if w == 0 {
            results[i][j] += 1
            results[j][i] += -1
        } else {
            results[i][j] += -1
            results[j][i] += 1
        }
    }
    finished <- true
}

func main() {
    pops := InitPops()
    results := make([][]int, popSize)
    for i := 0; i < popSize; i++ {
        results[i] = make([]int, popSize)
    }
    finished := make(chan bool)
    starts := make([]chan bool, 0)
    for i := 0; i < popSize; i++ {
        for j := i+1; j < popSize; j++ {
            for n := 0; n < nGames; n++ {
                starts = append(starts, make(chan bool))
                go PlayGame(pops, results, i, j, starts[len(starts)-1], finished)
            }
        }
    }
    PlayGames(finished, starts)
    fmt.Println(results)
}
