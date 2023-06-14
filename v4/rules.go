package main

import (
    //"encoding/json"
    "fmt"
    //"math/rand"
    "reflect"
    "strings"
    //T "gorgonia.org/tensor"
    //"github.com/go-test/deep"
)

var suits = []string{"Clubs", "Spades", "Hearts", "Diamonds"}
var ranks = []string{"6", "7", "8", "9", "10", "Jack", "Queen", "King", "Ace"}

func (this *Card) Beats(other *Card, trumpSuit string) bool {
    if this.Suit == trumpSuit && this.Suit != other.Suit {
        return true
    }
    return IndexOf(ranks, this.Rank) > IndexOf(ranks, other.Rank) && this.Suit == other.Suit
}

func (card *Card) ToStr() string {
    return fmt.Sprintf("%s of %s", card.Rank, card.Suit)
}

func (game *Game) BoardSize() int {
    return len(game.Board.Plays) + len(game.Board.Covers)
}

func (game *Game) GetByRole(role string) *Player {
    for i,p := range game.Players {
        if role == "Defender" && i == game.Defender {
            return p
        }
        if role == "Attacker" && i != game.Defender {
            return p
        }
    }
    return nil
}

func (game *Game) GetAttacker() *Player {
    return game.GetByRole("Attacker")
}

func (game *Game) GetDefender() *Player {
    return game.GetByRole("Defender")
}

func (game *Game) PlayerActions(p *Player) []*Action {
    if p == game.GetAttacker() {
        return game.AttackerActions()
    } else {
        return game.DefenderActions()
    }
}

func (game *Game) MaskedPlayers(pIdx int) []*Player {
    po := game.Players[pIdx]
    hand := make([]*Card, len(po.Hand))
    for i := 0; i < len(hand); i++ {
        hand[i] = &Card{Rank: "Unkown", Suit: "Unknown"}
    }
    p := Player{Idx: pIdx, Hand: hand}
    players := make([]*Player, 2)
    players[1-pIdx] = game.Players[1-pIdx]
    players[pIdx] = &p
    return players
}

func (board *Board) Covered() int {
    return Count(board.Covers, func (c *Card) bool {return c != nil})
}

func (game *Game) AttackerActions() []*Action {
    res := make([]*Action, 0)
    p := game.GetAttacker()
    if game.BoardSize() == 0 {
        for _,card := range p.Hand {
            act := Action{PlayerIdx: p.Idx, Verb: "Attack", Card: card}
            res = append(res, &act)
        }
    } else {
        for _,bc := range Cat(game.Board.Plays, NotNil(game.Board.Covers)) {
            for _,pc := range p.Hand {
                // Unique actions
                if bc != nil && bc.Rank == pc.Rank && IndexOfFn(res, func(act *Action) bool {return act.Card == pc}) == -1 {
                    act := Action{PlayerIdx: p.Idx, Verb: "Attack", Card: pc}
                    res = append(res, &act)
                }
            }
        }
    }
    if game.PickingUp || (game.Board.Covered() == len(game.Board.Plays) && len(game.Board.Plays) > 0) {
        act := Action{PlayerIdx: p.Idx, Verb: "Pass"}
        res = append(res, &act)
    }
    return res
}

func (board *Board) ReverseRank() string {
    if len(board.Plays) == 0 || len(NotNil(board.Covers)) > 0 {
        return ""
    }
    r := board.Plays[0].Rank
    for _,c := range Cat(board.Plays, NotNil(board.Covers)) {
        if c.Rank != r {
            return ""
        }
    }
    return r
}

func (game *Game) DefenderActions() []*Action {
    res := make([]*Action, 0)
    if game.PickingUp {
        return res
    }
    p := game.GetDefender()
    revRank := game.Board.ReverseRank()
    if revRank != "" {
        for _,pc := range p.Hand {
            if pc.Rank == revRank {
                act := Action{PlayerIdx: p.Idx, Verb: "Reverse", Card: pc}
                res = append(res, &act)
            }
        }
    }
    for i,bp := range game.Board.Plays {
        if game.Board.Covers[i] != nil {
            continue
        }
        for _,pc := range p.Hand {
            if pc.Beats(bp, game.Trump.Suit) {
                act := Action{PlayerIdx: p.Idx, Verb: "Defend", Card: pc, Cover: bp}
                res = append(res, &act)
            }
        }
    }
    // Get non-nil covers
    if len(game.Board.Plays) > 0 && game.Board.Covered() < len(game.Board.Plays) {
        act := Action{PlayerIdx: p.Idx, Verb: "Pickup"}
        res = append(res, &act)
    }
    return res
}

/*func (game *Game) ReverseRank() string {
    if len(game.Board.Plays) == 0 || len(game.Board.Covers) > 0 {
        return ""
    }
    r := game.Board.Plays[0].Rank
    for _,c := range Cat(game.Board.Plays, NotNil(game.Board.Covers)) {
        if c.Rank != r {
            return ""
        }
    }
    return r
}*/

func (game *Game) TakeAction(act *Action) *Update {
    valid := false
    p := game.Players[act.PlayerIdx]
    for _,a := range game.PlayerActions(p) {
        if reflect.DeepEqual(a, act) {
            valid = true
            break
        }
    }
    if !valid {
        fmt.Println("Not valid")
        return nil
    }
    //fmt.Println(act.Verb)
    switch act.Verb {
        case "Attack": {
            p.Hand = Remove(p.Hand, act.Card)
            game.Board.Plays = append(game.Board.Plays, act.Card)
            game.Board.Covers = append(game.Board.Covers, nil)
            game.memory.RemoveCard(p.Idx, act.Card)
        }
        case "Defend": {
            p.Hand = Remove(p.Hand, act.Card)
            idx := IndexOf(game.Board.Plays, act.Cover)
            game.Board.Covers[idx] = act.Card
            game.memory.RemoveCard(p.Idx, act.Card)
        }
        case "Pickup": {
            game.PickingUp = true
        }
        case "Reverse": {
            p.Hand = Remove(p.Hand, act.Card)
            game.Board.Plays = append(game.Board.Plays, act.Card)
            game.Board.Covers = append(game.Board.Covers, nil)
            game.Defender = 1-game.Defender
            game.memory.RemoveCard(p.Idx, act.Card)
        }
        case "Pass": {
            board := Cat(game.Board.Plays, NotNil(game.Board.Covers))
            if game.Board.Covered() < len(game.Board.Plays) {
                game.GetDefender().Hand = append(game.GetDefender().Hand, board...) 
                game.memory.AddCards(1-p.Idx, board)
            } else {
                game.memory.DiscardCards(board)
            }
            game.Board.Plays = make([]*Card,0)
            game.Board.Covers = make([]*Card,0)
            game.Deal(game.GetAttacker())
            game.Deal(game.GetDefender())
            game.Turn += 1
            if !game.PickingUp {
                game.Defender = 1-game.Defender
            }
            game.PickingUp = false
            game.memory.SetSizes(game.Players)
        }
    }
    // Only used for recording
    actions := make([][]*Action,0)
    for _,p := range game.Players {
        actions = append(actions, game.PlayerActions(p))
    }
    return &Update{
        Board: game.Board, 
        Deck: len(game.Deck), 
        Trump: game.Trump, 
        Players: game.Players,
        Actions: actions,
        Winner: game.CheckWinner(),
    }
}

func (game *Game) CheckWinner() int {
    for i,p := range game.Players {
        if len(p.Hand) == 0 && len(game.Deck) == 0 {
            return i
        }
    }
    return -1
}

/*func (act *Action) ToTensor(game *Game) T.Tensor {
    rankFeat := float64(IndexOf(ranks, act.Card.Rank))
    suitFeat := float64(Ternary(act.Card.Suit == game.Trump.Suit, 1, 0))
    boardSizeFeat := float64(game.BoardSize())
    handSizeFeat := float64(len(game.Players[act.PlayerIdx].Hand))
    back := []float64{rankFeat, suitFeat, boardSizeFeat, handSizeFeat}
    return T.New(T.WithShape(2), T.WithBacking(back))
}*/

func InitDeck() []*Card {
    cards := make([]*Card, 0)
    for _,r := range ranks {
        for _,s := range suits {
            cards = append(cards, &Card{Rank: r, Suit: s})
        }
    }
    Shuffle(cards)
    return cards
}

func InitPlayer(idx int) *Player {
    return &Player{Hand: make([]*Card, 0), Idx: idx}
}

func InitBoard() *Board {
    return &Board{Plays: make([]*Card,0), Covers: make([]*Card,0)}
}

func (game *Game) Draw() *Card {
    if len(game.Deck) == 0 {
        return nil
    }
    card := game.Deck[len(game.Deck)-1]
    game.Deck = game.Deck[:len(game.Deck)-1]
    return card
}

func (game *Game) Deal(p *Player) {
    for len(game.Deck) > 0 && len(p.Hand) < 6 {
        p.Hand = append(p.Hand, game.Draw())
    }
}

func (game *Game) DealAll() {
    for _,p := range game.Players {
        game.Deal(p)
    }
}

func (game *Game) ToStr() string {
    str := make([]string, 0)
    for _,p := range game.Players {
        h := Apply(p.Hand, func(c *Card) string {return c.ToStr()})
        str = append(str, strings.Join(h, ", "))
    }
    str = append(str, fmt.Sprintf("Trump: %s", game.Trump.ToStr()))
    return strings.Join(str, "\n")
}

func InitGame(key int, comp string) *Game {
    game := Game{
        Key: key,
        Deck: InitDeck(), 
        Board: InitBoard(), 
        Turn: 0, 
        Discard: make([]*Card, 0), 
        Players: []*Player{InitPlayer(0), InitPlayer(1)}, 
        Defender: 1, 
        PickingUp: false,
        Recording: make([]string, 0),
        Versus: comp, 
        joined: false,
        memory: InitMemory(2)}
    fmt.Println(comp)
    game.Trump = game.Deck[0]
    game.DealAll()
    if comp == "Easy" {
        go RandomLoop(&game)
    } else if comp == "Medium" {
        go MediumLoop(&game)
    }
    return &game
}

