ARG UBUNTU_VERSION=22.04
# This needs to generally match the container host's environment.
ARG CUDA_VERSION=12.0.0
# Target the CUDA build image
ARG BASE_CUDA_DEV_CONTAINER=nvidia/cuda:${CUDA_VERSION}-devel-ubuntu${UBUNTU_VERSION}
# Target the CUDA runtime image
ARG BASE_CUDA_RUN_CONTAINER=nvidia/cuda:${CUDA_VERSION}-runtime-ubuntu${UBUNTU_VERSION}

FROM nvidia/cuda:${CUDA_VERSION}-devel-ubuntu${UBUNTU_VERSION} AS build
WORKDIR /app
# Unless otherwise specified, we make a fat build.
ARG CUDA_DOCKER_ARCH=all
# Set nvcc architecture
ENV CUDA_DOCKER_ARCH=${CUDA_DOCKER_ARCH}
# Enable cuBLAS
ENV WHISPER_CUBLAS=1

#apt-get
RUN apt-get update && \
  apt-get install -y --no-install-recommends build-essential git gcc g++ wget \
  && rm -rf /var/lib/apt/lists/* /var/cache/apt/archives/*

# install golang
RUN wget --progress=dot:giga https://go.dev/dl/go1.22.10.linux-amd64.tar.gz
RUN rm -rf /usr/local/go && tar -C /usr/local -xzf go1.22.10.linux-amd64.tar.gz
ENV PATH ${PATH}:/usr/local/go/bin

# Ref: https://stackoverflow.com/a/53464012
ENV CUDA_MAIN_VERSION=12.0
ENV LD_LIBRARY_PATH /usr/local/cuda-${CUDA_MAIN_VERSION}/compat:$LD_LIBRARY_PATH

COPY ./ .
RUN make dependency && env && make build && \
  mv bin/go-whisper /bin/ && \
  rm -rf bin 

FROM nvidia/cuda:${CUDA_VERSION}-runtime-ubuntu${UBUNTU_VERSION} AS runtime
WORKDIR /app

LABEL maintainer="Bo-Yi Wu <appleboy.tw@gmail.com>" \
  org.label-schema.name="Speech-to-Text" \
  org.label-schema.vendor="Bo-Yi Wu" \
  org.label-schema.schema-version="1.0"

LABEL org.opencontainers.image.source=https://github.com/appleboy/go-whisper
LABEL org.opencontainers.image.description="Speech-to-Text."
LABEL org.opencontainers.image.licenses=MIT

RUN apt-get update && \
  apt-get install -y --no-install-recommends curl ffmpeg \
  && rm -rf /var/lib/apt/lists/* /var/cache/apt/archives/*

COPY --from=build /bin/go-whisper /bin/go-whisper
ENTRYPOINT ["/bin/go-whisper"]