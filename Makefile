IMPORT_PATH ?= github.com/asciishell/hse-calendar
BUILD_DIR ?= bin
PKG_DIR = .pkg
GOROOT ?= /usr/local/go

# Common constants
BINARIES_DIR := cmd
BINARIES := $$(find $(BINARIES_DIR) -maxdepth 1 \( ! -iname "$(BINARIES_DIR)" \) -type d -exec basename {} \;)
VERSION := $(shell git describe --long --tags --always --abbrev=8)
ifndef CI
	BRANCH = $(shell git rev-parse --abbrev-ref HEAD)
else
	BRANCH = $(CI_COMMIT_REF_NAME)
endif


DOCKER_BUILDER_FLAGS := --rm=true -u $$(id -u):$$(id -g) -v $(CURDIR):/go/src/$(IMPORT_PATH) -w /go/src/$(IMPORT_PATH)
DOCKER_BUILDER_IMAGE := golang:1.17

DOCKER_IMAGE_SPACE ?= asciishell
ifeq ($(BRANCH), master)
	 DOCKER_IMAGE_TAG = $(BRANCH)
else
	DOCKER_IMAGE_TAG = $(BRANCH)_$(VERSION)
endif

OSFLAG 				:=
ifeq ($(OS),Windows_NT)
	OSFLAG = "WIN"
else
	UNAME_S := $(shell uname -s)
	ifeq ($(UNAME_S),Linux)
		OSFLAG = LINUX
	endif
	ifeq ($(UNAME_S),Darwin)
		OSFLAG = OSX
	endif
endif

all:
	$(info $$OSFLAG is [${OSFLAG}])
	$(info $$BRANCH is [${BRANCH}])
	$(info $$VERSION is [${VERSION}])
	$(info $$DOCKER_IMAGE_TAG is [${DOCKER_IMAGE_TAG}])
# Build targets
$(BUILD_DIR):
	cp -rf $(GOROOT)/pkg/linux_amd64 $(CURDIR)/$(PKG_DIR) || true
	GOCACHE=`pwd`/.cache GOOS=linux GOARCH=amd64 CGO_ENABLED=0 GOBIN="" go install -pkgdir $(CURDIR)/$(PKG_DIR) ./...
	for bin in $(BINARIES); do \
		GOCACHE=`pwd`/.cache GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -pkgdir $(CURDIR)/$(PKG_DIR) -o $(BUILD_DIR)/$$bin $(IMPORT_PATH)/$(BINARIES_DIR)/$$bin;\
    done

.PHONY: clean
clean:
	@-rm -rf $(BUILD_DIR)
	@-rm -rf $(PKG_DIR)

# Docker targets
.PHONY: docker-build
docker-build: clean
	docker run $(DOCKER_BUILDER_FLAGS) $(DOCKER_BUILDER_IMAGE) make

.PHONY: docker-images
docker-images:
	for bin in $(BINARIES); do \
		docker build --rm --pull --tag $(DOCKER_IMAGE_SPACE)/$$bin:$(DOCKER_IMAGE_TAG) --file $(BINARIES_DIR)/$$bin/Dockerfile .;\
	done

.PHONY: docker-push
docker-push:
	for bin in $(BINARIES); do \
		docker tag $(DOCKER_IMAGE_SPACE)/$$bin:$(DOCKER_IMAGE_TAG) $(DOCKER_IMAGE_SPACE)/$$bin:$(DOCKER_IMAGE_TAG);\
		docker push $(DOCKER_IMAGE_SPACE)/$$bin:$(DOCKER_IMAGE_TAG);\
	done

.PHONY: docker-clean
docker-clean:
	for bin in $(BINARIES); do \
		docker rmi -f $$(docker images $(DOCKER_IMAGE_SPACE)/$$bin:$(DOCKER_IMAGE_TAG) -q);\
	done

.PHONY: lint
lint:
	golangci-lint run -c .golangci.yml ./...

.PHONY: test
test:
	if [ $(OSFLAG) = "WIN" ]; then \
		go test -v ./... ; \
	else \
		TIMEOUT_MULTIPLY=10 go test -v -race ./... ; \
	fi


.PHONY: ci-deploy
ci-deploy:
	ssh -t root@$$TARGET_HOST 'cd calendar-api && docker-compose stop'
	scp ./docker-compose.yml root@$$TARGET_HOST:calendar-api/docker-compose.yml
	ssh -t root@$$TARGET_HOST 'cd calendar-api && IMAGE_TAG=$(DOCKER_IMAGE_TAG) docker-compose pull && IMAGE_TAG=$(DOCKER_IMAGE_TAG) docker-compose up -d'
