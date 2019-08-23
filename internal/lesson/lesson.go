package lesson

import (
	"sort"
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

const Day = time.Hour * 24

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
		sort.Slice(result[k], func(i, j int) bool {
			return result[k][i].Begin.Before(result[k][j].Begin)
		})
		out = append(out, GroupedLessons{
			Date:    k,
			Lessons: result[k],
		})
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].Date.Before(out[j].Date)
	})
	return out
}
