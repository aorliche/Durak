
import (
    "reflect"
)

func AttackerActions(game *Game) []*Action {
    res := make([]*Action, 0)
    p := game.GetAttacker()
    if game.Turn == 0 {
        for _,card := range p.Cards {
            act := Action{Player: p, Mode: "Attack", Card: &card}
            res = append(res, &act)
        }
    } else {
        for _,bc := range Union(game.Board.Plays, game.Board.Covers) {
            for _,pc := range p.Cards {
                if bc.Rank == pc.Rank && IndexOf(res, func(act) {act.Card == pc}) == -1 {
                    act := Action{Player: p, Mode: "Attack", Card: &pc}
                    res = append(res, &act)
                }
            }
        }
    }
    if len(board.Plays) > 0 {
        act := Action{Player: p, Mode: "Pass"}
        res = append(res, &act)
    }
    return res
}

func DefenderActions(game *Game) []*Action {
    res := make([]*Action, 0)
    p := game.GetDefender()
    for _,bp := range game.Board.Plays {
        for _,pc := range p.Cards {
            if pc.Beats(bp, game.Trump) {
                act := Action{Player: p, Mode: "Defend", Card: &pc, Cover: &bp}
                res = append(res, &act)
            }
        }
    }
    if len(board.Plays) > 0 && len(board.Covers) < len(board.Plays) {
        act := Action{Player: p, Mode: "Pickup"}
        res = append(res, &act)
    }
    return res
}

func (w *Weights) Update(a *Weights, m *MultWeights, y float64, lr float64) {
    ww := reflect.ValueOf(w)
    aa := reflect.ValueOf(a)
    yhat := 0
    for i := 0; i < ww.NumField(); i++ {
        yhat += aa.Field(i).Interface().(float64)
    }
    err := lr*(y-yhat)
    for i := 0; i < ww.NumField(); i++ {
        b := aa.Field(i).Interface().(float64)
        c := ww.Field(i).Interface().(float64)
        d := c+err*b
        ww.Field(i).SetFloat(d)
    }
}
