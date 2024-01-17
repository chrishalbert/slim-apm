package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSlimApp(t *testing.T) {
	t.Run("should initialize Slim App", func(t *testing.T) {
		app := NewSlimApp()
		assert.NotNil(t, app.versions)
		assert.Nil(t, app.best)
		assert.Nil(t, app.worst)
		assert.False(t, app.shouldAggregate)
	})
}

func TestAddVersionMetric(t *testing.T) {
	app := NewSlimApp()
	version := "abc"
	metric := SlimMetric{Timestamp: 1, QueryTime: 123}
	t.Run("should create a new SlimVersion if adding a metric for a new version", func(t *testing.T) {
		_, ok := app.versions[version]
		assert.False(t, ok)

		app.AddVersionMetric(version, metric)
		assert.Nil(t, app.best)
		assert.Nil(t, app.worst)
		_, ok = app.versions[version]
		assert.True(t, ok)
		assert.True(t, app.shouldAggregate)
	})

	t.Run("should use existing SlimVersion if adding a metric for an existing version", func(t *testing.T) {
		appVersionExisting, ok := app.versions[version]
		assert.True(t, ok)

		app.AddVersionMetric(version, metric)
		assert.Nil(t, app.best)
		assert.Nil(t, app.worst)

		appVersion, ok := app.versions[version]
		assert.True(t, ok)
		assert.True(t, app.shouldAggregate)

		assert.Equal(t, appVersionExisting, appVersion)

		assert.Equal(t, 2, len(app.versions[version].timestamps))
		assert.Equal(t, 2, len(app.versions[version].queryTime))
	})
}
