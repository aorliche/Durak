
use crate::rules;

type GameState = rules::GameState;
type Action = rules::Action;
type Card = rules::Card;

// These are sometimes not representative of how long it takes to eval
// That's why we also measure time
// Aka iterative crappening

pub fn depth_limit(state: &GameState) -> usize {
    let ncards = state.plays.len() + rules::num_not_unk(&state.covers);
    if ncards > 18 {
        5
    } else if ncards > 12 {
        6
    } else if ncards > 10 {
        7
    } else if ncards > 8{
        8
    } else {
        12
    }
}

pub fn eval_node(orig: &GameState, working: Option<GameState>, me: usize, depth: usize, dlim: usize) -> (Vec<Action>, i32) {
    let mut cur : GameState;
    let mut dlim_adj : usize;
    if depth == 0 {
        assert!(working.is_none());
        cur = orig.clone();
        dlim_adj = depth_limit(orig);
        cur.start = Some(time::Instant::now());
    } else {
        cur = working.unwrap();
        dlim_adj = dlim;
    }
    if cur.start.unwrap().elapsed().as_seconds_f32() > 1.5 {
        dlim_adj -= 1
    }
    if depth > dlim_adj {
        return (Vec::new(), cur.eval(&orig, me))
    }
    return (Vec::new(), 0)
}


impl GameState {
    pub fn eval(&self, cur: &GameState, me: usize) -> i32 {
        let mut v = 0;
        // Cards you've played
        for i in 0..self.hands[me].len() {
            let card = self.hands[me][i];
            if cur.hands[me].contains(&card) {
                continue
            }
            v -= self.card_value(card);
        }
        // Cards opponent has played
        for i in 0..cur.plays.len() {
            let c1 = cur.plays[i];
            let c2 = cur.covers[i];
            if !self.hands[me].contains(&c1) {
                v += self.card_value(c1);
            }
            if c2 != rules::UNK_CARD {
                v += self.card_value(c2);
            }
        }
        v
    } 
    pub fn card_value(&self, card: Card) -> i32 {
        let mut v = 0;
        if card == rules::UNK_CARD {
            return v
        }
        if rules::card_suit(card) == rules::card_suit(self.trump) {
            v += 6;
        }
        v += rules::card_rank(card) as i32 - 4;
        v
    }
}

#[cfg(test)]
mod tests {
    use crate::rules;
    
    #[test]
    fn test_eval() {
        let mut g = rules::Game::new(0, "computer".to_string());
        let val = g.state.eval(&g.state.clone(), 0);
        assert_eq!(val, 0);
    }

    /*#[test]
    fn test_eval_2moves() {
        let mut g = rules::Game::new(0, "computer".to_string());
        g.take_action(&g.state.random_action());
        g.take_action(&g.state.random_action());
        let val = g.state.eval(&g.state.clone(), 0);
        assert_eq!(val, 1);
    }*/
}

/*func (state *GameState) EvalMystery(me int) float64 {
    if me == state.Defender {
        return state.SumValue(state.Plays) - state.SumValue(state.Covers)
    } else {
        return state.SumValue(state.Covers) - state.SumValue(state.Plays)
    }
}*/

/*func (state *GameState) Move(me int, depth int, dlim int) ([]*FastAction,float64) {
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
}*/
