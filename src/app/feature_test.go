package app

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type testFeatureA struct{ Feature }
type testFeatureB struct{ Feature }
type testFeatureC struct{ Feature }
type testFeatureD struct{ Feature }

func (f *testFeatureA) Init() {}
func (f *testFeatureB) Init() {
	f.AddFeature(&testFeatureC{})
}
func (f *testFeatureC) Init() {
	f.AddFeature(&testFeatureA{}).AddFeature(&testFeatureD{})
}
func (f *testFeatureD) Init() {}

func TestNestedFeature(t *testing.T) {
	assert := assert.New(t)

	var feature IFeature = &testFeatureA{}
	feature.Init()
	assert.Len(feature.GetAndInitNestedFeatures(), 0)

	feature = &testFeatureB{}
	feature.Init()
	assert.Len(feature.GetAndInitNestedFeatures(), 3)

	feature = &testFeatureC{}
	feature.Init()
	assert.Len(feature.GetAndInitNestedFeatures(), 2)

	feature = &testFeatureD{}
	feature.Init()
	assert.Len(feature.GetAndInitNestedFeatures(), 0)
}
