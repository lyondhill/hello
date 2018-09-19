package main

import (
	"fmt"
	"os"

	"github.com/wcharczuk/go-chart"
)

func main() {
	graph := chart.Chart{
		Series: []chart.Series{
			chart.ContinuousSeries{
				XValues: []float64{1.0, 2.0, 3.0, 4.0},
				YValues: []float64{1.0, 1.2, 3.1, 1.3},
			},
		},
	}

	file, err := os.Create("test.png")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	err = graph.Render(chart.PNG, file)
	fmt.Println(err)
}
