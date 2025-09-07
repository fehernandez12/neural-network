package server

import (
	"net/http"

	"github.com/fehernandez12/sonate"
)

func (s *Server) router() http.Handler {
	router := sonate.NewRouter()
	router.StrictSlash(true)
	router.Use(s.logger.RequestLoggerMiddleware)
	router.HandleFunc("/train", s.trainRoute).Methods(http.MethodPost)
	router.HandleFunc("/predict", s.predictRoute).Methods(http.MethodPost)
	return router
}
