package background

import (
	"context"
	"time"

	"github.com/asciishell/HSE_calendar/internal/schedulerimporter"

	"github.com/pkg/errors"

	"github.com/asciishell/HSE_calendar/internal/client"
	"github.com/asciishell/HSE_calendar/internal/storage"
	"github.com/asciishell/HSE_calendar/pkg/log"
)

type Background struct {
	logger   log.Logger
	storage  storage.Storage
	rerun    <-chan interface{}
	importer []schedulerimporter.Getter
}

const SleepTime = time.Hour

func (b Background) RunFetchDiff() {
	go func() {
		for {
			ctx, cancel := context.WithTimeout(context.Background(), SleepTime)
			select {
			case <-ctx.Done():
			case <-b.rerun:
			}
			cancel()
			clients, err := b.storage.GetClients()
			if err != nil {
				b.logger.WithError(errors.Wrapf(err, "can't fetch clients from storage"))
				continue
			}
			for _, c := range clients {
				go func(client2 client.Client) {

					if err := b.FetchUser(client2); err != nil {
						b.logger.WithError(err)
					}
				}(c)
			}

		}
	}()
}

func (b Background) FetchUser(c client.Client) error {
	for i := range b.importer{
		//b.importer[i].GetLessons(c, nil, nil, )
	}
	// В бэкграунд над
	return nil
}

func NewBackground(logger log.Logger, storage storage.Storage, rerun chan interface{}, importer []schedulerimporter.Getter) Background {
	result := Background{logger: logger, storage: storage, rerun: rerun, importer: importer}
	result.RunFetchDiff()
	return result
}
