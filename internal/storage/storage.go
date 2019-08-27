package storage

import (
	"time"

	"github.com/asciishell/hse-calendar/internal/client"
	"github.com/asciishell/hse-calendar/internal/lesson"
)

type Storage interface {
	// Create structure if not exists
	Migrate(index int) error
	// Create new client
	CreateClient(c *client.Client) error
	// Return list of clients
	GetClients() ([]client.Client, error)
	// Get list of lessons for client in particular date
	GetLessonsFor(c client.Client, day time.Time) (lesson.GroupedLessons, error)
	// Save actual lessons for client
	SetLessonsFor(c client.Client, groupedLessons lesson.GroupedLessons) error
	// Get unselected lessons for client between dates
	GetNewLessonsFor(c client.Client, start time.Time, end time.Time) ([]lesson.GroupedLessons, error)
}
