package postgresqldb

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/asciishell/hse-calendar/pkg/environment"

	"github.com/asciishell/hse-calendar/internal/client"
	"github.com/asciishell/hse-calendar/pkg/log"
)

func TestPostgresGormStorage_CreateClient(t *testing.T) {
	r := require.New(t)
	db, err := NewPostgresGormStorage(DBCredential{
		URL:        environment.GetStr("DB_URL_TEST", ""),
		Debug:      true,
		Migrate:    true,
		MigrateNum: 0,
	})
	if err != nil {
		log.New().Fatalf("can't use database:%s", err)
	}
	defer db.DB.Close()

	c := client.Client{
		ID:         0,
		Email:      environment.GetStr("TEST_EMAIL", ""),
		HSERuzID:   0,
		GoogleCode: "testtesttest",
	}
	defer func() {
		db.DB.Delete(c)
	}()
	err = db.CreateClient(&c)
	r.NoError(err)
	r.NotEqual(0, c.ID, "client ID not set %+v", c)
	fmt.Printf("client after insert:%+v\n%+v", err, c)
	cli, err := db.GetClients()
	fmt.Printf("client select err:%+v\n", err)
	for i := range cli {
		less, _ := db.GetLessonsFor(cli[i], time.Date(2019, 9, 2, 0, 0, 0, 0, time.UTC))

		fmt.Printf("client no:%d: %+v %+v\n", i, cli[i], err)
		for j := range less.Lessons {
			fmt.Printf("\t\tlesson:%+v\n", less.Lessons[j])

		}
	}
	less, err := db.GetLessonsFor(cli[0], time.Date(2019, 9, 2, 0, 0, 0, 0, time.UTC))
	err = db.SetLessonsFor(cli[0], less)
	fmt.Printf("set error:%+v\n", err)
}
