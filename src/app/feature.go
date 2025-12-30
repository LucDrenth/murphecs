package app

import (
	"fmt"
	"reflect"

	"github.com/lucdrenth/murphecs/src/utils"
)

type IFeature interface {
	Init()
	GetResources() []Resource
	GetSystems() []FeatureSystem

	// GetAndInitNestedFeatures recursively gets and Inits all features. It needs to init them so
	// that we can get its nested features.
	GetAndInitNestedFeatures() []IFeature
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
	features  []IFeature
}

type FeatureSystem struct {
	schedule Schedule
	system   System
	source   string
}

func (feature *Feature) AddSystem(schedule Schedule, system System) *Feature {
	feature.systems = append(feature.systems, FeatureSystem{schedule, system, utils.Caller(2, SystemErrorPackageDepth)})
	return feature
}

func (feature *Feature) AddResource(resource Resource) *Feature {
	feature.resources = append(feature.resources, resource)
	return feature
}

// AddFeature adds a nested feature.
//
// Systems of the nested feature will be added after that of the parent feature.
func (feature *Feature) AddFeature(f IFeature) *Feature {
	feature.features = append(feature.features, f)
	return feature
}

func (feature *Feature) GetResources() []Resource {
	return feature.resources
}

func (feature *Feature) GetSystems() []FeatureSystem {
	return feature.systems
}

func (feature *Feature) GetAndInitNestedFeatures() []IFeature {
	result := []IFeature{}

	for _, nestedFeature := range feature.features {
		nestedFeature.Init()
		result = append(result, nestedFeature)
		result = append(result, nestedFeature.GetAndInitNestedFeatures()...)
	}

	return result
}

func validateFeature(feature IFeature) error {
	initHasPointerReceiver, err := utils.MethodHasPointerReceiver(feature, "Init")
	if err != nil {
		return fmt.Errorf("failed to add feature %s: failed to validate: %v", reflect.TypeOf(feature).String(), err)
	}
	if !initHasPointerReceiver {
		return fmt.Errorf("failed to add feature %s: Init must be pointer receiver", reflect.TypeOf(feature).String())
	}

	return nil
}
