package main

import (
	"encoding/base64"
	"strings"
	"time"

	"github.com/asciishell/HSE_calendar/internal/background"

	"github.com/asciishell/HSE_calendar/internal/postgresqldb"
	"github.com/asciishell/HSE_calendar/pkg/environment"
	"github.com/asciishell/HSE_calendar/pkg/log"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

type config struct {
	DB          postgresqldb.DBCredential
	HTTPAddress string
	HTTPTimeout time.Duration
	MaxRequests int
	PrintConfig bool
}

func loadConfig() config {
	cfg := config{}
	cfg.DB.URL = environment.GetStr("BASE64_DB_URL", "")
	cfg.DB.Debug = environment.GetBool("DB_DEBUG", false)
	cfg.DB.Migrate = environment.GetBool("DB_MIGRATE", false)
	cfg.MaxRequests = environment.GetInt("MAX_REQUESTS", 100)
	cfg.HTTPAddress = environment.GetStr("ADDRESS", ":8000")
	cfg.HTTPTimeout = environment.GetDuration("HTTP_TIMEOUT", 500*time.Second)
	cfg.PrintConfig = environment.GetBool("PRINT_CONFIG", false)
	if cfg.PrintConfig {
		log.New().Infof("%+v", cfg)
	}
	if cfg.DB.URL != "" {
		dbURL, err := base64.StdEncoding.DecodeString(cfg.DB.URL)
		if err != nil {
			log.New().Fatalf("Can't parse base64 db url: %s", err)
		}
		cfg.DB.URL = strings.TrimSpace(string(dbURL))
	}

	return cfg
}
func main() {
	cfg := loadConfig()

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
	background.NewBackground(logger, db, rerunChan)

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Throttle(cfg.MaxRequests))
	r.Use(middleware.Timeout(cfg.HTTPTimeout))

	r.Route("/v1", func(r chi.Router) {
		r.Get("/diff", handler.GetDiff)

	})
}
