
use crate::rules;

type GameState = rules::GameState;
type Action = rules::Action;
type Card = rules::Card;
type Verb = rules::Verb;

// These are sometimes not representative of how long it takes to eval
// That's why we also measure time
// Aka iterative crappening

pub fn depth_limit(state: &GameState) -> usize {
    let ncards = state.hands.iter().map(|x| x.len()).sum::<usize>();
    if ncards > 18 {
        6
    } else if ncards > 12 {
        8
    } else if ncards > 10 {
        14
    } else if ncards > 8{
        12
    } else {
        16
    }
}

pub fn eval_node(orig: &GameState, working: Option<GameState>, me: usize, depth: usize, dlim: usize, empty_deck: bool) 
        -> (Option<Vec<Action>>, i32) {
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
    let elapsed = cur.start.unwrap().elapsed().as_seconds_f32();
    if elapsed > 20.0 {
        dlim_adj -= 1*(elapsed-20.0).floor() as usize;
    }
    if depth > dlim_adj {
        return (Some(Vec::new()), orig.eval(&cur, me, empty_deck))
    }
    let acts = cur.player_actions(me);
    // If you have no actions, return
    if acts.len() == 0 {
        return (None, 0)
    }
    let mut evals = vec![0; 2*acts.len()];
    let mut chains = vec![None; 2*acts.len()];
    let mut did_mystery = false;
    for i in 0..acts.len() {
        // Only allow defer as the first action of search
        // So AI can keep deferring on polling
        // But won't clog up the search stack
        if acts[i].verb == Verb::Defer && depth != 0 {
            continue
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
                if cur.picking_up {
                    let (c, r) = eval_node(orig, Some(s), me, depth+1, dlim_adj, true);
                    let mut cu = c.unwrap();
                    cu.push(acts[i]);
                    evals[2*i] = r;
                    evals[2*i+1] = 0;
                    chains[2*i] = Some(cu);
                    chains[2*i+1] = None;
                // Opponent successfully defended and will go next
                } else {
                    let (c, r) = eval_node(orig, Some(s), 1-me, depth+1, dlim_adj, true);
                    let mut cu = c.unwrap();
                    cu.push(acts[i]);
                    evals[2*i] = 0;
                    evals[2*i+1] = r;
                    chains[2*i+1] = Some(cu);
                    chains[2*i] = None;
                }
            // Go per-hand
            } else {
                return (Some(vec![acts[i]]), orig.eval_pass(&cur, me, false))
            }
        }
        // Unknown card play or cover
        // Only check one mystery card
        /*else if !did_mystery
                && (acts[i].verb == Verb::Play || acts[i].verb == Verb::Cover) 
                && acts[i].card == rules::UNK_CARD {
            did_mystery = true;
            evals[2*i] = orig.eval(&s, me, empty_deck);
            evals[2*i+1] = 0;
            chains[2*i] = Some(vec![acts[i]]);
            chains[2*i+1] = None;
        }*/
        // Pickup - Opponent's move will determine evaluation
        // Penalize taking cards with zero deck size (end of game)
        else if acts[i].verb == Verb::PickUp {
            evals[2*i] = 0;
            chains[2*i] = None;
            let (c, r) = eval_node(orig, Some(s), 1-me, depth+1, dlim_adj, empty_deck);
            let mut cu = c.unwrap();
            cu.push(acts[i]);
            evals[2*i+1] = r+3; // +3 penalty
            chains[2*i+1] = Some(cu);
        // Ordinary action
        } else {
            let (c, r) = eval_node(orig, Some(s.clone()), me, depth+1, dlim_adj, empty_deck);
            match c {
                Some(mut cu) => {
                    cu.push(acts[i]);
                    evals[2*i] = r;
                    chains[2*i] = Some(cu);
                }
                None => {
                    chains[2*i] = None;
                }
            }
            let (c, r) = eval_node(orig, Some(s), 1-me, depth+1, dlim_adj, empty_deck);
            match c {
                Some(mut cu) => {
                    cu.push(acts[i]);
                    evals[2*i+1] = r;
                    chains[2*i+1] = Some(cu);
                }
                None => {
                    chains[2*i+1] = None;
                }
            }
        }
    }
    let mut best = -10000;
    let mut besti = 0;
    for i in 0..acts.len() {
        if chains[2*i].is_none() && chains[2*i+1].is_none() {
            continue
        }
        if chains[2*i].is_none() {
            let mut e = -evals[2*i+1];
            if e > best {
                besti = 2*i+1;
                best = e;
            }
        } else if chains[2*i+1].is_none() {
            let mut e = evals[2*i];
            if e > best {
                besti = 2*i;
                best = e;
            }
        } else {
            let mut e1 = evals[2*i];
            let mut e2 = -evals[2*i+1];
            if e1 > best {
                besti = 2*i;
                best = e1;
            }
            if e2 > best {
                besti = 2*i+1;
                best = e2;
            }
        }
    }
    return (chains.swap_remove(besti), best)
}

impl GameState {
    pub fn eval(&self, cur: &GameState, me: usize, empty_deck: bool) -> i32 {
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
        if empty_deck {
            v += 2 - cur.hands[me].len() as i32 + cur.hands[1-me].len() as i32;
        }
        v
    } 
    pub fn eval_pass(&self, cur: &GameState, me: usize, empty_deck: bool) -> i32 {
        let mut bonus = 0;
        if cur.picking_up {
            bonus = cur.plays.len() as i32 + rules::num_not_unk(&cur.covers) as i32;
            if cur.defender == me {
                bonus *= -1;
            }
        }
        self.eval(cur, me, empty_deck) + bonus
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
    use crate::search;
    
    #[test]
    fn test_eval() {
        let mut g = rules::Game::new(0, "computer".to_string());
        let val = g.state.eval(&g.state.clone(), 0, false);
        assert_eq!(val, 0);
    }

    #[test]
    fn test_eval_node() {
        let mut g = rules::Game::new(0, "computer".to_string());
        let (c, r) = search::eval_node(&g.state, None, 0, 0, 0, false);
        assert_eq!(c.is_none(), false);
        assert_ne!(c.unwrap().len(), 0);
    }

    #[test]
    fn test_eval_node_endgame() {
        let mut g = rules::Game::new(0, "computer".to_string());
        g.deck = Vec::new();
        let (c, r) = search::eval_node(&g.state, None, 0, 0, 0, true);
        assert_eq!(c.is_none(), false);
        assert_ne!(c.unwrap().len(), 0);
    }
}
