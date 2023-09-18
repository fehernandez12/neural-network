package images

import (
	"image"
	"runtime"
	"sync"
)

func Invert(img image.Image) *image.NRGBA {
	src := NewScanner(img)
	dst := image.NewNRGBA(image.Rect(0, 0, src.W(), src.H()))
	Parallel(0, src.H(), func(ys <-chan int) {
		for y := range ys {
			i := y * dst.Stride
			src.Scan(0, y, src.W(), y+1, dst.Pix[i:i+src.W()*4])
			for x := 0; x < src.W(); x++ {
				d := dst.Pix[i : i+3 : i+3]
				d[0] = 255 - d[0]
				d[1] = 255 - d[1]
				d[2] = 255 - d[2]
				i += 4
			}
		}
	})
	return dst
}

// Parallel processes the data in separate goroutines.
func Parallel(start, stop int, fn func(<-chan int)) {
	count := stop - start
	if count < 1 {
		return
	}

	procs := runtime.GOMAXPROCS(0)
	if procs > count {
		procs = count
	}

	c := make(chan int, count)
	for i := start; i < stop; i++ {
		c <- i
	}
	close(c)

	var wg sync.WaitGroup
	for i := 0; i < procs; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			fn(c)
		}()
	}
	wg.Wait()
}
