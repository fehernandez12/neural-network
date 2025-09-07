package server

import "net/http"

func getStatusCode(operation string) int {
	switch operation {
	case "train":
		return http.StatusCreated
	case "predict":
		return http.StatusOK
	default:
		return http.StatusOK
	}
}
