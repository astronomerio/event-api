IMAGE_NAME ?= astronomerio/clickstream-ingestion-api

GIT_COMMIT=$(shell git rev-parse --short HEAD)
GIT_DIRTY=$(shell test -n "`git status --porcelain`" && echo "+CHANGES" || true)
GIT_DESCRIBE=$(shell git describe --tags --always)
GIT_IMPORT=github.com/astronomerio/clickstream-ingestion-api/pkg/version
GOLDFLAGS=-X $(GIT_IMPORT).GitCommit=$(GIT_COMMIT)$(GIT_DIRTY) -X $(GIT_IMPORT).GitDescribe=$(GIT_DESCRIBE)

VERSION ?= SNAPSHOT-$(GIT_COMMIT)

build:
	go build -ldflags '$(GOLDFLAGS)' -tags static -o server main.go

build-image:
	docker build -t $(IMAGE_NAME):$(VERSION) .

push-image:
	docker push $(IMAGE_NAME):$(VERSION)