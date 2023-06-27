package main

import (
    "fmt"
    "rand"
    "reflect"
    "strings"
    "time"
)

type Verb int
type Card int

var UNK_CARD = Card(-1)

const (
    PlayV int = iota
    CoverV 
    PassV
    ReverseV
    PickupV
    DeferV
)

var suits = []string{"Clubs", "Spades", "Hearts", "Diamonds"}
var ranks = []string{"6", "7", "8", "9", "10", "Jack", "Queen", "King", "Ace"}

func CardFromRankSuit(rank int, suit int) Card {
    return suit*9 + rank 
}

func (card Card) Rank() int {
    return card%9
}

func (card Card) Suit() int {
    return card/9
}

func (card Card) RankStr() string {
    return ranks[card.Rank()]
}

func (card Card) SuitStr() string {
    return suits[card.Suit()]
}

func (card Card) ToStr() string {
    return fmt.Sprintf("%s of %s", card.RankStr(), card.SuitStr()) 
}

func (card Card) Beats(other Card, trump Card) bool {
    if CardSuit(card) == CardSuit(trump) && CardSuit(card) != CardSuit(other) {
        return true
    }
    return CardRank(card) > CardRank(other) && CardSuit(card) == CardSuit(other)
}

func GenerateDeck() []Card {
    res := make([]Card, 0)
    for _,rank := range ranks {
        for _,suit := range suits {
            res = append(res, Card{Rank: rank, Suit: suit})
        }
    }
    rand.Shuffle(len(res), func(i, j int) {
        res[i], res[j] = res[j], res[i]
    })
    return res
}

func RemoveCard(cards *[]Card, c Card) bool {
    for i,card := range cards {
        if card == c {
            cards[i] = cards[len(cards)-1]
            *cards = cards[:len(cards)-1]
            return true
        }
    }
    return false
}

type Action struct {
    Player int,
    Verb Verb,
    Card Card,
    Covering Card,
}

type GameState struct {
    Defender int,
    PickingUp bool,
    Deferring bool,
    Trump Card,
    Plays []Card,
    Covers []Card,
    Hands [][]Card,
    Start time.Time,
}

func NumNotUnk(cards *[]Card) int {
    res := 0
    for _,c := range cards {
        if c != UNK_CARD {
            res++
        }
    }
    return res
}

func InitGameState(trump Card, hands [][]Card) *GameState {
    return &GameState{
        Defender: 1, 
        PickingUp: false, 
        Deferring: false,
        Trump: trump,
        Plays: make([]Card, 0),
        Covers: make([]Card, 0),
        Hands: hands,
        Start: nil,
    }
}
    
func (state *GameState) AttackerActions(player int) []Action {
    res := make([]Action, 0)
    if len(state.Plays) == 0 {
        for _,card := range state.hands[player] {
            res = append(res, Action{player, PlayVerb, card, UNK_CARD})
        }
        return res
    }
    for _,card := range state.hands[player] {
        // Allow play unknown card in search
        if card == UNK_CARD {
            res = append(res, Action{player, PlayVerb, card, UNK_CARD})
            continue
        }
        // Regular cards
        for i := 0; i < len(state.Plays); i++ {
            if card.Rank() == state.Plays[i].Rank() || (card != UNK_CARD && card.Rank() == state.Covers[i].Rank()) {
                res = append(res, Action{player, PlayVerb, card, UNK_CARD})
                break
            }
        }
    }
    if state.PickingUp || (NumNotUnk(&state.Covers) == len(state.Plays) && len(state.Plays) > 0) {
        res = append(res, Action{player, PassVerb, UNK_CARD, UNK_CARD})
    }
    // For AI to not throw trumps away
    if !state.PickingUp && len(state.Plays) > NumNotUnk(&state.Covers) {
        res = append(res, Action{player, DeferVerb, UNK_CARD, UNK_CARD})    
    }
    return res
}

func (state *GameState) ReverseRank() int {
    if len(state.Plays) == 0 {
        return -1
    }
    rank := state.Plays[0].Rank()
    for i := 0; i < len(state.Plays); i++ {
        if rank != state.Plays[i].Rank() {
            return -1
        }
        if state.Covers[i] != UNK_CARD && rank != state.Covers[i].Rank() {
            return -1
        }
    }
    return rank
}
    pub fn defender_actions(&self, pidx: usize) -> Vec<Action> {
        let mut res = Vec::new();
        if self.picking_up {
            return res
        }
        let rev_rank = self.reverse_rank();
        match rev_rank {
            Some(rank) => {
                for &card in &self.hands[pidx] {
                    if card_rank(card) == rank {
                        res.push(Action::new(pidx, Verb::Reverse, card, UNK_CARD))
                    }
                }
            },
            None => (),
        }
        for i in 0..self.plays.len() {
            for &card in self.hands[pidx].iter() {
                if self.covers[i] == UNK_CARD && beats(card, self.plays[i], self.trump) {
                    res.push(Action::new(pidx, Verb::Cover, card, self.plays[i]))
                }
                // For AI search allow cover with unknown card
                if card == UNK_CARD {
                    res.push(Action::new(pidx, Verb::Cover, card, self.plays[i]))
                }
            }
        }
        if self.plays.len() > 0 && self.covers.iter().filter(|&x| *x != UNK_CARD).count() < self.plays.len() {
            res.push(Action::new(pidx, Verb::PickUp, UNK_CARD, UNK_CARD))
        }
        res
    }
    pub fn reverse_rank(&self) -> Option<usize> {
        if self.plays.len() == 0 && self.covers.len() == 0 {
            return None
        }
        let mut card = self.plays[0];
        for &c in self.plays.iter().chain(self.covers.iter()) {
            if c != UNK_CARD && c != card {
                return None
            }
        }
        return Some(card_rank(card))
    }
        if self.plays.len() == 0 && self.covers.len() == 0 {
            return None
        }
        let mut card = self.plays[0];
        for &c in self.plays.iter().chain(self.covers.iter()) {
            if c != UNK_CARD && c != card {
                return None
            }
        }
        return Some(card_rank(card))

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

