package main

import (
    //"encoding/json"
    "fmt"
    "math/rand"
    "reflect"
    "strings"
    T "gorgonia.org/tensor"
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

func (game *Game) PlayerNames() []string {
    names := make([]string, len(game.Players))
    for i,p := range game.Players {
        names[i] = p.Name
    }
    return names
}

func (game *Game) PlayerFromName(name string) *Player {
    for _,p := range game.Players {
        if p.Name == name {
            return p
        }
    }
    return nil
}

func (game *Game) PlayerActions(p *Player) []*Action {
    if p == game.GetAttacker() {
        return game.AttackerActions()
    } else {
        return game.DefenderActions()
    }
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
        for _,bc := range Cat(game.Board.Plays, game.Board.Covers) {
            for _,pc := range p.Hand {
                // Unique actions
                if bc != nil && bc.Rank == pc.Rank && IndexOfFn(res, func(act *Action) bool {return act.Card == pc}) == -1 {
                    act := Action{PlayerIdx: p.Idx, Verb: "Attack", Card: pc}
                    res = append(res, &act)
                }
            }
        }
    }
    if len(game.Board.Plays) > 0 {
        act := Action{PlayerIdx: p.Idx, Verb: "Pass"}
        res = append(res, &act)
    }
    return res
}

func (game *Game) DefenderActions() []*Action {
    res := make([]*Action, 0)
    p := game.GetDefender()
    for _,bp := range game.Board.Plays {
        for _,pc := range p.Hand {
            if pc.Beats(bp, game.Trump.Suit) {
                act := Action{PlayerIdx: p.Idx, Verb: "Defend", Card: pc, Cover: bp}
                res = append(res, &act)
            }
        }
    }
    // Get non-nil covers
    cLen := Count(game.Board.Covers, func (c *Card) bool {return c != nil})
    if len(game.Board.Plays) > 0 && cLen < len(game.Board.Plays) {
        act := Action{PlayerIdx: p.Idx, Verb: "Pickup"}
        res = append(res, &act)
    }
    return res
}

func (game *Game) TakeAction(act *Action) error {
    valid := false
    p := game.Players[act.PlayerIdx]
    //fmt.Println("---")
    for _,a := range game.PlayerActions(p) {
        //jsn,_ := json.Marshal(a)
        //fmt.Printf("%s\n", jsn)
        if reflect.DeepEqual(a, act) {
            valid = true
            break
        }
    }
    if !valid {
        return fmt.Errorf("Invalid action")
    }
    switch act.Verb {
        case "Attack": {
            p.Hand = Remove(p.Hand, act.Card)
            for _,c := range p.Hand {
                fmt.Println(c.Rank, c.Suit)
            }
            game.Board.Plays = append(game.Board.Plays, act.Card)
            game.Board.Covers = append(game.Board.Covers, nil)
        }
    }
    return nil
}

func (act *Action) ToTensor(game *Game) T.Tensor {
    rankFeat := float64(IndexOf(ranks, act.Card.Rank))
    suitFeat := float64(Ternary(act.Card.Suit == game.Trump.Suit, 1, 0))
    boardSizeFeat := float64(game.BoardSize())
    handSizeFeat := float64(len(game.Players[act.PlayerIdx].Hand))
    back := []float64{rankFeat, suitFeat, boardSizeFeat, handSizeFeat}
    return T.New(T.WithShape(2), T.WithBacking(back))
}

func InitDeck() []*Card {
    cards := make([]*Card, 0)
    for _,r := range ranks {
        for _,s := range suits {
            cards = append(cards, &Card{Rank: r, Suit: s})
        }
    }
    rand.Shuffle(len(cards), func(i, j int) { cards[i], cards[j] = cards[j], cards[i] })
    return cards
}

func InitPlayer(idx int) *Player {
    return &Player{Hand: make([]*Card, 0), Idx: idx, Name: fmt.Sprint(idx)}
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

func InitGame() *Game {
    game := Game{Deck: InitDeck(), Board: InitBoard(), Turn: 0, Discard: make([]*Card, 0), Players: []*Player{InitPlayer(0), InitPlayer(1)}, Defender: 1}
    game.Trump = game.Deck[0]
    game.DealAll()
    return &game
}

