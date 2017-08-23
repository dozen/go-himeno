package main

import (
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"
)

//プリプロセス
const MIMAX = 0
const MJMAX = 0
const MKMAX = 0

var (
	p                = make([][][]float32, MIMAX, MIMAX)
	a                = make([][][][]float32, 4, 4)
	b                = make([][][][]float32, 3, 3)
	c                = make([][][][]float32, 3, 3)
	bnd              = make([][][]float32, MIMAX, MIMAX)
	wrk1             = make([][][]float32, MIMAX, MIMAX)
	wrk2             = make([][][]float32, MIMAX, MIMAX)
	imax, jmax, kmax int
	omega            float32
	concurrency      = 8
)

func init() {
	if len(os.Args) > 1 {
		if num, err := strconv.Atoi(os.Args[1]); err == nil && num > 0 {
			concurrency = num
		}
	}
	fmt.Printf("Max Goroutine: %d\n", concurrency)
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

	for i = 0; i < 4; i++ {
		a[i] = make([][][]float32, MIMAX, MIMAX)
	}

	for i = 0; i < 3; i++ {
		b[i] = make([][][]float32, MIMAX, MIMAX)
		c[i] = make([][][]float32, MIMAX, MIMAX)
	}

	for i = 0; i < MIMAX; i++ {
		a[0][i] = make([][]float32, MJMAX, MJMAX)
		a[1][i] = make([][]float32, MJMAX, MJMAX)
		a[2][i] = make([][]float32, MJMAX, MJMAX)
		a[3][i] = make([][]float32, MJMAX, MJMAX)
		b[0][i] = make([][]float32, MJMAX, MJMAX)
		b[1][i] = make([][]float32, MJMAX, MJMAX)
		b[2][i] = make([][]float32, MJMAX, MJMAX)
		c[0][i] = make([][]float32, MJMAX, MJMAX)
		c[1][i] = make([][]float32, MJMAX, MJMAX)
		c[2][i] = make([][]float32, MJMAX, MJMAX)
		p[i] = make([][]float32, MJMAX, MJMAX)
		wrk1[i] = make([][]float32, MJMAX, MJMAX)
		bnd[i] = make([][]float32, MJMAX, MJMAX)
		wrk2[i] = make([][]float32, MJMAX, MJMAX)
		for j = 0; j < MJMAX; j++ {
			a[0][i][j] = make([]float32, MKMAX, MKMAX)
			a[1][i][j] = make([]float32, MKMAX, MKMAX)
			a[2][i][j] = make([]float32, MKMAX, MKMAX)
			a[3][i][j] = make([]float32, MKMAX, MKMAX)
			b[0][i][j] = make([]float32, MKMAX, MKMAX)
			b[1][i][j] = make([]float32, MKMAX, MKMAX)
			b[2][i][j] = make([]float32, MKMAX, MKMAX)
			c[0][i][j] = make([]float32, MKMAX, MKMAX)
			c[1][i][j] = make([]float32, MKMAX, MKMAX)
			c[2][i][j] = make([]float32, MKMAX, MKMAX)
			p[i][j] = make([]float32, MKMAX, MKMAX)
			wrk1[i][j] = make([]float32, MKMAX, MKMAX)
			bnd[i][j] = make([]float32, MKMAX, MKMAX)
			wrk2[i][j] = make([]float32, MKMAX, MKMAX)

			for k = 0; k < MKMAX; k++ {
				a[0][i][j][k] = 0.0
				a[1][i][j][k] = 0.0
				a[2][i][j][k] = 0.0
				a[3][i][j][k] = 0.0
				b[0][i][j][k] = 0.0
				b[1][i][j][k] = 0.0
				b[2][i][j][k] = 0.0
				c[0][i][j][k] = 0.0
				c[1][i][j][k] = 0.0
				c[2][i][j][k] = 0.0
				p[i][j][k] = 0.0
				wrk1[i][j][k] = 0.0
				bnd[i][j][k] = 0.0
			}
		}
	}

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

		lock := sync.Mutex{}
		semaphore := make(chan struct{}, concurrency)
		ws := sync.WaitGroup{}

		for i := 1; i < imax-1; i++ {
			semaphore <- struct{}{}
			ws.Add(1)
			go func(i int) {
				defer func() {
					ws.Done()
				}()

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
				<-semaphore
				lock.Lock()
				//fmt.Printf("%.16f\n", ssxss)
				gosa += ssxss
				lock.Unlock()
			}(i)

		}
		ws.Wait()

		for i := 1; i < imax-1; i++ {
			for j := 1; j < jmax-1; j++ {
				for k := 1; k < kmax-1; k++ {
					p[i][j][k] = wrk2[i][j][k]
				}
			}
		}
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
