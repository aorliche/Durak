package main

import (
    "encoding/json"
    "fmt"
    "math/rand"
    "sync"
    "time"
    
    "github.com/gorilla/websocket"
)

type Verb int
type Card int

var UNK_CARD = Card(-1)

const (
    PlayVerb Verb = iota
    CoverVerb 
    ReverseVerb
    PassVerb
    PickupVerb
    DeferVerb
)

var suits = []string{"Clubs", "Spades", "Hearts", "Diamonds"}
var ranks = []string{"6", "7", "8", "9", "10", "Jack", "Queen", "King", "Ace"}
var verbs = []string{"Play", "Cover", "Reverse", "Pass", "Pickup", "Defer"}

func CardFromRankSuit(rank int, suit int) Card {
    return Card(suit*9 + rank)
}

func (card Card) Rank() int {
    return int(card)%9
}

func (card Card) Suit() int {
    return int(card)/9
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
    if card.Suit() == trump.Suit() && card.Suit() != other.Suit() {
        return true
    }
    return card.Rank() > other.Rank() && card.Suit() == other.Suit()
}

func GenerateDeck() []Card {
    res := make([]Card, 0)
    for suit := 0; suit < 4; suit++ {
        for rank := 0; rank < 9; rank++ {
            res = append(res, CardFromRankSuit(rank, suit))
        }
    }
    rand.Shuffle(len(res), func(i, j int) {
        res[i], res[j] = res[j], res[i]
    })
    return res
}

func RemoveCard(cards *[]Card, c Card) bool {
    for i,card := range *cards {
        if card == c {
            (*cards)[i] = (*cards)[len(*cards)-1]
            *cards = (*cards)[:len(*cards)-1]
            return true
        }
    }
    return false
}

type Action struct {
    Player int
    Verb Verb
    Card Card
    Covering Card
}

func (a Action) IsNull() bool {
    return a.Card == 0 && a.Covering == 0
}

func (a Action) ToStr() string {
    mp := map[string]any {
        "Player": a.Player,
        "Verb": verbs[a.Verb],
        "Card": a.Card,
        "Covering": a.Covering,
    }
    jsn, _ := json.Marshal(mp)
    return string(jsn)
}

type GameState struct {
    Defender int
    PickingUp bool
    Deferring bool
    Trump Card
    Plays []Card
    Covers []Card
    Hands [][]Card
    start time.Time
}

func NumNotUnk(cards []Card) int {
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
        start: time.Now(),
    }
}
    
func (state *GameState) AttackerActions(player int) []Action {
    res := make([]Action, 0)
    if len(state.Plays) == 0 {
        for _,card := range state.Hands[player] {
            res = append(res, Action{player, PlayVerb, card, UNK_CARD})
        }
        return res
    }
    for _,card := range state.Hands[player] {
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
    if state.PickingUp || (NumNotUnk(state.Covers) == len(state.Plays) && len(state.Plays) > 0) {
        res = append(res, Action{player, PassVerb, UNK_CARD, UNK_CARD})
    }
    // For AI to not throw trumps away
    if !state.PickingUp && len(state.Plays) > NumNotUnk(state.Covers) {
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

func (state *GameState) DefenderActions(player int) []Action {
    res := make([]Action, 0)
    if state.PickingUp {
        return res
    }
    revRank := state.ReverseRank()
    if revRank != -1 && NumNotUnk(state.Covers) == 0 {
        for _,card := range state.Hands[player] {
            if card.Rank() == revRank {
                res = append(res, Action{player, ReverseVerb, card, UNK_CARD})
            }
        }
    }
    for i := 0; i < len(state.Plays); i++ {
        for _,card := range state.Hands[player] {
            // For AI search allow cover with unknown card
            if card == UNK_CARD || (state.Covers[i] == UNK_CARD && card.Beats(state.Plays[i], state.Trump)) {
                res = append(res, Action{player, CoverVerb, card, state.Plays[i]})
            }
        }
    }
    if len(state.Plays) > 0 && NumNotUnk(state.Covers) < len(state.Plays) {
        res = append(res, Action{player, PickupVerb, UNK_CARD, UNK_CARD})
    }
    return res
}

func (state *GameState) PlayerActions(player int) []Action {
    if player == state.Defender {
        return state.DefenderActions(player)
    } else {
        return state.AttackerActions(player)
    }
}

func (state *GameState) RandomAction() Action {
    acts := append(make([]Action, 0), state.PlayerActions(0)...)
    acts = append(acts, state.PlayerActions(1)...)
    return acts[rand.Intn(len(acts))]
}

func (state *GameState) TakeAction(action Action) {
    switch action.Verb {
        case PlayVerb: {
            state.Plays = append(state.Plays, action.Card)
            state.Covers = append(state.Covers, UNK_CARD)
            RemoveCard(&state.Hands[action.Player], action.Card)
            state.Deferring = false
        }
        case CoverVerb: {
            for i := 0; i < len(state.Plays); i++ {
                if action.Covering == state.Plays[i] {
                    state.Covers[i] = action.Card
                }
            }
            RemoveCard(&state.Hands[action.Player], action.Card)
            state.Deferring = false
        }
        case ReverseVerb: {
            state.Plays = append(state.Plays, action.Card)
            state.Covers = append(state.Covers, UNK_CARD)
            RemoveCard(&state.Hands[action.Player], action.Card)
            state.Defender = 1-state.Defender
            state.Deferring = false
        }
        case PickupVerb: {
            state.PickingUp = true
            state.Deferring = false
        }
        case PassVerb: {
            if state.PickingUp {
                for i := 0; i < len(state.Plays); i++ {
                    state.Hands[state.Defender] = append(state.Hands[state.Defender], state.Plays[i])
                    if state.Covers[i] != UNK_CARD {
                        state.Hands[state.Defender] = append(state.Hands[state.Defender], state.Covers[i])
                    }
                }
            } else {
                state.Defender = 1-state.Defender
            }
            state.Plays = make([]Card, 0)
            state.Covers = make([]Card, 0)
            state.PickingUp = false
            state.Deferring = false
        }
        case DeferVerb: {
            state.Deferring = true
        }
    }
}

func (state *GameState) Clone() *GameState {
    hands := make([][]Card, 2)
    hands[0] = append(make([]Card, 0), state.Hands[0]...)
    hands[1] = append(make([]Card, 0), state.Hands[1]...)
    return &GameState {
        Defender: state.Defender,
        PickingUp: state.PickingUp,
        Deferring: state.Deferring,
        Trump: state.Trump,
        Plays: append(make([]Card, 0), state.Plays...),
        Covers: append(make([]Card, 0), state.Covers...),
        Hands: hands,
        start: state.start,
    }
}

func (state *GameState) ToStr() string {
    jsn, _ := json.Marshal(state)
    return string(jsn)
}

type Memory struct {
    Hands [][]Card
    Sizes []int
    Discard []Card
}

type Game struct {
    Key int
    State *GameState
    Deck []Card
    Memory *Memory
    Recording *Recording
    Versus string
    joined bool
    mutex sync.Mutex
    conns []*websocket.Conn
}

type Recording struct {
    Versus string
    Deck []Card
    Hands [][]Card
    Actions []Action
    Winner int
}

func InitGame(key int, versus string) *Game {
    deck := GenerateDeck()
    h0 := append(make([]Card, 0), deck[0:6]...)
    h1 := append(make([]Card, 0), deck[6:12]...)
    deck = append(make([]Card, 0), deck[12:]...)
    recording := &Recording{
        Versus: versus,
        Deck: append(make([]Card, 0), deck...),
        Hands: [][]Card{
            append(make([]Card, 0), h0...), 
            append(make([]Card, 0), h1...),
        },
        Actions: make([]Action, 0),
        Winner: -1,
    }
    return &Game{
        Key: key, 
        State: InitGameState(deck[0], [][]Card{h0, h1}),
        Deck: deck,
        Memory: &Memory {
            Hands: [][]Card{make([]Card, 0), make([]Card, 0)},
            Sizes: make([]int, 2),
            Discard: make([]Card, 0),
        },
        Recording: recording,
        Versus: versus,
        joined: false,
        conns: make([]*websocket.Conn, 2),
    }
}

func (game *Game) CheckWinner() int {
    if game.Recording.Winner != -1 {
        return game.Recording.Winner
    }
    if len(game.Deck) > 0 {
        return -1
    }
    for i := 0; i < len(game.State.Hands); i++ {
        if len(game.State.Hands[i]) == 0 {
            game.Recording.Winner = i
            return i
        }
    }
    return -1
}

func (game *Game) Deal(player int) {
   for len(game.Deck) > 0 && len(game.State.Hands[player]) < 6 {
       game.State.Hands[player] = append(game.State.Hands[player], game.Deck[len(game.Deck)-1])
       game.Deck = game.Deck[:len(game.Deck)-1]
   } 
   game.Memory.Sizes[player] = len(game.State.Hands[player])
}

func (game *Game) TakeAction(action Action) bool {
    // Check that action is still legal
    acts := game.State.PlayerActions(action.Player)
    found := false
    for _, act := range acts {
        if act == action {
            found = true
        }
    }
    if !found {
        return false
    }
    // No strings of Defer verbs in recordings
    if action.Verb == DeferVerb {
        if len(game.Recording.Actions) > 0 && game.Recording.Actions[len(game.Recording.Actions)-1].Verb != DeferVerb {
            game.Recording.Actions = append(game.Recording.Actions, action)
        }
    } else {
        game.Recording.Actions = append(game.Recording.Actions, action)
    }
    switch action.Verb {
        case PlayVerb, CoverVerb, ReverseVerb: {
            game.State.TakeAction(action);
            RemoveCard(&game.Memory.Hands[action.Player], action.Card)
            game.Memory.Sizes[action.Player] -= 1
        }
        case PassVerb: {
            if game.State.PickingUp {
                for i := 0; i < len(game.State.Plays); i++ {
                    game.Memory.Hands[game.State.Defender] = append(game.Memory.Hands[game.State.Defender], game.State.Plays[i])
                    game.Memory.Sizes[game.State.Defender] += 1
                    if game.State.Covers[i] != UNK_CARD {
                        game.Memory.Hands[game.State.Defender] = append(game.Memory.Hands[game.State.Defender], game.State.Covers[i])
                        game.Memory.Sizes[game.State.Defender] += 1
                    }
                }
            } else {
                for i := 0; i < len(game.State.Plays); i++ {
                    game.Memory.Discard = append(game.Memory.Discard, game.State.Plays[i])
                    game.Memory.Discard = append(game.Memory.Discard, game.State.Covers[i])
                }
            }
            game.State.TakeAction(action);
            game.Deal(1-game.State.Defender)
            game.Deal(game.State.Defender)
        }
        case PickupVerb, DeferVerb: {
            game.State.TakeAction(action);
        }
    }
    return true
}
