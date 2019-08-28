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
const FetchDuration = time.Hour * 24 * 30 * 2

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
func (b Background) FetchClient(c client.Client) error {
	start := time.Now()
	end := start.Add(FetchDuration)
	lessons, err := b.importer.GetLessons(c, start, end)
	if err != nil {
		return errors.Wrapf(err, "can't get lessons for %+v", c)
	}
	grouped := lesson.GroupLessons(lessons)
	for i := range grouped {
		newLessons := grouped[i]
		oldLessons, err := b.storage.GetLessonsFor(c, grouped[i].Day)
		if err != nil {
			return errors.Cause(err)
		}
		if newLessons.Equal(oldLessons) {
			continue
		}

		if err := b.storage.SetLessonsFor(c, newLessons); err != nil {
			return errors.Cause(err)
		}
	}
	return nil
}
func (b Background) FetchAllClients() error {

	clients, err := b.storage.GetClients()
	if err != nil {
		return errors.Wrapf(err, "can't fetch clients from storage")
	}
	for _, c := range clients {
		go func(client2 client.Client) {

			if err := b.FetchClient(client2); err != nil {
				b.logger.WithError(err)
			}
			b.logger.Infof("client %v handled successfully", client2)
		}(c)
	}
	return nil
}
func NewBackground(logger log.Logger, storage storage.Storage, rerun chan interface{}, importer schedulerimporter.Getter) Background {
	result := Background{logger: logger, storage: storage, rerun: rerun, importer: importer}
	result.RunFetchDiff()
	return result
}
