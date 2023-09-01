package network

type MultiLayerNetwork struct {
	Lrate                  float64
	NeuralLayers           []*NeuralLayer
	TransferFunc           TransferFunc
	TransferFuncDerivative TransferFunc
}

func NewMultiLayerNetwork(l []int, lRate float64, tFunc, tFuncDerivative TransferFunc) *MultiLayerNetwork {
	mlp := MultiLayerNetwork{
		Lrate:                  lRate,
		TransferFunc:           tFunc,
		TransferFuncDerivative: tFuncDerivative,
	}
	mlp.NeuralLayers = make([]*NeuralLayer, len(l))
	for i, q := range l {
		if i != 0 {
			mlp.NeuralLayers[i] = NewNeuralLayer(q, l[i-1])
		} else {
			mlp.NeuralLayers[i] = NewNeuralLayer(q, 0)
		}
	}
	return &mlp
}
