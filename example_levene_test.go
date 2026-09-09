package statstest_test

import (
	"fmt"

	"github.com/motah-fard/statstest"
)

func ExampleLevenesTest() {
	g1 := []float64{8.1, 8.3, 7.9, 8.0, 8.2, 8.5}
	g2 := []float64{7.5, 9.0, 6.8, 9.5, 7.0, 8.8}
	g3 := []float64{8.0, 8.1, 7.9, 8.0, 8.1, 7.95}

	res, err := statstest.LevenesTest(g1, g2, g3)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	fmt.Printf("W = %.3f, p = %.6f\n", res.Statistic, res.PValue)

	// Output:
	// W = 38.825, p = 0.000001
}
