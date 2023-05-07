package main

import (
    "math"
    clone "github.com/huandu/go-clone"
)

func (b *Board) Size() int {
    return len(b.Cards()) 
}

func (b *Board) Cards() []*Card {
    return Cat(b.Plays, NotNil(b.Covers)) 
}

func (b *Board) ReverseRank() string {
    if len(b.Plays) == 0 {
        return ""
    }
    r := b.Plays[0].Rank
    for _,c := range b.Cards() {
        if c.Rank != r {
            return ""
        }
    }
    return r
}

type GameState struct {
    Defender int
    PickingUp bool
    Trump string
    DeckSize int
    Board *Board
    Hands [][]*Card
}

func (state *GameState) Clone() *GameState {
    return clone.Clone(state).(*GameState) 
}

func InitGameState(game *Game) *GameState {
    hands := make([][]*Card, len(game.Players))
    for i,p := range game.Players {
        hands[i] = p.Hand
    }
    return &GameState{
        Defender: game.Defender,
        PickingUp: false,
        Trump: game.Trump.Suit,
        DeckSize: len(game.Deck),
        Board: clone.Clone(game.Board).(*Board),
        Hands: hands,
    }
}

func (state *GameState) AttackerActions(pIdx int) []*Action {
    res := make([]*Action, 0)
    if state.Board.Size() == 0 {
        for _,card := range state.Hands[pIdx] {
            act := Action{PlayerIdx: pIdx, Verb: "Attack", Card: card}
            res = append(res, &act)
        }
    } else {
        for _,bc := range state.Board.Cards() {
            for _,pc := range state.Hands[pIdx] {
                // Unique actions
                if bc != nil && bc.Rank == pc.Rank && IndexOfFn(res, func(act *Action) bool {return act.Card == pc}) == -1 {
                    act := Action{PlayerIdx: pIdx, Verb: "Attack", Card: pc}
                    res = append(res, &act)
                }
            }
        }
    }
    if state.PickingUp || (state.Board.Covered() == len(state.Board.Plays) && len(state.Board.Plays) > 0) {
        act := Action{PlayerIdx: pIdx, Verb: "Pass"}
        res = append(res, &act)
    }
    return res
}

func (state *GameState) DefenderActions(pIdx int) []*Action {
    res := make([]*Action, 0)
    revRank := state.Board.ReverseRank()
    if revRank != "" {
        for _,pc := range state.Hands[pIdx] {
            if pc.Rank == revRank {
                act := Action{PlayerIdx: pIdx, Verb: "Reverse", Card: pc}
                res = append(res, &act)
            }
        }
    }
    for i,bp := range state.Board.Plays {
        if state.Board.Covers[i] != nil {
            continue
        }
        for _,pc := range state.Hands[pIdx] {
            if pc.Beats(bp, state.Trump) {
                act := Action{PlayerIdx: pIdx, Verb: "Defend", Card: pc, Cover: bp}
                res = append(res, &act)
            }
        }
    }
    // Get non-nil covers
    if len(state.Board.Plays) > 0 && state.Board.Covered() < len(state.Board.Plays) {
        act := Action{PlayerIdx: pIdx, Verb: "Pickup"}
        res = append(res, &act)
    }
    return res
}

func (state *GameState) Move(me int, opp int, cur int) (*Action, float64) {
    var acts []*Action
    var r float64
    var bestAct *Action
    if cur == state.Defender {
        acts = state.DefenderActions(cur)
    } else {
        acts = state.AttackerActions(cur)
    }
    evals := make([]float64, 2*len(acts))
    for i,act := range acts {
        s := state.Clone()
        s.TakeAction(act)
        // Check win
        if s.DeckSize == 0 && len(s.Hands[cur]) == 0 {
            return act, math.Inf(Ternary(cur == me, 1, -1))
        }
        // Check mystery card played
        // TODO mystery eval or propagation
        if act.Card != nil && act.Card.Suit == "?" {
            return act, 0
        }
        if act.Verb == "Pass" {
            return act, s.EvalPass(me) 
        }
        if act.Verb == "Pickup" {
            // Opponent's move will determine evaluation
            evals[2*i] = math.Inf(-1)
            _,r = s.Move(me, opp, 1-cur)
            evals[2*i+1] = -r
        } else {
            _,r = s.Move(me, opp, cur)
            evals[2*i] = r
            _,r = s.Move(me, opp, 1-cur)
            evals[2*i+1] = -r
        }
    }
    bestAct = nil
    best := math.Inf(-1)
    for i:=0; i<len(acts); i++ {
        e := Ternary(evals[2*i] < evals[2*i+1], evals[2*i], evals[2*i+1])
        if e > best {
            best = e
            bestAct = acts[i]
        }
    }
    return bestAct, best
}

func (state *GameState) TakeAction(act *Action) {
    switch act.Verb {
        case "Attack": {
            state.Hands[act.PlayerIdx] = Remove(state.Hands[act.PlayerIdx], act.Card)
            state.Board.Plays = append(state.Board.Plays, act.Card)
            state.Board.Covers = append(state.Board.Covers, nil)
        } 
        case "Defend": {
            state.Hands[act.PlayerIdx] = Remove(state.Hands[act.PlayerIdx], act.Card)
            idx := IndexOf(state.Board.Plays, act.Cover)
            state.Board.Covers[idx] = act.Card
        }
        case "Pickup": {
            state.PickingUp = true
        }
        case "Reverse": {
            state.Hands[act.PlayerIdx] = Remove(state.Hands[act.PlayerIdx], act.Card)
            state.Board.Plays = append(state.Board.Plays, act.Card)
            state.Board.Covers = append(state.Board.Covers, nil)
            state.Defender = 1-state.Defender
        }
        case "Pass": {
            // Skip, handled in Move code
        }
    }
}

// Simple evaluation function
// Assume you're the defender
func (state *GameState) EvalPass(me int) float64 {
    var val float64
    if state.PickingUp {
        val = 0
        // Trumps and Aces positive,
        // other cards negative
        for _,c := range state.Board.Cards() {
            if c.Suit == state.Trump || c.Rank == "Ace" {
                val += 1
            } else {
                val -= 1
            }
        }
        if me != state.Defender {
            val *= -1
        }
    } else {
        // Positive for low covers and high plays (assuming you're the defender)
        // Advantage for defending
        if me == state.Defender {
            val = 2 
            for _,c := range state.Board.Covers {
                if c.Suit == state.Trump || c.Rank == "Ace" {
                    val -= 1
                } 
            }
        } else {
            val = 2
            for _,c := range state.Board.Plays {
                if c.Suit == state.Trump || c.Rank == "Ace" {
                    val -= 1
                } 
            }
            for _,c := range state.Board.Covers {
                if c.Suit == state.Trump || c.Rank == "Ace" {
                    val += 1
                }
            }
        }
    }
    return val
}
