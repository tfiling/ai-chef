LAST_TAG?=$(shell git describe --tags 2>/dev/null || echo 'latest')
IMAGE_TAG?=$(LAST_TAG)-$(shell git branch --show-current)
LOCAL_REPO := "ldg"

export DOCKER_BUILDKIT:=1

.PHONY: build
build: build-webserver

.PHONY: build-webserver
build-webserver:
	@echo "====================== building webserver ======================"
	docker build --target webserver -t $(LOCAL_REPO)/webserver:$(IMAGE_TAG) .
	@echo "====================== building webserver completed ======================"

.PHOMY: run
run: test lint
	@echo "====================== Running Local Dev Env ======================"
	@TAG=${IMAGE_TAG} docker compose -f dev_utils/local_dev_env/compose.yaml up -d

.PHONY: stop
stop:
	@echo "====================== Stopping Local Dev Env ======================"
	@TAG=${IMAGE_TAG} docker compose -f dev_utils/local_dev_env/compose.yaml down --remove-orphans -t 0

.PHONY: test
test: build
	@echo "====================== Running Tests ======================"
	docker build . --target unit-test --tag $(LOCAL_REPO)/webserver-tests:latest
	@echo "====================== Completed Running Tests ======================"


.PHONY: lint
lint: build
	@echo "====================== Running Linter ======================"
	docker build . --target linter --tag $(LOCAL_REPO)/webserver-tests:latest
	@echo "====================== Completed Running Linter ======================"

