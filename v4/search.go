package main

import (
    "encoding/json"
    //"fmt"
    "math"
    "time"
    //clone "github.com/huandu/go-clone"
)

func ChainString(chain []*FastAction) string {
    cc := make([]*Action, len(chain))
    for i,a := range chain {
        cc[i] = a.ToAction()
    }
    jsn, _ := json.Marshal(cc) 
    return string(jsn) 
}

func HandString(h []*Card) string {
    jsn, _ := json.Marshal(h) 
    return string(jsn)
}

func FastHandString(h []int) string {
    hh := make([]*Card, len(h))
    for i,c := range h {
        hh[i] = &Card{Rank: ranks[c%9], Suit: suits[c/9]}
    }
    return HandString(hh)
}

func (a *Action) String() string {
    jsn, _ := json.Marshal(a)
    return string(jsn)
}

func (c *Card) ToFastCard() int {
    if c == nil {
        return -1
    }
    if c.Rank == "?" || c.Suit == "?" {
        return 36
    }
    return IndexOf(ranks, c.Rank) + 9*IndexOf(suits, c.Suit)
}

func FastCardToCard(c int) *Card {
    if c == -1 {
        return nil
    }
    if c == 36 {
        return &Card{Rank: "?", Suit: "?"}
    }
    return &Card{Rank: ranks[c%9], Suit: suits[c/9]}
}

func (a *FastAction) ToAction() *Action {
    if a == nil {
        return nil
    }
    aa := &Action{
        PlayerIdx: a.PlayerIdx, 
        Verb: verbs[a.Verb], 
        Card: FastCardToCard(a.Card),
        Cover: FastCardToCard(a.Cover),
    }
    return aa
}

func (a *FastAction) String() string {
    return a.ToAction().String()
}

var AttackV = 0
var DefendV = 1
var PassV = 2
var ReverseV = 3
var PickupV = 4
var DeferV = 5
var verbs = []string{"Attack", "Defend", "Pass", "Reverse", "Pickup", "Defer"}

type FastAction struct {
    PlayerIdx int
    Verb int 
    Card int
    Cover int
}

func (state *GameState) BoardSize() int {
    return len(state.CardsOnBoard()) 
}

func (state *GameState) CardsOnBoard() []int {
    return Cat(state.Plays, FastNotNil(state.Covers)) 
}

func (state *GameState) ReverseRank() int {
    if len(state.Plays) == 0 || FastNumNotNil(state.Covers) > 0 {
        return -1
    }
    r := state.Plays[0]%9
    for _,c := range state.CardsOnBoard() {
        if c%9 != r {
            return -1
        }
    }
    return r
}

func FastBeats(card1 int, card2 int, trumpSuit int) bool {
    if card1/9 == trumpSuit && card1/9 != card2/9 {
        return true
    }
    return card1%9 > card2%9 && card1/9 == card2/9
}

func FastNotNil(sl []int) []int {
    res := make([]int, 0)
    for _,s := range sl {
        if s != -1 {
            res = append(res, s)
        } 
    }
    return res
}

func FastNumNotNil(sl []int) int {
    n := 0
    for _,s := range sl {
        if s != -1 {
            n++
        } 
    }
    return n
}

type GameState struct {
    Defender int
    PickingUp bool
    Trump int           // Card
    DeckSize int
    Plays []int
    Covers []int
    Hands [][]int
    Defering bool
    Start *time.Time
}

func copyIntSlice(sl []int) []int {
    res := make([]int, len(sl))
    copy(res, sl)
    return res
}

func (state *GameState) Clone() *GameState {
    pCopy := make([]int, len(state.Plays))
    cCopy := make([]int, len(state.Covers))
    hCopy := make([][]int, len(state.Hands))
    hCopy[0] = make([]int, len(state.Hands[0]))
    hCopy[1] = make([]int, len(state.Hands[1]))
    copy(pCopy, state.Plays)
    copy(cCopy, state.Covers)
    copy(hCopy[0], state.Hands[0])
    copy(hCopy[1], state.Hands[1])
    return &GameState{
        Defender: state.Defender,
        PickingUp: state.PickingUp,
        Trump: state.Trump,
        DeckSize: state.DeckSize,
        Plays: pCopy,
        Covers: cCopy,
        Hands: hCopy, 
        Defering: state.Defering,
        Start: state.Start,
    }
}

/*func (state *GameState) Clone() *GameState {
    return clone.Clone(state).(*GameState) 
}*/

func InitGameState(game *Game) *GameState {
    pCopy := make([]int, len(game.Board.Plays))
    cCopy := make([]int, len(game.Board.Covers))
    hCopy := make([][]int, len(game.Players))
    hCopy[0] = make([]int, len(game.Players[0].Hand))
    hCopy[1] = make([]int, len(game.Players[1].Hand))
    for i,p := range game.Players {
        for j,c := range p.Hand {
            hCopy[i][j] = c.ToFastCard()
        }
    }
    for i,p := range game.Board.Plays {
        pCopy[i] = p.ToFastCard()
    }
    for i,c := range game.Board.Covers {
        cCopy[i] = c.ToFastCard()
    }
    return &GameState{
        Defender: game.Defender,
        PickingUp: game.PickingUp,
        Trump: game.Trump.ToFastCard(),
        DeckSize: len(game.Deck),
        Plays: pCopy,
        Covers: cCopy,
        Hands: hCopy, 
        Defering: false,
    }
}

func (state *GameState) AttackerActions(pIdx int) []*FastAction {
    res := make([]*FastAction, 0)
    if state.Defering {
        return res
    }
    if state.BoardSize() == 0 {
        for _,card := range state.Hands[pIdx] {
            act := FastAction{PlayerIdx: pIdx, Verb: AttackV, Card: card, Cover: -1}
            res = append(res, &act)
        }
    } else {
        for _,pc := range state.Hands[pIdx] {
            for _,bc := range state.CardsOnBoard() {
                if bc%9 == pc%9 || pc == 36 {
                    // Unique actions
                    act := FastAction{PlayerIdx: pIdx, Verb: AttackV, Card: pc, Cover: -1}
                    res = append(res, &act)
                    break
                }
            }
        }
    }
    if state.PickingUp || (FastNumNotNil(state.Covers) == len(state.Plays) && len(state.Plays) > 0) {
        act := FastAction{PlayerIdx: pIdx, Verb: PassV, Card: -1, Cover: -1}
        res = append(res, &act)
    }
    // For AI to not throw trumps away
    if len(state.Plays) > FastNumNotNil(state.Covers) {
        act := FastAction{PlayerIdx: pIdx, Verb: DeferV, Card: -1, Cover: -1}
        res = append(res, &act)
    }
    return res
}

func (state *GameState) DefenderActions(pIdx int) []*FastAction {
    res := make([]*FastAction, 0)
    if state.PickingUp {
        return res
    }
    revRank := state.ReverseRank()
    if revRank != -1 {
        for _,pc := range state.Hands[pIdx] {
            if pc%9 == revRank {
                act := FastAction{PlayerIdx: pIdx, Verb: ReverseV, Card: pc, Cover: -1}
                res = append(res, &act)
            }
        }
    }
    for i,bp := range state.Plays {
        if state.Covers[i] != -1 {
            continue
        }
        for _,pc := range state.Hands[pIdx] {
            if FastBeats(pc, bp, state.Trump/9) || pc == 36 {
                act := FastAction{PlayerIdx: pIdx, Verb: DefendV, Card: pc, Cover: bp}
                res = append(res, &act)
            }
        }
    }
    // Get non-nil covers
    if len(state.Plays) > 0 && FastNumNotNil(state.Covers) < len(state.Plays) {
        act := FastAction{PlayerIdx: pIdx, Verb: PickupV, Card: -1, Cover: -1}
        res = append(res, &act)
    }
    return res
}

func StartChain(act *FastAction) []*FastAction {
    return []*FastAction{act}
}

// These are sometimes not representative of how long it takes to eval
// That's why we also measure time
// TODO iterative deepening
func (state *GameState) DepthLimit() int {
    nCards := len(Cat(state.Hands[0], state.Hands[1]))
    if nCards > 18 {
        return 5
    } else if nCards > 12 {
        return 6
    } else if nCards > 10 {
        return 7
    } else if nCards > 8{
        return 8
    } else {
        return 12
    }
}

func (state *GameState) Move(me int, depth int, dlim int) ([]*FastAction,float64) {
    if depth == 0 {
        dlim = state.DepthLimit()
        start := time.Now()
        state.Start = &start
    }
    dlimAdjusted := dlim
    tDiff := time.Now().Sub(*state.Start)
    if tDiff.Seconds() > 2 {
        dlimAdjusted -= 1
    }
    if depth > dlimAdjusted {
        return StartChain(nil), state.EvalMystery(me)
    }
    var acts []*FastAction
    if me == state.Defender {
        acts = state.DefenderActions(me)
    } else {
        acts = state.AttackerActions(me)
    }
    // If you have no actions, return
    if len(acts) == 0 {
        return nil, 0 
    }
    evals := make([]float64, 2*len(acts))
    chains := make([][]*FastAction, 2*len(acts))
    endGame := state.DeckSize <= 1
    didMystery := false
    for i,act := range acts {
        if act.Verb == DeferV {
            // If end of game, treat as pass and keep going
            if endGame {
                state.Defering = true
            } else {
                return StartChain(act), state.EvalPass(me) 
            }
        }
        s := state.Clone()
        s.TakeAction(act, endGame)
        // Check win
        // Infinite value can confuse later action selection?
        if s.DeckSize == 0 && len(s.Hands[me]) == 0 {
            return StartChain(act), 1000
        }
        // End hand
        if act.Verb == PassV {
            if endGame {
            // Go to end of game
                if state.PickingUp {
                    c,r := s.Move(me, depth+1, dlim)
                    evals[2*i] = r
                    chains[2*i] = append(c, act)
                    chains[2*i+1] = nil
                } else {
                    c,r := s.Move(1-me, depth+1, dlim)
                    evals[2*i+1] = r
                    chains[2*i+1] = append(c, act)
                    chains[2*i] = nil
                }
            } else {
            // Go per-hand
                return StartChain(act), s.EvalPass(me) 
            }
        }
        // Mystery card
        // Only check one mystery card
        // TODO incorporate possibility of reverse
        if !didMystery && act.Card == 36 {
            //return StartChain(act), s.EvalMystery(me)
            chains[2*i] = StartChain(act)
            evals[2*i] = s.EvalMystery(me)
            chains[2*i+1] = nil
            didMystery = true
        // Pickup - Opponent's move will determine evaluation
        // Penalize high hand count
        // Penalize taking cards with zero deck size (end of game)
        } else if act.Verb == PickupV {
            chains[2*i] = nil
            c,r := s.Move(1-me, depth+1, dlim)
            evals[2*i+1] = s.PickupPenalty(me) + r
            chains[2*i+1] = Ternary(c == nil, nil, append(c, act))
        // Regular known move
        } else {
            c,r := s.Move(me, depth+1, dlim)
            evals[2*i] = r
            chains[2*i] = Ternary(c == nil, nil, append(c, act))
            c,r = s.Move(1-me, depth+1, dlim)
            evals[2*i+1] = r
            chains[2*i+1] = Ternary(c == nil, nil, append(c, act))
        }
    }
    // Find likely action
    var bestChain []*FastAction
    best := math.Inf(-1)
    for i,_ := range acts {
        if chains[2*i] == nil && chains[2*i+1] == nil {
            continue
        }
        if chains[2*i] == nil {
            e := -evals[2*i+1]
            if e > best {
                bestChain = chains[2*i+1]
                best = e
            }
        } else if chains[2*i+1] == nil {
            e := evals[2*i]
            if e > best {
                bestChain = chains[2*i]
                best = e
            }
        } else {
            e1 := evals[2*i]
            e2 := -evals[2*i+1]
            if e1 > best {
                bestChain = chains[2*i]
                best = e1
            }
            if e2 > best {
                bestChain = chains[2*i+1]
                best = e2
            }
        }
    }
    return bestChain, best
}

// EndGame means we need to clear the board
func (state *GameState) TakeAction(act *FastAction, endGame bool) {
    if state.Defering && act.PlayerIdx == state.Defender {
        state.Defering = false
    }
    switch act.Verb {
        case AttackV: {
            state.Hands[act.PlayerIdx] = Remove(state.Hands[act.PlayerIdx], act.Card)
            state.Plays = append(state.Plays, act.Card)
            state.Covers = append(state.Covers, -1)
        } 
        case DefendV: {
            state.Hands[act.PlayerIdx] = Remove(state.Hands[act.PlayerIdx], act.Card)
            idx := IndexOf(state.Plays, act.Cover)
            state.Covers[idx] = act.Card
        }
        case PickupV: {
            state.PickingUp = true
        }
        case ReverseV: {
            state.Hands[act.PlayerIdx] = Remove(state.Hands[act.PlayerIdx], act.Card)
            state.Plays = append(state.Plays, act.Card)
            state.Covers = append(state.Covers, -1)
            state.Defender = 1-state.Defender
        }
        case PassV: {
            // Skip, handled in Move code
            // ...unless it's the endgame
            if endGame {
                state.Plays = make([]int, 0)
                state.Covers = make([]int, 0)
                if !state.PickingUp {
                    state.Defender = 1-state.Defender
                }
                state.PickingUp = false
            }
        }
    }
}

func (state *GameState) SumValue(cards []int) float64 {
    res := float64(0)
    for _,c := range cards {
        if c != -1 && c != 36 {
            res += float64(c%9 - 4)
            if c/9 == state.Trump/9 {
                res += 9
            }
        }
    }
    return res
}

// Simple evaluation function
func (state *GameState) EvalPass(me int) float64 {
    var val float64
    if state.PickingUp {
        val = state.SumValue(Cat(state.Plays, state.Covers)) - 4
        if me != state.Defender {
            val *= -1
        }
    } else {
        return state.EvalMystery(me)
    }
    return val
}

func (state *GameState) EvalMystery(me int) float64 {
    if me == state.Defender {
        return state.SumValue(state.Plays) - state.SumValue(state.Covers)
    } else {
        return state.SumValue(state.Covers) - state.SumValue(state.Plays)
    }
}

func (state *GameState) PickupPenalty(me int) float64 {
    val := 0
    if state.DeckSize < 6 {
        val += 6 - state.DeckSize
    }
    if len(state.Hands[me]) > 6 {
        val += len(state.Hands[me])-6
    }
    val += FastNumNotNil(Cat(state.Plays, state.Covers))
    return float64(val)
}
