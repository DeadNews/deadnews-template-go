# Build the application from source.
FROM golang:1.23.0-alpine@sha256:d0b31558e6b3e4cc59f6011d79905835108c919143ebecc58f35965bf79948f4 AS go-builder

ENV GOCACHE="/cache/go-build" \
    # Disable CGO to build a static binary.
    CGO_ENABLED=0

WORKDIR /app

COPY go.mod go.sum cmd ./
RUN --mount=type=cache,target=${GOCACHE} \
    go build -o /app/dist/deadnews-template-go ./...

# Deploy the application binary into a lean image.
FROM gcr.io/distroless/static-debian12:debug@sha256:0744faa82b4de871151f62d4c776f222ce5df373fb278c8810a94a8459994059 AS runtime
LABEL maintainer="DeadNews <deadnewsgit@gmail.com>"

ENV SERVICE_PORT=8000

COPY --from=go-builder /app/dist/deadnews-template-go /bin/deadnews-template-go

RUN ["/busybox/sh", "-c", "ln -s /busybox/sh /bin/sh"]

USER nonroot:nonroot
EXPOSE ${SERVICE_PORT}
HEALTHCHECK NONE

ENTRYPOINT [ "/bin/deadnews-template-go" ]
