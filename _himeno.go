package main

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"sync"
	"time"
)

//プリプロセス
const MIMAX = 0
const MJMAX = 0
const MKMAX = 0

var (
	p                [MIMAX][MJMAX][MKMAX]float32
	a                [4][MIMAX][MJMAX][MKMAX]float32
	b                [3][MIMAX][MJMAX][MKMAX]float32
	c                [3][MIMAX][MJMAX][MKMAX]float32
	bnd              [MIMAX][MJMAX][MKMAX]float32
	wrk1             [MIMAX][MJMAX][MKMAX]float32
	wrk2             [MIMAX][MJMAX][MKMAX]float32
	imax, jmax, kmax int
	omega            float32
	concurrency      = runtime.NumCPU()
	copyConcurrency  = concurrency
	mainJobChan      = make(chan int, MIMAX)
	gosaChan         = make(chan float32, MIMAX)
	sumJobChan       = make(chan int, MIMAX)
	ws               = sync.WaitGroup{}
)

func init() {
	if len(os.Args) > 1 {
		if num, err := strconv.Atoi(os.Args[1]); err == nil && num > 0 {
			concurrency = num
		}
	}
	if len(os.Args) > 2 {
		if num, err := strconv.Atoi(os.Args[2]); err == nil && num > 0 {
			copyConcurrency = num
		}
	}
	fmt.Printf("Max Goroutine: %d\n", concurrency)

	for i := 0; i < concurrency; i++ {
		go JacobiMainWorker()
	}

	for i := 0; i < copyConcurrency; i++ {
		go JacobiSumWorker()
	}
}

func main() {
	var (
		nn                    int
		gosa                  float32
		cpu, cpu0, cpu1, flop float64
		target                = 60.0
	)
	imax = MIMAX - 1
	jmax = MJMAX - 1
	kmax = MKMAX - 1
	omega = 0.8

	initmt()
	fmt.Printf("mimax = %d mjmax = %d mkmax = %d\n", MIMAX, MJMAX, MKMAX)
	fmt.Printf("imax = %d jmax = %d kmax =%d\n", imax, jmax, kmax)

	nn = 3
	fmt.Printf(" Start rehearsal measurement process.\n")
	fmt.Printf(" Measure the performance in %d times.\n\n", nn)

	cpu0 = second()
	gosa = jacobi(nn)
	cpu1 = second()
	cpu = cpu1 - cpu0

	flop = fflop(imax, jmax, kmax)

	fmt.Printf(" MFLOPS: %f time(s): %f %e\n\n",
		mflops(nn, cpu, flop), cpu, gosa)

	nn = int(target / (cpu / 3.0))

	fmt.Printf(" Now, start the actual measurement process.\n")
	fmt.Printf(" The loop will be excuted in %d times\n", nn)
	fmt.Printf(" This will take about one minute.\n")
	fmt.Printf(" Wait for a while\n\n")

	/*
	* Start measuring
	 */
	cpu0 = second()
	gosa = jacobi(nn)
	cpu1 = second()

	cpu = cpu1 - cpu0

	fmt.Printf(" Loop executed for %d times\n", nn)
	fmt.Printf(" Gosa : %e \n", gosa)
	fmt.Printf(" MFLOPS measured : %f\tcpu : %f\n", mflops(nn, cpu, flop), cpu)
	fmt.Printf(" Score based on Pentium III 600MHz : %f\n",
		mflops(nn, cpu, flop)/82)
}

func initmt() {
	var i, j, k int

	for i = 0; i < imax; i++ {
		for j = 0; j < jmax; j++ {
			for k = 0; k < kmax; k++ {
				a[0][i][j][k] = 1.0
				a[1][i][j][k] = 1.0
				a[2][i][j][k] = 1.0
				a[3][i][j][k] = 1.0 / 6.0
				b[0][i][j][k] = 0.0
				b[1][i][j][k] = 0.0
				b[2][i][j][k] = 0.0
				c[0][i][j][k] = 1.0
				c[1][i][j][k] = 1.0
				c[2][i][j][k] = 1.0
				p[i][j][k] = float32(i*i) / float32((imax-1)*(imax-1))
				wrk1[i][j][k] = 0.0
				bnd[i][j][k] = 1.0
			}
		}
	}
}

func jacobi(nn int) float32 {
	var gosa float32

	for n := 1; n < nn+1; n++ {
		gosa = 0.0

		go func() {
			for i := 1; i < imax-1; i++ {
				mainJobChan <- i
			}
		}()
		for i := 1; i < imax-1; i++ {
			gosa += <-gosaChan
		}

		ws.Add(imax - 2)
		for i := 1; i < imax-1; i++ {
			sumJobChan <- i
		}
		ws.Wait()
	}

	return gosa
}

func fflop(mx, my, mz int) float64 {
	return float64(mz-2) * float64(my-2) * float64(mx-2) * 34.0
}

func mflops(nn int, cpu, flop float64) float64 {
	return flop / cpu * 1.e-6 * float64(nn)
}

var (
	baseTime = time.Time{}
)

func second() float64 {
	now := time.Now()

	if (baseTime == time.Time{}) {
		baseTime = now
		return 0.0
	} else {
		sub := now.Sub(baseTime)
		return float64(sub.Seconds())
	}
}

func JacobiMainWorker() {
	var i int
	for {
		i = <-mainJobChan

		var ssxss float32
		for j := 1; j < jmax-1; j++ {
			for k := 1; k < kmax-1; k++ {
				var s0, ss float32
				s0 = a[0][i][j][k]*p[i+1][j][k] +
					a[1][i][j][k]*p[i][j+1][k] +
					a[2][i][j][k]*p[i][j][k+1] +
					b[0][i][j][k]*(p[i+1][j+1][k]-p[i+1][j-1][k]-p[i-1][j+1][k]+p[i-1][j-1][k]) +
					b[1][i][j][k]*(p[i][j+1][k+1]-p[i][j-1][k+1]-p[i][j+1][k-1]+p[i][j-1][k-1]) +
					b[2][i][j][k]*(p[i+1][j][k+1]-p[i-1][j][k+1]-p[i+1][j][k-1]+p[i-1][j][k-1]) +
					c[0][i][j][k]*p[i-1][j][k] + c[1][i][j][k]*p[i][j-1][k] + c[2][i][j][k]*p[i][j][k-1] +
					wrk1[i][j][k]

				ss = (s0*a[3][i][j][k] - p[i][j][k]) * bnd[i][j][k]
				//fmt.Printf("%.16f\n", ssxss)

				ssxss += ss * ss
				/* gosa= (gosa > ss*ss) ? a : b; */
				wrk2[i][j][k] = p[i][j][k] + omega*ss
			}
		}
		gosaChan <- ssxss
	}
}

func JacobiSumWorker() {
	var i int
	for {
		i = <-sumJobChan

		for j := 1; j < jmax-1; j++ {
			for k := 1; k < kmax-1; k++ {
				p[i][j][k] = wrk2[i][j][k]
			}
		}
		ws.Done()
	}
}
