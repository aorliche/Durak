package main

import (
    "testing"
)

func TestSearch(t *testing.T) {
    game := InitGame(0, false)
    state := InitGameState(game)
    state.Move(0, 1, 0)
}
