package statstest_test

import (
	"fmt"

	"github.com/motah-fard/statstest"
)

func ExampleAdjustPValues() {
	p := []float64{0.001, 0.01, 0.03, 0.2}

	adj, err := statstest.AdjustPValues(p, statstest.BenjaminiHochberg)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	fmt.Printf("%.3f %.3f %.3f %.3f\n", adj[0], adj[1], adj[2], adj[3])

	// Output:
	// 0.004 0.020 0.040 0.200
}
