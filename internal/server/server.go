package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"github.com/sandeep-jaiswar/jaiswar-securities/internal/session"
	"go.uber.org/zap"
)

type Server struct {
	router         *mux.Router
	logger         *zap.Logger
	handlers       *Handlers
	sessionManager *session.SessionManager
	httpServer     *http.Server
	sessionClient  *session.SessionManager
}

type Handlers struct {
	logger         *zap.Logger
	sessionManager *session.SessionManager
	sessionClient  *session.SessionManager
}

func NewServer(logger *zap.Logger, port string, sessionClient *session.SessionManager) *Server {
	router := mux.NewRouter()
	sessionManager := session.NewSessionManager()
	handlers := &Handlers{
		logger:         logger,
		sessionManager: sessionManager,
		sessionClient:  sessionClient,
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
		sessionClient:  sessionClient,
	}
	s.InitializeRoutes()
	return s
}

func (s *Server) InitializeRoutes() {
	s.router.HandleFunc("/api/v1/login", s.handlers.LoginHandler).Methods(http.MethodGet)
	s.router.HandleFunc("/api/v1/token", s.handlers.TokenHandler).Methods(http.MethodGet)
	s.router.HandleFunc("/api/v1/profile", s.handlers.ProfileHandler).Methods(http.MethodGet)
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
}

func (h *Handlers) TokenHandler(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("Token handler called")
}

func (h *Handlers) ProfileHandler(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("Profile handler called")
}
