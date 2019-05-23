package main

import (
	"errors"
	"fmt"
	"math"
)

func main() {
	fmt.Println("Welcome to Lauren's years of things cheat program ")
	fmt.Println("press ctrl+c to quit")

	var startYear int
	var stopYear int

	var err error
	for {
		fmt.Println("please insert a starting year:")
		startYear, err = readYear()
		if err != nil {
			fmt.Println("invalid year: try '1232bc'")
			continue
		}
		break
	}

	for {
		fmt.Println("Please insert a end year:")
		stopYear, err = readYear()
		if err != nil {
			fmt.Println("invalid year: try '1232ad'")
			continue
		}
		break
	}

	totalYears := math.Abs(float64(startYear - stopYear))
	fmt.Println(totalYears, "---")

	for {
		// read
		fmt.Printf("Year:")
		year, err := readYear()
		if err != nil {
			fmt.Println("invalid year: try '1232ad'")
			continue
		}

		// parse
		floatYear := math.Abs(float64(startYear - year))

		// quit
		fmt.Printf("%.4f/100\n", (floatYear/totalYears)*100)

		// output
	}
}

func readYear() (int, error) {
	var yearNumber int
	var yearStamp string
	n, err := fmt.Scanf("%d%s\n", &yearNumber, &yearStamp)
	if err != nil {
		return n, err
	}

	switch yearStamp {
	case "ad":
		return yearNumber, nil
	case "bc":
		return -yearNumber, nil
	}
	return n, errors.New("bad yearstamp")
}
