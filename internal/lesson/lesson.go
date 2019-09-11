package lesson

import (
	"encoding/json"
	"sort"
	"time"

	"github.com/asciishell/hse-calendar/internal/tz"
)

const GoogleDateFormat = "January 02, 2006 15:04:05 MST"
const dateLayout = "2006-01-02"

type Lesson struct {
	ID         int             `json:"id" gorm:"PRIMARY_KEY"`
	Begin      time.Time       `json:"begin" gorm:"NOT NULL"`
	End        time.Time       `json:"end" gorm:"NOT NULL"`
	Name       string          `json:"name" gorm:"NOT NULL"`
	Building   string          `json:"building"`
	Auditorium string          `json:"auditorium"`
	Lecturer   string          `json:"lecturer"`
	KindOfWork string          `json:"kindOfWork"`
	Stream     string          `json:"stream"`
	CreatedAt  time.Time       `json:"created_at" gorm:"NOT NULL"`
	Grouped    *GroupedLessons `json:"-" gorm:"foreignkey:GroupedID"`
	GroupedID  uint
}

type lessonJSON struct {
	Begin      string `json:"begin"`
	End        string `json:"end"`
	Name       string `json:"name"`
	Building   string `json:"building"`
	Auditorium string `json:"auditorium"`
	Lecturer   string `json:"lecturer"`
	KindOfWork string `json:"kindOfWork"`
	Stream     string `json:"stream"`
	CreatedAt  string `json:"created_at"`
}

func (l Lesson) MarshalJSON() ([]byte, error) {
	return json.Marshal(lessonJSON{
		Begin:      l.Begin.Format(GoogleDateFormat),
		End:        l.End.Format(GoogleDateFormat),
		Name:       l.Name,
		Building:   l.Building,
		Auditorium: l.Auditorium,
		Lecturer:   l.Lecturer,
		KindOfWork: l.KindOfWork,
		Stream:     l.Stream,
		CreatedAt:  tz.GetTime(l.CreatedAt).Format(time.RFC3339),
	})
}

const Day = time.Hour * 24

func (l Lesson) Equal(l2 Lesson) bool {
	return l.Begin.Equal(l2.Begin) &&
		l.End.Equal(l2.End) &&
		l.Name == l2.Name &&
		l.Building == l2.Building &&
		l.Auditorium == l2.Auditorium &&
		l.Lecturer == l2.Lecturer &&
		l.KindOfWork == l2.KindOfWork &&
		l.Stream == l2.Stream
}

// Group lessons by date
//
// Return ordered by date list of structures,
// where every element ordered by begin time
func GroupLessons(lessons []Lesson) []GroupedLessons {
	if len(lessons) == 0 {
		return nil
	}
	result := make(map[string][]Lesson)
	minDate := lessons[0].Begin
	maxDate := lessons[0].Begin
	for i := range lessons {
		day := lessons[i].Begin.Format(dateLayout)
		result[day] = append(result[day], lessons[i])
		if lessons[i].Begin.Before(minDate) {
			minDate = lessons[i].Begin
		} else if lessons[i].Begin.After(maxDate) {
			maxDate = lessons[i].Begin
		}
	}
	for d := minDate; d.Before(maxDate); d = d.Add(time.Hour * 24) {
		day := d.Format(dateLayout)
		if _, ok := result[day]; !ok {
			result[day] = make([]Lesson, 0)
		}
	}
	out := make([]GroupedLessons, 0, len(result))
	for k := range result {
		slice := result[k]
		sort.Slice(slice, func(i, j int) bool {
			return slice[i].Begin.Before(slice[j].Begin)
		})
		result[k] = slice
		d, _ := time.Parse(dateLayout, k)
		out = append(out, GroupedLessons{
			Day:     d,
			Lessons: result[k],
		})
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].Day.Before(out[j].Day)
	})
	return out
}
