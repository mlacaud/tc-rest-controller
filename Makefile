GO_APP=tc-rest-controller
GO_MAIN=tc-rest-controller

GOPATH=${PWD}

BUILD_PACKAGE=tc-rest-controller
BUILD_OUTPUT=bin/
BUILD_BIN=tc-rest-controller

DOCKER_REPO=mlacaud
DOCKER_IMAGE=tc-rest-controller
DOCKER_TAG=latest

EXEC=build


all: $(EXEC)

get: 
	go get $(GO_MAIN)

install: get 
	go install $(GO_MAIN)

build: get 
	go build -tags netgo -a -v -o $(BUILD_OUTPUT)$(BUILD_BIN) $(BUILD_PACKAGE)

docker-build: build
	docker build -t $(DOCKER_REPO)/$(DOCKER_IMAGE):$(DOCKER_TAG) .

docker-push: docker-build
	docker push $(DOCKER_REPO)/$(DOCKER_IMAGE):$(DOCKER_TAG)

clean:
	rm -r $(BUILD_OUTPUT)

clean-go:
	rm -r pkg src/github.com

clean-docker:
	docker rmi -f $(DOCKER_REPO)/$(DOCKER_IMAGE):$(DOCKER_TAG)

mrproper: clean-go clean clean-docker
