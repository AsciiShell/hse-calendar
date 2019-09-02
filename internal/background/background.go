package background

import (
	"context"
	"time"

	"github.com/asciishell/hse-calendar/internal/lesson"

	"github.com/asciishell/hse-calendar/internal/schedulerimporter"

	"github.com/pkg/errors"

	"github.com/asciishell/hse-calendar/internal/client"
	"github.com/asciishell/hse-calendar/internal/storage"
	"github.com/asciishell/hse-calendar/pkg/log"
)

type Background struct {
	logger   log.Logger
	storage  storage.Storage
	rerun    <-chan interface{}
	importer schedulerimporter.Getter
}

const SleepTime = time.Hour
const FetchDuration = time.Hour * 24 * 30 * 3

func (b Background) RunFetchDiff() {
	go func() {
		for {
			if err := b.FetchAllClients(); err != nil {
				b.logger.WithError(err)
			}
			b.waitSignal()
		}
	}()
}
func (b Background) waitSignal() {
	ctx, cancel := context.WithTimeout(context.Background(), SleepTime)
	defer cancel()
	select {
	case <-ctx.Done():
	case <-b.rerun:
	}
}
func (b Background) FetchClient(c client.Client, nextSignal chan interface{}) {
	start := time.Now()
	end := start.Add(FetchDuration)
	lessons, err := b.importer.GetLessons(c, start, end, nextSignal)
	if err != nil {
		b.logger.Errorf("can't get lessons for %+v: %+v", c, err)
		return
	}
	grouped := lesson.GroupLessons(lessons)
	for i := range grouped {
		newLessons := grouped[i]
		oldLessons, err := b.storage.GetLessonsFor(c, grouped[i].Day)
		if err != nil {
			b.logger.Errorf("can't get lessons from storage for %v: %+v", c, err)
		}
		if newLessons.Equal(oldLessons) {
			continue
		}
		if err := b.storage.SetLessonsFor(c, newLessons); err != nil {
			b.logger.Errorf("can't set lessons for %v: %+v", c, err)
		}
	}
	b.logger.Infof("client %v handled successfully", c)
}
func (b Background) FetchAllClients() error {

	clients, err := b.storage.GetClients()
	if err != nil {
		return errors.Wrapf(err, "can't fetch clients from storage")
	}
	nextSignal := make(chan interface{}, 1)
	for i := range clients {
		go b.FetchClient(clients[i], nextSignal)
		<-nextSignal
	}
	return nil
}
func NewBackground(logger log.Logger, storage storage.Storage, rerun chan interface{}, importer schedulerimporter.Getter) Background {
	result := Background{logger: logger, storage: storage, rerun: rerun, importer: importer}
	result.RunFetchDiff()
	return result
}
