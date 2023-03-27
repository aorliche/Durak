
type Game struct {
    Players []*Player
    Board *Board
    Deck []*Card
    Discard []*Card
    Trump *Card
    Turn int
}

type Board struct {
    Plays []*Card
    Covers []*Card
}

type Player struct {
    Name string
    Hand []*Card
    Attacking bool
}

// TODO Remember player having cards (e.g. from pickup or process of elimination at the end of deck)
type Card struct {
    Visible bool
    Rank string
    Suit string
}

// No card matching, only predefined actions
type Action struct {
    Player *Player
    Mode string     // Attack Defend Pickup Pass (Reverse later)
    Card *Card
    Cover *Card
}

// Card suit weight i.e. card is trump
// Card ranks are just ordinals 1,2,3... for six, seven, eight...
type Weights struct {
    CardRankW float64
    CardSuitW float64
    BoardRankW float64
    BoardSuitW float64
    BoardSizeW float64
    OppHandSizeW float64
}

// Multiplies above weights
type MultWeights struct {
    DeckSizeW float64       
}

type Goal struct {
    Name string
    AttackW *Weights
    DefendW *Weights
    PickupW *Weights
    PassW *Weights
    // Feedback
    Feedback FeedbackFn
}

type FeedbackFn func(*Action, float64)

// Match card with respect to goal
// Get rid of small values, acquire trumps, get rid of cards in hand, acquire cards for opponent, get rid of trumps for opponent
// Goals change with time (goal weights change with game state... get rid of cards in hand at end of game)
// Search action space to see how goals affected
// 1. Hard coded goals, can be turned on or off
// 1a. Action search shallow
// 1b. Learning what actions help these goals: free to pick naive action but result of hand gives feedback signal to goals
// 2. Tell the computer it did a bad action (and tell it which goal it hindered)

func FeedbackFn(act *Action, goal *Goal, val float64) {

}
