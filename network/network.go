package network

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
	"gonum.org/v1/gonum/mat"
)

// Network is a neural network with 3 layers
type Network struct {
	Inputs        int
	Hiddens       int
	Outputs       int
	HiddenWeights *mat.Dense
	OutputWeights *mat.Dense
	LearningRate  float64
}

// NewNetwork creates a neural network with random weights
func NewNetwork(input, hidden, output int, rate float64) (net *Network) {
	net = &Network{
		Inputs:       input,
		Hiddens:      hidden,
		Outputs:      output,
		LearningRate: rate,
	}
	net.HiddenWeights = mat.NewDense(net.Hiddens, net.Inputs, randomArray(net.Inputs*net.Hiddens, float64(net.Inputs)))
	net.OutputWeights = mat.NewDense(net.Outputs, net.Hiddens, randomArray(net.Hiddens*net.Outputs, float64(net.Hiddens)))
	return
}

// Train the neural network
func (net *Network) Train(inputData []float64, targetData []float64) {
	// feedforward
	inputs := mat.NewDense(len(inputData), 1, inputData)
	hiddenInputs := dot(net.HiddenWeights, inputs)
	hiddenOutputs := apply(sigmoid, hiddenInputs)
	finalInputs := dot(net.OutputWeights, hiddenOutputs)
	finalOutputs := apply(sigmoid, finalInputs)

	// find errors
	targets := mat.NewDense(len(targetData), 1, targetData)
	outputErrors := subtract(targets, finalOutputs)
	hiddenErrors := dot(net.OutputWeights.T(), outputErrors)

	// backpropagate
	net.OutputWeights = add(net.OutputWeights,
		scale(net.LearningRate,
			dot(multiply(outputErrors, sigmoidPrime(finalOutputs)),
				hiddenOutputs.T()))).(*mat.Dense)

	net.HiddenWeights = add(net.HiddenWeights,
		scale(net.LearningRate,
			dot(multiply(hiddenErrors, sigmoidPrime(hiddenOutputs)),
				inputs.T()))).(*mat.Dense)
}

// Predict uses the neural network to predict the value given input data
func (net *Network) Predict(inputData []float64) mat.Matrix {
	// feedforward
	inputs := mat.NewDense(len(inputData), 1, inputData)
	hiddenInputs := dot(net.HiddenWeights, inputs)
	hiddenOutputs := apply(sigmoid, hiddenInputs)
	finalInputs := dot(net.OutputWeights, hiddenOutputs)
	finalOutputs := apply(sigmoid, finalInputs)
	return finalOutputs
}

func (net *Network) Save() error {
	logrus.WithField("step", "saving weights").Info("training network")
	h, err := os.Create("./data/hweights.model")
	if err != nil {
		return fmt.Errorf("error creating the weights file: %v", err)
	}
	if err == nil {
		net.HiddenWeights.MarshalBinaryTo(h)
	}
	logrus.WithField("path", h.Name()).Info("training network")
	defer h.Close()
	o, err := os.Create("./data/oweights.model")
	if err == nil {
		net.OutputWeights.MarshalBinaryTo(o)
	}
	defer o.Close()
	return nil
}

// load a neural network from file
func (net *Network) Load() {
	h, err := os.Open("./data/hweights.model")
	if err == nil {
		logrus.Info("Loading hidden weights")
		net.HiddenWeights.Reset()
		net.HiddenWeights.UnmarshalBinaryFrom(h)
	}
	defer h.Close()
	o, err := os.Open("./data/oweights.model")
	if err == nil {
		logrus.Info("Loading output weights")
		net.OutputWeights.Reset()
		net.OutputWeights.UnmarshalBinaryFrom(o)
	}
	defer o.Close()
}

// predict a number from an image
// image should be 28 x 28 PNG file
func (net *Network) PredictFromImage(path string) int {
	input := DataFromImage(path)
	output := net.Predict(input)
	matrixPrint(output)
	best := 0
	highest := 0.0
	for i := 0; i < net.Outputs; i++ {
		if output.At(i, 0) > highest {
			best = i
			highest = output.At(i, 0)
		}
	}
	return best
}

func (net *Network) MnistTrain(ep int) error {
	rand.NewSource(time.Now().UTC().UnixNano())
	t1 := time.Now()

	for epochs := 0; epochs < ep; epochs++ {
		testFile, err := os.Open("./mnist_dataset/mnist_train.csv")
		if err != nil {
			logrus.Errorf("error opening the training file: %v", err)
			return err
		}
		r := csv.NewReader(bufio.NewReader(testFile))
		for {
			record, err := r.Read()
			if err == io.EOF {
				break
			}

			inputs := make([]float64, net.Inputs)
			for i := range inputs {
				x, _ := strconv.ParseFloat(record[i], 64)
				inputs[i] = (x / 255.0 * 0.999) + 0.001
			}

			targets := make([]float64, 10)
			for i := range targets {
				targets[i] = 0.001
			}
			x, _ := strconv.Atoi(record[0])
			targets[x] = 0.999

			net.Train(inputs, targets)
		}
		testFile.Close()
	}
	elapsed := time.Since(t1)
	fmt.Printf("\nTime taken to train: %s\n", elapsed)
	return nil
}

func (net *Network) MnistPredict() {
	t1 := time.Now()
	checkFile, _ := os.Open("mnist_dataset/mnist_test.csv")
	defer checkFile.Close()

	score := 0
	r := csv.NewReader(bufio.NewReader(checkFile))
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		inputs := make([]float64, net.Inputs)
		for i := range inputs {
			if i == 0 {
				inputs[i] = 1.0
			}
			x, _ := strconv.ParseFloat(record[i], 64)
			inputs[i] = (x / 255.0 * 0.999) + 0.001
		}
		outputs := net.Predict(inputs)
		best := 0
		highest := 0.0
		for i := 0; i < net.Outputs; i++ {
			if outputs.At(i, 0) > highest {
				best = i
				highest = outputs.At(i, 0)
			}
		}
		target, _ := strconv.Atoi(record[0])
		if best == target {
			score++
		}
	}

	elapsed := time.Since(t1)
	fmt.Printf("Time taken to check: %s\n", elapsed)
	fmt.Println("score:", score)
}
