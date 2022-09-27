package api

import (
	"database/sql"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/sirupsen/logrus"

	"github.com/jordanknott/monitor/internal/config"
	"github.com/jordanknott/monitor/internal/db"
	"github.com/jordanknott/monitor/internal/graph"
)

// TaskcafeHandler contains all the route handlers
type MonitorApi struct {
	Data   db.Data
	Config config.AppConfig
}

// NewRouter creates a new router for chi
func NewRouter(dbConnection *sql.DB, appConfig config.AppConfig) (chi.Router, error) {
	formatter := new(logrus.TextFormatter)
	formatter.TimestampFormat = "02-01-2006 15:04:05"
	formatter.FullTimestamp = true

	routerLogger := logrus.New()
	routerLogger.SetLevel(logrus.InfoLevel)
	routerLogger.Formatter = formatter
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	// r.Use(logger.NewStructuredLogger(routerLogger))
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Cookie", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	repository := db.NewData(dbConnection)
	monitorApi := MonitorApi{*repository, appConfig}
	r.Group(func(mux chi.Router) {
		mux.Post("/auth/login", monitorApi.AuthLogin)
		mux.Post("/auth/validate", monitorApi.AuthValidate)
		mux.Handle("/__graphql", graph.NewPlaygroundHandler("/graphql"))

	})
	auth := AuthenticationMiddleware{*repository}
	r.Group(func(mux chi.Router) {
		mux.Use(auth.Middleware)
		mux.Mount("/graphql", graph.NewHandler(*repository, appConfig))
		mux.Post("/snapshots/create", monitorApi.SnapshotCreate)
		mux.Post("/installs/sync/preview", monitorApi.InstallSyncPreview)
		mux.Post("/installs/sync", monitorApi.InstallSync)
	})

	return r, nil
}
