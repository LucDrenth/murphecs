package app

import (
	"github.com/lucdrenth/murph_engine/src/log"
)

type ID int

type SubApp interface {
	Run(exitChannel <-chan struct{}, isDoneChannel chan<- bool)
	Logger() log.Logger
	AddSystem(Schedule, System)
	AddSchedule(Schedule, ScheduleType)
	AddResource(Resource)
}
