mod durak;

fn main() {
    let c = durak::card_from_suit_rank(1,2);
    println!("Hello, world!");
    println!("{}", durak::card_to_string(c));
    let a = durak::Action::new(0, durak::Verb::PickUp, 36, 36);
    let mut j = serde_json::to_string(&a).unwrap();
    println!("{}", j);
    let h = vec![Vec::new(), Vec::new()];
    let gs = durak::GameState::new(c, h);
    j = serde_json::to_string(&gs).unwrap();
    println!("{}", j);
}
