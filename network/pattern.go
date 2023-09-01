package network

import (
	"encoding/csv"
	"io"
	"math/rand"
	"os"
	"strings"
)

const (
	SCALING_FACTOR = 0.0000000000001
)

// Pattern represents a single input to the Perceptron struct.
// It technically is a pattern with certain dimensions and desired value.
//
// * Features is a slice of float64 values that represent the input dimensions.
//
// * SingleRawExpectation is a string and is filled by the parser with the input classification (in terms of belonging class)
//
// * SingleExpectation is a float64 representation of the class to which the Pattern belongs.
//
// * MultipleExpectation is a slice of float64 values and is used for multiple classification cases.
type Pattern struct {
	Features             []float64
	SingleRawExpectation string
	SingleExpectation    float64
	MultipleExpectation  []float64
}

func LoadFromCSV(path string) ([]*Pattern, []string, error) {
	var patterns []*Pattern
	fileContents, err := os.ReadFile(path)
	if err != nil {
		return nil, nil, err
	}
	reader := csv.NewReader(strings.NewReader(string(fileContents)))
	lineCounter := 0
	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, nil, err
		}
		floatingValues := StringToFloat(line, 1, -1.0)
		patterns = append(patterns, &Pattern{
			Features:             floatingValues,
			SingleRawExpectation: line[len(line)-1],
		})
		lineCounter++
	}
	mapped := RawConvert(patterns)
	return patterns, mapped, nil
}

func RawConvert(patterns []*Pattern) []string {
	var mapped []string
	for _, pattern := range patterns {
		check, _ := StringInSlice(pattern.SingleRawExpectation, mapped)
		if !check {
			mapped = append(mapped, pattern.SingleRawExpectation)
		}
	}
	for i := range patterns {
		for mappedIndex, mappedValue := range mapped {
			if strings.Compare(mappedValue, patterns[i].SingleRawExpectation) == 0 {
				patterns[i].SingleExpectation = float64(mappedIndex)
			}
		}
	}
	return mapped
}

func NewRandomInt(d int) int64 {
	return rand.Int63n(int64(2 ^ d))
}

func NewRandomPatternArray(d int, k int) []*Pattern {
	var patterns []*Pattern
	for i := 0; i < k; i++ {
		a := NewRandomInt(d)
		b := NewRandomInt(d)
		c := a + b
		ab := IntToBinary(a, d)
		bb := IntToBinary(b, d)
		ab = append(ab, bb...)

		patterns = append(patterns, &Pattern{
			Features:            ab,
			MultipleExpectation: IntToBinary(c, d+1),
		})
	}
	return patterns
}
