package statstest_test

import (
	"fmt"

	"github.com/motah-fard/statstest"
)

func ExampleKruskalWallis() {
	g1 := []float64{2.9, 3.0, 2.5, 2.6, 3.2}
	g2 := []float64{3.8, 2.7, 4.0, 2.4, 3.9}
	g3 := []float64{2.8, 3.4, 3.7, 2.2, 2.0}

	res, err := statstest.KruskalWallis(g1, g2, g3)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	fmt.Printf("H = %.3f, p = %.4f\n", res.H, res.PValue)

	// Output:
	// H = 1.860, p = 0.3946
}
