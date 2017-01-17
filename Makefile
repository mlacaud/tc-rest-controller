DOCKER_REPO=mlacaud
DOCKER_IMAGE=tc-rest-controller
DOCKER_TAG=latest

BUILD_OUTPUT=bin/
BUILD_BIN=tc-rest-controller

EXEC=docker-build


all: $(EXEC)

build:
	go build -tags netgo -a -v -o $(BUILD_OUTPUT)$(BUILD_BIN)

docker-build: build
	docker build -t $(DOCKER_REPO)/$(DOCKER_IMAGE):$(DOCKER_TAG) .

docker-push: build docker-build
	docker push $(DOCKER_REPO)/$(DOCKER_IMAGE):$(DOCKER_TAG)

clean:
	rm -r $(BUILD_OUTPUT)


clean-docker:
	docker rmi -f $(DOCKER_REPO)/$(DOCKER_IMAGE):$(DOCKER_TAG)

mrproper: clean-docker clean
