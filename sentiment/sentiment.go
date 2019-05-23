package main

import (
	"os"
	"fmt"
	"strings"

	sent "gopkg.in/vmarkovtsev/BiDiSentiment.v1"
	"github.com/cdipaolo/sentiment"
)

func main() {
	model, err := sentiment.Restore()
	if err != nil {
		panic(err)
	}

	stuff := strings.Join(os.Args[1:], " ")

	analysis := model.SentimentAnalysis(stuff, sentiment.English)

	fmt.Printf("1: %+v\n", analysis)

	session, _ := sent.OpenSession()
  defer session.Close()
  result, _ := sent.Evaluate(
    []string{stuff},
    session)
  fmt.Println("2:", (result[0])
}