package schedulerimporter

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"github.com/pkg/errors"

	"github.com/asciishell/hse-calendar/internal/client"
	"github.com/asciishell/hse-calendar/internal/lesson"
)

const timeOut = time.Second * 60
const Location = "Europe/Moscow"

type RuzOld struct {
	Mutex *sync.Mutex
}

func NewRuzOld() *RuzOld {
	return &RuzOld{Mutex: &sync.Mutex{}}
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
	Url			string `json:"url1"`
}

func (r RuzOld) GetLessons(client client.Client, start time.Time, end time.Time) ([]lesson.Lesson, error) {
	const SourceURL = "http://ruz2019.hse.ru/ruzservice.svc/personlessons?language=1&receivertype=0&email=%s&fromdate=%s&todate=%s"
	const DateFormat = "2006.1.2"
	httpClient := &http.Client{
		Timeout: timeOut,
	}
	url := fmt.Sprintf(SourceURL, client.Email, start.Format(DateFormat), end.Format(DateFormat))
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, errors.Wrap(err, "can't create request")
	}
	var resp *http.Response
	r.Mutex.Lock()
	resp, err = httpClient.Do(req)
	r.Mutex.Unlock()

	if err != nil {
		return nil, errors.Wrap(err, "can't do request")
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "can't read bytes")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("server returned %v %s:\n%s", resp.StatusCode, resp.Status, body))
	}
	var result []ruzOldJSON
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, errors.Wrap(err, "can't parse ruz old json")
	}
	var lessons []lesson.Lesson
	for i := range result {
		less, err := result[i].Convert()
		if err != nil {
			return nil, errors.Cause(err)
		}
		lessons = append(lessons, less)
	}
	return lessons, nil
}

func (r ruzOldJSON) Convert() (lesson.Lesson, error) {
	const timeLayout = "2006.01.02 15:04"
	loc, err := time.LoadLocation(Location)
	if err != nil {
		return lesson.Lesson{}, errors.Wrapf(err, "timezone %s not found", Location)
	}
	name := r.Discipline
	if len(r.KindOfWork) > 0 {
		name = string([]rune(r.KindOfWork)[0]) + "." + name
	}
	start, err := time.ParseInLocation(timeLayout, r.Date+" "+r.BeginLesson, loc)
	if err != nil {
		return lesson.Lesson{}, errors.Wrapf(err, "can't parse time %s %s", r.Date, r.BeginLesson)
	}
	end, err := time.ParseInLocation(timeLayout, r.Date+" "+r.EndLesson, loc)
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
		Stream:     r.Stream + "\n" + r.Url}, nil
}
