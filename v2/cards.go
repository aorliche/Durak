package main

import (
    "fmt"
)

var suits = []string{"Clubs", "Spades", "Hearts", "Diamonds"}
var ranks = []string{"6", "7", "8", "9", "10", "Jack", "Queen", "King", "Ace"}

func MakeCardObject(rank string, suit string) *Object {
    obj := Object{Type: "Card", Props: make(map[string]*Object)}
    obj.Props["rank"] = MakeStringObject(rank)
    obj.Props["suit"] = MakeStringObject(suit)
    return &obj
}

func MakePlayerObject(name string, hand *Object) *Object {
    obj := Object{Type: "Player", Props: make(map[string]*Object)}
    obj.Props["name"] = MakeStringObject(name)
    obj.Props["hand"] = hand
    return &obj
}

func MakeGameObject(trump *Object) *Object {
    obj := Object{Type: "Game", Props: make(map[string]*Object)}
    obj.Props["trump"] = trump
    return &obj
}

func CardStr(card *Object) string {
    rank := card.Props["rank"].Val.(string)
    suit := card.Props["suit"].Val.(string)
    str := fmt.Sprintf("%s of %s", rank, suit)
    return str
}

func PlayerStr(p *Object) string {
    name := p.Props["name"].Val.(string)
    hand := p.Props["hand"].ToStr()
    str := fmt.Sprintf("Player \"%s\" Hand: %s", name, hand)
    return str
}

func GameStr(game *Object) string {
    str := fmt.Sprintf("Game Trump: %s", game.Props["trump"].ToStr())
    return str
}
