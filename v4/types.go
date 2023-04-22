package main

import (
    //T "gorgonia.org/tensor"
    "sync"
)

/*type AI struct {
    Weights map[string]*T.Tensor
}*/

type Game struct {
    Key int
    Players []*Player
    Board *Board
    Deck []*Card
    Discard []*Card
    Trump *Card
    Defender int
    Turn int
    PickingUp bool
    Recording []*Record
    Versus string
    mutex sync.Mutex
}

type Record struct {
    Action *Action
    Actions [][]*Action
    // On hand boundaries
    Update *GameUpdate
}

type ActionResponse struct {
    Success bool
    Actions []*Action
}

/*type ActionsUpdate struct {
    PlayerIdx int
    Actions []*Action
}*/

type GameUpdate struct {
    Board *Board
    Deck int
    Trump *Card
    Players []*Player
    Actions [][]*Action
    Winner int
}

type Board struct {
    Plays []*Card
    Covers []*Card
}

type Player struct {
    Name string     `json:"omitempty"`
    Idx int
    Hand []*Card
}

// TODO Remember player having cards (e.g. from pickup or process of elimination at the end of deck)
type Card struct {
    Rank string     
    Suit string     
}

// No card matching, only predefined actions
type Action struct {
    PlayerIdx int
    Verb string     // Attack Defend Pickup Pass (Reverse later)
    Card *Card      // When covering Card covers Cover
    Cover *Card   
}

// Match card with respect to goal
// Get rid of small values, acquire trumps, get rid of cards in hand, acquire cards for opponent, get rid of trumps for opponent
// Goals change with time (goal weights change with game state... get rid of cards in hand at end of game)
// Search action space to see how goals affected
// 1. Hard coded goals, can be turned on or off
// 1a. Action search shallow
// 1b. Learning what actions help these goals: free to pick naive action but result of hand gives feedback signal to goals
// 2. Tell the computer it did a bad action (and tell it which goal it hindered)
