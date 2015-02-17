package main

import (
	"crypto/md5"
	"flag"
	"fmt"
	"math/rand"
	"reflect"
	"runtime"
	"unsafe"
)

var workers = flag.Int("w", runtime.NumCPU(), "the number of workers")

type result struct {
	C    []float64
	fill float64
}

func main() {
	flag.Parse()
	runtime.GOMAXPROCS(*workers)

	// The checksum of the expected answer.
	checksum := computeChecksum(data.C)

	problem := make(chan result)

	for i := 0; i < *workers; i++ {
		go func() {
			for {
				C := make([]float64, data.m*data.n)

				// Fill in with a number specific to the current goroutine.
				fill := rand.Float64()
				for j := range C {
					C[j] = fill
				}

				// Multiply A by B and store the result in C.
				multiply(data.A, data.B, C, data.m, data.p, data.n)

				// Check the result against the expected answer.
				if checksum != computeChecksum(C) {
					problem <- result{C: C, fill: fill}
				}
			}
		}()
	}

	bad := <-problem

	// Sometimes the program reaches this point and reports that some of the
	// entries of bad.C have not been touched: they are equal to bad.fill.
	fmt.Printf("Fill: %.20e\n", bad.fill)
	for i := range data.C {
		if data.C[i] != bad.C[i] {
			fmt.Printf("%6d %30.20e %30.20e\n", i, data.C[i], bad.C[i])
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

func computeChecksum(data []float64) [16]byte {
	var bytes []byte

	header := (*reflect.SliceHeader)(unsafe.Pointer(&bytes))
	header.Data = (*reflect.SliceHeader)(unsafe.Pointer(&data)).Data
	header.Len = 8 * len(data)
	header.Cap = 8 * len(data)

	return md5.Sum(bytes)
}
