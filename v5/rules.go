package main

import (
    "encoding/json"
    "fmt"
    //"log"
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
    Attacker int
    Defender int
    PickingUp bool
    Deferring []bool
    Passed []bool
    Won []bool
    Trump Card
    Plays []Card
    Covers []Card
    Hands [][]Card
    Dir int
    start time.Time
    gamePtr *Game
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
        Attacker: 0,
        Defender: 1, 
        PickingUp: false, 
        Deferring: make([]bool, len(hands)),
        Passed: make([]bool, len(hands)),
        Won: make([]bool, len(hands)),
        Trump: trump,
        Plays: make([]Card, 0),
        Covers: make([]Card, 0),
        Hands: hands,
        Dir: 1,
        start: time.Now(),
        gamePtr: nil,
    }
}
    
func (state *GameState) AttackerActions(player int) []Action {
    res := make([]Action, 0)
    // Player has already passed
    if state.Passed[player] {
        return res
    }
    // Player has already won
    if state.Won[player] {
        return []Action{Action{player, PassVerb, UNK_CARD, UNK_CARD}}
    }
    if len(state.Plays) == 0 {
        // Only initial attacker may play first
        if state.Attacker != player {
            return res
        }
        for _,card := range state.Hands[player] {
            res = append(res, Action{player, PlayVerb, card, UNK_CARD})
        }
        return res
    }
    // Don't allow playing more cards than defender can defend
    unmetCards := len(state.Plays) - NumNotUnk(state.Covers)
    if len(state.Hands[state.Defender]) - unmetCards > 0 {
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
    // Only allow reverse when defender can potentially meet it
    if revRank != -1 && NumNotUnk(state.Covers) == 0 && len(state.Plays)+1 <= len(state.Hands[state.Attacker]) {
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

func IndexOf[T comparable](slice []T, val T) int {
    for i, v := range slice {
        if v == val {
            return i
        }
    }
    return -1
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

func (state *GameState) NextRole(player int) int {
    for i := 0; i < len(state.Hands); i++ {
        p := player+((i+1)*state.Dir)
        if p < 0 {
            p += len(state.Hands)
        }
        if p >= len(state.Hands) {
            p -= len(state.Hands)
        }
        if !state.Won[p] {
            return p
        }
    }
    panic("NextRole failed")
}

func (state *GameState) AllPassed() bool {
    for i := 0; i < len(state.Hands); i++ {
        if !state.Passed[i] && state.Defender != i && !state.Won[i] {
            return false
        }
    }
    return true
}

func (state *GameState) CheckGameOver() bool {
    if len(state.gamePtr.Deck) > 0 {
        return false
    }
    n := 0
    for i := 0; i < len(state.Hands); i++ {
        if len(state.Hands[i]) == 0 {
            state.Won[i] = true
            n++
        }
    }
    return n == len(state.Hands)-1
}

func (state *GameState) TakeAction(action Action) {
    switch action.Verb {
        case PlayVerb: {
            state.Plays = append(state.Plays, action.Card)
            state.Covers = append(state.Covers, UNK_CARD)
            RemoveCard(&state.Hands[action.Player], action.Card)
            // Sets to false
            state.Deferring = make([]bool, len(state.Hands))
            state.CheckGameOver()
        }
        case CoverVerb: {
            for i := 0; i < len(state.Plays); i++ {
                if action.Covering == state.Plays[i] {
                    state.Covers[i] = action.Card
                }
            }
            RemoveCard(&state.Hands[action.Player], action.Card)
            state.Deferring = make([]bool, len(state.Hands))
            state.CheckGameOver()
        }
        case ReverseVerb: {
            state.Plays = append(state.Plays, action.Card)
            state.Covers = append(state.Covers, UNK_CARD)
            RemoveCard(&state.Hands[action.Player], action.Card)
            state.Attacker, state.Defender = state.Defender, state.Attacker
            state.Deferring = make([]bool, len(state.Hands))
            state.Dir *= -1
            state.CheckGameOver()
        }
        case PickupVerb: {
            state.PickingUp = true
            state.Deferring = make([]bool, len(state.Hands))
        }
        case PassVerb: {
            // Handled in Game.TakeAction 
            // Need to be in both places for Medium search
            state.Passed[action.Player] = true
            if !state.AllPassed() {
                break
            }
            if state.PickingUp {
                for i := 0; i < len(state.Plays); i++ {
                    state.Hands[state.Defender] = append(state.Hands[state.Defender], state.Plays[i])
                    if state.Covers[i] != UNK_CARD {
                        state.Hands[state.Defender] = append(state.Hands[state.Defender], state.Covers[i])
                    }
                }
                state.Attacker = state.NextRole(state.Defender) 
                state.Defender = state.NextRole(state.Attacker)
            } else {
                state.Attacker = state.NextRole(state.Attacker) 
                state.Defender = state.NextRole(state.Attacker)
            }
            state.Plays = make([]Card, 0)
            state.Covers = make([]Card, 0)
            state.PickingUp = false
            state.Deferring = make([]bool, len(state.Hands))
            state.Passed = make([]bool, len(state.Hands))
        }
        case DeferVerb: {
            state.Deferring[action.Player] = true
        }
    }
}

func (state *GameState) Clone() *GameState {
    hands := make([][]Card, len(state.Hands))
    for i := 0; i < len(hands); i++ {
        hands[i] = append(make([]Card, 0), state.Hands[i]...)
    }
    return &GameState {
        Defender: state.Defender,
        Attacker: state.Attacker,
        PickingUp: state.PickingUp,
        Deferring: append(make([]bool, 0), state.Deferring...),
        Passed: append(make([]bool, 0), state.Passed...),
        Won: append(make([]bool, 0), state.Won...),
        Trump: state.Trump,
        Plays: append(make([]Card, 0), state.Plays...),
        Covers: append(make([]Card, 0), state.Covers...),
        Hands: hands,
        Dir: state.Dir,
        gamePtr: state.gamePtr,
        start: state.start,
    }
}

func (state *GameState) ToStr() string {
    jsn, _ := json.Marshal(state.Clone())
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
    Players []string
    Names []string
    joined []bool
    mutex sync.Mutex
    conns []*websocket.Conn
}

type Recording struct {
    Players []string
    Deck []Card
    Hands [][]Card
    Actions []Action
    Winners []int
}

func InitGame(key int, players []string) *Game {
    deck := GenerateDeck()
    numPlayers := len(players)
    handsState := make([][]Card, numPlayers)
    handsRec := make([][]Card, numPlayers)
    handsMemory := make([][]Card, numPlayers)
    for i := 0; i < numPlayers; i++ {
        handsState[i] = append(make([]Card, 0), deck[i*6:(i+1)*6]...)
        handsRec[i] = append(make([]Card, 0), deck[i*6:(i+1)*6]...)
        handsMemory[i] = make([]Card, 0)
    }
    deck = append(make([]Card, 0), deck[numPlayers*6:]...)
    recording := &Recording{
        Players: players,
        Deck: append(make([]Card, 0), deck...),
        Hands: handsRec,
        Actions: make([]Action, 0),
        Winners: make([]int, 0),
    }
    game := &Game{
        Key: key,
        State: InitGameState(deck[0], handsState),
        Deck: deck,
        Memory: &Memory {
            Hands: handsMemory,
            Sizes: make([]int, numPlayers),
            Discard: make([]Card, 0),
        },
        Recording: recording,
        Players: players,
        Names: make([]string, numPlayers),
        joined: make([]bool, numPlayers),
        conns: make([]*websocket.Conn, numPlayers),
    }
    game.State.gamePtr = game
    return game
}

func (game *Game) CheckGameOver() bool {
    if len(game.Deck) > 0 {
        return false
    }
    if len(game.Recording.Winners) == len(game.Players)-1 {
        return true
    }
    for i := 0; i < len(game.State.Hands); i++ {
        if len(game.State.Hands[i]) == 0 {
            if IndexOf(game.Recording.Winners, i) == -1 {
                game.Recording.Winners = append(game.Recording.Winners, i)
            }
        }
    }
    return len(game.Recording.Winners) == len(game.Players)-1 
}

// TODO Deal in correct order
func (game *Game) Deal(player int) {
   for len(game.Deck) > 0 && len(game.State.Hands[player]) < 6 {
       game.State.Hands[player] = append(game.State.Hands[player], game.Deck[len(game.Deck)-1])
       game.Deck = game.Deck[:len(game.Deck)-1]
   } 
   game.Memory.Sizes[player] = len(game.State.Hands[player])
   // If some player has zero cards and we don't deal them any
   // Adjust winners, attackers, defenders
   game.CheckGameOver()
   if IndexOf(game.Recording.Winners, game.State.Attacker) != -1 {
       game.State.Attacker = game.State.NextRole(game.State.Attacker)
       game.State.Defender = game.State.NextRole(game.State.Attacker)
   }
   if IndexOf(game.Recording.Winners, game.State.Defender) != -1 {
       game.State.Defender = game.State.NextRole(game.State.Defender)
       game.State.Attacker = game.State.NextRole(game.State.Defender)
   }
}

func (game *Game) TakeAction(action Action) bool {
    // Check that game isn't over
    if game.CheckGameOver() {
        return false
    }
    // Check that action is still legal
    acts := game.State.PlayerActions(action.Player)
    if IndexOf(acts, action) == -1 {
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
            game.State.Passed[action.Player] = true
            if !game.State.AllPassed() {
                break
            }
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
            for i := 0; i < len(game.State.Hands); i++ {
                game.Deal(i)
            }
        }
        case PickupVerb, DeferVerb: {
            game.State.TakeAction(action);
        }
    }
    return true
}
