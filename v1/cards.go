package main

import (
    "fmt"
)

func MakeGame(trump *Object) *Object {
    g := MakeObject("game")
    g.SetProp("trump", "card", trump)
    g.SetProp("board", "[]card", make([]*Object, 0))
    return g
}

func MakePlayer(name string) *Object {
    p := MakeObject("player")
    p.SetProp("hand", "[]card", make([]*Object, 0))
    return p
}

func MakeCard(rank string, suit string) *Object {
    c := MakeObject("card")
    c.SetProp("rank", "string", rank)
    c.SetProp("suit", "string", suit)
    return c
}

var ranks = []string{"6", "7", "8", "9", "10", "Jack", "Queen", "King", "Ace"}

func CardStr(c *Object) string {
    return fmt.Sprintf("%s of %s", c.Props["rank"].(string), c.Props["suit"].(string))
}
