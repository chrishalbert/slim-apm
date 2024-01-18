# SlimAPM
Slim-APM is a golang library that can be used to analyze healthcheck responses, which could be used for an APM.

## Running
Options:
1. You can simply run `go run .`. (you need `go` installed)
1. You can run with `make slimapm`

## Developing
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
