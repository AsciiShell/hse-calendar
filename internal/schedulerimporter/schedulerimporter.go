package schedulerimporter

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/pkg/errors"

	"github.com/asciishell/HSE_calendar/internal/client"
	"github.com/asciishell/HSE_calendar/internal/lesson"
)

type Getter interface {
	GetLessons(client client.Client, start time.Time, end time.Time) ([]lesson.Lesson, error)
}

const TimeOut = time.Second * 15

type RuzOld struct{}

func (RuzOld) GetLessons(client client.Client, start time.Time, end time.Time) ([]lesson.Lesson, error) {
	const SourceURL = "http://ruz2019.hse.ru/ruzservice.svc/personlessons?language=1&receivertype=0&email=%s&fromdate=%s&todate=%s"
	const DateFormat = "2006.1.2"
	httpClient := &http.Client{
		Timeout: TimeOut,
	}
	url := fmt.Sprintf(SourceURL, client.Email, start.Format(DateFormat), end.Format(DateFormat))
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, errors.Wrap(err, "can't create request")
	}

	resp, err := httpClient.Do(req)
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

type RuzWeb struct{}

func (RuzWeb) GetLessons(client client.Client, start time.Time, end time.Time) ([]lesson.Lesson, error) {
	panic("implement me")
}