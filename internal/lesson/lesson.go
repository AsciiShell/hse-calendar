package lesson

import (
	"time"

	"github.com/asciishell/hse-calendar/internal/client"
)

type Lesson struct {
	ID         int            `json:"id" gorm:"PRIMARY_KEY;AUTO_INCREMENT"`
	Begin      time.Time      `json:"begin" gorm:"NOT NULL"`
	End        time.Time      `json:"end" gorm:"NOT NULL"`
	Name       string         `json:"name" gorm:"NOT NULL"`
	Building   string         `json:"building"`
	Auditorium string         `json:"auditorium"`
	Lecturer   string         `json:"lecturer"`
	KindOfWork string         `json:"kindOfWork"`
	Stream     string         `json:"stream"`
	Owner      *client.Client `json:"-" gorm:"NOT NULL"`
	CreatedAt  time.Time      `json:"created_at" gorm:"NOT NULL"`
}
