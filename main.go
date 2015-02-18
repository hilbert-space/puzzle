package main

import (
	"flag"
	"fmt"
	"math/rand"
	"runtime"
)

type result struct {
	C    []float64
	fill float64
}

func main() {
	var m, p, n, w int

	flag.IntVar(&m, "m", 20, "the first dimension")
	flag.IntVar(&p, "p", 2, "the second dimension")
	flag.IntVar(&n, "n", 10000, "the third dimension")
	flag.IntVar(&w, "w", runtime.NumCPU(), "the number of workers")
	flag.Parse()

	fmt.Printf("m = %d, p = %d, n = %d, w = %d\n", m, p, n, w)

	runtime.GOMAXPROCS(w)

	// Construct two matrices.
	A := make([]float64, m*p)
	B := make([]float64, p*n)
	for i := range A {
		A[i] = 42
	}
	for i := range B {
		B[i] = 42
	}

	// Compute the product of A and B.
	expectedC := make([]float64, m*n)
	multiply(A, B, expectedC, m, p, n)

	problem := make(chan result)

	for i := 0; i < w; i++ {
		go func() {
			for {
				C := make([]float64, m*n)

				// Fill in with a number specific to the current goroutine.
				fill := rand.Float64()
				for j := range C {
					C[j] = fill
				}

				// Multiply A by B and store the result in C.
				multiply(A, B, C, m, p, n)

				// Check the result against the expected answer.
				for j := range expectedC {
					if expectedC[j] != C[j] {
						problem <- result{C: C, fill: fill}
						break
					}
				}
			}
		}()
	}

	bad := <-problem

	// Sometimes the program reaches this point and reports that some of the
	// entries of bad.C have not been touched: they are equal to bad.fill.
	fmt.Printf("Fill: %.20e\n", bad.fill)
	for i := range expectedC {
		if expectedC[i] != bad.C[i] {
			fmt.Printf("%6d %30.20e %30.20e\n", i, expectedC[i], bad.C[i])
		}
	}
}

func multiply(A, B, C []float64, m, p, n int) {
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			sum := 0.0
			for k := 0; k < p; k++ {
				sum += A[k*m+i] * B[j*p+k]
			}
			C[j*m+i] = sum
		}
	}
}
