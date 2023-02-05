package main

import (
    "fmt"
)

func makeGame(trump *object) *object {
    g := makeObject("game")
    g.setProp("trump", "*object", trump)
    return g
}

func makePlayer(name string) *object {
    p := makeObject("player")
    p.setProp("hand", "[]*object", make([]*object, 0))
    return p
}

func makeCard(rank string, suit string) *object {
    c := makeObject("card")
    c.setProp("rank", "string", rank)
    c.setProp("suit", "string", suit)
    return c
}

var ranks = []string{"6", "7", "8", "9", "10", "Jack", "Queen", "King", "Ace"}

func cardStr(c *object) string {
    return fmt.Sprintf("%s of %s", c.props["rank"].(string), c.props["suit"].(string))
}
