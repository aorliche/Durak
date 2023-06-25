
use crate::rules;

type GameState = rules::GameState;
type Action = rules::Action;
type Card = rules::Card;
type Verb = rules::Verb;

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

pub fn eval_node(orig: &GameState, working: Option<GameState>, me: usize, depth: usize, dlim: usize, empty_deck: bool) -> (Option<Vec<Action>>, i32) {
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
        return (Some(Vec::new()), cur.eval(&orig, me))
    }
    let acts = cur.player_actions(me);
    // If you have no actions, return
    if acts.len() == 0 {
        return (None, 0)
    }
    let mut evals = Vec::new();
    let mut chains = Vec::new();
    let mut did_mystery = false;
    for i in 0..acts.len() {
        if empty_deck && acts[i].verb as usize == Verb::Defer as usize {
            return (Some(Vec::new()), cur.eval_pass(&orig, me))
        }
        let mut s = cur.clone();
        s.take_action(&acts[i]);
        // Check win
        if empty_deck && s.hands[me].len() == 0 {
            return (Some(vec![acts[i]]), 1000)
        }
        // End hand
        if acts[i].verb == Verb::Pass {
            // Go to the end of game
            if empty_deck {
                // You passed with opponent picking up
                // You will play next turn
                if s.picking_up {
                    let (c, r) = eval_node(orig, Some(s), me, depth+1, dlim, true);
                    c.push(acts[i]);
                    evals[2*i] = r;
                    evals[2*i+1] = 0;
                    chains[2*i] = c;
                    chains[2*i+1] = None;
                // Opponent successfully defended and will go next
                } else {
                    let (c, r) = eval_node(orig, Some(s), 1-me, depth+1, dlim, true);
                    c.push(acts[i]);
                    evals[2*i] = 0;
                    evals[2*i+1] = r;
                    chains[2*i+1] = c;
                    chains[2*i] = None;
                }
            // Go per-hand
            } else {
                return (Some(Vec::new()), s.eval_pass(&orig, me))
            }
        }
        // Unknown card play or cover
        // Only check one mystery card
        else if (acts[i].verb as usize == Verb::Play as usize 
            || acts[i].verb as usize == Verb::Cover as usize) 
            && !did_mystery 
            && acts[i].card == rules::UNK_CARD {
                did_mystery = true;
                evals[2*i] = s.eval(&orig, me);
                evals[2*i+1] = 0;
                chains[2*i] = Some(vec![acts[i]]);
                chains[2*i+1] = None;
            }
        }
        // Pickup - Opponent's move will determine evaluation
        // Penalize taking cards with zero deck size (end of game)
        else if acts[i].verb as usize == Verb::Pickup as usize {
            evals[2*i] = 0;
            chains[2*i] = None;
            let (c, r) = eval_node(orig, Some(s), 1-me, depth+1, dlim, empty_deck);
            c.push(acts[i]);
            chains[2*i+1] = r;
            chains[2*i+1] = Some(c);
        } else {
            let mut (c, r) = eval_node(orig, Some(s), me, depth+1, dlim, empty_deck);
            c.push(acts[i]);
            evals[2*i] = r;
            chains[2*i] = c;
            (c, r) = eval_node(orig, Some(s), 1-me, depth+1, dlim, empty_deck);
            c.push(acts[i]);
            evals[2*i+1] = r;
            chains[2*i+1] = c;
        }
    }
    let mut best = -10000;
    let mut besti = 0;
    for i in 0..chains.len() {
        
    }
    return (Some(Vec::new()), 0)
}

/*
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
    pub fn eval_pass(&self, orig: &GameState, me: usize) -> i32 {
        let bonus = 0;
        if self.picking_up {
            bonus = self.plays.len() as i32 + rules::num_not_unk(self.covers.len()) as i32;
            if self.defender == me {
                bonus *= -1;
            }
        }
        self.eval(orig, me) + bonus
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
