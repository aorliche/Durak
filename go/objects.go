package main

import "fmt"

type object struct {
    typ string
    props map[string] interface{}
}

/*func (o object) valtyp() string {
    switch o.typ {
        case "rank": 
        case "suit": 
            return "string"
    }
    return "none"
}*/

func appendList(obj *object, prop string, item *object) {
    obj.props[prop] = append(obj.props[prop].([]*object), item)
}

func makeObject(typ string) *object {
    obj := object{typ: typ, props: make(map[string]interface{})}
    return &obj
}

func makePlayer(name string) *object {
    p := makeObject("player")
    p.props["hand"] = make([]*object, 0)
    return p
}

func makeCard(rank string, suit string) *object {
    c := makeObject("card")
    c.props["rank"] = rank
    c.props["suit"] = suit
    return c
}

func getListItem(obj *object, prop string, idx int) *object {
    return obj.props[prop].([]*object)[idx]
}

func main() {
    c := makeCard("8", "Hearts")
    p := makePlayer("Guy")
    appendList(p, "hand", c)
    c = getListItem(p, "hand", 0)
    fmt.Println(c.props["rank"])
}
