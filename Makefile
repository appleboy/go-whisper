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
GOLDFLAGS ?= -X 'main.Version=$(VERSION)'
INCLUDE_PATH := $(abspath third_party/whisper.cpp):$(INCLUDE_PATH)
LIBRARY_PATH := $(abspath third_party/whisper.cpp):$(LIBRARY_PATH)

ifdef WHISPER_CUBLAS
	CGO_CFLAGS      += -DGGML_USE_CUBLAS -I/usr/local/cuda/include -I/opt/cuda/include -I$(CUDA_PATH)/targets/$(UNAME_M)-linux/include
	CGO_CXXFLAGS    += -DGGML_USE_CUBLAS -I/usr/local/cuda/include -I/opt/cuda/include -I$(CUDA_PATH)/targets/$(UNAME_M)-linux/include
	EXTLDFLAGS      = -extldflags "-lcuda -lcublas -lculibos -lcudart -lcublasLt -lpthread -ldl -lrt -L/usr/local/cuda/lib64 -L/opt/cuda/lib64 -L$(CUDA_PATH)/targets/$(UNAME_M)-linux/lib"

build: $(EXECUTABLE)

$(EXECUTABLE): $(GOFILES)
	CGO_CXXFLAGS=${CGO_CXXFLAGS} CGO_CFLAGS=${CGO_CFLAGS} C_INCLUDE_PATH=${INCLUDE_PATH} LIBRARY_PATH=${LIBRARY_PATH} $(GO) build -v -tags '$(TAGS)' -ldflags '$(EXTLDFLAGS)-s -w $(GOLDFLAGS)' -o bin/$@
endif

all: build

clone:
	@[ -d third_party/whisper.cpp ] || git clone https://github.com/appleboy/whisper.cpp.git third_party/whisper.cpp

dependency: clone
	@echo Build whisper
	@make -C third_party/whisper.cpp libwhisper.a

test:
	@C_INCLUDE_PATH=${INCLUDE_PATH} LIBRARY_PATH=${LIBRARY_PATH} $(GO) test -v -cover -coverprofile coverage.txt ./... && echo "\n==>\033[32m Ok\033[m\n" || exit 1

install: $(GOFILES)
	C_INCLUDE_PATH=${INCLUDE_PATH} LIBRARY_PATH=${LIBRARY_PATH} $(GO) install -v -tags '$(TAGS)' -ldflags '$(EXTLDFLAGS)-s -w $(GOLDFLAGS)'

build: $(EXECUTABLE)

$(EXECUTABLE): $(GOFILES)
	C_INCLUDE_PATH=${INCLUDE_PATH} LIBRARY_PATH=${LIBRARY_PATH} $(GO) build -v -tags '$(TAGS)' -ldflags '$(EXTLDFLAGS)-s -w $(GOLDFLAGS)' -o bin/$@

clean:
	$(GO) clean -x -i ./...
	rm -rf coverage.txt $(EXECUTABLE) $(DIST)

version:
	@echo $(VERSION)
