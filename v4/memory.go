package main

func InitMemory(np int) *Memory {
    mem := Memory{Hands: make([][]*Card, np), Sizes: make([]int, np), Discard: make([]*Card, 0)}
    for i,_ := range mem.Sizes {
        mem.Hands[i] = make([]*Card, 0)
        mem.Sizes[i] = 6
    }
    return &mem
}

func (mem *Memory) RemoveCard(pIdx int, c *Card) {
    mem.Hands[pIdx] = Remove(mem.Hands[pIdx], c)
    mem.Sizes[pIdx] -= 1
}

func (mem *Memory) AddCards(pIdx int, cards []*Card) {
    mem.Hands[pIdx] = append(mem.Hands[pIdx], cards...)
    mem.Sizes[pIdx] += len(cards)
}

func (mem *Memory) DiscardCards(cards []*Card) {
    mem.Discard = append(mem.Discard, cards...)
}

func (mem *Memory) SetSizes(players []*Player) {
    for i,p := range players {
        mem.Sizes[i] = len(p.Hand)
    }
}

func (mem *Memory) SetKnownCards(state *GameState, me int, opp int) {
    if state.DeckSize <= 1 {
        state.Hands[opp] = mem.GuessFinalCards(state, me)
    } else {
        for i,mc := range mem.Hands[opp] {
            c := state.Hands[opp][i] 
            c.Rank = mc.Rank
            c.Suit = mc.Suit
        }
    }
}

func (mem *Memory) GuessFinalCards(state *GameState, me int) []*Card {
    cards := InitDeck()
    notit := make([]bool, len(cards))
    oppCards := make([]*Card, 0)
    for i,c := range cards {
        for _,mc := range mem.Discard {
            if c.Rank == mc.Rank && c.Suit == mc.Suit {
                notit[i] = true
                break
            }
        }
        for _,mc := range state.Hands[me] {
            if c.Rank == mc.Rank && c.Suit == mc.Suit {
                notit[i] = true
                break
            }
        }
        if state.DeckSize == 1 {
            if c.Rank == state.TrumpRank && c.Suit == state.Trump {
                notit[i] = true
            }
        }
    }
    for i,val := range notit {
        if !val {
            oppCards = append(oppCards, cards[i])
        }
    }
    return oppCards
}
