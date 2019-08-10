package schedulerimporter

import (
	"testing"
	"time"

	"github.com/asciishell/HSE_calendar/internal/client"
	"github.com/stretchr/testify/require"
)

func TestSourceRuzOld_GetLessons(t *testing.T) {
	r := require.New(t)
	const testMail = "aepodchezertsev@edu.hse.ru"
	const testDuration = time.Hour * 24 * 7 * 2
	c := client.Client{Email: testMail}
	start := time.Date(2019, 9, 1, 0, 0, 0, 0, time.UTC)
	end := start.Add(testDuration)
	lessons, err := SourceRuzOld{}.GetLessons(c, start, end)
	r.NoError(err, "Error during fetching data")
	r.True(len(lessons) > 0, "Server return no data")
}
