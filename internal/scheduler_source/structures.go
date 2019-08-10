package scheduler_source

import (
	"time"

	"github.com/pkg/errors"

	"github.com/asciishell/HSE_calendar/internal/lesson"
)

type ruzConverter interface {
	Convert() lesson.Lesson
}
type ruzOldJSON struct {
	Auditorium  string `json:"auditorium"`
	BeginLesson string `json:"beginLesson"`
	Building    string `json:"building"`
	Date        string `json:"date"`
	Discipline  string `json:"discipline"`
	EndLesson   string `json:"endLesson"`
	KindOfWork  string `json:"kindOfWork"`
	Lecturer    string `json:"lecturer"`
	Stream      string `json:"stream"`
}

func (r ruzOldJSON) Convert() (lesson.Lesson, error) {
	const timeLayout = "2006.01.02 15:04"
	name := string(r.KindOfWork[0]) + "." + r.Discipline
	start, err := time.Parse(timeLayout, r.Date+" "+r.BeginLesson)
	if err != nil {
		return lesson.Lesson{}, errors.Wrapf(err, "can't parse time %s %s", r.Date, r.BeginLesson)
	}
	end, err := time.Parse(timeLayout, r.Date+" "+r.EndLesson)
	if err != nil {
		return lesson.Lesson{}, errors.Wrapf(err, "can't parse time %s %s", r.Date, r.BeginLesson)
	}

	return lesson.Lesson{Name: name,
		KindOfWork: r.KindOfWork,
		Begin:      start,
		End:        end,
		Auditorium: r.Auditorium,
		Building:   r.Building,
		Lecturer:   r.Lecturer,
		Stream:     r.Stream}, nil
}
