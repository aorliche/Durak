package main

import (
    "fmt"
    "testing"
    //clone "github.com/huandu/go-clone"
)

/*func TestSearch(t *testing.T) {
    game := InitGame(0, false)
    state := InitGameState(game)
    act,val := state.Move(0, 0)
    fmt.Println(state.Trump)
    fmt.Println(ChainString(act),val)
    fmt.Println(HandString(state.Hands[0]))
    fmt.Println(HandString(state.Hands[1]))
}*/

/*func TestMystery(t *testing.T) {
    game := InitGame(0, false)
    state := InitGameState(game)
    for i,_ := range state.Hands[1] {
        c := state.Hands[1][i] 
        c.Rank = "?"
        c.Suit = "?"
    }
    act,val := state.Move(0, 0)
    fmt.Println(state.Trump)
    fmt.Println(ChainString(act),val)
    fmt.Println(HandString(state.Hands[0]))
}*/

/*func TestPartialMystery(t *testing.T) {
    game := InitGame(0, false)
    state := InitGameState(game)
    for i,_ := range state.Hands[1][:4] {
        c := state.Hands[1][i] 
        c.Rank = "?"
        c.Suit = "?"
    }
    act,val := state.Move(0, 0)
    fmt.Println(state.Trump)
    fmt.Println(ChainString(act),val)
    fmt.Println(HandString(state.Hands[0]))
    fmt.Println(HandString(state.Hands[1]))
}*/

/*func TestPassAfterPickup(t *testing.T) {
    game := InitGame(0, "Human")
    state := InitGameState(game)
    for i,_ := range state.Hands[1][:4] {
        c := state.Hands[1][i] 
        c.Rank = "?"
        c.Suit = "?"
    }
    act,who,val := state.Move(0, 0)
    fmt.Println(state.Trump)
    fmt.Println(ChainString(act),who,val)
    fmt.Println(HandString(state.Hands[0]))
    fmt.Println(HandString(state.Hands[1]))
}*/

/*func TestPlaySecondTrumpBug(t *testing.T) {
    game := InitGame(0, "Human")
    state := InitGameState(game)
    state.Hands[1][0].Rank = "6"
    state.Hands[1][0].Suit = "Clubs"
    state.Hands[1][1].Rank = "6"
    state.Hands[1][1].Suit = "Hearts"
    state.Trump = "Hearts"
    state.Defender = 0
    act,val := state.Move(1, 0)
    fmt.Println(state.Trump)
    fmt.Println(ChainString(act),val)
    fmt.Println(HandString(state.Hands[0]))
    fmt.Println(HandString(state.Hands[1]))
}*/

/*func TestGuessFinalCards(t *testing.T) {
    game := InitGame(0, "Human")
    mem := InitMemory(2)
    gameClone := clone.Clone(game).(*Game)
    state := InitGameState(game)
    for _,c := range state.Hands[0] {
        c.Rank = "?"
        c.Suit = "?"
    }
    mem.Discard = game.Deck
    game.Deck = make([]*Card,0) //game.Deck[:1]
    state.DeckSize = 0
    state.Hands[0] = mem.GuessFinalCards(state, 1)
    fmt.Println(HandString(gameClone.Players[0].Hand))
    fmt.Println(HandString(state.Hands[0]))
}*/

/*func TestSetKnownCards(t *testing.T) {
    game := InitGame(0, "Human")
    mem := InitMemory(2)
    gameClone := clone.Clone(game).(*Game)
    state := InitGameState(game)
    for _,c := range state.Hands[0] {
        c.Rank = "?"
        c.Suit = "?"
    }
    mem.AddCards(0, game.Players[0].Hand[:3])
    mem.Sizes[0] = 6
    mem.SetKnownCards(state, 1, 0)
    fmt.Println(HandString(mem.Hands[0]))
    fmt.Println(HandString(gameClone.Players[0].Hand))
    fmt.Println(HandString(state.Hands[0]))
}*/

func TestEndgame(t *testing.T) {
    game := InitGame(0, "Human")
    mem := InitMemory(2)
    state := InitGameState(game)
    for _,c := range state.Hands[0] {
        c.Rank = "?"
        c.Suit = "?"
    }
    mem.Discard = game.Deck
    game.Deck = make([]*Card,0) //game.Deck[:1]
    state.DeckSize = 0
    state.Defender = 0
    mem.SetKnownCards(state, 1, 0)
    act,val := state.Move(1, 0)
    fmt.Println(state.Trump)
    fmt.Println(ChainString(act),val)
    fmt.Println(HandString(state.Hands[0]))
    fmt.Println(HandString(state.Hands[1]))
}
