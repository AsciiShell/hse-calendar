package storage

import (
	"time"

	"github.com/asciishell/hse-calendar/internal/client"
	"github.com/asciishell/hse-calendar/internal/lesson"
)

type Storage interface {
	// Create structure if not exists
	Migrate()
	// Return list of clients
	GetClients() ([]client.Client, error)
	// Return list of lessons for client in particular date
	GetLessonsFor(c client.Client, day time.Time) (lesson.GroupedLessons, error)
	// Add calculated diff in lessons to database
	AddDiff(c client.Client, lessons []lesson.Lesson) error
	// Fetch diff for client and delete them
	GetDiffBetween(c client.Client, start *time.Time, end *time.Time) ([]lesson.Lesson, error)
	// Save actual lessons for client
	SetLessonsFor(c client.Client, groupedLessons []lesson.GroupedLessons) error
}
