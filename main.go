package main

import (
	"flag"
	"fmt"
	"gonn/network"
	"gonn/utils"
)

func main() {
	// 784 inputs - 28 x 28 pixels, each pixel is an input
	// 100 hidden nodes - an arbitrary number
	// 10 outputs - digits 0 to 9
	// 0.1 is the learning rate
	net := network.NewNetwork(784, 200, 10, 0.1)

	mnist := flag.String("mnist", "", "Either train or predict to evaluate neural network")
	file := flag.String("file", "", "File name of 28 x 28 PNG file to evaluate")
	invert := flag.Bool("invert", false, "Invert the image before prediction")
	flag.Parse()

	// train or mass predict to determine the effectiveness of the trained network
	switch *mnist {
	case "train":
		net.MnistTrain()
		net.Save()
	case "predict":
		net.Load()
		net.MnistPredict()
	default:
		// don't do anything
	}

	// predict individual digit images
	if *file != "" {
		// print the image out nicely on the terminal
		utils.PrintImage(utils.GetImage(*file), *invert, *file)
		// load the neural network from file
		net.Load()
		// predict which number it is
		fmt.Println("prediction:", net.PredictFromImage(*file))
	}

}
