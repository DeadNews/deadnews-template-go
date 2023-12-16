# Build the application from source.
FROM golang:1.21.5-alpine@sha256:feceecc0e1d73d085040a8844de11a2858ba4a0c58c16a672f1736daecc2a4ff AS go-builder

WORKDIR /app
COPY go.mod go.sum cmd ./
RUN go build -o /app/dist/deadnews-template-go ./...

# Deploy the application binary into a lean image.
FROM gcr.io/distroless/static-debian12:latest@sha256:4a2c1a51ae5e10ec4758a0f981be3ce5d6ac55445828463fce8dff3a355e0b75 AS runtime
LABEL maintainer "DeadNews <aurczpbgr@mozmail.com>"

COPY --from=go-builder /app/dist/deadnews-template-go /usr/local/bin/deadnews-template-go

USER nonroot:nonroot
EXPOSE 1271
HEALTHCHECK NONE

CMD ["deadnews-template-go"]
