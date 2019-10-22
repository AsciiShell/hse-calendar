package main

import (
	"net/http"
	"time"

	"github.com/asciishell/hse-calendar/internal/background"
	"github.com/asciishell/hse-calendar/internal/postgresqldb"
	"github.com/asciishell/hse-calendar/internal/schedulerimporter"
	"github.com/asciishell/hse-calendar/internal/tz"
	"github.com/asciishell/hse-calendar/pkg/environment"
	"github.com/asciishell/hse-calendar/pkg/log"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

type config struct {
	DB          postgresqldb.DBCredential
	HTTPAddress string
	HTTPTimeout time.Duration
	MaxRequests int
	PrintConfig bool
	Timezone    string
}

func loadConfig() config {
	cfg := config{}
	cfg.DB.URL = environment.GetStr("DB_URL", "")
	cfg.DB.Debug = environment.GetBool("DB_DEBUG", false)
	cfg.DB.Migrate = environment.GetBool("DB_MIGRATE", false)
	cfg.DB.MigrateNum = environment.GetInt("DB_MIGRATE_NUM", -1)
	cfg.MaxRequests = environment.GetInt("MAX_REQUESTS", 100)
	cfg.HTTPAddress = environment.GetStr("ADDRESS", ":8000")
	cfg.HTTPTimeout = environment.GetDuration("HTTP_TIMEOUT", 500*time.Second)
	cfg.PrintConfig = environment.GetBool("PRINT_CONFIG", false)
	cfg.Timezone = environment.GetStr("TZ", "")
	if cfg.PrintConfig {
		log.New().Infof("%+v", cfg)
	}
	return cfg
}
func main() {
	cfg := loadConfig()
	if err := tz.SetTimezone(cfg.Timezone); err != nil {
		log.New().Fatalf("can't set timezone:%s", err)
	}
	db, err := postgresqldb.NewPostgresGormStorage(cfg.DB)
	if err != nil {
		log.New().Fatalf("can't use database:%s", err)
	}
	defer func() {
		_ = db.DB.Close()
	}()
	logger := log.New()

	rerunChan := make(chan interface{})
	handler := NewHandler(logger, db, rerunChan)
	background.NewBackground(logger, db, rerunChan, schedulerimporter.NewRuzOld())

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Throttle(cfg.MaxRequests))
	r.Use(middleware.Timeout(cfg.HTTPTimeout))

	r.Route("/v1", func(r chi.Router) {
		r.Post("/diff/{email}/{id}", handler.GetDiff)
		r.Get("/run", handler.Rerun)
		r.Route("/client/{email}/{id}", func(r chi.Router) {
			r.Post("/", handler.CreateClient)
			r.Delete("/", handler.DeleteClient)
		})
	})
	if err := http.ListenAndServe(cfg.HTTPAddress, r); err != nil {
		logger.Fatalf("server error:%s", err)
	}
}
