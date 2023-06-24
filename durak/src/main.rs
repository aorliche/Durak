mod rules;
mod search;

fn random_game() {
    let mut n = 0;
    let mut g = rules::Game::new(0, "computer".to_string());
    while !g.is_over() {
        n += 1;
        if n > 4 {
            break;
        }
        let a = g.state.random_action();
        g.take_action(&a);
        println!("{}", serde_json::to_string(&a).unwrap());
        println!("{}", serde_json::to_string(&g).unwrap());
        println!("{}", g.state.eval(&g.state, 0));
    }
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
    random_game();
}
