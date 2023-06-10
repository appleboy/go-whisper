FROM jrottenberg/ffmpeg:4.4-alpine312

LABEL maintainer="Bo-Yi Wu <appleboy.tw@gmail.com>" \
  org.label-schema.name="Speech-to-Text" \
  org.label-schema.vendor="Bo-Yi Wu" \
  org.label-schema.schema-version="1.0"

LABEL org.opencontainers.image.source=https://github.com/appleboy/go-whisper
LABEL org.opencontainers.image.description="Speech-to-Text."
LABEL org.opencontainers.image.licenses=MIT

RUN apk update && apk add --no-cache make git gcc g++ && \
  rm -rf /var/cache/apk/*

RUN make whisper && make build

COPY bin/go-whisper /bin/

WORKDIR /app

ENTRYPOINT ["/bin/go-whisper"]
