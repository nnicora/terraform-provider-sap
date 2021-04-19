SRC=main.go
APP=terraform-provider-sap
BUILD_DIR=${CURDIR}/dist
VERSION=$(shell cat VERSION)

.PHONY: all prepare-dist linux darwin clean

all: build

deps:
	GOSUMDB=off go mod vendor
	go mod verify

prepare-dist:
	if [ ! -d "$(BUILD_DIR)" ]; then mkdir -p $(BUILD_DIR)/linux $(BUILD_DIR)/darwin; fi

build: clean darwin

darwin: deps prepare-dist
	env GOOS=darwin GOARCH=amd64 go build \
	-a \
	-o $(BUILD_DIR)/darwin/$(APP) \
	-ldflags "-X main.version=$(VERSION) -X main.buildTime=`date -u +%Y%m%d.%H%M%S` -X main.revision=`git rev-parse HEAD`" \
	$(SRC)

linux: deps prepare-dist
	env GOOS=linux GOARCH=amd64 go build \
	-a \
	-o $(BUILD_DIR)/linux/$(APP) \
	-ldflags "-X main.version=$(VERSION) -X main.buildTime=`date -u +%Y%m%d.%H%M%S` -X main.revision=`git rev-parse HEAD`" \
	$(SRC)

windows: deps prepare-dist
	env GOOS=windows GOARCH=amd64 go build \
	-a \
	-o $(BUILD_DIR)/windows/$(APP).exe \
	-ldflags "-X main.version=$(VERSION) -X main.buildTime=`date -u +%Y%m%d.%H%M%S` -X main.revision=`git rev-parse HEAD`" \
	$(SRC)

clean:
	rm -rf *.out
	rm -rf dist/