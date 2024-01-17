CUR_DIR=$(shell pwd)
ARTIFACTS_DIR=${CUR_DIR}/.artifacts

.PHONY: slimapm
slimapm:
	go run .

.PHONY: ci
ci: lint test

.PHONY: lint
lint:
	docker run --rm -v ${CUR_DIR}:/app -w /app golangci/golangci-lint:v1.54 \
		golangci-lint run -v

.PHONY: test
test:
	@rm -rf ${ARTIFACTS_DIR} || true
	@mkdir ${ARTIFACTS_DIR}
	go test -v ./. -coverprofile ${ARTIFACTS_DIR}/cover.out || true
	go tool cover -html ${ARTIFACTS_DIR}/cover.out -o ${ARTIFACTS_DIR}/cover.html