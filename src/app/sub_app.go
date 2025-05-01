package app

type ID int

type SubApp interface {
	Run(exitChannel <-chan struct{}, isDoneChannel chan<- bool)
	AddSystem(Schedule, System) SubApp
	AddSchedule(Schedule, ScheduleType) SubApp
	AddResource(Resource) SubApp
}
