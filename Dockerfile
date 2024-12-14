# Build the application from source.
FROM golang:1.23.4-alpine@sha256:6c5c9590f169f77c8046e45c611d3b28fe477789acd8d3762d23d4744de69812 AS go-builder

ENV GOCACHE="/cache/go-build" \
    # Disable CGO to build a static binary.
    CGO_ENABLED=0

WORKDIR /app

COPY go.mod go.sum cmd ./
RUN --mount=type=cache,target=${GOCACHE} \
    go build -o /app/dist/deadnews-template-go ./...

# Deploy the application binary into a lean image.
FROM gcr.io/distroless/static-debian12:debug@sha256:0eb1b021dc83cd103a4f5bec6d5292907b7f649137a200e178aa42d768d93ba2 AS runtime
LABEL maintainer="DeadNews <deadnewsgit@gmail.com>"

ENV SERVICE_PORT=8000

COPY --from=go-builder /app/dist/deadnews-template-go /bin/deadnews-template-go

RUN ["/busybox/sh", "-c", "ln -s /busybox/sh /bin/sh"]

USER nonroot:nonroot
EXPOSE ${SERVICE_PORT}
HEALTHCHECK NONE

ENTRYPOINT ["/bin/deadnews-template-go"]
