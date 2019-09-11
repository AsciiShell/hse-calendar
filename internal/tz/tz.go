package tz

import (
	"time"

	"github.com/pkg/errors"
)

// nolint:gochecknoglobals
var loc *time.Location

func SetTimezone(tz string) error {
	location, err := time.LoadLocation(tz)
	if err != nil {
		return errors.Wrapf(err, "can't set timezone to %s", tz)
	}
	loc = location
	return nil
}

func GetTime(t time.Time) time.Time {
	return t.In(loc)
}

func GetTimezone() *time.Location {
	return loc
}
