package main

import (
    "encoding/json"
    //"fmt"
    "math"
    clone "github.com/huandu/go-clone"
)

func ChainString(chain []*Action) string {
    jsn, _ := json.Marshal(chain) 
    return string(jsn) 
}

func HandString(h []*Card) string {
    jsn, _ := json.Marshal(h) 
    return string(jsn)
}

func (a *Action) String() string {
    jsn, _ := json.Marshal(a)
    return string(jsn)
}

func (b *Board) Size() int {
    return len(b.Cards()) 
}

func (b *Board) Cards() []*Card {
    return Cat(b.Plays, NotNil(b.Covers)) 
}

func (b *Board) ReverseRank() string {
    if len(b.Plays) == 0 || len(NotNil(b.Covers)) > 0 {
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
        hands[i] = clone.Clone(p.Hand).([]*Card)
    }
    return &GameState{
        Defender: game.Defender,
        PickingUp: game.PickingUp,
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
                if bc != nil && (bc.Rank == pc.Rank || pc.Rank == "?")&& IndexOfFn(res, func(act *Action) bool {return act.Card == pc}) == -1 {
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
    if state.PickingUp {
        return res
    }
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
            if pc.Beats(bp, state.Trump) || pc.Rank == "?" {
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

func StartChain(act *Action) []*Action {
    return []*Action{act}
}

func (state *GameState) Move(me int, depth int) ([]*Action,float64) {
    if depth > 7 {
        return StartChain(nil), state.EvalMystery(me)
    }
    var acts []*Action
    if me == state.Defender {
        acts = state.DefenderActions(me)
    } else {
        acts = state.AttackerActions(me)
    }
    // If you have no actions, return
    if len(acts) == 0 {
        return nil,0 
    }
    evals := make([]float64, 2*len(acts))
    chains := make([][]*Action, 2*len(acts))
    for i,act := range acts {
        s := state.Clone()
        s.TakeAction(act)
        // Check win
        // Infinite value can confuse later action selection?
        if s.DeckSize == 0 && len(s.Hands[me]) == 0 {
            return StartChain(act), 1000
        }
        // End hand
        if act.Verb == "Pass" {
            return StartChain(act), s.EvalPass(me) 
        }
        // Mystery card
        // TODO incorporate possibility of reverse
        if act.Card != nil && act.Card.Rank == "?" {
            //return StartChain(act), s.EvalMystery(me)
            chains[2*i] = StartChain(act)
            evals[2*i] = s.EvalMystery(me)
            chains[2*i+1] = nil
        // Pickup - Opponent's move will determine evaluation
        // Penalize high hand count
        // Penalize taking cards with zero deck size (end of game)
        } else if act.Verb == "Pickup" {
            chains[2*i] = nil
            c,r := s.Move(1-me, depth+1)
            evals[2*i+1] = s.PickupPenalty(me) + r
            chains[2*i+1] = Ternary(c == nil, nil, append(c, act))
        // Regular known move
        } else {
            c,r := s.Move(me, depth+1)
            evals[2*i] = r
            chains[2*i] = Ternary(c == nil, nil, append(c, act))
            c,r = s.Move(1-me, depth+1)
            evals[2*i+1] = r
            chains[2*i+1] = Ternary(c == nil, nil, append(c, act))
        }
    }
    // Find likely action
    var bestChain []*Action
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

func (state *GameState) SumValue(cards []*Card) float64 {
    res := float64(0)
    for _,c := range cards {
        if c != nil && c.Rank != "?" {
            res += float64(IndexOf(ranks, c.Rank) - 4)
            if c.Suit == state.Trump {
                res += 7
            }
        }
    }
    return res
}

// Simple evaluation function
func (state *GameState) EvalPass(me int) float64 {
    var val float64
    if state.PickingUp {
        val = state.SumValue(Cat(state.Board.Plays, state.Board.Covers))
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
        return state.SumValue(state.Board.Plays) - state.SumValue(state.Board.Covers)
    } else {
        return state.SumValue(state.Board.Covers) - state.SumValue(state.Board.Plays)
    }
}

func (state *GameState) PickupPenalty(me int) float64 {
    val := 0
    if state.DeckSize < 2 {
        val += 8
    }
    hsize := len(state.Hands[me])
    if hsize > 4 {
        val += hsize-4
    }
    val += len(NotNil(Cat(state.Board.Plays, state.Board.Covers)))
    return float64(val)
}
