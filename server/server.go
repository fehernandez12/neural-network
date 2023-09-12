package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"gonn/cache"
	"gonn/models"
	"gonn/network"
	"gonn/utils"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/disintegration/imaging"
	"github.com/fehernandez12/sonate"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"gonum.org/v1/gonum/mat"
	"gopkg.in/yaml.v3"
)

type ServerConfig struct {
	Timeout      int
	OpenAIAPIKey string
}

type Server struct {
	config  *ServerConfig
	addr    string
	logger  *Logger
	network *network.Network
}

func (s *Server) Config() *ServerConfig {
	return s.config
}

func (s *Server) Addr() string {
	return s.addr
}

func (s *Server) Logger() *Logger {
	return s.logger
}

func (s *Server) Network() *network.Network {
	return s.network
}

func NewServer(addr string) (*Server, error) {
	if addr == "" {
		return nil, errors.New("addr cannot be empty")
	}
	config, err := readServerConfig()
	if err != nil {
		return nil, err
	}
	cacheRep := cache.NewRedisCacheRepository()
	cache.SetRepository(cacheRep)
	return &Server{
		addr:    addr,
		config:  config,
		logger:  NewLogger(),
		network: network.NewNetwork(784, 200, 10, 0.1),
	}, nil
}

func (s *Server) Start(stop <-chan struct{}) error {
	s.network.Load()
	srv := &http.Server{
		Addr:    s.addr,
		Handler: s.router(),
	}
	go func() {
		logrus.WithField("addr", s.addr).Info("starting server")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logrus.Fatalf("listen: %s\n", err)
		}
	}()
	<-stop
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(s.config.Timeout)*time.Millisecond)
	defer cancel()
	logrus.WithField("timeout", s.config.Timeout).Info("shutting down server")
	return srv.Shutdown(ctx)
}

func (s *Server) router() http.Handler {
	router := sonate.NewRouter()
	// router.Use(middleware.RequestLogger)
	router.HandleFunc("/", s.defaultRoute).Methods(http.MethodPost)
	return router
}

func (s *Server) defaultRoute(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	if err := r.ParseMultipartForm(2 << 20); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		s.logger.Error(http.StatusBadRequest, r.URL.Path, err)
		sendErrorResponse(w, http.StatusBadRequest, err)
		return
	}
	var resp models.Response
	switch r.FormValue("operation") {
	case "train":
		resp = s.TrainNetwork()
	case "predict":
		if _, err := os.Stat("./tmp"); err != nil {
			if os.IsNotExist(err) {
				os.Mkdir("./tmp", 0755)
			} else {
				w.WriteHeader(http.StatusInternalServerError)
				s.logger.Error(http.StatusInternalServerError, r.URL.Path, err)
				sendErrorResponse(w, http.StatusInternalServerError, err)
				return
			}
		}
		var status int
		var err error
		resp, status, err = s.PredictNetwork(r)
		if err != nil {
			w.WriteHeader(status)
			s.logger.Error(status, r.URL.Path, err)
			sendErrorResponse(w, status, err)
			return
		}
		os.RemoveAll("./tmp")
	}
	response, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
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

func (s *Server) TrainNetwork() *models.TrainResponse {
	start := time.Now()
	resp := &models.TrainResponse{}
	resp.Operation = "train"
	s.network.MnistTrain()
	s.network.Save()
	resp.Time = time.Since(start).String()
	resp.Message = "Training complete"
	resp.Success = true
	return resp
}

func (s *Server) PredictNetwork(r *http.Request) (*models.PredictResponse, int, error) {
	start := time.Now()
	resp := &models.PredictResponse{}
	resp.Operation = "predict"
	// get the file from the request and
	// save it as a temporary image file
	_, h, err := r.FormFile("image")
	if err != nil {
		return nil, http.StatusBadRequest, err
	}
	// create file
	fileName, err := makeFile(h, r)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}
	isCached, cachedResult, err := checkCache(fileName)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	if isCached {
		logrus.Info("Retrieved from cache")
		resp.Prediction = cachedResult.Prediction
		resp.Results = cachedResult.Results
		resp.Accuracy = cachedResult.Accuracy
		resp.Time = time.Since(start).String()
	} else {
		prediction, results, accuracy := s.Predict(fileName)
		resp.Prediction = prediction
		resp.Results = makeResultsMap(results)
		resp.Accuracy = accuracy
		resp.Time = time.Since(start).String()
		logrus.Info("Saving to cache")
		cacheValue, err := json.Marshal(resp)
		if err != nil {
			return nil, http.StatusInternalServerError, err
		}
		err = cache.Put(utils.GetSHA256Checksum(fileName), string(cacheValue))
		if err != nil {
			return nil, http.StatusInternalServerError, err
		}
	}
	resp.Success = true
	return resp, http.StatusOK, nil
}

func makeResultsMap(results []float64) map[string]float64 {
	m := make(map[string]float64)
	for i := 0; i < len(results); i++ {
		m[strconv.Itoa(i)] = results[i]
	}
	return m
}

func makeFile(h *multipart.FileHeader, r *http.Request) (string, error) {
	name := uuid.New().String()
	saveFile(h, name)
	fileName := fmt.Sprintf("./tmp/%s_raw.png", name)
	img := utils.GetImage(fmt.Sprintf("./tmp/%s_raw.png", name))
	invert, err := strconv.ParseBool(r.FormValue("invert"))
	if err != nil {
		return "", err
	}
	if invert {
		inverted := imaging.Invert(img)
		imaging.Save(inverted, fmt.Sprintf("./tmp/%s_inverted.png", name))
		fileName = fmt.Sprintf("./tmp/%s_inverted.png", name)
	}
	return fileName, nil
}

func (s *Server) Predict(fileName string) (int, []float64, float64) {
	input := network.DataFromImage(fileName)
	output := s.network.Predict(input)
	results := MatrixToSlice(output)
	best := 0
	highest := 0.0
	for i := 0; i < s.network.Outputs; i++ {
		if output.At(i, 0) > highest {
			best = i
			highest = output.At(i, 0)
		}
	}
	return best, results, float64(int(highest*10000)) / 100
}

func checkCache(fileName string) (bool, *models.PredictResponse, error) {
	checkSum := utils.GetSHA256Checksum(fileName)
	logrus.Info("Checking cache for ", checkSum)
	result, err := cache.Get(checkSum)
	if err == redis.Nil {
		logrus.Info("Not found in cache")
		return false, nil, nil
	}
	if result != "" {
		logrus.Info("Found in cache")
		resp := &models.PredictResponse{}
		err = json.Unmarshal([]byte(result), &resp)
		if err != nil {
			return false, nil, err
		}
		return true, resp, nil
	}
	return false, nil, nil
}

func MatrixToSlice(m mat.Matrix) []float64 {
	r, _ := m.Dims()
	s := make([]float64, r)
	for i := 0; i < r; i++ {
		s[i] = float64(int(m.At(i, 0)*10000)) / 100
	}
	return s
}

func saveFile(h *multipart.FileHeader, name string) {
	file, err := h.Open()
	if err != nil {
		return
	}
	defer file.Close()
	out, err := os.Create(fmt.Sprintf("./tmp/%s_raw.%s", name, getExtension(h.Filename)))
	if err != nil {
		return
	}
	defer out.Close()
	_, err = io.Copy(out, file)
	if err != nil {
		return
	}
}

func sendErrorResponse(w http.ResponseWriter, status int, err error) {
	response, err := json.Marshal(models.NewErrorResponse(status, err))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(status)
	w.Write(response)
}

func getExtension(filename string) string {
	return filename[len(filename)-3:]
}

func readServerConfig() (*ServerConfig, error) {
	yamlFile, err := os.ReadFile("./server/config.yml")
	if err != nil {
		return nil, err
	}
	config := &ServerConfig{}
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		return nil, err
	}
	config.OpenAIAPIKey = os.Getenv("OPENAI_API_KEY")
	return config, nil
}
