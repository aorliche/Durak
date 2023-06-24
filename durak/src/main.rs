mod durak;

fn random_game() {
    use rand::{thread_rng, seq::IteratorRandom};
    let mut rng = thread_rng();

    let mut n = 0;
    let mut g = durak::Game::new(0, "computer".to_string());
    while !g.is_over() {
        n += 1;
        if n > 100 {
            break;
        }
        let p0a = g.state.player_actions(0);
        let p1a = g.state.player_actions(1);
        let a = p0a.iter().chain(p1a.iter()).choose(&mut rng).unwrap();
        g.take_action(&a);
        println!("{}", serde_json::to_string(&g).unwrap());
    }

}

fn main() {
    /*let c = durak::card_from_suit_rank(1,2);
    println!("Hello, world!");
    println!("{}", durak::card_to_string(c));
    let a = durak::Action::new(0, durak::Verb::PickUp, 36, 36);
    let mut j = serde_json::to_string(&a).unwrap();
    println!("{}", j);
    let h = vec![Vec::new(), Vec::new()];
    let gs = durak::GameState::new(c, h);
    j = serde_json::to_string(&gs).unwrap();
    println!("{}", j);
    let d = durak::generate_deck();
    for card in d {
        println!("{}", durak::card_to_string(card));
    }
    let mut g = durak::Game::new(0, "computer".to_string());
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
