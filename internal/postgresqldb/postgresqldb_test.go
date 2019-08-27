package postgresqldb

import (
	"fmt"
	"testing"
	"time"

	"github.com/asciishell/hse-calendar/internal/lesson"

	"github.com/stretchr/testify/require"

	"github.com/asciishell/hse-calendar/pkg/environment"

	"github.com/asciishell/hse-calendar/internal/client"
)

func TestPostgresGormStorage_CreateClient(t *testing.T) {
	r := require.New(t)
	db, err := NewPostgresGormStorage(DBCredential{
		URL:        environment.GetStr("DB_URL_TEST", ""),
		Debug:      true,
		Migrate:    true,
		MigrateNum: 0,
	})
	r.NoError(err, "can't use database")
	defer db.DB.Close()

	c := client.Client{
		ID:         0,
		Email:      environment.GetStr("TEST_EMAIL", ""),
		HSERuzID:   0,
		GoogleCode: "testtesttest",
	}
	err = db.CreateClient(&c)
	r.NoError(err, "can't create client")
	defer func() {
		db.DB.Delete(c)
	}()
	r.NotEqual(0, c.ID, "client ID not set %+v", c)
	fmt.Printf("client after insert:%+v", c)

	g := lesson.GroupedLessons{
		Client:     c,
		ClientID:   c.ID,
		Day:        time.Date(2019, 9, 2, 0, 0, 0, 0, time.Local),
		IsSelected: false,
	}
	g.Lessons = append(g.Lessons, lesson.Lesson{
		Begin: time.Date(2019, 9, 2, 8, 0, 0, 0, time.Local),
		End:   time.Date(2019, 9, 2, 9, 0, 0, 0, time.Local),
		Name:  "Test name",
	})
	err = db.SetLessonsFor(c, g)
	r.NoError(err, "can't lessons group")
	cli, err := db.GetClients()
	r.NoError(err, "can't get all clients")
	for i := range cli {
		less, err := db.GetLessonsFor(cli[i], time.Date(2019, 9, 2, 0, 0, 0, 0, time.UTC))
		r.NoError(err, "can't get lessons for client")
		r.True(len(less.Lessons) > 0, "no lessons")
		err = db.SetLessonsFor(cli[i], less)
		r.NoError(err, "can't update lessons")
	}
}
