package main

import (
	"fmt"
	"math/big"
	"os"
	"strings"

	"github.com/scottsquatch/adventofcode2019/utils/fileutils"
)

type cardoperation func(*CardDeck)
type fastcardoperation func(index, size int) int
type fastinversecardoperation func(card, size *big.Int) *big.Int

func runPartA() {
	decksize := 10007
	cd := NewCardDeck(decksize)
	for _, line := range strings.Split(fileutils.ReadFile(os.Args[1]), "\n") {
		cardop := getcardoperation(strings.Trim(line, "\r"))
		cardop(cd)
	}
	for i, c := range cd.cards {
		if c == 2019 {
			fmt.Printf("Position of card 2019 is %d\n", i)
		}
	}
}

func getcardoperation(line string) cardoperation {
	var n int
	if r, _ := fmt.Sscanf(line, "cut %d", &n); r == 1 {
		return func(cd *CardDeck) { cd.Cut(n) }
	} else if r, _ := fmt.Sscanf(line, "deal with increment %d", &n); r == 1 {
		return func(cd *CardDeck) { cd.Deal(n) }
	} else if strings.EqualFold(line, "deal into new stack") {
		return func(cd *CardDeck) { cd.DealIntoNewStack() }
	}

	panic(line + " is not a valid card operation")
}

func getfastinversecardoperation(line string) fastinversecardoperation {
	var n int
	if r, _ := fmt.Sscanf(line, "cut %d", &n); r == 1 {
		return func(card, size *big.Int) *big.Int {
			return big.NewInt(0).Mod(big.NewInt(0).Add(big.NewInt(0).Add(card, big.NewInt(int64(n))), size), size)
		}
	} else if r, _ := fmt.Sscanf(line, "deal with increment %d", &n); r == 1 {
		return func(card, size *big.Int) *big.Int {
			return big.NewInt(0).Mod(big.NewInt(0).Mul(card, big.NewInt(0).ModInverse(big.NewInt(int64(n)), size)), size)
		}
	} else if strings.EqualFold(line, "deal into new stack") {
		return func(card, size *big.Int) *big.Int {
			return big.NewInt(0).Sub(big.NewInt(0).Sub(size, big.NewInt(1)), card)
		}
	}

	panic(line + " is not a valid card operation")
}

func getfastcardoperation(line string) fastcardoperation {
	var n int
	if r, _ := fmt.Sscanf(line, "cut %d", &n); r == 1 {
		return func(index, size int) int {
			ip := (index - n) % size
			if ip < 0 {
				ip += size
			}

			return ip
		}
	} else if r, _ := fmt.Sscanf(line, "deal with increment %d", &n); r == 1 {
		return func(index, size int) int {
			return (index * n) % size
		}
	} else if strings.EqualFold(line, "deal into new stack") {
		return func(index, size int) int {
			ip := (-1 - index) % size
			if ip < 0 {
				ip += size
			}

			return ip
		}
	}

	panic(line + " is not a valid card operation")
}

func runPartB() {
	// Using method from https://www.reddit.com/r/adventofcode/comments/ee0rqi/2019_day_22_solutions/fbnifwk/
	decksize := big.NewInt(119315717514047)
	inverseoperations := []fastinversecardoperation{}
	for _, line := range strings.Split(fileutils.ReadFile(os.Args[1]), "\n") {
		op := getfastinversecardoperation(strings.Trim(line, "\r"))
		inverseoperations = append([]fastinversecardoperation{op}, inverseoperations...)

	}

	pos := big.NewInt(2020)
	card := big.NewInt(2020)
	for _, op := range inverseoperations {
		card = op(card, decksize)
	}
	cardofcard := card
	for _, op := range inverseoperations {
		cardofcard = op(cardofcard, decksize)
	}

	A := big.NewInt(0).Mod(big.NewInt(0).Mul(big.NewInt(0).Sub(card, cardofcard), big.NewInt(0).ModInverse(big.NewInt(0).Add(big.NewInt(0).Sub(pos, card), decksize), decksize)), decksize)
	B := big.NewInt(0).Mod(big.NewInt(0).Sub(card, big.NewInt(0).Mul(A, pos)), decksize)

	n := big.NewInt(101741582076661)
	An := big.NewInt(0).Exp(A, n, decksize)
	An2 := big.NewInt(0).Sub(An, big.NewInt(1))
	first := big.NewInt(0).Mul(An, pos)

	mi := big.NewInt(0).ModInverse(big.NewInt(0).Sub(A, big.NewInt(1)), decksize)
	second := big.NewInt(0).Mul(big.NewInt(0).Mul(An2, mi), B)

	result := big.NewInt(0).Mod(big.NewInt(0).Add(first, second), decksize)

	fmt.Printf("Value of index 2020 is %v\n", result)
}

func main() {
	runPartA()
	runPartB()
}
