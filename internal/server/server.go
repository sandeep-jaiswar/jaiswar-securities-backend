package server

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"github.com/sandeep-jaiswar/jaiswar-securities/internal/paytm"
	"github.com/sandeep-jaiswar/jaiswar-securities/internal/session"
	"go.uber.org/zap"
)

type Server struct {
	router         *mux.Router
	logger         *zap.Logger
	handlers       *Handlers
	sessionManager *session.SessionManager
	httpServer     *http.Server
	paytmClient    *paytm.PaytmMoneyClient
}

type Handlers struct {
	logger         *zap.Logger
	sessionManager *session.SessionManager
	paytmClient    *paytm.PaytmMoneyClient
}

func NewServer(logger *zap.Logger, port string, paytmClient *paytm.PaytmMoneyClient) *Server {
	router := mux.NewRouter()
	sessionManager := session.NewSessionManager()
	handlers := &Handlers{
		logger:         logger,
		sessionManager: sessionManager,
		paytmClient:    paytmClient,
	}

	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: router,
	}

	s := &Server{
		router:         router,
		logger:         logger,
		handlers:       handlers,
		sessionManager: sessionManager,
		httpServer:     httpServer,
		paytmClient:    paytmClient,
	}
	s.InitializeRoutes()
	return s
}

func (s *Server) InitializeRoutes() {
	s.router.HandleFunc("/api/v1/login", s.handlers.LoginHandler).Methods(http.MethodGet)
	s.router.HandleFunc("/api/v1/token", s.handlers.TokenHandler).Methods(http.MethodGet)
}

func (s *Server) Start() {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)

	go func() {
		<-signalChan
		s.logger.Info("Received shutdown signal, shutting down gracefully...")
		s.Shutdown()
	}()

	s.logger.Info("Starting server", zap.String("address", s.httpServer.Addr))
	if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		s.logger.Fatal("Failed to start server", zap.Error(err))
	}
}

func (s *Server) Shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.httpServer.Shutdown(ctx); err != nil {
		s.logger.Fatal("Server forced to shutdown:", zap.Error(err))
	}

	s.logger.Info("Server exited gracefully")
}

func (h *Handlers) LoginHandler(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("Login handler called")

	stateKey := r.URL.Query().Get("stateKey")
	if stateKey == "" {
		h.logger.Error("State key is missing")
		http.Error(w, "State key is required", http.StatusBadRequest)
		return
	}

	loginURL, err := h.paytmClient.Login(stateKey)
	if err != nil {
		h.logger.Error("Error during Paytm login", zap.Error(err))
		http.Error(w, "Failed to log in", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, loginURL, http.StatusFound)
}

func (h *Handlers) TokenHandler(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	success := queryParams.Get("success")
	requestToken := queryParams.Get("requestToken")
	state := queryParams.Get("state")

	if success != "true" || requestToken == "" || state != "token" {
		h.logger.Error("Invalid request parameters for token handler")
		http.Error(w, "Invalid request parameters", http.StatusBadRequest)
		return
	}

	accessTokenResponse, err := h.paytmClient.GenerateAccessToken(requestToken)
	if err != nil {
		h.logger.Error("Error generating access token", zap.Error(err))
		http.Error(w, "Failed to generate access token", http.StatusInternalServerError)
		return
	}

	// Read the response body
	responseBody, err := ioutil.ReadAll(accessTokenResponse.Body)
	if err != nil {
		h.logger.Error("Failed to read response body", zap.Error(err))
		http.Error(w, "Failed to read response body", http.StatusInternalServerError)
		return
	}
	defer accessTokenResponse.Body.Close()

	// Log the response body
	h.logger.Info("Access token response received", zap.String("response_data", string(responseBody)))
}
