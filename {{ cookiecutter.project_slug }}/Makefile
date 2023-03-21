DATE       ?= $(shell TZ=Asia/Shanghai date +%FT%T%z)
VERSION    ?= $(shell git rev-parse HEAD)
GO_VERSION ?= $(shell go version | cut -d ' ' -f 3-)

# show help
help:
	@echo ''
	@echo 'Usage:'
	@echo ' make [target]'
	@echo ''
	@echo 'Targets:'
	@awk '/^[a-zA-Z\-\_0-9]+:/ { \
	helpMessage = match(lastLine, /^# (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")-1); \
			helpMessage = substr(lastLine, RSTART + 2, RLENGTH); \
			printf "\033[36m%-22s\033[0m %s\n", helpCommand,helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)


.PHONY: build
# build
build:
	@ mkdir -p bin/
	@ go build -v -ldflags "-s -w -X {{ cookiecutter.project_slug }}.BuildTime=${DATE} -X {{ cookiecutter.project_slug }}.Version=${VERSION} -X '{{ cookiecutter.project_slug }}.GoVersion=${GO_VERSION}'" -o ./bin/{{ cookiecutter.project_slug }} ./cmd/{{ cookiecutter.project_slug }}

.PHONY: run
# run server
run:
	@ go run -ldflags "-s -w -X {{ cookiecutter.project_slug }}.BuildTime=${DATE} -X {{ cookiecutter.project_slug }}.Version=${VERSION} -X '{{ cookiecutter.project_slug }}.GoVersion=${GO_VERSION}'" ./cmd/{{ cookiecutter.project_slug }} run


.PHONY: test
# test
test:
	@ go test ./... -coverprofile=coverage.out


.PHONY: gen
# gen gorm-gen
gen:
	@ rm -rf ./internal/dal/query
	@ go run ./cmd/gen/main.go

lint:
	golines -m 180 -w --reformat-tags .

clear-tag:
	@ git tag | grep t- | tail -r | tail -n +2 | tail -r | xargs git push --delete origin  
	@ git tag | grep t- | tail -r | tail -n +2 | tail -r | xargs git tag -d
