package main

// CardDeck represents a deck of cards
type CardDeck struct {
	size  int
	cards []int
}

// NewCardDeck creates a new deck of cards with the given sizes
func NewCardDeck(size int) *CardDeck {
	cards := make([]int, size)
	for i := 0; i < size; i++ {
		cards[i] = i
	}

	return &CardDeck{size: size, cards: cards}
}

// DealIntoNewStack deal into a new card deck and return
func (cd *CardDeck) DealIntoNewStack() {
	newcards := []int{}
	for _, c := range cd.cards {
		newcards = append([]int{c}, newcards...)
	}

	copy(cd.cards, newcards)
}

// Cut cuts given number of cards
func (cd *CardDeck) Cut(n int) {
	if n < 0 {
		n = len(cd.cards) + n
	}
	cd.cards = append(cd.cards[n:], cd.cards[0:n]...)
}

// Deal card with given increment
func (cd *CardDeck) Deal(increment int) {
	board := make([]int, cd.size)
	for i, j := 0, 0; i < cd.size; i, j = i+1, (j+increment)%cd.size {
		board[j] = cd.cards[i]
	}

	copy(cd.cards, board)
}
