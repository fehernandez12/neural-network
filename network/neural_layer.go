package network

// A NeuralLayer represents a single layer of the Perceptron.
type NeuralLayer struct {
	NeuronUnits []*NeuronUnit
	Length      int
}

func NewNeuralLayer(n, p int) *NeuralLayer {
	layer := &NeuralLayer{
		NeuronUnits: make([]*NeuronUnit, n),
		Length:      n,
	}
	for i := 0; i < n; i++ {
		layer.NeuronUnits[i] = NewRandomNeuronUnit(p)
	}
	return layer
}
