package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"neural-network/cache"
	"neural-network/logger"
	"neural-network/models"
	"neural-network/network"
	"neural-network/utils"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/handlers"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"gonum.org/v1/gonum/mat"
)

type ServerConfig struct {
	Timeout int
	Addr    string
}

type Server struct {
	config  *ServerConfig
	logger  *logger.Logger
	network *network.Network
}

func (s *Server) Config() *ServerConfig {
	return s.config
}

func (s *Server) Network() *network.Network {
	return s.network
}

func StartServer() error {
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	stopper := make(chan struct{})
	go func() {
		<-done
		close(stopper)
	}()
	server, err := newServer()
	if err != nil {
		return err
	}
	return server.Start(stopper)
}

func newServer() (*Server, error) {
	config, err := readServerConfig()
	if err != nil {
		return nil, err
	}
	cacheRep := cache.NewRedisCacheRepository()
	cache.SetRepository(cacheRep)
	return &Server{
		config:  config,
		logger:  logger.NewLogger(),
		network: network.NewNetwork(784, 200, 10, 0.1),
	}, nil
}

func (s *Server) Start(stop <-chan struct{}) error {
	corsObj := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"POST", "GET", "OPTIONS", "PUT", "DELETE"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	)
	// load the neural network weights from file
	s.network.Load()
	srv := &http.Server{
		Addr:    s.config.Addr,
		Handler: corsObj(s.router()),
	}
	go func() {
		s.logger.WithField("addr", s.config.Addr)
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

func (s *Server) TrainNetwork() *models.TrainResponse {
	start := time.Now()
	// check if weights are already trained
	h, err := os.Stat("./data/hweights.model")
	o, err1 := os.Stat("./data/oweights.model")
	if (err1 == nil && o.Size() > 0) && (err == nil && h.Size() > 0) {
		logrus.WithField("step", "skipping training").Info("training network")
		return &models.TrainResponse{
			OperationResponse: models.OperationResponse{
				Operation: "train",
				ApiResponse: models.ApiResponse{
					Success: true,
					Time:    "0s",
				},
			},
			Message: "Weights already trained, skipping training",
		}
	}
	logrus.WithField("step", "starting training").Info("training network")
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
		inverted := utils.InvertImage(img)
		utils.Save(inverted, fmt.Sprintf("./tmp/%s_inverted.png", name))
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
	w.Header().Set("Content-Type", "application/json")
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
	config := &ServerConfig{}
	timeout, err := strconv.Atoi(os.Getenv("SERVER_TIMEOUT"))
	if err != nil {
		timeout = 30000
	}
	config.Timeout = timeout
	addr := os.Getenv("SERVER_ADDR")
	if addr == "" {
		addr = ":8080"
	}
	config.Addr = addr
	return config, nil
}
