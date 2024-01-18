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

func TestAggregateSlimVersion(t *testing.T) {
	var slimVersionScenarios = []struct {
		testCase        string
		shouldAggregate bool
		queryTime       []uint16
		expectedMin     uint16
		expectedMax     uint16
		expectedAvg     float32
	}{
		{
			testCase:        "should not aggregate if flag set",
			shouldAggregate: false,
			queryTime:       []uint16{1, 2, 3},
			expectedMin:     0,
			expectedMax:     0,
			expectedAvg:     0.0,
		},
		{
			testCase:        "should aggregate if flag set",
			shouldAggregate: true,
			queryTime:       []uint16{1, 2, 3},
			expectedMin:     1,
			expectedMax:     3,
			expectedAvg:     2.0,
		},
	}

	for _, testCase := range slimVersionScenarios {
		version := &SlimVersion{
			shouldAggregate: testCase.shouldAggregate,
			queryTime:       testCase.queryTime,
		}
		t.Run(testCase.testCase, func(t *testing.T) {
			version.aggregate()
			assert.Equal(t, testCase.expectedMin, version.min)
			assert.Equal(t, testCase.expectedMax, version.max)
			assert.Equal(t, testCase.expectedAvg, version.avg)
			// After running aggregate(), the flag should be set to false
			assert.False(t, version.shouldAggregate)
		})
	}
}

func TestAggregateSlimApp(t *testing.T) {
	var slimAppScenarios = []struct {
		testCase        string
		shouldAggregate bool
		versions        map[string]*SlimVersion
		expectedBest    string
		expectedWorst   string
	}{
		{
			testCase:        "should not aggregate if flag is false",
			shouldAggregate: false,
			versions: map[string]*SlimVersion{
				"abc": {
					hash:            "abc",
					timestamps:      []uint32{2, 4, 6},
					queryTime:       []uint16{2, 4, 6},
					shouldAggregate: true,
				},
			},
			expectedBest:  "",
			expectedWorst: "",
		},
		{
			testCase:        "should aggregate versions logically",
			shouldAggregate: true,
			versions: map[string]*SlimVersion{
				"abc": {
					hash:            "abc",
					timestamps:      []uint32{1, 2, 3},
					queryTime:       []uint16{1, 2, 3},
					shouldAggregate: true,
				},
				"cde": {
					hash:            "cde",
					timestamps:      []uint32{4, 5, 6},
					queryTime:       []uint16{4, 5, 6},
					shouldAggregate: true,
				},
				"fgh": {
					hash:            "fgh",
					timestamps:      []uint32{7, 8, 9},
					queryTime:       []uint16{7, 8, 9},
					shouldAggregate: true,
				},
			},
			expectedBest:  "abc",
			expectedWorst: "fgh",
		},
	}

	for _, testCase := range slimAppScenarios {
		slimApp := &SlimApp{
			versions:        testCase.versions,
			shouldAggregate: testCase.shouldAggregate,
		}

		t.Run(testCase.testCase, func(t *testing.T) {
			var actualBest, actualWorst string
			slimApp.aggregate()
			if slimApp.best != nil {
				actualBest = slimApp.best.hash
			}
			if slimApp.worst != nil {
				actualWorst = slimApp.worst.hash
			}
			assert.Equal(t, testCase.expectedBest, actualBest)
			assert.Equal(t, testCase.expectedWorst, actualWorst)
			// After running aggregate(), the flag should be set to false
			assert.False(t, slimApp.shouldAggregate)
		})
	}
}

func TestGetVersions(t *testing.T) {
	var slimAppScenarios = []struct {
		testCase        string
		shouldAggregate bool
		versions        map[string]*SlimVersion
		expectedMin     uint16
		expectedMax     uint16
		expectedAvg     float32
	}{
		{
			testCase: "should return a slice of aggregated versions",
			versions: map[string]*SlimVersion{
				"abc": {
					hash:            "abc",
					timestamps:      []uint32{2, 4, 6},
					queryTime:       []uint16{2, 4, 6},
					shouldAggregate: true,
				},
			},
			expectedMin: 2,
			expectedMax: 6,
			expectedAvg: 4.0,
		},
	}
	for _, testCase := range slimAppScenarios {
		slimApp := &SlimApp{
			versions:        testCase.versions,
			shouldAggregate: true,
		}

		versions := slimApp.GetVersions()

		t.Run(testCase.testCase, func(t *testing.T) {
			assert.Equal(t, testCase.expectedMin, versions[0].min)
			assert.Equal(t, testCase.expectedMax, versions[0].max)
			assert.Equal(t, testCase.expectedAvg, versions[0].avg)
			// After running aggregate(), the flag should be set to false
			assert.False(t, slimApp.shouldAggregate)
		})
	}
}

func TestIncludeMetrics(t *testing.T) {
	t.Run("should add metrics and set to aggregate", func(t *testing.T) {
		version := &SlimVersion{}
		version.IncludeMetrics(SlimMetric{Timestamp: 123, QueryTime: 456})
		assert.Contains(t, version.timestamps, uint32(123))
		assert.Contains(t, version.queryTime, uint16(456))
		assert.True(t, version.shouldAggregate)
	})
}

func TestStringSlimApp(t *testing.T) {
	t.Run("should return the best and worst releases", func(t *testing.T) {
		app := NewSlimApp()
		app.best = &SlimVersion{shouldAggregate: false, hash: "thebest"}
		app.worst = &SlimVersion{shouldAggregate: false, hash: "theworst"}
		actual := app.String()
		assert.Equal(t, "Best Release: thebest\nWorst Release: theworst", actual)
	})
}

func TestStringSlimVersion(t *testing.T) {
	t.Run("should return the best and worst releases", func(t *testing.T) {
		version := &SlimVersion{
			shouldAggregate: false,
			hash:            "hash",
			avg:             3.0,
			min:             1,
			max:             5,
		}
		actual := version.String()
		assert.Equal(t, "VERSION: hash\n  max: 5\n  min: 1\n  avg: 3.00", actual)
	})
}

func TestGetReleaseHistory(t *testing.T) {
	t.Run("should return a map organized by initial version date", func(t *testing.T) {
		app := NewSlimApp()
		app.AddVersionMetric("abc", SlimMetric{Timestamp: 1})
		app.AddVersionMetric("ghi", SlimMetric{Timestamp: 5})
		app.AddVersionMetric("def", SlimMetric{Timestamp: 2})
		app.AddVersionMetric("def", SlimMetric{Timestamp: 3})
		app.AddVersionMetric("def", SlimMetric{Timestamp: 4})
		releaseHistory := app.GetReleaseHistory()
		assert.Equal(t, struct {
			hash  string
			start uint32
		}{hash: "abc", start: 1}, (*releaseHistory)[0])
		assert.Equal(t, struct {
			hash  string
			start uint32
		}{hash: "def", start: 2}, (*releaseHistory)[1])
		assert.Equal(t, struct {
			hash  string
			start uint32
		}{hash: "ghi", start: 5}, (*releaseHistory)[2])
		assert.Equal(t, 3, len(*releaseHistory))
	})
}
