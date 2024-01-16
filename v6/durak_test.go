package durak

import (
    "encoding/json"
    "log"
    "testing"
    "time"
)

func TestBeats(t *testing.T) {
    if Card(10).Beats(Card(11), Card(20)) {
        t.Errorf("Card(10).Beats(Card(11), Card(20))")
    }
    if Card(4).Beats(Card(17), Card(11)) {
        t.Errorf("Card(4).Beats(Card(17), Card(1))")
    }
}

func TestInitGame(t *testing.T) {
    game := InitGame(0, []string{"Human", "Medium"})
    if game == nil {
        t.Errorf("InitGame failed")
    }
}

func TestGetActions(t *testing.T) {
    game := InitGame(0, []string{"Easy", "Medium"})
    acts0 := game.State.PlayerActions(0)
    acts1 := game.State.PlayerActions(1)
    if len(acts0) == 0 {
        t.Errorf("No actions for player 0")
    }
    if len(acts1) != 0 {
        t.Errorf("Actions for player 1")
    }
}

func TestTakeAction(t *testing.T) {
    game := InitGame(1, []string{"Human", "Human", "Medium"})
    game.TakeAction(game.State.RandomAction())
    acts1 := game.State.PlayerActions(1)
    if len(acts1) == 0 {
        t.Errorf("No actions for player 1")
    }
}

func TestSearchStart(t *testing.T) {
    game := InitGame(0, []string{"Human", "Medium"})
    c, _ := game.State.EvalNode(0, nil)
    if len(c) == 0 {
        t.Errorf("No action chain for search")
    }
}

func TestSearchEnd(t *testing.T) {
    game := InitGame(0, []string{"Human", "Medium"})
    game.Deck = make([]Card, 0)
    c, _ := game.State.EvalNode(0, nil)
    if len(c) == 0 {
        t.Errorf("No action chain for search")
    }
}

/*func TestSearchEndSimple1(t *testing.T) {
    game := InitGame(0, []string{"Human", "Medium"})
    game.Deck = make([]Card, 0)
    game.State.Trump = 1
    game.State.Hands[1] = []Card{Card(1), Card(2)}
    game.State.Hands[0] = []Card{Card(5), Card(6)}
    game.MaskUnknownCards(0)
    c, _ := game.State.EvalNode(nil, 0, 0, 0, len(game.Deck))
    if len(c) == 0 {
        t.Errorf("No action chain for search")
    }
}*/

func TestMaskUnkownCardStart(t *testing.T) {
    game := InitGame(0, []string{"Human", "Medium"})
    state := game.MaskUnknownCards(0)
    for _, card := range state.Hands[1] {
        if card != UNK_CARD {
            t.Errorf("Card not UNK_CARD")
        }
    }
    for _, card := range state.Hands[0] {
        if card == UNK_CARD {
            t.Errorf("Card UNK_CARD")
        }
    }
}

func TestMaskUnknownCard_WithKnown(t *testing.T) {
    game := InitGame(0, []string{"Human", "Medium"})
    game.Memory.Hands[0] = []Card{game.State.Hands[0][0], game.State.Hands[0][1]}
    state := game.MaskUnknownCards(1)
    nKnown := 0
    for _, card := range state.Hands[0] {
        if card != UNK_CARD {
            nKnown++
        }
    }
    if nKnown != 2 {
        t.Errorf("Wrong number of known cards")
    }
}

/*func TestThreeEasyComputerRandomGame(t *testing.T) {
    game := InitGame(0, []string{"Easy", "Easy", "Easy"})
    game.StartComputer("Easy", 0)
    game.StartComputer("Easy", 1)
    game.StartComputer("Easy", 2)
    for !game.CheckGameOver() {
        time.Sleep(100 * time.Millisecond)
        acts := make([]int, 3)
        for i := 0; i < len(acts); i++ {
            acts[i] = len(game.State.PlayerActions(i))
        }
        log.Println(len(game.Deck), acts, game.State.ToStr(), game.Recording.Winners)
    }
    jsn, _ := json.Marshal(game.Recording)
    log.Println(string(jsn))
}*/

func TestThreeMediumComputerRandomGame(t *testing.T) {
    log.Println("TestThreeMediumComputerRandomGame")
    game := InitGame(0, []string{"Medium", "Medium", "Medium"})
    game.StartComputer("Medium", 0, nil, nil, nil)
    game.StartComputer("Medium", 1, nil, nil, nil)
    game.StartComputer("Medium", 2, nil, nil, nil)
    for !game.CheckGameOver() {
        time.Sleep(1000 * time.Millisecond)
        acts := make([]int, 3)
        for i := 0; i < len(acts); i++ {
            acts[i] = len(game.State.PlayerActions(i))
        }
        log.Println(len(game.Deck), acts, game.State.ToStr(), game.Recording.Winners)
    }
    jsn, _ := json.Marshal(game.Recording)
    log.Println(string(jsn))
}

/*func TestThreeMediumComputerOneWon(t *testing.T) {
    game := InitGame(0, []string{"Medium", "Medium", "Medium"})
    game.Deck = make([]Card, 0)
    game.Recording.Winners = []int{0}
    game.State.Hands[0] = make([]Card, 0)
    game.State.Won[0] = true
    game.State.Attacker = 1
    game.State.Defender = 2
    game.StartComputer("Medium", 0)
    game.StartComputer("Medium", 1)
    game.StartComputer("Medium", 2)
    for !game.CheckGameOver() {
        time.Sleep(1000 * time.Millisecond)
        acts := make([]int, 3)
        for i := 0; i < len(acts); i++ {
            acts[i] = len(game.State.PlayerActions(i))
        }
        log.Println(len(game.Deck), acts, game.State.ToStr(), game.Recording.Winners)
    }
    jsn, _ := json.Marshal(game.Recording)
    log.Println(string(jsn))
}*/

/*func TestTwoMediumComputerRandomGame(t *testing.T) {
    game := InitGame(0, []string{"Medium", "Medium"})
    game.StartComputer("Medium", 0)
    game.StartComputer("Medium", 1)
    for !game.CheckGameOver() {
        time.Sleep(1000 * time.Millisecond)
        acts := make([]int, 2)
        for i := 0; i < len(acts); i++ {
            acts[i] = len(game.State.PlayerActions(i))
        }
        log.Println(len(game.Deck), acts, game.State.ToStr())
    }
    jsn, _ := json.Marshal(game.Recording)
    log.Println(string(jsn))
}*/

/*func TestPlayerTwoMediumAttacker(t *testing.T) {
    game := InitGame(0, []string{"Medium", "Medium"})
    game.State.Attacker = 1
    game.State.Defender = 0
    c,r := game.State.EvalNode(nil, 1, 0, 0, len(game.Deck))
    log.Println(c, r)
}*/
