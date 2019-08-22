package background

import (
	"context"
	"fmt"
	"time"

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
const FetchDuration = time.Hour * 24 * 30

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
	for i := range lessons {
		fmt.Printf("%+v\n", lessons[i])
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
		}(c)
	}
	return nil
}
func NewBackground(logger log.Logger, storage storage.Storage, rerun chan interface{}, importer schedulerimporter.Getter) Background {
	result := Background{logger: logger, storage: storage, rerun: rerun, importer: importer}
	result.RunFetchDiff()
	return result
}
