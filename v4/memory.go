package main

func InitMemory(np int) *Memory {
    mem := Memory{Hands: make([][]*Card, np), Sizes: make([]int, np), Discard: make([]*Card, 0)}
    for i,_ := range mem.Sizes {
        mem.Hands[i] = make([]*Card, 0)
        mem.Sizes[i] = 6
    }
    return &mem
}

func (mem *Memory) RemoveCard(p *Player, c *Card) {
    mem.Hands[p.Idx] = Remove(mem.Hands[p.Idx], c)
    mem.Sizes[p.Idx] -= 1
}

func (mem *Memory) AddCards(p *Player, cards []*Card) {
    mem.Hands[1-p.Idx] = append(mem.Hands[1-p.Idx], cards...)
    mem.Sizes[1-p.Idx] += len(cards)
}

func (mem *Memory) DiscardCards(cards []*Card) {
    mem.Discard = append(mem.Discard, cards...)
}

func (mem *Memory) SetSizes(players []*Player) {
    for i,p := range players {
        mem.Sizes[i] = len(p.Hand)
    }
}
