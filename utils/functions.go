package utils

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/png"
	"os"
	"runtime"
	"sync"
	"sync/atomic"

	"github.com/disintegration/imaging"
)

var maxProcs int64

func SetMaxProcs(value int) {
	atomic.StoreInt64(&maxProcs, int64(value))
}

func Parallel(start, stop int, fn func(<-chan int)) {
	count := stop - start
	if count < 1 {
		return
	}
	procs := runtime.GOMAXPROCS(0)
	limit := int(atomic.LoadInt64(&maxProcs))
	if procs > limit && limit > 0 {
		procs = limit
	}
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

// print out image on iTerm2; equivalent to imgcat on iTerm2
func PrintImage(img image.Image, invert bool, file string) {
	var buf bytes.Buffer
	fmt.Printf("Should invert: %v\n", invert)
	if invert {
		inverted := imaging.Invert(img)
		imaging.Save(inverted, file+"_inverted.png")
		png.Encode(&buf, inverted)
	} else {
		png.Encode(&buf, img)
	}
	imgBase64Str := base64.StdEncoding.EncodeToString(buf.Bytes())
	fmt.Printf("\x1b]1337;File=inline=1:%s\a\n", imgBase64Str)
}

// get the file as an image
func GetImage(filePath string) image.Image {
	imgFile, err := os.Open(filePath)
	defer imgFile.Close()
	if err != nil {
		fmt.Println("Cannot read file:", err)
	}
	img, _, err := image.Decode(imgFile)
	if err != nil {
		fmt.Println("Cannot decode file:", err)
	}
	return img
}