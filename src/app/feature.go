package app

type IFeature interface {
	Init()
	getResources() []Resource
	getSystems() []FeatureSystem
}

// Feature is a set of resources and systems that will be initialized and added to an app before the
// app runs. This is useful because, in contrast to adding systems directly to an app, resources that
// are used in system params won't have to be added before adding the system.
//
// Added systems and resources will not be directly verified. It will be done once the app processes
// the features, which is done before running the app.
type Feature struct {
	resources []Resource
	systems   []FeatureSystem
}

type FeatureSystem struct {
	schedule Schedule
	system   System
}

func (feature *Feature) AddSystem(schedule Schedule, system System) *Feature {
	feature.systems = append(feature.systems, FeatureSystem{schedule, system})
	return feature
}

func (feature *Feature) AddResource(resource Resource) *Feature {
	feature.resources = append(feature.resources, resource)
	return feature
}

func (feature *Feature) getResources() []Resource {
	return feature.resources
}

func (feature *Feature) getSystems() []FeatureSystem {
	return feature.systems
}
