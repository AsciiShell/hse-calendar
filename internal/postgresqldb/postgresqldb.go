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
var migrations = []string{`CREATE TABLE clients
(
    id          SERIAL PRIMARY KEY,
    email       VARCHAR(50) NOT NULL,
    hse_ruz_id  INTEGER,
    google_code VARCHAR(50) NOT NULL UNIQUE
);
CREATE TABLE grouped
(
    id          SERIAL PRIMARY KEY,
    client_id   INTEGER NOT NULL REFERENCES clients (ID) ON DELETE CASCADE,
    day         date    NOT NULL,
    is_selected BOOLEAN NOT NULL DEFAULT FALSE,
    UNIQUE (client_id, day)
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
    grouped_id   INTEGER                  NOT NULL REFERENCES grouped (ID) ON DELETE CASCADE,
    created_at   TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now()
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

func (p *PostgresGormStorage) CreateClient(c *client.Client) error {
	if err := p.DB.Create(c).Error; err != nil {
		return errors.Wrapf(err, "can't get clients")
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
	if err := p.DB.Where("grouped_id IN (SELECT id FROM grouped WHERE day::date = ? and client_id = ?)", day.Format("2006-1-2"), c.ID).Find(&result).Error; err != nil {
		return lesson.GroupedLessons{}, errors.Wrapf(err, "can't get lessons")
	}
	return lesson.GroupedLessons{
		Day:     day,
		Lessons: result,
	}, nil
}

func (p *PostgresGormStorage) SetLessonsFor(c client.Client, groupedLessons lesson.GroupedLessons) error {
	t := p.DB.Begin()
	defer t.Commit()
	if err := p.DB.Where("day::date = ? and client_id = ?", groupedLessons.Day.Format("2006-1-2"), c.ID).Delete(lesson.GroupedLessons{}).Error; err != nil {
		t.Rollback()
		return errors.Wrapf(err, "can't delete old lessons")
	}
	//groupedLessons.Client = c
	groupedLessons.ClientID = c.ID
	if err := p.DB.Create(&groupedLessons).Error; err != nil {
		t.Rollback()
		return errors.Wrapf(err, "can't create new lessons group")
	}
	return nil
}

func (p *PostgresGormStorage) GetNewLessonsFor(c client.Client, start time.Time, end time.Time) ([]lesson.GroupedLessons, error) {
	var result []lesson.GroupedLessons
	if err := p.DB.
		Where("client_id = ? AND (day::date BETWEEN ? AND ?) AND NOT is_selected",
			c.ID,
			start.Format("2006-1-2"),
			end.Format("2006-1-2")).
		Find(&result).Error; err != nil {
		return nil, errors.Wrapf(err, "can't read new lessons")
	}
	if err := p.DB.Model(lesson.GroupedLessons{}).
		Where("client_id = ? AND (day::date BETWEEN ? AND ?) AND NOT is_selected",
			c.ID,
			start.Format("2006-1-2"),
			end.Format("2006-1-2")).
		UpdateColumn("is_selected", true).Error; err != nil {
		return nil, errors.Wrapf(err, "can't read new lessons")
	}
	return result, nil
}
