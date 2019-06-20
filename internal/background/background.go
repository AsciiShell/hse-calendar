package background

import (
	"context"
	"time"

	"github.com/asciishell/HSE_calendar/internal/storage"
	"github.com/asciishell/HSE_calendar/pkg/log"
)

type Background struct {
	logger  log.Logger
	storage storage.Storage
	rerun   <-chan interface{}
}

const Timeout = time.Hour

func (b Background) RunFetchDiff() {
	go func() {
		for {
			ctx, cancel := context.WithTimeout(context.Background(), Timeout)
			select {
			case <-ctx.Done():
			case <-b.rerun:
			}
			cancel()
			// TODO do update
		}
	}()
}
func NewBackground(logger log.Logger, storage storage.Storage, rerun chan interface{}) Background {
	result := Background{logger: logger, storage: storage, rerun: rerun}
	result.RunFetchDiff()
	return result
}
