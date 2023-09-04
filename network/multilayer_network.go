package network

import (
	"math"
	"math/rand"

	log "github.com/sirupsen/logrus"
)

type MultiLayerNetwork struct {
	Lrate                  float64
	NeuralLayers           []*NeuralLayer
	TransferFunc           TransferFunc
	TransferFuncDerivative TransferFunc
}

func NewMultiLayerNetwork(
	l []int,
	lRate float64,
	tFunc, tFuncDerivative TransferFunc,
) *MultiLayerNetwork {
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

func NewElmanNetwork(
	i, h, o int,
	lRate float64,
	tFunc, tFuncDerivative TransferFunc,
) *MultiLayerNetwork {
	mlp := NewMultiLayerNetwork([]int{i, h, o}, lRate, tFunc, tFuncDerivative)
	return mlp
}

func (mlp *MultiLayerNetwork) Execute(pattern *Pattern, options ...int) []float64 {
	nv := 0.0
	r := make([]float64, mlp.NeuralLayers[len(mlp.NeuralLayers)-1].Length)
	for i, feature := range pattern.Features {
		mlp.NeuralLayers[0].NeuronUnits[i].Value = feature
	}
	for i := len(pattern.Features); i < mlp.NeuralLayers[0].Length; i++ {
		mlp.NeuralLayers[0].NeuronUnits[i].Value = 0.5
	}
	for i := 1; i < len(mlp.NeuralLayers); i++ {
		for j := 0; j < mlp.NeuralLayers[i].Length; j++ {
			nv = 0.0
			for k := 0; k < mlp.NeuralLayers[i-1].Length; k++ {
				nv += mlp.NeuralLayers[i].NeuronUnits[j].Weights[k] * mlp.NeuralLayers[i-1].NeuronUnits[k].Value
			}
			nv += mlp.NeuralLayers[i].NeuronUnits[j].Bias
			mlp.NeuralLayers[i].NeuronUnits[j].Value = mlp.TransferFunc(nv)
			if i == 1 && len(options) == 0 && options[0] == 1 {
				for x := len(pattern.Features); x < mlp.NeuralLayers[0].Length; x++ {
					mlp.NeuralLayers[0].NeuronUnits[x].Value = mlp.NeuralLayers[i].NeuronUnits[x-len(pattern.Features)].Value
				}
			}
		}
	}
	for i := 0; i < mlp.NeuralLayers[len(mlp.NeuralLayers)-1].Length; i++ {
		r[i] = mlp.NeuralLayers[len(mlp.NeuralLayers)-1].NeuronUnits[i].Value
	}
	return r
}

func (mlp *MultiLayerNetwork) BackPropagate(p *Pattern, o []float64, options ...int) float64 {
	var no []float64
	var r float64
	if len(options) == 1 {
		no = mlp.Execute(p, options[0])
	} else {
		no = mlp.Execute(p)
	}

	e := 0.0
	for i := 0; i < mlp.NeuralLayers[len(mlp.NeuralLayers)-1].Length; i++ {
		e = o[i] - no[i]
		mlp.NeuralLayers[len(mlp.NeuralLayers)-1].NeuronUnits[i].Delta = e * mlp.TransferFuncDerivative(no[i])
	}
	for j := len(mlp.NeuralLayers) - 2; j >= 0; j-- {
		for k := 0; k < mlp.NeuralLayers[j].Length; k++ {
			e = 0.0
			for l := 0; l < mlp.NeuralLayers[j+1].Length; l++ {
				e += mlp.NeuralLayers[j+1].NeuronUnits[l].Delta * mlp.NeuralLayers[j+1].NeuronUnits[l].Weights[k]
			}
			mlp.NeuralLayers[j].NeuronUnits[k].Delta = e * mlp.TransferFuncDerivative(mlp.NeuralLayers[j].NeuronUnits[k].Value)
		}
		for l := 0; l < mlp.NeuralLayers[j+1].Length; l++ {
			for m := 0; m < mlp.NeuralLayers[j].Length; m++ {
				mlp.NeuralLayers[j+1].NeuronUnits[l].Weights[m] +=
					mlp.Lrate * mlp.NeuralLayers[j+1].NeuronUnits[l].Delta * mlp.NeuralLayers[j].NeuronUnits[m].Value
			}
			mlp.NeuralLayers[j+1].NeuronUnits[l].Bias += mlp.Lrate * mlp.NeuralLayers[j+1].NeuronUnits[l].Delta
		}
		if j == 1 && len(options) > 0 && options[0] == 1 {
			for i := len(p.Features); i < mlp.NeuralLayers[0].Length; i++ {
				mlp.NeuralLayers[0].NeuronUnits[i].Value = mlp.NeuralLayers[j].NeuronUnits[i-len(p.Features)].Value
			}
		}
	}
	for i := 0; i < len(o); i++ {
		r = math.Abs(no[i] - o[i])
	}
	return r / float64(len(o))
}

func (mlp *MultiLayerNetwork) Train(patterns []*Pattern, mapped []string, epochs int) {
	output := make([]float64, len(mapped))
	for epoch := 0; epoch <= epochs; epoch++ {
		for _, pattern := range patterns {
			for io, _ := range output {
				output[io] = 0.0
			}
			output[int(pattern.SingleExpectation)] = 1.0
			mlp.BackPropagate(pattern, output)
		}
	}
}

func (mlp *MultiLayerNetwork) ElmanTrain(patterns []*Pattern, epochs int) {
	for epoch := 0; epoch <= epochs; epoch++ {
		pir := rand.Intn(len(patterns))
		for i, pattern := range patterns {
			mlp.BackPropagate(pattern, pattern.MultipleExpectation, 1)
			if epoch%100 == 0 && i == pir {
				o_out := mlp.Execute(pattern, 1)
				for _, ov := range o_out {
					ov = Round(ov, .5, 0)
					log.WithFields(log.Fields{
						"SUM": "  ==========================",
					}).Info()
					log.WithFields(log.Fields{
						"a_n_1": BinaryToInt(pattern.Features[0:int(len(pattern.Features)/2)]),
						"a_n_2": pattern.Features[0:int(len(pattern.Features)/2)],
					}).Info()
					log.WithFields(log.Fields{
						"b_n_1": BinaryToInt(pattern.Features[int(len(pattern.Features)/2):]),
						"b_n_2": pattern.Features[int(len(pattern.Features)/2):],
					}).Info()
					log.WithFields(log.Fields{
						"sum_1": BinaryToInt(pattern.MultipleExpectation),
						"sum_2": pattern.MultipleExpectation,
					}).Info()
					log.WithFields(log.Fields{
						"sum_1": BinaryToInt(o_out),
						"sum_2": o_out,
					}).Info()
					log.WithFields(log.Fields{
						"END": "  ==========================",
					}).Info()
				}
			}
		}
	}

}
