FROM golang:1.15-buster as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN make bin/pc

FROM debian:buster-slim

ARG VERSION=0.0.1
ARG BUILD_DATE=2021-05-1

LABEL \
  org.opencontainers.image.created="$BUILD_DATE" \
  org.opencontainers.image.authors="edvin.norling@gmail.com" \
  org.opencontainers.image.homepage="https://github.com/NissesSenap/falcosidekick-tekton" \
  org.opencontainers.image.documentation="https://github.com/NissesSenap/falcosidekick-tekton" \
  org.opencontainers.image.source="https://github.com/NissesSenap/falcosidekick-tekton" \
  org.opencontainers.image.version="$VERSION" \
  org.opencontainers.image.vendor="GitHub" \
  org.opencontainers.image.licenses="MIT" \
  summary="Kubernetes response Engine deletes pods after commands from falcosidekick" \
  description="Kubernetes response Engine delete pods is meant to be run inside k8s and delete pods after getting instructions from falcosidekick." \
  name="podDeleter"

RUN apt-get update && apt-get upgrade -y \
 && rm -rf /var/lib/apt/lists/*

USER 1001

WORKDIR /app

COPY --from=builder /app/bin/poddeleter .

CMD ["/app/poddeleter"]
