# Build the application from source.
FROM golang:1.21.1-alpine@sha256:1c9cc949513477766da12bfa80541c4f24957323b0ee00630a6ff4ccf334b75b AS go-builder

WORKDIR /app
COPY go.mod go.sum cmd ./
RUN go build -o /app/dist/deadnews-template-go ./...

# Deploy the application binary into a lean image.
FROM gcr.io/distroless/static-debian12:latest@sha256:98e138282ba524ff4f5124fec603f82ee2331df4ba981d169b3ded8bcd83ca52 AS runtime
LABEL maintainer "DeadNews <aurczpbgr@mozmail.com>"

COPY --from=go-builder /app/dist/deadnews-template-go /usr/local/bin/deadnews-template-go

USER nonroot:nonroot
EXPOSE 1271
HEALTHCHECK NONE

CMD ["deadnews-template-go"]
