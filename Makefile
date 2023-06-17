EXECUTABLE := go-whisper
GO ?= go
GOFILES := $(shell find . -name "*.go" -type f)
HAS_GO = $(shell hash $(GO) > /dev/null 2>&1 && echo "GO" || echo "NOGO" )

ifneq ($(shell uname), Darwin)
	EXTLDFLAGS = -extldflags "-static" $(null)
else
	EXTLDFLAGS =
endif

ifeq ($(HAS_GO), GO)
	GOPATH ?= $(shell $(GO) env GOPATH)
	export PATH := $(GOPATH)/bin:$(PATH)

	CGO_EXTRA_CFLAGS := -DSQLITE_MAX_VARIABLE_NUMBER=32766
	CGO_CFLAGS ?= $(shell $(GO) env CGO_CFLAGS) $(CGO_EXTRA_CFLAGS)
endif

ifeq ($(OS), Windows_NT)
	GOFLAGS := -v -buildmode=exe
	EXECUTABLE ?= $(EXECUTABLE).exe
else ifeq ($(OS), Windows)
	GOFLAGS := -v -buildmode=exe
	EXECUTABLE ?= $(EXECUTABLE).exe
else
	GOFLAGS := -v
	EXECUTABLE ?= $(EXECUTABLE)
endif

ifneq ($(DRONE_TAG),)
	VERSION ?= $(DRONE_TAG)
else
	VERSION ?= $(shell git describe --tags --always || git rev-parse --short HEAD)
endif

TAGS ?=
LDFLAGS ?= -X 'main.Version=$(VERSION)'
INCLUDE_PATH := $(abspath third_party/whisper.cpp)
LIBRARY_PATH := $(abspath third_party/whisper.cpp)

all: build

clone:
	@[ -d third_party/whisper.cpp ] || git clone https://github.com/appleboy/whisper.cpp.git third_party/whisper.cpp

dependency: clone
	@echo Build whisper
	@make -C third_party/whisper.cpp libwhisper.a

test:
	@C_INCLUDE_PATH=${INCLUDE_PATH} LIBRARY_PATH=${LIBRARY_PATH} $(GO) test -v -cover -coverprofile coverage.txt ./... && echo "\n==>\033[32m Ok\033[m\n" || exit 1

install: $(GOFILES)
	C_INCLUDE_PATH=${INCLUDE_PATH} LIBRARY_PATH=${LIBRARY_PATH} $(GO) install -v -tags '$(TAGS)' -ldflags '$(EXTLDFLAGS)-s -w $(LDFLAGS)'

build: $(EXECUTABLE)

$(EXECUTABLE): $(GOFILES)
	C_INCLUDE_PATH=${INCLUDE_PATH} LIBRARY_PATH=${LIBRARY_PATH} $(GO) build -v -tags '$(TAGS)' -ldflags '$(EXTLDFLAGS)-s -w $(LDFLAGS)' -o bin/$@

clean:
	$(GO) clean -x -i ./...
	rm -rf coverage.txt $(EXECUTABLE) $(DIST)

version:
	@echo $(VERSION)
