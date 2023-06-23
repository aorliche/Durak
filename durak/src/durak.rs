use serde::{Deserialize, Serialize};
use serde_json::Value;

#[derive(Serialize, Deserialize, Debug)]
pub enum Verb {
    Play = 1,
    Cover,
    Reverse,
    Pass,
    PickUp,
    Defer,
}

type Card = usize;
const unknown_card : Card = 36;

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

#[derive(Serialize, Deserialize, Debug)]
pub struct Action {
    player: usize,
    verb: Verb,
    card: Card,
    covering: Card,
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
    defender: usize,
    picking_up: bool,
    trump: Card,
    deck_size: usize,
    plays: Vec<Card>,
    covers: Vec<Card>,
    hands: Vec<Vec<Card>>,
    defering: bool,
    #[serde(skip)]
    start: Option<time::Instant>,
}

pub fn num_not_nil(v: &Vec<Card>) -> usize {
    v.iter().filter(|&x| *x != unknown_card).count()
}

impl GameState {
    pub fn new(trump: Card, hands: Vec<Vec<Card>>) -> GameState {
        GameState {
            defender: 1,
            picking_up: false,
            trump,
            deck_size: 24,
            plays: Vec::new(),
            covers: Vec::new(),
            hands,
            defering: false,
            start: None,
        } 
    }

    pub fn attacker_actions(&self, pidx: usize) -> Vec<Action> {
        let mut res = Vec::new();
        if self.defering {
            return res
        }
        if self.plays.len() == 0 {
           for &card in &self.hands[pidx] {
               res.push(Action::new(pidx, Verb::Play, card, unknown_card))
           } 
        } else {
            for &card in &self.hands[pidx] {
                for &board_card in self.plays.iter().chain(self.covers.iter()) {
                    if board_card == unknown_card {
                        continue;
                    }
                    // Card equality with unknown card for AI search
                    if card_rank(card) == card_rank(board_card) || card == unknown_card {
                        res.push(Action::new(pidx, Verb::Play, card, unknown_card));
                        break;
                    }
                }
            }
        }
        if self.picking_up || (num_not_nil(&self.covers) == self.plays.len() && self.plays.len() > 0) {
            res.push(Action::new(pidx, Verb::Pass, unknown_card, unknown_card))
        }
        // For AI to not throw trumps away
        if self.plays.len() > num_not_nil(&self.covers) {
            res.push(Action::new(pidx, Verb::Defer, unknown_card, unknown_card)) 
        }
        res
    }
}

/*pub struct Game {

}*/
