BUILD_FILES = $(shell go list -f '{{range .GoFiles}}{{$$.Dir}}/{{.}}\
{{end}}' ./...)

GH_VERSION ?= $(shell git describe --tags 2>/dev/null || git rev-parse --short HEAD)
DATE_FMT = +%Y-%m-%d
BUILD_DATE = $(shell date "$(DATE_FMT)")
IMAGE_REPO = "quay.io/nissessenap/poddeleter"

print-%  : ; @echo $* = $($*)

bin/pc: $(BUILD_FILES)
	@go build -trimpath -o bin/poddeleter ./main.go

bin/container:
	docker build --build-arg BUILD_DATE=$(BUILD_DATE) --build-arg VERSION=$(GH_VERSION) . -t $(IMAGE_REPO):$(GH_VERSION)
.PHONY: bin/container

bin/push:
	docker push $(IMAGE_REPO):$(GH_VERSION)

clean:
	rm -rf ./bin ./share
.PHONY: clean

test:
	go test ./...
.PHONY: test
