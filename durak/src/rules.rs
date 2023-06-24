use serde::{Deserialize, Serialize};
use serde_json::Value;

use rand::seq::SliceRandom;
use rand::{thread_rng, seq::IteratorRandom};

#[derive(Serialize, Deserialize, Clone, Copy, Debug)]
pub enum Verb {
    Play = 1,
    Cover,
    Reverse,
    Pass,
    PickUp,
    Defer,
}

pub type Card = usize;
pub const UNK_CARD : Card = 36;

pub fn new_card(card: usize) -> Card {
    card
}

pub fn card_from_suit_rank(suit: usize, rank: usize) -> Card {
    suit * 9 + rank
}

pub fn card_suit(card: Card) -> usize {
    card / 9
}

pub fn card_rank(card: Card) -> usize {
    card % 9
}

pub fn card_suit_string(card: Card) -> String {
    match card_suit(card) {
        0 => "spades".to_string(),
        1 => "hearts".to_string(),
        2 => "diamonds".to_string(),
        3 => "clubs".to_string(),
        _ => "unknown".to_string(),
    }
}

pub fn card_rank_string(card: Card) -> String {
    match card_rank(card) {
        0 => "6".to_string(),
        1 => "7".to_string(),
        2 => "8".to_string(),
        3 => "9".to_string(),
        4 => "10".to_string(),
        5 => "J".to_string(),
        6 => "Q".to_string(),
        7 => "K".to_string(),
        8 => "A".to_string(),
        _ => "unknown".to_string(),
    }
}

pub fn card_to_string(card: Card) -> String {
    format!("{} of {}", card_rank_string(card), card_suit_string(card))
}

pub fn beats(card1: Card, card2: Card, trump: Card) -> bool {
    if card_suit(card1) == card_suit(trump) && card_suit(card1) != card_suit(card2) {
        return true
    }
    card_rank(card1) > card_rank(card2) && card_suit(card1) == card_suit(card2)
}
    
pub fn generate_deck() -> Vec<Card> {
    let mut res = Vec::new();
    for suit in 0..4 {
        for rank in 0..9 {
            res.push(card_from_suit_rank(suit, rank));
        }
    }
    res.shuffle(&mut thread_rng());
    res
}

pub fn remove_card(cards: &mut Vec<Card>, card: Card) -> bool {
    let idx = cards.iter().position(|x| *x == card);
    match idx {
        Some(i) => {
            cards.remove(i);
            true
        },
        None => false,
    }
}

#[derive(Serialize, Deserialize, Clone, Copy, Debug)]
pub struct Action {
    pub player: usize,
    pub verb: Verb,
    pub card: Card,
    pub covering: Card,
}

impl Action {
    pub fn new(player: usize, verb: Verb, card: Card, covering: Card) -> Action {
        Action {
            player,
            verb,
            card,
            covering,
        }
    }
}

#[derive(Serialize, Deserialize, Debug)]
pub struct GameState {
    pub defender: usize,
    pub picking_up: bool,
    pub trump: Card,
    pub plays: Vec<Card>,
    pub covers: Vec<Card>,
    pub hands: Vec<Vec<Card>>,
    pub defering: bool,
    #[serde(skip)]
    pub start: Option<time::Instant>,
}

pub fn num_not_unk(v: &Vec<Card>) -> usize {
    v.iter().filter(|&x| *x != UNK_CARD).count()
}

impl GameState {
    pub fn new(trump: Card, hands: Vec<Vec<Card>>) -> GameState {
        GameState {
            defender: 1,
            picking_up: false,
            trump,
            plays: Vec::new(),
            covers: Vec::new(),
            hands,
            defering: false,
            start: None,
        } 
    }
    pub fn attacker_actions(&self, pidx: usize) -> Vec<Action> {
        let mut res = Vec::new();
        if self.plays.len() == 0 {
           for &card in &self.hands[pidx] {
               res.push(Action::new(pidx, Verb::Play, card, UNK_CARD))
           } 
        } else {
            for &card in &self.hands[pidx] {
                for &board_card in self.plays.iter().chain(self.covers.iter()) {
                    if board_card == UNK_CARD {
                        continue;
                    }
                    // Card equality with unknown card for AI search
                    if card_rank(card) == card_rank(board_card) || card == UNK_CARD {
                        res.push(Action::new(pidx, Verb::Play, card, UNK_CARD));
                        break;
                    }
                }
            }
        }
        if self.picking_up || (num_not_unk(&self.covers) == self.plays.len() && self.plays.len() > 0) {
            res.push(Action::new(pidx, Verb::Pass, UNK_CARD, UNK_CARD))
        }
        // For AI to not throw trumps away
        if !self.picking_up && self.plays.len() > num_not_unk(&self.covers) {
            res.push(Action::new(pidx, Verb::Defer, UNK_CARD, UNK_CARD)) 
        }
        res
    }
    pub fn reverse_rank(&self) -> Option<usize> {
        if self.plays.len() == 0 && self.covers.len() == 0 {
            return None
        }
        let mut card = self.plays[0];
        for &c in self.plays.iter().chain(self.covers.iter()) {
            if c != UNK_CARD && c != card {
                return None
            }
        }
        return Some(card_rank(card))
    }
    pub fn defender_actions(&self, pidx: usize) -> Vec<Action> {
        let mut res = Vec::new();
        if self.picking_up {
            return res
        }
        let rev_rank = self.reverse_rank();
        match rev_rank {
            Some(rank) => {
                for &card in &self.hands[pidx] {
                    if card_rank(card) == rank {
                        res.push(Action::new(pidx, Verb::Reverse, card, UNK_CARD))
                    }
                }
            },
            None => (),
        }
        for i in 0..self.plays.len() {
            for &card in self.hands[pidx].iter() {
                if self.covers[i] == UNK_CARD && beats(card, self.plays[i], self.trump) {
                    res.push(Action::new(pidx, Verb::Cover, card, self.plays[i]))
                }
            }
        }
        if self.plays.len() > 0 && self.covers.iter().filter(|&x| *x != UNK_CARD).count() < self.plays.len() {
            res.push(Action::new(pidx, Verb::PickUp, UNK_CARD, UNK_CARD))
        }
        res
    }
    pub fn player_actions(&self, pidx: usize) -> Vec<Action> {
        if pidx == self.defender {
            self.defender_actions(pidx)
        } else {
            self.attacker_actions(pidx)
        }
    }
    pub fn random_action(&self) -> Action {
        self.player_actions(0).iter().chain(self.player_actions(1).iter()).choose(&mut thread_rng()).unwrap().clone()
    }
    pub fn take_action(&mut self, action: &Action) {
        match action.verb {
            Verb::Play => {
                self.plays.push(action.card);
                self.covers.push(UNK_CARD);
                remove_card(&mut self.hands[action.player], action.card);
            }, 
            Verb::Cover => {
                for i in 0..self.plays.len() {
                    if action.covering == self.plays[i] {
                        self.covers[i] = action.card;
                    }
                }
                remove_card(&mut self.hands[action.player], action.card);
            },
            Verb::Reverse => {
                self.plays.push(action.card);
                self.covers.push(UNK_CARD);
                remove_card(&mut self.hands[action.player], action.card);
                self.defender = 1-self.defender;
                self.defering = false;
            },
            Verb::PickUp => {
                self.picking_up = true;
            },
            Verb::Pass => {
                if self.picking_up {
                    for &card in self.plays.iter().chain(self.covers.iter()) {
                        if card != UNK_CARD {
                            self.hands[self.defender].push(card);
                        } 
                    } 
                }
                self.plays.clear();
                self.covers.clear();
                self.picking_up = false;
                self.defering = false;
            },
            Verb::Defer => {
                self.defering = true;
            },
            _ => (),
        }
    }
    pub fn clone(&self) -> GameState {
        GameState {
            defender: self.defender,
            picking_up: self.picking_up,
            trump: self.trump,
            plays: self.plays.clone(),
            covers: self.covers.clone(),
            hands: vec![self.hands[0].clone(), self.hands[1].clone()],
            defering: self.defering,
            start: self.start,
        }
    }
}

#[derive(Serialize, Deserialize, Debug)]
pub struct Memory {
    pub hands: Vec<Vec<Card>>,
    pub sizes: Vec<usize>,
    pub discard: Vec<Card>,
}

#[derive(Serialize, Deserialize, Debug)]
pub struct Game {
    pub key: usize,
    pub state: GameState,
    pub deck: Vec<Card>,
    pub memory: Memory,
    pub recording: Vec<String>,
    pub versus: String,
    pub winner: Option<usize>,
    pub joined: bool,
}

impl Game {
    pub fn new(key: usize, versus: String) -> Game {
        let deck = generate_deck();
        Game {
            key,
            state: GameState::new(deck[0], vec![deck[30..36].to_vec(), deck[24..30].to_vec()]),
            deck: deck[..24].to_vec(),
            memory: Memory {
                hands: vec![Vec::new(), Vec::new()],
                sizes: vec![6, 6],
                discard: Vec::new(),
            },
            recording: Vec::new(),
            versus,
            winner: None,
            joined: false,
        }
    }
    pub fn take_action(&mut self, action: &Action) {
        self.recording.push(serde_json::to_string(&action).unwrap());
        match action.verb {
            Verb::Play | Verb::Cover => {
                self.state.take_action(&action);
                remove_card(&mut self.memory.hands[action.player], action.card);
                self.memory.sizes[action.player] -= 1;
            },
            Verb::Reverse => {
                self.state.take_action(&action);
                remove_card(&mut self.memory.hands[action.player], action.card);
                self.memory.sizes[action.player] -= 1;
            },
            Verb::Pass => {
                if self.state.picking_up {
                    for &card in self.state.plays.iter().chain(self.state.covers.iter()) {
                        if card != UNK_CARD {
                            self.memory.hands[self.state.defender].push(card);
                            self.memory.sizes[self.state.defender] += 1;
                        }
                    }
                } else {
                    for &card in self.state.plays.iter().chain(self.state.covers.iter()) {
                        assert!(card != UNK_CARD);
                        self.memory.discard.push(card);
                    }
                }
                self.state.take_action(&action);
                self.deal(1-self.state.defender);
                self.deal(self.state.defender);
            },
            Verb::PickUp | Verb::Defer => {
                self.state.take_action(&action);
            },
        }
    }
    pub fn deal(&mut self, player: usize) {
        while self.deck.len() > 0 && self.state.hands[player].len() < 6 {
            self.state.hands[player].push(self.deck.pop().unwrap());
        }
        self.memory.sizes[player] = self.state.hands[player].len();
    }
    pub fn is_over(&self) -> bool {
        if self.deck.len() == 0 && self.state.hands.iter().any(|x| x.len() == 0) {
            return true
        }
        false
    }
}
