# Build the application from source.
FROM golang:1.24.2-alpine@sha256:7772cb5322baa875edd74705556d08f0eeca7b9c4b5367754ce3f2f00041ccee AS go-builder

ENV GOCACHE="/cache/go-build" \
    # Disable CGO to build a static binary.
    CGO_ENABLED=0

WORKDIR /app

COPY go.mod go.sum cmd ./
RUN --mount=type=cache,target=${GOCACHE} \
    go build -o /app/dist/deadnews-template-go ./...

# Deploy the application binary into a lean image.
FROM gcr.io/distroless/static-debian12:debug@sha256:5c474275684bbf271ec40502ab50158b2f9826de5877d8feec27e22e8d6ee3d2 AS runtime
LABEL maintainer="DeadNews <deadnewsgit@gmail.com>"

ENV SERVICE_PORT=8000

COPY --from=go-builder /app/dist/deadnews-template-go /bin/deadnews-template-go

RUN ["/busybox/sh", "-c", "ln -s /busybox/sh /bin/sh"]

USER nonroot:nonroot
EXPOSE ${SERVICE_PORT}
HEALTHCHECK NONE

ENTRYPOINT ["/bin/deadnews-template-go"]
