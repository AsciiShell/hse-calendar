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
	URL        string
	Debug      bool
	Migrate    bool
	MigrateNum int
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
		if err := result.Migrate(credential.MigrateNum); err != nil {
			defer db.Close()
			return nil, errors.Wrapf(err, "can't apply migration number %v", credential.MigrateNum)
		}
		logger.Info("Migration complete")
	}
	return &result, nil
}

//nolint:gochecknoglobals
var migrations = []string{`create table clients
(
    id          SERIAL PRIMARY KEY,
    email       VARCHAR(50) NOT NULL,
    hse_ruz_id  INTEGER,
    google_code VARCHAR(50) NOT NULL UNIQUE
);
CREATE TABLE lessons
(
    id           SERIAL PRIMARY KEY,
    begin        TIMESTAMP WITH TIME ZONE NOT NULL,
    "end"        TIMESTAMP WITH TIME ZONE NOT NULL,
    name         TEXT                     NOT NULL,
    building     TEXT,
    auditorium   TEXT,
    lecturer     TEXT,
    kind_of_work TEXT,
    stream       TEXT,
    owner_refer  INTEGER                  NOT NULL REFERENCES clients (ID),
    created_at   TIMESTAMP WITH TIME ZONE NOT NULL
);`}

func (p *PostgresGormStorage) Migrate(index int) error {
	if index < 0 || index > len(migrations) {
		return errors.New("migration with such index does not exist")
	}
	if err := p.DB.Exec(migrations[index]).Error; err != nil {
		return errors.Cause(err)
	}
	return nil
}

func (p *PostgresGormStorage) GetClients() ([]client.Client, error) {
	var result []client.Client
	if err := p.DB.Find(&result).Error; err != nil {
		return nil, errors.Wrapf(err, "can't get clients")
	}
	return result, nil
}

func (p *PostgresGormStorage) GetLessonsFor(c client.Client, day time.Time) (lesson.GroupedLessons, error) {
	var result []lesson.Lesson
	p.DB.Where("begin::date = ?", day.Format("2006-1-2"))
	return lesson.GroupedLessons{
		Date:    day,
		Lessons: result,
	}, nil
}

func (p *PostgresGormStorage) GetDiffBetween(c client.Client, start *time.Time, end *time.Time) ([]lesson.Lesson, error) {
	panic("implement me")
}

func (p *PostgresGormStorage) AddDiff(c client.Client, lessons []lesson.Lesson) error {
	panic("implement me")
}

func (p *PostgresGormStorage) SetLessonsFor(c client.Client, groupedLessons []lesson.GroupedLessons) error {
	panic("implement me")
}
