package main

import (
	"crypto/md5"
	"flag"
	"fmt"
	"reflect"
	"runtime"
	"unsafe"
)

var workers = flag.Int("w", runtime.NumCPU(), "the number of workers")

func main() {
	if data.checksum != computeChecksum(data.C) {
		panic("the checksum is incorrect")
	}

	flag.Parse()

	runtime.GOMAXPROCS(*workers)

	problem := make(chan []float64)

	for i := 0; i < *workers; i++ {
		go func() {
			for {
				someC := multiply(data.A, data.B, data.m, data.p, data.n)
				if data.checksum != computeChecksum(someC) {
					problem <- someC
				}
			}
		}()
	}

	badC := <-problem

	for i := range data.C {
		if badC[i] != data.C[i] {
			fmt.Printf("%6d %30.20e %30.20e\n", i, data.C[i], badC[i])
		}
	}
}

func multiply(A, B []float64, m, p, n int) []float64 {
	C := make([]float64, m*n)

	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			sum := 0.0
			for k := 0; k < p; k++ {
				sum += A[k*m+i] * B[j*p+k]
			}
			C[j*m+i] = sum
		}
	}

	return C
}

func computeChecksum(data []float64) [16]byte {
	var bytes []byte

	header := (*reflect.SliceHeader)(unsafe.Pointer(&bytes))
	header.Data = (*reflect.SliceHeader)(unsafe.Pointer(&data)).Data
	header.Len = 8 * len(data)
	header.Cap = 8 * len(data)

	return md5.Sum(bytes)
}
