package durak

import (
    //"encoding/json"
    //"log"
    "time" 
)

// Emprically chosen
func (state *GameState) DepthLimit() int {
    nCards := 0
    for i := 0; i < len(state.Hands); i++ {
        nCards += len(state.Hands[i])
    }
    if nCards > 20 {
        return 6
    } else if nCards > 14 {
        return 9
    } else if nCards > 12 {
        return 10
    } else if nCards > 10 {
        return 11
    } else {
        return 16
    }
}

// Tunable parameters
type EvalParams struct {
    CardValueTrumpBonus int
    CardValueCardOffset int
    HandSizePickingUpMult float64
    HandSizeSmallDeckLimit int
    HandSizeSmallDeckMult float64
    SmallDeckLimit int
    NotLastWinnerValue int
    HandMult float64
    KnownMult float64
}

// Default parameters for eval
var DefaultEvalParams = EvalParams {6, 4, 2.0, 3, 2.0, 2, 20, 2.0, 1.0}

func Contains(card Card, cards []Card) bool {
    for _, c := range cards {
        if c == card {
            return true
        }
    }
    return false
}

func (state *GameState) CardValue(card Card, params *EvalParams) float64 {
    if card == UNK_CARD {
        return 0
    }
    if card.Suit() == state.Trump.Suit() {
        return float64(params.CardValueTrumpBonus + card.Rank())
    }
    return float64(card.Rank() - params.CardValueCardOffset)
}

func (orig *GameState) EvalTable(cur *GameState, me int, params *EvalParams) float64 {
    v := float64(0)
    // Cards you've played
    for _, card := range orig.Hands[me] {
        if Contains(card, cur.Hands[me]) {
            continue
        }
        v -= float64(orig.CardValue(card, params))
    }
    // Cards opponents have played on the table
    for i, _ := range cur.Plays {
        c1 := cur.Plays[i]
        c2 := cur.Covers[i]
        if !Contains(c1, orig.Hands[me]) {
            v += float64(orig.CardValue(c1, params))
        }
        if c2 != UNK_CARD && !Contains(c2, orig.Hands[me]) {
            v += float64(orig.CardValue(c2, params))
        }
    }
    return v
}

func (orig *GameState) EvalHandSizes(cur *GameState, me int, params *EvalParams) float64 {
    v := float64(0)
    div := float64(len(cur.Hands)-1)
    for i := 0; i < len(cur.Hands); i++ {
        if i == me {
            v -= float64(len(cur.Hands[me]))
        } else {
            v += float64(len(cur.Hands[i]))/div
        }
    }
    if cur.PickingUp {
        u := float64(len(cur.Plays) + NumNotUnk(cur.Covers))
        u *= params.HandSizePickingUpMult
        if cur.Defender == me {
            v -= u
        } else {
            v += u/div
        }
    }
    if len(cur.gamePtr.Deck) < params.HandSizeSmallDeckLimit {
        v *= params.HandSizeSmallDeckMult
    }
    return v
}

// Will only work once cards are known
func (orig *GameState) EvalKnownCards(cur *GameState, me int, params *EvalParams) float64 {
    v := float64(0)
    div := float64(len(cur.Hands)-1)
    for i := 0; i < len(cur.Hands); i++ {
        for _, card := range cur.Hands[i] {
            if i == me {
                v += float64(orig.CardValue(card, params))
            } else {
                v -= float64(orig.CardValue(card, params))/div
            }
        }
    }
    return v
}

func (orig *GameState) Eval(cur *GameState, me int, params *EvalParams) float64 {
    a := params.KnownMult*orig.EvalKnownCards(cur, me, params)
    b := params.HandMult*orig.EvalHandSizes(cur, me, params)
    c := orig.EvalTable(cur, me, params)
    return a + b + c
}

func (orig *GameState) EvalNode(me int, params *EvalParams) ([]Action, float64) {
    var chain []Action
    var res float64
    node := orig.Clone()
    node.start = time.Now() 
    for depth := 1; depth <= 20; depth++ {
        acts, v, ok := node.EvalNodeDeep(node.Clone(), me, depth, params)
        if !ok {
            break
        } else {
            if acts != nil {
                chain = acts
                res = v
            }
        }
    }
    return chain, res
}

func (orig *GameState) EvalNodeDeep(cur *GameState, me int, depth int, params *EvalParams) ([]Action, float64, bool) {
    if params == nil {
        params = &DefaultEvalParams
    }
    deckSize := len(cur.gamePtr.Deck)
    // Exceeded time limit
    elapsed := time.Now().Sub(orig.start)
    if elapsed.Seconds() > 5 {
        return nil, 0, false
    }
    // You've already won and don't take actions
    if cur.Won[me] {
        return nil, 0, true
    }
    if depth == 0 {
        return make([]Action, 0), orig.Eval(cur, me, params), true
    }
    acts := cur.PlayerActions(me)
    // If you have no actions, return
    if len(acts) == 0 {
        return nil, 0, true
    }
    // If it's depth 0 and you have 1 action, just take that action (don't waste time evaluating)
    // Note that oftentimes you have a Defer action in addition to bad ones
    if depth == 0 && len(acts) == 1 {
        return []Action{acts[0]}, 0, true
    }
    // Default values should be 0 and nil
    np := len(cur.Hands)
    evals := make([]float64, np*len(acts))
    chains := make([][]Action, np*len(acts))
    playedUnkCard := false
    for i, act := range acts {
        s := cur.Clone();
        s.TakeAction(act);
        // Check win (getting rid of cards)
        if deckSize == 0 && len(s.Hands[me]) == 0 {
            // If everyone else has already won, the final player with
            // cards should feel bad
            if s.CheckGameOver() {
                return []Action{act}, 1000, true
            }
            return []Action{act}, float64(params.NotLastWinnerValue), true
        }
        // You don't get actions but opponents do
        if act.Verb == DeferVerb {
            for j := 0; j < np; j++ {
                if j == me {
                    continue
                }
                c, r, finished := orig.EvalNodeDeep(s, j, depth-1, params)
                if !finished {
                    return nil, 0, false
                }
                if c != nil {
                    evals[np*i+j] = r
                    chains[np*i+j] = append(c, act)
                }
            }
        // End hand
        } else if act.Verb == PassVerb {
            // Go to the end of game
            if deckSize <= params.SmallDeckLimit {
                // Defender picking up
                // All others can play
                // NOTE: Action is Pass, so cur.PickingUp is valid
                // NOTE: Original code did not take into account more than 2 players?
                if cur.PickingUp {
                    for j := 0; j < np; j++ {
                        if j == cur.Defender {
                            continue
                        }
                        c, r, finished := orig.EvalNodeDeep(s, j, depth-1, params)
                        if !finished {
                            return nil, 0, false
                        }
                        if c != nil {
                            evals[np*i+j] = r
                            chains[np*i+j] = append(c, act)
                        }
                    }
                // Defender successfully defended and will go next
                // But other players can also maybe go
                } else {
                    for j := 0; j < np; j++ {
                        if j != me {
                            c, r, finished := orig.EvalNodeDeep(s, j, depth-1, params)
                            if !finished {
                                return nil, 0, false
                            }
                            if c != nil {
                                evals[np*i+j] = r
                                chains[np*i+j] = append(c, act)
                            }
                        }
                    }
                }
            // Go per-hand
            } else {
                return []Action{act}, orig.Eval(cur, me, params), true
            }
        // Unknown card play or cover
        // Only check one mystery card
        } else if act.Card == UNK_CARD && (act.Verb == PlayVerb || act.Verb == CoverVerb) {
            if !playedUnkCard {
                playedUnkCard = true
                evals[np*i+me] = orig.Eval(s, me, params)
                chains[np*i+me] = []Action{act}
            } 
        // Pickup - Opponents' moves will determine evaluation
        // Penalize taking cards with zero deck size (end of game)
        } else if act.Verb == PickupVerb {
            for j := 0; j < np; j++ {
                if j == me {
                    continue
                }
                c, r, finished := orig.EvalNodeDeep(s, j, depth-1, params)
                if !finished {
                    return nil, 0, false
                }
                if c != nil {
                    evals[np*i+j] = r
                    chains[np*i+j] = append(c, act)
                }
            }
        // Ordinary action
        } else {
            for j := 0; j < np; j++ {
                c, r, finished := orig.EvalNodeDeep(s, j, depth-1, params)
                if !finished {
                    return nil, 0, false
                }
                if c != nil {
                    evals[np*i+j] = r
                    chains[np*i+j] = append(c, act)
                }
            }
        }
    }
    best := -10000.0
    besti := 0
    for i, _ := range acts {
        for j := 0; j < np; j++ {
            if chains[np*i+j] == nil {
                continue
            }
            e := evals[np*i+j]
            if j != me {
                e *= -1
            }
            if e > best {
                besti = np*i+j
                best = e
            }
        }
    }
    return chains[besti], best, true
}

