package app

type ID int

type SubApp interface {
	// Run runs the sub app on a separate goroutine.
	Run(exitChannel <-chan struct{}, isDoneChannel chan<- bool)

	// AddSchedule adds a new schedule to which systems can be added.
	AddSchedule(Schedule, ScheduleType) SubApp

	// AddSystem adds a new system that runs on the specified schedule. If the given system is not valid, an error will be
	// logged.
	AddSystem(Schedule, System) SubApp

	// AddResource adds a new resource to the app. There can only exist 1 resource per type per app. The resource can then by
	// used as a system param, either by reference or by value.
	AddResource(Resource) SubApp

	// AddFeature adds a feature that will be processed before running this app. Features are useful because feature resources
	// are added before the feature systems are added. This way we don't have to worry about adding resources before adding
	// systems that use them.
	AddFeature(IFeature) SubApp
}
