package postgresqldb

import (
	"time"

	"github.com/asciishell/hse-calendar/internal/client"
	"github.com/asciishell/hse-calendar/internal/lesson"

	"github.com/asciishell/hse-calendar/pkg/log"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"

	// Registry postgres
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type PostgresGormStorage struct {
	DB *gorm.DB
}

type DBCredential struct {
	URL     string
	Debug   bool
	Migrate bool
}

func NewPostgresGormStorage(credential DBCredential) (*PostgresGormStorage, error) {
	var err error
	var db *gorm.DB
	logger := log.New()
	db, err = gorm.Open("postgres", credential.URL)
	if err != nil {
		return nil, errors.Wrapf(err, "can't connect to database, dsn %s", credential.URL)
	}
	if err = db.DB().Ping(); err != nil {
		return nil, errors.Wrapf(err, "can't ping database, dsn %s", credential.URL)
	}
	db.LogMode(credential.Debug)
	result := PostgresGormStorage{DB: db}
	if credential.Migrate {
		result.Migrate()
		logger.Info("Migration complete")
	}
	return &result, nil
}

func (p *PostgresGormStorage) Migrate() {
	p.DB.AutoMigrate(&client.Client{}, &lesson.Lesson{})
}

func (p *PostgresGormStorage) GetClients() ([]client.Client, error) {
	panic("implement me")
}

func (p *PostgresGormStorage) GetLessonsFor(c client.Client, start *time.Time, end *time.Time) ([]lesson.Lesson, error) {
	panic("implement me")
}

func (p *PostgresGormStorage) AddDiff([]lesson.Lesson) error {
	panic("implement me")
}

func (p *PostgresGormStorage) GetDiffBetween(c client.Client, start *time.Time, end *time.Time) ([]lesson.Lesson, error) {
	panic("implement me")
}

func (p *PostgresGormStorage) SetLessonsFor(c client.Client, date time.Time, lessons []lesson.Lesson) error {
	panic("implement me")
}
