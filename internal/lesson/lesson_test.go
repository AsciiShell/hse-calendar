package lesson

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

const testLessons = "test.json"
const testLessonsGrouped = "test_grouped.json"

func TestGroupLessons(t *testing.T) {
	r := require.New(t)

	plan, err := ioutil.ReadFile(filepath.Clean(testLessons))
	r.NoError(err)
	var testLessonsData []Lesson
	err = json.Unmarshal(plan, &testLessonsData)
	r.NoError(err)
	plan, err = ioutil.ReadFile(filepath.Clean(testLessonsGrouped))
	r.NoError(err)
	var testGroupedLessons []GroupedLessons
	err = json.Unmarshal(plan, &testGroupedLessons)
	r.NoError(err)

	grouped := GroupLessons(testLessonsData)
	eq := reflect.DeepEqual(grouped, testGroupedLessons)
	if !eq {
		fmt.Printf("%+v\n", grouped)
		fmt.Printf("%+v\n", testGroupedLessons)
	}
	r.True(eq)
}
