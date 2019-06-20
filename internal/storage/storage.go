package storage

import (
	"time"

	"github.com/asciishell/HSE_calendar/internal/client"
	"github.com/asciishell/HSE_calendar/internal/lesson"
)

type Storage interface {
	// Create structure if not exists
	Migrate()
	// Return list of clients
	GetClients() ([]client.Client, error)
	// Return list of lessons for client between two dates or from begin/to end if null
	GetLessonsFor(c client.Client, start *time.Time, end *time.Time) ([]lesson.Lesson, error)
	// Add calculated diff in lessons to database
	AddDiff([]lesson.Lesson) error
	// Fetch diff for client and delete them
	GetDiffBetween(c client.Client, start *time.Time, end *time.Time) ([]lesson.Lesson, error)
	// Save actual lessons for client
	SetLessonsFor(c client.Client, date time.Time, lessons []lesson.Lesson) error
}
