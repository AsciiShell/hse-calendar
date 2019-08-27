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

func (g GroupedLessons) Equal(g2 GroupedLessons) bool {
	if len(g.Lessons) != len(g2.Lessons) {
		return false
	}

baseLoop:
	for i := range g.Lessons {
		for j := range g2.Lessons {
			if g.Lessons[i].Equal(g2.Lessons[j]) {
				continue baseLoop
			}
		}
		return false
	}
	return true
}
