package main

import (
    "fmt"
    "math/rand"
    //"reflect"
    "strings"
    T "gorgonia.org/tensor"
)

var suits = []string{"Clubs", "Spades", "Hearts", "Diamonds"}
var ranks = []string{"6", "7", "8", "9", "10", "Jack", "Queen", "King", "Ace"}

func (this *Card) Beats(other *Card, trumpSuit string) bool {
    if this.Suit == trumpSuit && this.Suit != other.Suit {
        return true
    }
    return IndexOf(ranks, this.Rank) > IndexOf(ranks, other.Rank) && other.Suit != trumpSuit
}

func (card *Card) ToStr() string {
    return fmt.Sprintf("%s of %s", card.Rank, card.Suit)
}

func (game *Game) BoardSize() int {
    return len(game.Board.Plays) + len(game.Board.Covers)
}

func (game *Game) GetByRole(role string) *Player {
    for i,p := range game.Players {
        if role == "Defender" && i != game.Attacker {
            return p
        }
        if role == "Attacker" && i == game.Attacker {
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

func (game *Game) AttackerActions() []*Action {
    res := make([]*Action, 0)
    p := game.GetAttacker()
    if game.BoardSize() == 0 {
        for _,card := range p.Hand {
            act := Action{Player: p, Mode: "Attack", Card: card}
            res = append(res, &act)
        }
    } else {
        for _,bc := range Cat(game.Board.Plays, game.Board.Covers) {
            for _,pc := range p.Hand {
                if bc.Rank == pc.Rank && IndexOfFn(res, func(act *Action) bool {return act.Card == pc}) == -1 {
                    act := Action{Player: p, Mode: "Attack", Card: pc}
                    res = append(res, &act)
                }
            }
        }
    }
    if len(game.Board.Plays) > 0 {
        act := Action{Player: p, Mode: "Pass"}
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
                act := Action{Player: p, Mode: "Defend", Card: pc, Cover: bp}
                res = append(res, &act)
            }
        }
    }
    if len(game.Board.Plays) > 0 && len(game.Board.Covers) < len(game.Board.Plays) {
        act := Action{Player: p, Mode: "Pickup"}
        res = append(res, &act)
    }
    return res
}

func (act *Action) ToTensor(game *Game) T.Tensor {
    rankFeat := float64(IndexOf(ranks, act.Card.Rank))
    suitFeat := float64(Ternary(act.Card.Suit == game.Trump.Suit, 1, 0))
    boardSizeFeat := float64(game.BoardSize())
    handSizeFeat := float64(len(act.Player.Hand))
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

func InitPlayer() *Player {
    return &Player{Hand: make([]*Card, 0)}
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
    game := Game{Deck: InitDeck(), Board: InitBoard(), Turn: 0, Discard: make([]*Card, 0), Players: []*Player{InitPlayer(), InitPlayer()}}
    game.Trump = game.Deck[0]
    game.DealAll()
    return &game
}

/*func main() {
    game := InitGame()
    fmt.Println(game.ToStr())
}*/

