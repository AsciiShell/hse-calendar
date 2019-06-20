package lesson

import "time"

type GroupedLessons struct {
	Date    time.Time `json:"date"`
	Lessons []Lesson  `json:"lessons"`
}
