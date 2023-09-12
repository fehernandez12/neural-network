package models

type Response interface {
	GetOperation() string
}

type ApiResponse struct {
	Success bool   `json:"success"`
	Time    string `json:"process_time"`
}

type OperationResponse struct {
	ApiResponse
	Operation string `json:"operation"`
}

type ErrorResponse struct {
	ApiResponse
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func NewErrorResponse(status int, err error) *ErrorResponse {
	return &ErrorResponse{
		Status:  status,
		Message: err.Error(),
	}
}

type PredictResponse struct {
	OperationResponse
	Results    map[string]float64 `json:"results"`
	Prediction int                `json:"prediction"`
	Accuracy   float64            `json:"accuracy"`
}

func (r *PredictResponse) GetOperation() string {
	return r.Operation
}

type TrainResponse struct {
	OperationResponse
	Message string `json:"message"`
}

func (r *TrainResponse) GetOperation() string {
	return r.Operation
}
