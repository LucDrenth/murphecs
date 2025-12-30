package app

import (
	"fmt"

	"github.com/lucdrenth/murphecs/src/ecs"
)

type Schedule string

type Scheduler struct {
	systems map[Schedule]*ScheduleSystems
	order   []Schedule
}

func NewScheduler() Scheduler {
	return Scheduler{
		systems: map[Schedule]*ScheduleSystems{},
		order:   []Schedule{},
	}
}

func (s *Scheduler) AddSchedule(schedule Schedule, scheduleSystemsId ScheduleSystemsId) error {
	if _, exists := s.systems[schedule]; exists {
		return fmt.Errorf("schedule already exists")
	}

	s.systems[schedule] = &ScheduleSystems{id: scheduleSystemsId}
	s.order = append(s.order, schedule)

	return nil
}

func (s *Scheduler) AddSystem(schedule Schedule, system System, source string, world *ecs.World, outerWorlds *map[ecs.WorldId]*ecs.World, logger Logger, resources *resourceStorage, eventStorage *EventStorage) error {
	scheduleSystems, exists := s.systems[schedule]
	if !exists {
		return fmt.Errorf("schedule %s does not exist", schedule)
	}

	return scheduleSystems.add(system, source, world, outerWorlds, logger, resources, eventStorage)
}

func (s *Scheduler) GetScheduleSystems() ([]*ScheduleSystems, error) {
	if len(s.order) != len(s.systems) {
		return nil, fmt.Errorf("order of length %d does not match schedules of length %d", len(s.order), len(s.systems))
	}

	result := make([]*ScheduleSystems, len(s.order))

	for i, schedule := range s.order {
		scheduleSystems, ok := s.systems[schedule]
		if !ok {
			return nil, fmt.Errorf("schedule %s from schedule order does not exist", schedule)
		}

		result[i] = scheduleSystems
	}

	return result, nil
}

func (s *Scheduler) NumberOfSystems() uint {
	result := uint(0)
	for _, scheduleSystems := range s.systems {
		for _, systemGroup := range scheduleSystems.systemGroups {
			result += uint(len(systemGroup.systems))
		}
	}
	return result
}
