# Challenge Summary
We have a monitoring service which regularly pings an Order Management service and records metrics that indicate the health of that service. The Order Management service responds with the following fields:

- timestamp: the unix time when pinging the service, in seconds
- version: the git SHA of the code run by the micro service
- query_time: how long it took the micro service to generate its response, in nanoseconds

# Assumptions
The monitoring service pings the services every hour on the hour

# Sample Data
```json
[
    {
        "timestamp": 1536051600,
        "version": "356a192b7913b04c54574d18c28d46e6395428ab",
        "query_time": 189
    },
    {
        "timestamp": 1536832800,
        "version": "77de68daecd823babbb58edb1c8e14d7106e83bb",
        "query_time": 124
    }
]
```
# Deliverables
1. Find the minimum, average and maximum query times by version. 
2. Find the best and worst performing releases.
3. Using the health data, reconstruct the release history of the service.
4. Print output to stdout. 
5. Be able to provide the completed assessment via a publicly accessible code repository or a compressed file that includes the source code.

# SlimAPM Interface
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