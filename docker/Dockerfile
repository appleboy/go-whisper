FROM golang:1.20-alpine3.18 as SOURCE

WORKDIR /app

RUN apk update && apk add --no-cache make git gcc g++ && \
  rm -rf /var/cache/apk/*

COPY . .

RUN make dependency && make build && \
  mv bin/go-whisper /bin/ && \
  rm -rf bin && \
  apk del make git gcc g++

FROM jrottenberg/ffmpeg:4.4-alpine312

LABEL maintainer="Bo-Yi Wu <appleboy.tw@gmail.com>" \
  org.label-schema.name="Speech-to-Text" \
  org.label-schema.vendor="Bo-Yi Wu" \
  org.label-schema.schema-version="1.0"

LABEL org.opencontainers.image.source=https://github.com/appleboy/go-whisper
LABEL org.opencontainers.image.description="Speech-to-Text."
LABEL org.opencontainers.image.licenses=MIT

RUN apk update && apk add --no-cache curl && \
  rm -rf /var/cache/apk/*

WORKDIR /app

COPY --from=SOURCE /bin/go-whisper /bin/go-whisper

ENTRYPOINT [ "/bin/go-whisper" ]
