package durak

func (game *Game) MaskUnknownCards(me int) *GameState {
    state := game.State.Clone()
    for opp := 0; opp < len(game.State.Hands); opp++ {
        if opp == me {
            continue
        }
        for i, card := range game.State.Hands[opp] {
            if !Contains(card, game.Memory.Hands[opp]) {
                state.Hands[opp][i] = UNK_CARD
            }
        }
    }
    return state
}
