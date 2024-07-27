# Build the application from source.
FROM golang:1.22.5-alpine@sha256:0d3653dd6f35159ec6e3d10263a42372f6f194c3dea0b35235d72aabde86486e AS go-builder

ENV GOCACHE="/cache/go-build" \
    # Disable CGO to build a static binary.
    CGO_ENABLED=0

WORKDIR /app

COPY go.mod go.sum cmd ./
RUN --mount=type=cache,target=${GOCACHE} \
    go build -o /app/dist/deadnews-template-go ./...

# Deploy the application binary into a lean image.
FROM gcr.io/distroless/static-debian12:latest@sha256:ce46866b3a5170db3b49364900fb3168dc0833dfb46c26da5c77f22abb01d8c3 AS runtime
LABEL maintainer "DeadNews <deadnewsgit@gmail.com>"

ENV SERVICE_PORT=8000

COPY --from=go-builder /app/dist/deadnews-template-go /usr/local/bin/deadnews-template-go

USER nonroot:nonroot
EXPOSE ${SERVICE_PORT}
HEALTHCHECK NONE

CMD ["deadnews-template-go"]
