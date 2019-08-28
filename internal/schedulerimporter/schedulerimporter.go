package schedulerimporter

import (
	"time"

	"github.com/asciishell/hse-calendar/internal/client"
	"github.com/asciishell/hse-calendar/internal/lesson"
)

type Getter interface {
	GetLessons(client client.Client, start time.Time, end time.Time, endSignal chan<- interface{}) ([]lesson.Lesson, error)
}
