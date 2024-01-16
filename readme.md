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