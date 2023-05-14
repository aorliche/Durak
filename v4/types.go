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
    joined bool
    mutex sync.Mutex
    memory *Memory
}

type Record struct {
    Action *Action
    Update *Update
}

type Update struct {
    Key int
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
    Idx int
    Hand []*Card
}

type Memory struct {
    Hands [][]*Card 
    Sizes []int
    Discard []*Card
}

type Card struct {
    Rank string     
    Suit string     
}

type Action struct {
    PlayerIdx int
    Verb string     // Attack Defend Pickup Pass Reverse
    Card *Card      
    Cover *Card     // When covering, Card covers Cover, otherwise nil
}
