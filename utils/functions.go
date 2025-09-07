package utils

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"image"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"neural-network/images"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"

	"golang.org/x/image/bmp"
	"golang.org/x/image/tiff"
)

var maxProcs int64

// EncodeOption sets an optional parameter for the Encode and Save functions.
type EncodeOption func(*encodeConfig)

type encodeConfig struct {
	jpegQuality         int
	gifNumColors        int
	gifQuantizer        draw.Quantizer
	gifDrawer           draw.Drawer
	pngCompressionLevel png.CompressionLevel
}

var defaultEncodeConfig = encodeConfig{
	jpegQuality:         95,
	gifNumColors:        256,
	gifQuantizer:        nil,
	gifDrawer:           nil,
	pngCompressionLevel: png.DefaultCompression,
}

// Format is an image file format.
type Format int

// Image file formats.
const (
	JPEG Format = iota
	PNG
	GIF
	TIFF
	BMP
)

var formatExts = map[string]Format{
	"jpg":  JPEG,
	"jpeg": JPEG,
	"png":  PNG,
	"gif":  GIF,
	"tif":  TIFF,
	"tiff": TIFF,
	"bmp":  BMP,
}

type fileSystem interface {
	Create(string) (io.WriteCloser, error)
	Open(string) (io.ReadCloser, error)
}

type localFS struct{}

func (localFS) Create(name string) (io.WriteCloser, error) { return os.Create(name) }
func (localFS) Open(name string) (io.ReadCloser, error)    { return os.Open(name) }

var fs fileSystem = localFS{}

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
		inverted := InvertImage(img)
		Save(inverted, file+"_inverted.png")
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
	if err != nil {
		fmt.Println("Cannot read file:", err)
		return nil
	}
	defer imgFile.Close()
	img, _, err := image.Decode(imgFile)
	if err != nil {
		fmt.Println("Cannot decode file:", err)
		return nil
	}
	return img
}

// GetSHA256Checksum gets the checksum of an uploaded file.
// This is used to check if the file has been uploaded before
// and retrieve its result from cache. It is also used to
// save results in cache.
func GetSHA256Checksum(filePath string) string {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Cannot read file:", err)
	}
	defer file.Close()
	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		fmt.Println("Cannot copy file:", err)
	}
	return hex.EncodeToString(hash.Sum(nil))
}

func InvertImage(img image.Image) *image.NRGBA {
	return images.Invert(img)
}

// FormatFromExtension parses image format from filename extension:
// "jpg" (or "jpeg"), "png", "gif", "tif" (or "tiff") and "bmp" are supported.
func FormatFromExtension(ext string) (Format, error) {
	if f, ok := formatExts[strings.ToLower(strings.TrimPrefix(ext, "."))]; ok {
		return f, nil
	}
	return -1, errors.New("unknown image format: " + ext)
}

// FormatFromFilename parses image format from filename:
// "jpg" (or "jpeg"), "png", "gif", "tif" (or "tiff") and "bmp" are supported.
func FormatFromFilename(filename string) (Format, error) {
	ext := filepath.Ext(filename)
	return FormatFromExtension(ext)
}

// Encode writes the image img to w in the specified format (JPEG, PNG, GIF, TIFF or BMP).
func Encode(w io.Writer, img image.Image, format Format, opts ...EncodeOption) error {
	cfg := defaultEncodeConfig
	for _, option := range opts {
		option(&cfg)
	}

	switch format {
	case JPEG:
		if nrgba, ok := img.(*image.NRGBA); ok && nrgba.Opaque() {
			rgba := &image.RGBA{
				Pix:    nrgba.Pix,
				Stride: nrgba.Stride,
				Rect:   nrgba.Rect,
			}
			return jpeg.Encode(w, rgba, &jpeg.Options{Quality: cfg.jpegQuality})
		}
		return jpeg.Encode(w, img, &jpeg.Options{Quality: cfg.jpegQuality})

	case PNG:
		encoder := png.Encoder{CompressionLevel: cfg.pngCompressionLevel}
		return encoder.Encode(w, img)

	case GIF:
		return gif.Encode(w, img, &gif.Options{
			NumColors: cfg.gifNumColors,
			Quantizer: cfg.gifQuantizer,
			Drawer:    cfg.gifDrawer,
		})

	case TIFF:
		return tiff.Encode(w, img, &tiff.Options{Compression: tiff.Deflate, Predictor: true})

	case BMP:
		return bmp.Encode(w, img)
	}

	return errors.New("unknown image format")
}

// Save saves the image to file with the specified filename.
// The format is determined from the filename extension:
// "jpg" (or "jpeg"), "png", "gif", "tif" (or "tiff") and "bmp" are supported.
//
// Examples:
//
//	// Save the image as PNG.
//	err := imaging.Save(img, "out.png")
//
//	// Save the image as JPEG with optional quality parameter set to 80.
//	err := imaging.Save(img, "out.jpg", imaging.JPEGQuality(80))
func Save(img image.Image, filename string, opts ...EncodeOption) (err error) {
	f, err := FormatFromFilename(filename)
	if err != nil {
		return err
	}
	file, err := fs.Create(filename)
	if err != nil {
		return err
	}
	err = Encode(file, img, f, opts...)
	errc := file.Close()
	if err == nil {
		err = errc
	}
	return err
}
