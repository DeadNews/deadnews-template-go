# Build the application from source.
FROM golang:1.22.1-alpine@sha256:6f179eca0d49ec57ed6d64067d3d2c8c77fb4ca134b687f31cf1666e467cd1a9 AS go-builder

ENV GOCACHE="/cache/go-build" \
    # Disable CGO to build a static binary.
    CGO_ENABLED=0

WORKDIR /app

COPY go.mod go.sum cmd ./
RUN --mount=type=cache,target=${GOCACHE} \
    go build -o /app/dist/deadnews-template-go ./...

# Deploy the application binary into a lean image.
FROM gcr.io/distroless/static-debian12:latest@sha256:6dcc833df2a475be1a3d7fc951de90ac91a2cb0be237c7578b88722e48f2e56f AS runtime
LABEL maintainer "DeadNews <aurczpbgr@mozmail.com>"

ENV GO_PORT=1271

COPY --from=go-builder /app/dist/deadnews-template-go /usr/local/bin/deadnews-template-go

USER nonroot:nonroot
EXPOSE ${GO_PORT}
HEALTHCHECK NONE

CMD ["deadnews-template-go"]
