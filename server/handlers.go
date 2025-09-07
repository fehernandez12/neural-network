package server

import (
	"encoding/json"
	"net/http"
	"neural-network/models"
	"os"
	"time"
)

func (s *Server) trainRoute(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	if err := r.ParseMultipartForm(2 << 20); err != nil {
		s.logger.Error(http.StatusBadRequest, r.URL.Path, err)
		sendErrorResponse(w, http.StatusBadRequest, err)
		return
	}
	resp := s.TrainNetwork()
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
