package validations

import "neural-network/network"

type Validator struct {
	Scores    []float64
	Actual    []float64
	Predicted []float64
	Train     []*network.Pattern
	Test      []*network.Pattern
}

func (v *Validator) RandomSubsamplingValidation(
	neuron *network.NeuronUnit,
	patterns []*network.Pattern,
	percentage float64,
	epochs, folds int,
	shuffle bool,
) {
	v.Scores = make([]float64, folds)

	for i := 0; i < folds; i++ {
		train, test := TrainTestPatternsSplit(patterns, percentage, shuffle)
		neuron.Train(train, epochs, true)
		for _, pattern := range test {
			v.Actual = append(v.Actual, pattern.SingleExpectation)
			v.Predicted = append(v.Predicted, neuron.Predict(pattern))
		}
		_, correctPerc := network.Accuracy(v.Actual, v.Predicted)
		v.Scores[i] = correctPerc
	}
}

func (v *Validator) KFoldValidation(
	neuron *network.NeuronUnit,
	patterns []*network.Pattern,
	epochs, k int,
	shuffle bool,
) {
	v.Scores = make([]float64, k)
	folds := KFoldPatternsSplit(patterns, k, shuffle)

	for i := 0; i < k; i++ {
		train := make([]*network.Pattern, 0)
		for j := 0; j < k; j++ {
			if i != j {
				train = append(train, folds[j]...)
			}
		}
		neuron.Train(train, epochs, true)
		for _, pattern := range folds[i] {
			v.Actual = append(v.Actual, pattern.SingleExpectation)
			v.Predicted = append(v.Predicted, neuron.Predict(pattern))
		}
		_, correctPerc := network.Accuracy(v.Actual, v.Predicted)
		v.Scores[i] = correctPerc
	}
}

func (v *Validator) MLPRandomSubsamplingValidation(
	mlp *network.MultiLayerNetwork,
	patterns []*network.Pattern,
	epochs, k int,
	shuffle bool,
	mapped []string,
) {
}
