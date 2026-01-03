package app

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestScheduleOrder(t *testing.T) {
	const (
		schedule1 Schedule = "schedule1"
		schedule2 Schedule = "schedule2"
		schedule3 Schedule = "schedule3"
	)

	t.Run("ScheduleLast", func(t *testing.T) {
		assert := assert.New(t)
		var err error

		order := ScheduleLast{}
		schedules := []Schedule{}

		schedules, err = order.insert(schedule1, schedules)
		assert.NoError(err)
		assert.Equal([]Schedule{schedule1}, schedules)

		schedules, err = order.insert(schedule2, schedules)
		assert.NoError(err)
		assert.Equal([]Schedule{schedule1, schedule2}, schedules)

		schedules, err = order.insert(schedule3, schedules)
		assert.NoError(err)
		assert.Equal([]Schedule{schedule1, schedule2, schedule3}, schedules)
	})

	t.Run("ScheduleBefore", func(t *testing.T) {
		assert := assert.New(t)

		order := ScheduleBefore{Other: schedule2}
		var err error
		var schedules []Schedule

		schedules, err = order.insert(schedule1, []Schedule{})
		assert.Error(err) // Other does not exist

		schedules, err = order.insert(schedule1, []Schedule{schedule2})
		assert.NoError(err)
		assert.Equal([]Schedule{schedule1, schedule2}, schedules)
	})

	t.Run("ScheduleAfter", func(t *testing.T) {
		assert := assert.New(t)

		order := ScheduleAfter{Other: schedule2}
		var err error
		var schedules []Schedule

		schedules, err = order.insert(schedule3, []Schedule{})
		assert.Error(err) // Other does not exist

		schedules, err = order.insert(schedule3, []Schedule{schedule2})
		assert.NoError(err)
		assert.Equal([]Schedule{schedule2, schedule3}, schedules)
	})
}
