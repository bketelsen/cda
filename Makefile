GIT_SUMMARY := $(shell git describe --tags --dirty --always)
REPO=bketelsen/cda.ms
DOCKER_IMAGE := $(REPO):$(GIT_SUMMARY)

default: publish

repo:
	@echo $(DOCKER_IMAGE)

build:
	@go install github.com/bketelsen/cda
	@GOOS=linux CGO_ENABLE=0 go build -o cda main.go
	@docker build -t $(DOCKER_IMAGE) .
	@docker tag $(DOCKER_IMAGE) $(REPO)

push:
	@docker push $(DOCKER_IMAGE)
	@docker push $(REPO)

clean:
	@rm -rf dist/

install:
	@go install github.com/bketelsen/cda

release-snapshot: clean
	@goreleaser --snapshot

release: clean
	@github-release-notes -org bketelsen -repo cda.ms -include-commits > .releasenotes
	@goreleaser --release-notes=.releasenotes

release-notes:
	@github-release-notes -org bketelsen -repo cda.ms -include-commits
