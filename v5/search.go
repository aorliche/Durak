package main

import (
    //"fmt"
    "time" 
)

func (state *GameState) DepthLimit() int {
    nCards := len(state.Hands[0]) + len(state.Hands[1]);
    if nCards > 18 {
        return 6
    } else if nCards > 12 {
        return 8
    } else if nCards > 10 {
        return 10
    } else if nCards > 8{
        return 12
    } else {
        return 16
    }
}

func Contains(card Card, cards []Card) bool {
    for _, c := range cards {
        if c == card {
            return true
        }
    }
    return false
}

func (state *GameState) CardValue(card Card) int {
    v := 0
    if card == UNK_CARD {
        return v
    }
    if card.Suit() == state.Trump.Suit() {
        v += 6
    }
    v += card.Rank() - 4
    return v
}

func (orig *GameState) Eval(cur *GameState, me int, emptyDeck bool) int {
    v := 0
    // Cards you've played
    for _, card := range orig.Hands[me] {
        if Contains(card, cur.Hands[me]) {
            continue
        }
        v -= orig.CardValue(card);
    }
    // Cards opponent has played
    for i, _ := range cur.Plays {
        c1 := cur.Plays[i]
        c2 := cur.Covers[i]
        if !Contains(c1, orig.Hands[me]) {
            v += orig.CardValue(c1)
        }
        if c2 != UNK_CARD && !Contains(c2, orig.Hands[me]) {
            v += orig.CardValue(c2)
        }
    }
    if emptyDeck {
        v += 2 - len(cur.Hands[me]) + len(cur.Hands[1-me])
    }
    return v
}
    
func (orig *GameState) EvalPass(cur *GameState, me int, emptyDeck bool) int {
    bonus := 0
    if cur.PickingUp {
        bonus = len(cur.Plays) + NumNotUnk(cur.Covers)
        if cur.Defender == me {
            bonus *= -1
        }
    }
    return orig.Eval(cur, me, emptyDeck) + bonus
}

func (orig *GameState) EvalNode(cur *GameState, me int, depth int, dlim int, emptyDeck bool) ([]Action, int) {
    dlimAdj := dlim
    if depth == 0 {
        cur = orig.Clone()
        dlimAdj = orig.DepthLimit()
        orig.start = time.Now()
    } 
    elapsed := time.Now().Sub(orig.start)
    if elapsed.Seconds() > 2 {
        dlimAdj -= int(elapsed.Seconds()) - 2
    }
    if depth > dlimAdj {
        return make([]Action, 0), orig.Eval(cur, me, emptyDeck)
    }
    acts := cur.PlayerActions(me)
    // If you have no actions, return
    if len(acts) == 0 {
        return nil, 0
    }
    // Default values should be 0 and nil
    evals := make([]int, 2*len(acts))
    chains := make([][]Action, 2*len(acts))
    playedUnkCard := false
    for i, act := range acts {
        // Only allow defer as the first action of search
        // So AI can keep deferring on polling
        // But won't clog up the search stack
        if act.Verb == DeferVerb && depth != 0 {
            continue
        }
        s := cur.Clone();
        s.TakeAction(act);
        // Check win
        if emptyDeck && len(s.Hands[me]) == 0 {
            return []Action{act}, 1000
        }
        // End hand
        if act.Verb == PassVerb {
            // Go to the end of game
            if emptyDeck {
                // You passed with opponent picking up
                // You will play next turn
                if cur.PickingUp {
                    c, r := orig.EvalNode(s, me, depth+1, dlimAdj, true)
                    evals[2*i] = r
                    evals[2*i+1] = 0
                    chains[2*i] = append(c, act)
                    chains[2*i+1] = nil
                // Opponent successfully defended and will go next
                } else {
                    c, r := orig.EvalNode(s, 1-me, depth+1, dlimAdj, true)
                    evals[2*i] = 0
                    evals[2*i+1] = r
                    chains[2*i] = nil
                    chains[2*i+1] = append(c, act)
                }
            // Go per-hand
            } else {
                return []Action{act}, orig.EvalPass(cur, me, false)
            }
        // Unknown card play or cover
        // Only check one mystery card
        } else if act.Card == UNK_CARD && (act.Verb == PlayVerb || act.Verb == CoverVerb) {
            if !playedUnkCard {
                playedUnkCard = true
                evals[2*i] = orig.Eval(s, me, emptyDeck)
                evals[2*i+1] = 0
                chains[2*i] = []Action{act}
                chains[2*i+1] = nil
            } 
        // Pickup - Opponent's move will determine evaluation
        // Penalize taking cards with zero deck size (end of game)
        } else if act.Verb == PickupVerb {
            c, r := orig.EvalNode(s, 1-me, depth+1, dlimAdj, emptyDeck)
            evals[2*i] = 0
            evals[2*i+1] = r+3 // +3 penalty
            chains[2*i] = nil
            chains[2*i+1] = append(c, act)
        // Ordinary action
        } else {
            c, r := orig.EvalNode(s, me, depth+1, dlimAdj, emptyDeck)
            if c != nil {
                evals[2*i] = r
                chains[2*i] = append(c, act)
            }
            c, r = orig.EvalNode(s, 1-me, depth+1, dlimAdj, emptyDeck)
            if c != nil {
                evals[2*i+1] = r
                chains[2*i+1] = append(c, act)
            }
        }
    }
    best := -10000
    besti := 0
    for i, _ := range acts {
        if chains[2*i] == nil && chains[2*i+1] == nil {
            continue
        }
        if chains[2*i] == nil {
            e := -evals[2*i+1]
            if e > best {
                besti = 2*i+1
                best = e
            }
        } else if chains[2*i+1] == nil {
            e := evals[2*i]
            if e > best {
                besti = 2*i
                best = e
            }
        } else {
            e1 := evals[2*i]
            e2 := -evals[2*i+1]
            if e1 > best {
                besti = 2*i
                best = e1
            }
            if e2 > best {
                besti = 2*i+1
                best = e2
            }
        }
    }
    return chains[besti], best
}

