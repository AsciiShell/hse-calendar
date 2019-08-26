package lesson

import (
	"time"

	"github.com/asciishell/hse-calendar/internal/client"
)

type GroupedLessons struct {
	ID         int           `json:"-" gorm:"PRIMARY_KEY"`
	Client     client.Client `json:"-" gorm:"foreignkey:ClientID"`
	ClientID   uint          `json:"-" gorm:"NOT NULL"`
	Day        time.Time     `json:"day" gorm:"NOT NULL"`
	IsSelected bool          `json:"-" gorm:"DEFAULT FALSE"`
	Lessons    []Lesson      `json:"lessons" gorm:"foreignkey:GroupedID"`
}
