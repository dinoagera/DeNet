package api

import (
	"denettest/internal/config"
	"denettest/internal/handlers"
	"denettest/internal/logger"
	"denettest/internal/middleware"
	"denettest/internal/repository/postgres"
	"denettest/internal/service"
	"log/slog"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
)

type API struct {
	log     *slog.Logger
	cfg     *config.Config
	router  *mux.Router
	db      *pgxpool.Pool
	handler *handlers.Handler
}

func InitAPI() *API {
	log := logger.InitLogger()
	cfg := config.InitConfig(log)
	db, err := postgres.New(cfg.DBPath)
	if err != nil {
		log.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	authRepository := postgres.NewAuthStorage(db)
	userRepository := postgres.NewUserStorage(db)
	taskRepository := postgres.NewTaskStorage(db)
	service := service.New(log, authRepository, userRepository, taskRepository)
	handler := handlers.New(log, service, service, service)
	api := &API{
		log:     log,
		cfg:     cfg,
		router:  mux.NewRouter(),
		db:      db,
		handler: handler,
	}

	return api
}

func (api *API) StartServer() {
	api.setupRoutes()

	server := &http.Server{
		Handler:      api.router,
		Addr:         api.cfg.HTTPAddress,
		ReadTimeout:  api.cfg.HTTPIdleTimeout,
		WriteTimeout: api.cfg.HTTPReadTimeout,
		IdleTimeout:  api.cfg.HTTPIdleTimeout,
	}

	api.log.Info("starting server", "address", server.Addr)

	if err := server.ListenAndServe(); err != nil {
		api.log.Error("server failed", "error", err)
		os.Exit(1)
	}
}
func (api *API) setupRoutes() {
	//without middleware
	publicRouter := api.router.PathPrefix("").Subrouter()
	publicRouter.HandleFunc("/register", api.handler.Register).Methods("POST")
	publicRouter.HandleFunc("/login", api.handler.Login).Methods("POST")
	//with middleware
	privateRouter := api.router.PathPrefix("").Subrouter()
	privateRouter.Use(middleware.New(api.log, api.cfg.SecretKey))
	privateRouter.HandleFunc("/users/{id}/status", api.handler.GetStatus).Methods("GET")
	privateRouter.HandleFunc("/users/leaderboard", api.handler.GetLeaderboard).Methods("GET")
	privateRouter.HandleFunc("/users/{id}/task/complete", api.handler.CompleteTask).Methods("POST")
	privateRouter.HandleFunc("/users/{id}/referrer", api.handler.SetReferrer).Methods("POST")
}
