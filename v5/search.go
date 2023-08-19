package main

import (
    //"encoding/json"
    //"log"
    "time" 
)

// Empriically chosen
func (state *GameState) DepthLimit() int {
    nCards := len(state.Hands[0]) + len(state.Hands[1]);
    if nCards > 18 {
        return 6
    } else if nCards > 12 {
        return 9
    } else if nCards > 10 {
        return 10
    } else if nCards > 8{
        return 11
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
    if card == UNK_CARD {
        return 0
    }
    if card.Suit() == state.Trump.Suit() {
        return 6 + card.Rank()
    }
    return card.Rank() - 4
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
        v += 2
        // Hands modifier
        for i := 0; i < len(orig.Hands); i++ {
            if i == me {
                v -= len(cur.Hands[me])
            } else {
                v += len(cur.Hands[i])
            }
        }
    }
    return v
}

// Applied on end of midgame hand or unknown card play
func (cur *GameState) HandsPenalty(me int) int {
    v := 0
    for i := 0; i < len(cur.Hands); i++ {
        if len(cur.Hands[i]) <= 4 {
            continue
        }
        if i == me {
            v -= len(cur.Hands[me])-4
        } else {
            v += len(cur.Hands[i])-4
        }
    }
    return v
}

func (orig *GameState) EvalNode(cur *GameState, me int, depth int, dlim int, deckSize int) ([]Action, int) {
    dlimAdj := dlim
    if depth == 0 {
        cur = orig.Clone()
        dlimAdj = orig.DepthLimit()
        orig.start = time.Now()
    } 
    // Iterative crappening
    elapsed := time.Now().Sub(orig.start)
    if elapsed.Seconds() > 2 {
        dlimAdj -= (int(elapsed.Seconds()) - 2)/2
    }
    if depth > dlimAdj {
        return make([]Action, 0), orig.Eval(cur, me, deckSize < 3)
    }
    acts := cur.PlayerActions(me)
    // If you have no actions, return
    if len(acts) == 0 {
        return nil, 0
    }
    // If everyone else has already won, you lost
    if cur.gamePtr.CheckGameOver() {
        return make([]Action, 0), -1000
    }
    // Default values should be 0 and nil
    np := len(cur.Hands)
    evals := make([]int, np*len(acts))
    chains := make([][]Action, np*len(acts))
    playedUnkCard := false
    for i, act := range acts {
        s := cur.Clone();
        s.TakeAction(act);
        // Check win (getting rid of cards)
        if deckSize == 0 && len(s.Hands[me]) == 0 {
            return []Action{act}, 200
        }
        // You don't get actions but opponents do
        if act.Verb == DeferVerb {
            for j := 0; j < np; j++ {
                if j == me {
                    continue    
                }
                c, r := orig.EvalNode(s, j, depth+1, dlimAdj, deckSize)
                evals[np*i+j] = r
                chains[np*i+j] = append(c, act)
            }
        // End hand
        } else if act.Verb == PassVerb {
            // Go to the end of game
            if deckSize == 0 {
                // Defender picking up
                // All others can play
                // NOTE: Action is Pass, so cur.PickingUp is valid
                if cur.PickingUp {
                    for j := 0; j < np; j++ {
                        if j == cur.Defender {
                            continue
                        }
                        c, r := orig.EvalNode(s, s.Attacker, depth+1, dlimAdj, 0)
                        evals[np*i+j] = r + s.HandsPenalty(me)
                        chains[np*i+j] = append(c, act)
                    }
                // Defender successfully defended and will go next
                } else {
                    c, r := orig.EvalNode(s, cur.Defender, depth+1, dlimAdj, 0) 
                    evals[np*i+cur.Defender] = r + s.HandsPenalty(cur.Defender)
                    chains[np*i+cur.Defender] = append(c, act)
                }
            // Go per-hand
            } else {
                return []Action{act}, orig.Eval(cur, me, deckSize < 3) + s.HandsPenalty(me)
            }
        // Unknown card play or cover
        // Only check one mystery card
        } else if act.Card == UNK_CARD && (act.Verb == PlayVerb || act.Verb == CoverVerb) {
            if !playedUnkCard {
                playedUnkCard = true
                evals[np*i+me] = orig.Eval(s, me, deckSize < 3) + s.HandsPenalty(me)
                chains[np*i+me] = []Action{act}
            } 
        // Pickup - Opponents' moves will determine evaluation
        // Penalize taking cards with zero deck size (end of game)
        } else if act.Verb == PickupVerb {
            for j := 0; j < np; j++ {
                if j == me {
                    continue
                }
                c, r := orig.EvalNode(s, j, depth+1, dlimAdj, deckSize)
                if len(s.Hands[me]) > 6 {
                    r += len(s.Plays)+NumNotUnk(s.Covers)+len(s.Hands[me])-6
                }
                evals[np*i+j] = r
                chains[np*i+j] = append(c, act)
            }
        // Ordinary action
        } else {
            for j := 0; j < np; j++ {
                c, r := orig.EvalNode(s, j, depth+1, dlimAdj, deckSize)
                if c != nil {
                    evals[np*i+j] = r
                    chains[np*i+j] = append(c, act)
                }
            }
        }
    }
    best := -10000
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
    return chains[besti], best
}

