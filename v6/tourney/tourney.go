// Genetic algorithm to find best eval function parameters for Durak AI

package main

import (
    //"encoding/json"
    "fmt"
    "math/rand"
    "sort"
    "time"

    . "github.com/aorliche/durak"
)

var popSize = 12
var keepSize = 5
var numGen = 20
var nGames = 2 // Each player plays each other player this many times
var nSimulGames = 16 // We have 12 cores (2 players each game) but seems we can overcommit
var nSeconds = 360 // 6 minute timeout per game

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
    return rand.Float64() < 0.4
}

// 50% Chance of mutating each particular parameter
// Sometimes a successful mutation mutates ints by zero
func Mutate(params *EvalParams) *EvalParams {
    npar := new(EvalParams)
    *npar = *params
    params = npar
    if change() {
        params.CardValueTrumpBonus += IntFrom(-2, 2)
    }
    if change() {
        params.CardValueCardOffset += IntFrom(-2, 2)
    }
    if change() {
        params.HandSizePickingUpMult += FloatFrom(-2.0, 2.0)
    }
    if change() {
        params.HandSizeSmallDeckLimit += IntFrom(-2, 2)
    }
    if change() {
        params.HandSizeSmallDeckMult += FloatFrom(-2.0, 2.0)
    }
    if change() {
        params.SmallDeckLimit += IntFrom(-2, 2)
    }
    if change() {
        params.NotLastWinnerValue += IntFrom(-10, 10)
    }
    if change() {
        params.HandMult += FloatFrom(-1, 1)
    }
    if change() {
        params.KnownMult += FloatFrom(-2.0, 2.0)
    }
    return params
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
    nRem := 0
    for i := 0; i < len(games); i++ {
        if i < nSimulGames {
            nRem++
            games[i] <- true
            continue
        }
        <-finished
        games[i] <- true
    }
    for i := 0; i < nRem; i++ {
        <-finished
    }
}

// 6 min game time limit
func PlayGame(pops []*EvalParams, results [][]int, i int, j int, start chan bool, finished chan bool) {
    game := InitGame(i, []string{"Medium", "Medium"})
    <-start
    kill1 := game.StartComputer("Medium", 0, pops[i], nil, nil)
    kill2 := game.StartComputer("Medium", 1, pops[j], nil, nil)
    n := 0
    for !game.CheckGameOver() {
        time.Sleep(1000 * time.Millisecond)
        n++
        if n > nSeconds {
            *kill1 = true
            *kill2 = true
            break
        }
    }
    fmt.Println(game.Recording.Winners)
    if len(game.Recording.Winners) > 0 {
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

type ArgSlice struct {
    sort.Interface
    Idcs []int
}

func (s ArgSlice) Swap(i, j int) {
   s.Interface.Swap(i, j)
   s.Idcs[i], s.Idcs[j] = s.Idcs[j], s.Idcs[i]
}

func Argsort(vals []int) []int {
    s := ArgSlice{sort.IntSlice(vals), make([]int, len(vals))}
    for i := 0; i < len(vals); i++ {
        s.Idcs[i] = i
    }
    sort.Sort(s)
    return s.Idcs
}

func main() {
    pops := InitPops()
    for g := 0; g < numGen; g++ {
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
        // Synchronize
        PlayGames(finished, starts)
        fitness := make([]int, popSize)
        for i := 0; i < popSize; i++ {
            for j := 0; j < popSize; j++ {
                fitness[i] += results[i][j]
            }
        }
        idcs := Argsort(fitness)
        fmt.Println(idcs)
        fmt.Println(results)
        //idcs := []int{2, 3, 1, 0}
        newPops := make([]*EvalParams, popSize)
        // Keep the best unchanged
        for i := 0; i < keepSize; i++ {
            newPops[i] = pops[idcs[len(idcs)-i-1]]
        }
        // Mutate the rest from the best
        for i := keepSize; i < popSize; i++ {
            newPops[i] = newPops[IntFrom(0, keepSize-1)]
            newPops[i] = Mutate(newPops[i])
        }
        for i := 0; i < popSize; i++ {
            fmt.Println(*newPops[i])
        }
        pops = newPops
    }
}
