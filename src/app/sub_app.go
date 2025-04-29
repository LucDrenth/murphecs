package app

import (
	"github.com/lucdrenth/murph_engine/src/log"
)

type ID int

type SubApp interface {
	Run(exitChannel <-chan struct{}, isDoneChannel chan<- bool)
	Logger() log.Logger
	AddStartupSystem(schedule Schedule, system System)
	AddStartupSchedule(schedule Schedule)
	AddSystem(schedule Schedule, system System)
	AddSchedule(schedule Schedule)
	AddCleanupSystem(schedule Schedule, system System)
	AddCleanupSchedule(schedule Schedule)
	AddResource(resource Resource)
}
