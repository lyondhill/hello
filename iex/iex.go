package main

import (
	"fmt"
	"net/http"

	"github.com/timpalpant/go-iex"
)

func main() {
	client := iex.NewClient(&http.Client{})

	quotes, err := client.GetTOPS(nil)
	if err != nil {
		panic(err)
	}

	for _, quote := range quotes {
		fmt.Printf("%v: bid $%.02f (%v shares), ask $%.02f (%v shares) [as of %v]\n",
			quote.Symbol, quote.BidPrice, quote.BidSize,
			quote.AskPrice, quote.AskSize, quote.LastUpdated)
	}
	fmt.Println("tops: %d", len(quotes))
}
