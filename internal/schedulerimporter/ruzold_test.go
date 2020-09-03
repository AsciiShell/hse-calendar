package schedulerimporter

import (
	"testing"
	"time"

	"github.com/asciishell/hse-calendar/pkg/environment"

	"github.com/stretchr/testify/require"

	"github.com/asciishell/hse-calendar/internal/client"
)

func TestRuzOld_GetLessons(t *testing.T) {
	r := require.New(t)
	testMail := environment.GetStr("TEST_EMAIL", "")
	const testDuration = time.Hour * 24 * 7 * 2
	c := client.Client{Email: testMail}
	start := time.Date(2019, 9, 1, 0, 0, 0, 0, time.UTC)
	end := start.Add(testDuration)
	_, err := NewRuzOld().GetLessons(c, start, end)
	r.NoError(err, "Error during fetching data")
	//r.True(len(lessons) > 0, "Server return no data")
}
