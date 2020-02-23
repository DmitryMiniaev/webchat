BUILDPATH=$(CURDIR)
GO=$(shell which go)

PROJNAME=webchat

.PHONY: all
all: build

.PHONY: build
build:
	$(GO) install

.PHONY: clean
clean:
	$(GO) clean

.PHONY: docker-build
docker-build:
	docker build -t webchat .