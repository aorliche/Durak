package main

import (
    
)

func (game *Game) MaskUnknownCards(me int) *GameState {
    opp := 1 - me
    state := game.State.Clone()
    for i, card := range game.State.Hands[opp] {
        if !Contains(card, game.Memory.Hands[opp]) {
            state.Hands[opp][i] = UNK_CARD
        }
    }
    return state
}
