package network

import "math"

type TransferFunc func(float64) float64

func HeavySideTransfer(x float64) float64 {
	if x >= 0.0 {
		return 1.0
	}
	return 0.0
}

func HeavySideTransferDerivative(x float64) float64 {
	return 1.0
}

func SigmoidalTransfer(x float64) float64 {
	return 1.0 / (1.0 + math.Exp(-x))
}

func SigmoidalTransferDerivative(x float64) float64 {
	return x * (1.0 - x)
}

func HyperbolicTransfer(x float64) float64 {
	return math.Tanh(x)
}

func HyperbolicTransferDerivative(x float64) float64 {
	return 1.0 - x*x
}
