package validations

import "neural-network/network"

func RandomSubsamplingValidation(
	neuron *network.NeuronUnit,
	patterns []*network.Pattern,
	percentage float64,
	epochs, folds int,
	shuffle bool,
) []float64 {
	var actual, predicted []float64
	scores := make([]float64, folds)

	for i := 0; i < folds; i++ {
		train, test := TrainTestPatternsSplit(patterns, percentage, shuffle)
		neuron.Train(train, epochs, true)
		for _, pattern := range test {
			actual = append(actual, pattern.SingleExpectation)
			predicted = append(predicted, neuron.Predict(pattern))
		}
		_, correctPerc := network.Accuracy(actual, predicted)
		scores[i] = correctPerc
	}

	return scores
}

func KFoldValidation(
	neuron *network.NeuronUnit,
	patterns []*network.Pattern,
	epochs, k int,
	shuffle bool,
) []float64 {

}
