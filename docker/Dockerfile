FROM jrottenberg/ffmpeg:4.4-alpine312

RUN apk update && apk add --no-cache curl && \
  rm -rf /var/cache/apk/*

RUN curl -LJO https://huggingface.co/ggerganov/whisper.cpp/resolve/main/ggml-base.bin \
  --output /models/ggml-base.bin && \
  curl -LJO https://huggingface.co/ggerganov/whisper.cpp/resolve/main/ggml-large.bin \
  --output /models/ggml-large.bin

CMD [ "/bin/sh" ]
