package main

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
}

// Allows for outputing the details (deliverable 4.)
// func (version *SlimVersion) String() string

// aggregate is called within any accessor to compute aggregate values
// func (version *SlimVersion) aggregate()

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
// func (app *SlimApp) GetVersions() []SlimVersion

// GetReleaseHistory returns a Pointer to a map of Times pointing to SlimVersions (deliverable 3.)
// func (app *SlimApp) GetReleaseHistory() map[uint32]SlimVersion

// aggregate is called within the accessors to rebuild if needed
// func (app *SlimApp) aggregate()
