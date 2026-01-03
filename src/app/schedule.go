package app

import (
	"fmt"
	"slices"

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

func (s *Scheduler) AddSchedule(schedule Schedule, scheduleSystemsId ScheduleSystemsId, order ScheduleOrder, isPaused bool) (err error) {
	if _, exists := s.systems[schedule]; exists {
		return ErrScheduleAlreadyExists
	}

	scheduleSystems := &ScheduleSystems{id: scheduleSystemsId}
	if isPaused {
		scheduleSystems.isPaused.Store(true)
	}
	s.systems[schedule] = scheduleSystems

	s.order, err = order.insert(schedule, s.order)
	if err != nil {
		return fmt.Errorf("ScheduleOrder failed to insert: %w", err)
	}

	return nil
}

func (s *Scheduler) AddSystem(schedule Schedule, system System, source string, world *ecs.World, outerWorlds *map[ecs.WorldId]*ecs.World, logger Logger, resources *resourceStorage, eventStorage *EventStorage) error {
	scheduleSystems, exists := s.systems[schedule]
	if !exists {
		return fmt.Errorf("%w: %s", ErrScheduleNotFound, schedule)
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

type ScheduleOrder interface {
	insert(Schedule, []Schedule) ([]Schedule, error)
}

var (
	_ ScheduleOrder = &ScheduleLast{}
	_ ScheduleOrder = &ScheduleBefore{}
	_ ScheduleOrder = &ScheduleAfter{}
)

type ScheduleLast struct{}

func (scheduleOrder ScheduleLast) insert(schedule Schedule, schedules []Schedule) ([]Schedule, error) {
	return append(schedules, schedule), nil
}

type ScheduleBefore struct {
	Other Schedule
}

func (scheduleOrder ScheduleBefore) insert(schedule Schedule, schedules []Schedule) ([]Schedule, error) {
	i := slices.Index(schedules, scheduleOrder.Other)
	if i == -1 {
		return schedules, fmt.Errorf("%w: '%s'", ErrScheduleNotFound, scheduleOrder.Other)
	}

	return slices.Insert(schedules, i, schedule), nil
}

type ScheduleAfter struct {
	Other Schedule
}

func (scheduleOrder ScheduleAfter) insert(schedule Schedule, schedules []Schedule) ([]Schedule, error) {
	i := slices.Index(schedules, scheduleOrder.Other)
	if i == -1 {
		return schedules, fmt.Errorf("%w: '%s'", ErrScheduleNotFound, scheduleOrder.Other)
	}

	return slices.Insert(schedules, i+1, schedule), nil
}
