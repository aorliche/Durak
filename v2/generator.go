package main

import (
    "math/rand"
)

func RandomCard() *Object {
    rank := ranks[rand.Intn(len(ranks))]
    suit := suits[rand.Intn(len(suits))]    
    return Card(rank, suit)
}

func RandomGame() *Object {
    return MakeGameObject(RandomCard())
}

func Beats(a *Object, b *Object, game *Object) bool {
    c1 := a.Props["suit"].Val.(string) == b.Props["suit"].Val.(string) && 
        IndexOf(ranks, a.Props["rank"].Val.(string)) > IndexOf(ranks, b.Props["rank"].Val.(string))
    c2 := game.Props["trump"].Props["suit"].Val.(string) == a.Props["suit"].Val.(string) && 
        game.Props["trump"].Props["suit"].Val.(string) != b.Props["suit"].Val.(string)
    return c1 || c2
}

func MakeBeatsExample(n int, handSize int) ([][]*Object, []*Object, []*Object, [][]*Object) {
    inp := make([][]*Object, n)
    args := make([]*Object, n)
    games := make([]*Object, n)
    out := make([][]*Object, n)
    for i:=0; i<n; i++ {
        b := RandomCard()
        g := RandomGame()
        inp[i] = make([]*Object, handSize)
        args[i] = b
        games[i] = g
        out[i] = make([]*Object, 0)
        for j:=0; j<handSize; j++ {
            a := RandomCard()
            inp[i][j] = a
            if Beats(a, b, g) {
                out[i] = append(out[i], a)
            }
        }
        // Since solver generates unique output
        inp[i] = Unique(inp[i])
        out[i] = Unique(out[i])
    }
    return inp, args, games, out
}
