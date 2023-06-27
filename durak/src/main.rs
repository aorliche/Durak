mod rules;
mod search;

fn random_game() {
    let mut n = 0;
    let mut g = rules::Game::new(0, "computer".to_string());
    let orig_state = g.state.clone();
    while !g.is_over() {
        n += 1;
        if n > 4 {
            break;
        }
        let a = g.state.random_action();
        g.take_action(&a);
        println!("{}", serde_json::to_string(&a).unwrap());
        println!("{}", serde_json::to_string(&g).unwrap());
        println!("{}", orig_state.eval(&g.state, 1, false));
    }
}

fn random_game_search() {
    let mut g = rules::Game::new(0, "computer".to_string());
    let (c, r) = search::eval_node(&g.state.clone(), None, 0, 0, 0, false);
    println!("{}", r);
    println!("{}", serde_json::to_string(&c).unwrap());
    println!("{}", serde_json::to_string(&g.state).unwrap());
}

fn random_game_end_search() {
    let mut g = rules::Game::new(0, "computer".to_string());
    g.deck = Vec::new();
    let (c, r) = search::eval_node(&g.state.clone(), None, 0, 0, 0, true);
    println!("{}", r);
    println!("{}", serde_json::to_string(&c).unwrap());
    let mut cu = c.unwrap();
    g.take_action(&cu.last().unwrap());
    let acts = g.state.possible_actions();
    println!("{}", serde_json::to_string(&acts).unwrap());
    println!("{}", serde_json::to_string(&g.state).unwrap());
}

fn main() {
    /*let c = rules::card_from_suit_rank(1,2);
    println!("Hello, world!");
    println!("{}", rules::card_to_string(c));
    let a = rules::Action::new(0, rules::Verb::PickUp, 36, 36);
    let mut j = serde_json::to_string(&a).unwrap();
    println!("{}", j);
    let h = vec![Vec::new(), Vec::new()];
    let gs = rules::GameState::new(c, h);
    j = serde_json::to_string(&gs).unwrap();
    println!("{}", j);
    let d = rules::generate_deck();
    for card in d {
        println!("{}", rules::card_to_string(card));
    }
    let mut g = rules::Game::new(0, "computer".to_string());
    j = serde_json::to_string(&g).unwrap();
    println!("{}", j);
    let mut acts = g.state.attacker_actions(0);
    for act in &acts {
        println!("{}", serde_json::to_string(&act).unwrap());
    }
    g.take_action(&acts[0]);
    println!("{}", serde_json::to_string(&g).unwrap());
    acts = g.state.defender_actions(1);
    for act in &acts {
        println!("{}", serde_json::to_string(&act).unwrap());
    }
    g.take_action(&acts[0]);
    println!("{}", serde_json::to_string(&g).unwrap());*/
    //random_game();
    //random_game_search();
    random_game_end_search();
}
