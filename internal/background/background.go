package background

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/asciishell/HSE_calendar/internal/storage"
	"github.com/asciishell/HSE_calendar/pkg/log"
)

type Background struct {
	logger  log.Logger
	storage storage.Storage
	rerun   <-chan interface{}
}

const SleepTime = time.Hour
const TimeOut = time.Second * 10

func (b Background) RunFetchDiff() {
	go func() {
		for {
			ctx, cancel := context.WithTimeout(context.Background(), SleepTime)
			select {
			case <-ctx.Done():
			case <-b.rerun:
			}
			cancel()

			client := &http.Client{Timeout: TimeOut}
			b.storage
			req, err := http.NewRequest("GET", "https://postman-echo.com/get", nil)
			if err != nil {
				log.Fatal(err)
			}

			resp, err := client.Do(req)
			// always handle errors
			/**	if err != nil {
				log.Fatal(err)
			}**/
			defer resp.Body.Close()

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("%s", body)

		}
	}()
}
func NewBackground(logger log.Logger, storage storage.Storage, rerun chan interface{}) Background {
	result := Background{logger: logger, storage: storage, rerun: rerun}
	result.RunFetchDiff()
	return result
}
