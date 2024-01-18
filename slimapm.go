package main

import "fmt"

type SlimMetric struct {
	Timestamp uint32 `json:"timestamp"`
	QueryTime uint16 `json:"query_time"`
}

// SlimVersion stores the hash of the version along with aggregates (deliverable 1.)
type SlimVersion struct {
	hash            string
	max             uint16
	min             uint16
	avg             float32
	timestamps      []uint32
	queryTime       []uint16
	shouldAggregate bool
}

// includeMetrics is used to append each Metric
func (version *SlimVersion) IncludeMetrics(metrics SlimMetric) {
	version.timestamps = append(version.timestamps, metrics.Timestamp)
	version.queryTime = append(version.queryTime, metrics.QueryTime)
	version.shouldAggregate = true
}

// String for outputing the details (deliverable 4.)
func (version SlimVersion) String() string {
	version.aggregate()
	return fmt.Sprintf("VERSION: %s\n  max: %d\n  min: %d\n  avg: %.2f", version.hash, version.max, version.min, version.avg)
}

// aggregate is called within any accessor to compute aggregate values
func (version *SlimVersion) aggregate() {
	if !version.shouldAggregate {
		return
	}
	var sum float32
	count := 0
	for _, time := range version.queryTime {
		sum += float32(time)
		count++
		if version.min == 0.0 || time < version.min {
			version.min = time
		}
		if version.max == 0.0 || time > version.max {
			version.max = time
		}
	}
	version.avg = sum / float32(count)
	version.shouldAggregate = false
}

// SlimApp contains a slice of SlimVersions, along w/ the best and worst (deliverable 2.)
type SlimApp struct {
	versions        map[string]*SlimVersion
	best            *SlimVersion
	worst           *SlimVersion
	shouldAggregate bool
}

// NewSlimApp loads the healthcheck metrics into a SlimApp and returns a ptr to the object
func NewSlimApp() *SlimApp {
	return &SlimApp{versions: make(map[string]*SlimVersion)}
}

// AddRaw will construct the SlimVersions slice and check for the best and worst
func (app *SlimApp) AddVersionMetric(version string, metric SlimMetric) error {
	if _, ok := app.versions[version]; !ok {
		app.versions[version] = &SlimVersion{hash: version}
	}
	app.versions[version].IncludeMetrics(metric)
	app.shouldAggregate = true
	return nil
}

// GetVersions returns a slice of SlimVersions (deliverable 1.)
func (app *SlimApp) GetVersions() []SlimVersion {
	app.aggregate()
	cap := len(app.versions)
	versions := make([]SlimVersion, 0, cap)
	for _, version := range app.versions {
		versions = append(versions, *version)
	}
	return versions
}

// GetReleaseHistory returns a Pointer to a map of Times pointing to SlimVersions (deliverable 3.)
// func (app *SlimApp) GetReleaseHistory() map[uint32]SlimVersion

// aggregate is called within the accessors to rebuild if needed
func (app *SlimApp) aggregate() {
	if !app.shouldAggregate {
		return
	}
	for _, version := range app.versions {
		version.aggregate()
		if app.best == nil || version.avg < app.best.avg {
			app.best = version
		}
		if app.worst == nil || version.avg > app.worst.avg {
			app.worst = version
		}
	}
	app.shouldAggregate = false
}

// String returns the high level overview about the best and worst release
func (app SlimApp) String() string {
	app.aggregate()
	return fmt.Sprintf("Best Release: %v\nWorst Release: %v", app.best.hash, app.worst.hash)
}
