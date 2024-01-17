
# SlimAPM

## Setup
### As an end user...
Options:
1. You can simply run `go run .`. (you need `go` installed)
1. You can run with `make slimapm`

### As a developer...
Dependencies: You need `make`, `go`, and `docker` - docker is used for golangci-lint, per
maintainer recommendations to avoid installing as a golang mod

1. Verify code passes standards (linting and tests):
    ```bash
    make ci 
    ```
    _(this runs `make lint` and `make test`)_

1. Verify code coverage:
    ```bash
    make test
    ```
    _(then open up `./.artifacts/cover.html` in your browser)_


## Original Specs
```go
// SlimMetric allows for unmarshalling before including in a version details
type SlimMetric struct {
    Timestamp Time `json:"timestamp"`
    QueryTime uint16 `json:"query_time"`
}

// SlimVersion stores the hash of the version along with aggregates (deliverable 1.)
type SlimVersion struct {
    hash string
    max uint16
    min uint16
    avg float32
    timestamps []Time
    queryTime []uint16
    shouldAggregate bool
}

// includeMetrics is used to append each Metric
(version *SlimVersion) IncludeMetrics(metrics SlimMetric)

// Allows for outputing the details (deliverable 4.)
(version *SlimVersion) String() string

// aggregate is called within any accessor to compute aggregate values
(version *SlimVersion) aggregate()

// SlimApp contains a slice of SlimVersions, along w/ the best and worst (deliverable 2.)
type SlimApp {
    versions []SlimVersions
    best *SlimVersion
    worst *SlimVersion
    shouldAggregate bool
}

// New loads the healthcheck metrics into a SlimApp and returns a ptr to the object
New(fileBytes []byte) *SlimApp

// AddRaw will construct the SlimVersions slice and check for the best and worst
(app *SlimApp) AddRaw(healthcheck []byte) error

// GetVersions returns a slice of SlimVersions (deliverable 1.)
(app *SlimApp) GetVersions() []SlimVersion

// GetReleaseHistory returns a Pointer to a map of Times pointing to SlimVersions (deliverable 3.)
(app *SlimApp) GetReleaseHistory() map[Time]SlimVersion

// aggregate is called within the accessors to rebuild if needed
(app *SlimApp) aggregate()
```