package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"github.com/sandeep-jaiswar/jaiswar-securities/internal/session"
)

type Server struct {
	router         *mux.Router
	logger         *zap.Logger
	handlers       *Handlers
	sessionManager *session.SessionManager
	httpServer     *http.Server
}

type Handlers struct {
	logger         *zap.Logger
	sessionManager *session.SessionManager
}

func NewServer(logger *zap.Logger, port string) *Server {
	router := mux.NewRouter()
	sessionManager := session.NewSessionManager()
	handlers := &Handlers{
		logger:         logger,
		sessionManager: sessionManager,
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
	fmt.Fprintln(w, "Login endpoint")
}

func (h *Handlers) TokenHandler(w http.ResponseWriter, r *http.Request) {
	token := "some_generated_token"
	userID := "user123"

	h.logger.Info("Token handler called")
	h.logger.Info("Storing token", zap.String("userID", userID), zap.String("token", token))
	h.sessionManager.StoreToken(userID, token)

	fmt.Fprintf(w, "Token stored for user %s: %s", userID, token)
}
