package network

import "math/rand"

// A NeuronUnit represents a single computation unit. It corresponds
// to the simple binary perceptron as defined by Rosenblatt.
//
// * Weights is a slice of float64 values that represent the way each dimension of
// the pattern is modulated.
//
// * Bias is a float64 value that represents the propensity of the NeuronUnit to spread signal.
//
// * Lrate is a float64 that represents the learning rate of the neuron.
//
// * Value is a float64 that represents the desired value when I load the input pattern.
//
// * Delta is a float64 that represents the difference between the desired value and the actual value (the error).
type NeuronUnit struct {
	Weights []float64
	Bias    float64
	Lrate   float64
	Value   float64
	Delta   float64
}

func NewRandomNeuronUnit(dim int) *NeuronUnit {
	neuron := &NeuronUnit{
		Weights: make([]float64, dim),
		Bias:    rand.NormFloat64() * SCALING_FACTOR,
		Lrate:   rand.NormFloat64() * SCALING_FACTOR,
		Value:   rand.NormFloat64() * SCALING_FACTOR,
		Delta:   rand.NormFloat64() * SCALING_FACTOR,
	}

	for i := range neuron.Weights {
		neuron.Weights[i] = rand.NormFloat64() * SCALING_FACTOR
	}

	return neuron
}

func (n *NeuronUnit) UpdateWeights(pattern *Pattern) (float64, float64) {
	predictedValue, prevError, postError := n.Predict(pattern), 0.0, 0.0
	prevError = pattern.SingleExpectation - predictedValue
	n.Bias += (n.Lrate * prevError)
	for i := range n.Weights {
		n.Weights[i] += (n.Lrate * prevError * pattern.Features[i])
	}
	predictedValue = n.Predict(pattern)
	postError = pattern.SingleExpectation - predictedValue
	return prevError, postError
}

func (n *NeuronUnit) Train(patterns []*Pattern, epochs int, init bool) {
	if init {
		n.Weights = make([]float64, len(patterns[0].Features))
		n.Bias = 0.0
	}
	squaredPrevError, squaredPostError := 0.0, 0.0
	for i := 0; i < epochs; i++ {
		for _, pattern := range patterns {
			prevError, postError := n.UpdateWeights(pattern)
			squaredPrevError += prevError * prevError
			squaredPostError += postError * postError
		}
	}
}

func (n *NeuronUnit) Predict(pattern *Pattern) float64 {
	if (ScalarProduct(n.Weights, pattern.Features) + n.Bias) < 0 {
		return 0.0
	}
	return 1.0
}
