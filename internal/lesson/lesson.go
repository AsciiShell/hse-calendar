package lesson

import (
	"fmt"
	"sort"
	"time"
)

type Lesson struct {
	ID         int            `json:"id" gorm:"PRIMARY_KEY"`
	Begin      time.Time      `json:"begin" gorm:"NOT NULL"`
	End        time.Time      `json:"end" gorm:"NOT NULL"`
	Name       string         `json:"name" gorm:"NOT NULL"`
	Building   string         `json:"building"`
	Auditorium string         `json:"auditorium"`
	Lecturer   string         `json:"lecturer"`
	KindOfWork string         `json:"kindOfWork"`
	Stream     string         `json:"stream"`
	CreatedAt  time.Time      `json:"created_at" gorm:"NOT NULL"`
	Grouped    GroupedLessons `json:"-" gorm:"foreignkey:GroupedID"`
	GroupedID  uint
}

const Day = time.Hour * 24

// Take lessons should be processing
func Handle(lessons []Lesson) {
	go func() {
		grouped := GroupLessons(lessons)
		// Get data from database
		// Compare
		// Save new lessons
		// Save diff
		fmt.Printf("%+v", grouped)
	}()
}

// Group lessons by date
//
// Return ordered by date list of structures,
// where every element ordered by begin time
func GroupLessons(lessons []Lesson) []GroupedLessons {
	result := make(map[time.Time][]Lesson)
	for i := range lessons {
		day := lessons[i].Begin.Truncate(Day)
		result[day] = append(result[day], lessons[i])
	}
	out := make([]GroupedLessons, 0, len(result))
	for k := range result {
		slice := result[k]
		sort.Slice(slice, func(i, j int) bool {
			return slice[i].Begin.Before(slice[j].Begin)
		})
		result[k] = slice
		out = append(out, GroupedLessons{
			Day:     k,
			Lessons: result[k],
		})
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].Day.Before(out[j].Day)
	})
	return out
}
