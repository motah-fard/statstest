package statstest_test

import (
	"fmt"

	"github.com/motah-fard/statstest"
)

func ExampleShapiroWilk() {
	x := []float64{5.1, 4.9, 5.3, 5.0, 4.8, 5.2, 5.4, 4.7, 5.1, 5.0}

	res, err := statstest.ShapiroWilk(x)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	fmt.Printf("W = %.3f, p = %.4f\n", res.W, res.PValue)

	// Output:
	// W = 0.984, p = 0.9829
}
