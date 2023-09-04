package network

import (
	"math"
	"strconv"
)

func StringToFloat(values []string, mode int, defaultValue float64) []float64 {
	var floatValues []float64
	if mode == 0 {
		floatValues = make([]float64, len(values))
	}
	for i, value := range values {
		floatValue, err := strconv.ParseFloat(value, 64)
		if mode == 0 {
			if err == nil {
				floatValues[i] = floatValue
			} else {
				floatValues[i] = defaultValue
			}
		} else {
			if err == nil {
				floatValues = append(floatValues, floatValue)
			}
		}
	}
	return floatValues
}

func StringInSlice(el string, slice []string) (bool, int) {
	for i, v := range slice {
		if v == el {
			return true, i
		}
	}
	return false, -1
}

func IntToBinary(n int64, d int) []float64 {
	bs := strconv.FormatInt(n, 2)
	bi := make([]float64, d)
	zn := d - len(bs)
	if zn < 0 {
		bi = make([]float64, len(bs))
		zn = 0
	}
	for zi := 0; zi < d; zi++ {
		if zi < zn {
			bi[zi] = 0.0
		} else {
			bi[zi], _ = strconv.ParseFloat(string(bs[zi-zn]), 64)
		}
	}
	return bi
}

func BinaryToInt(n []float64) int {
	bi := 0
	for i := len(n) - 1; i >= 0; i-- {
		bi += int(n[i]) * int(math.Pow(float64(2), float64(len(n)-i-1)))
	}
	return bi
}

func ScalarProduct(a, b []float64) float64 {
	var sum float64
	for i := range a {
		sum += a[i] * b[i]
	}
	return sum
}

func Accuracy(actual, predicted []float64) (int, float64) {
	if len(actual) != len(predicted) {
		return -1, -1.0
	}
	correct := 0
	for i := range actual {
		if actual[i] == predicted[i] {
			correct++
		}
	}
	return correct, float64(correct) / float64(len(actual)) * 100.0
}

func Round(val, roundOn float64, places int) float64 {
	result := 0.0
	round := 0.0
	pow := math.Pow(10, float64(places))
	digit := pow * val
	_, div := math.Modf(digit)
	if div >= roundOn {
		round = math.Ceil(digit)
	} else {
		round = math.Floor(digit)
	}
	result = round / pow
	return result
}
