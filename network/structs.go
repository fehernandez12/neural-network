package network

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

// A NeuralLayer represents a single layer of the Perceptron.
type NeuralLayer struct {
	NeuronUnits []NeuronUnit
	Length      int
}

type MultiLayerNetwork struct {
	Lrate                  float64
	NeuralLayers           []NeuralLayer
	TransferFunc           tFunc
	TransferFuncDerivative tFunc
}
