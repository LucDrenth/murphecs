package app

import (
	"fmt"

	"github.com/lucdrenth/murph_engine/src/ecs"
	"github.com/lucdrenth/murph_engine/src/log"
)

type Schedule string

type Scheduler struct {
	schedules map[Schedule]*SystemSet
	order     []Schedule
}

func NewScheduler() Scheduler {
	return Scheduler{
		schedules: map[Schedule]*SystemSet{},
		order:     []Schedule{},
	}
}

func (s *Scheduler) AddSchedule(schedule Schedule) error {
	if _, exists := s.schedules[schedule]; exists {
		return fmt.Errorf("schedule already exists")
	}

	s.schedules[schedule] = &SystemSet{}
	s.order = append(s.order, schedule)

	return nil
}

func (s *Scheduler) AddSystem(schedule Schedule, system System, world *ecs.World, logger log.Logger) error {
	systemSet, exists := s.schedules[schedule]
	if !exists {
		return fmt.Errorf("schedule %s does not exist", schedule)
	}

	return systemSet.add(system, world, logger)
}

func (s *Scheduler) GetSystemSets() ([]*SystemSet, error) {
	if len(s.order) != len(s.schedules) {
		return nil, fmt.Errorf("order of length %d does not match schedules of length %d", len(s.order), len(s.schedules))
	}

	result := make([]*SystemSet, len(s.order))

	for i, schedule := range s.order {
		systemSet, ok := s.schedules[schedule]
		if !ok {
			return nil, fmt.Errorf("schedule %s from schedule order does not exist", schedule)
		}

		result[i] = systemSet
	}

	return result, nil
}
