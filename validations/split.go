package validations

import (
	"math/rand"
	"neural-network/network"
)

func TrainTestPatternsSplit(
	patterns []*network.Pattern,
	perc float64,
	shuffle bool,
) ([]*network.Pattern, []*network.Pattern) {
	splitPivot := int(float64(len(patterns)) * perc)
	train := make([]*network.Pattern, splitPivot)
	test := make([]*network.Pattern, len(patterns)-splitPivot)

	if shuffle {
		perm := rand.Perm(len(patterns))

		for i := 0; i < splitPivot; i++ {
			train[i] = patterns[perm[i]]
		}

		for i := 0; i < len(patterns)-splitPivot; i++ {
			test[i] = patterns[perm[i]]
		}
	} else {
		train = patterns[:splitPivot]
		test = patterns[splitPivot:]
	}
	return train, test
}

func KFoldPatternsSplit(
	patterns []*network.Pattern,
	k int,
	shuffle bool,
) [][]*network.Pattern {
	size := int(len(patterns) / k)
	remainder := int(len(patterns) % k)
	folds := make([][]*network.Pattern, k)
	var perm []int
	if shuffle {
		perm = rand.Perm(len(patterns))
	}
	currSize := 0
	start := 0
	curr := 0
	for i := 0; i < k; i++ {
		curr = start
		currSize = size
		if i < remainder {
			currSize++
		}
		folds[i] = make([]*network.Pattern, currSize)
		for j := 0; j < currSize; j++ {
			if shuffle {
				folds[i][j] = patterns[perm[curr]]
			} else {
				folds[i][j] = patterns[curr]
			}
			curr++
		}
		start = curr
	}
	return folds
}
