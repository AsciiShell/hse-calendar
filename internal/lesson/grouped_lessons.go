package lesson

import (
	"time"

	"github.com/asciishell/hse-calendar/internal/client"
)

type GroupedLessons struct {
	ID         int           `json:"-" gorm:"PRIMARY_KEY"`
	Client     client.Client `json:"-" gorm:"foreignkey:ClientID"`
	ClientID   int           `json:"-" gorm:"NOT NULL"`
	Day        time.Time     `json:"date" gorm:"NOT NULL"`
	IsSelected bool          `json:"-" gorm:"DEFAULT FALSE"`
	Lessons    []Lesson      `json:"lessons" gorm:"foreignkey:GroupedID"`
}

func (GroupedLessons) TableName() string {
	return "grouped"
}
