package server

import (
	"encoding/json"
	"net/http"
	"neural-network/models"
	"os"
	"time"
)

type ErrorResponse struct {
	Messages []string `json:messages`
}

func (s *Server) handleError(w http.ResponseWriter, statusCode int, route string, cause error) {
	errorResponse := &ErrorResponse{}
	errorResponse.Messages = append(errorResponse.Messages, cause.Error())
	response, err := json.Marshal(errorResponse)
	if err != nil {
		s.handleError(w, http.StatusInternalServerError, route, err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(response)
	s.logger.Error(statusCode, route, cause)
}

func (s *Server) trainRoute(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	var req models.TrainRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&req)
	if err != nil {
		s.handleError(w, http.StatusBadRequest, "/train", err)
	}
	resp, err := s.TrainNetwork(&req)
	if err != nil {
		s.handleError(w, http.StatusBadRequest, "/train", err)
	}
	response, err := json.Marshal(resp)
	if err != nil {
		s.logger.Error(http.StatusInternalServerError, r.URL.Path, err)
		sendErrorResponse(w, http.StatusInternalServerError, err)
		return
	}
	statusCode := getStatusCode(resp.GetOperation())
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(response)
	s.logger.Info(statusCode, r.URL.Path, start)
}

func (s *Server) predictRoute(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	if err := r.ParseMultipartForm(2 << 20); err != nil {
		s.logger.Error(http.StatusBadRequest, r.URL.Path, err)
		sendErrorResponse(w, http.StatusBadRequest, err)
		return
	}
	var resp models.Response
	if _, err := os.Stat("./tmp"); err != nil {
		if os.IsNotExist(err) {
			os.Mkdir("./tmp", 0755)
		} else {
			s.logger.Error(http.StatusInternalServerError, r.URL.Path, err)
			sendErrorResponse(w, http.StatusInternalServerError, err)
			return
		}
	}
	var status int
	var err error
	resp, status, err = s.PredictNetwork(r)
	if err != nil {
		s.logger.Error(status, r.URL.Path, err)
		sendErrorResponse(w, status, err)
		return
	}
	os.RemoveAll("./tmp")
	response, err := json.Marshal(resp)
	if err != nil {
		s.logger.Error(http.StatusInternalServerError, r.URL.Path, err)
		sendErrorResponse(w, http.StatusInternalServerError, err)
		return
	}
	statusCode := getStatusCode(resp.GetOperation())
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(response)
	s.logger.Info(statusCode, r.URL.Path, start)
}
