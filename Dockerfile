# Build the application from source.
FROM golang:1.21.3-alpine@sha256:926f7f7e1ab8509b4e91d5ec6d5916ebb45155b0c8920291ba9f361d65385806 AS go-builder

WORKDIR /app
COPY go.mod go.sum cmd ./
RUN go build -o /app/dist/deadnews-template-go ./...

# Deploy the application binary into a lean image.
FROM gcr.io/distroless/static-debian12:latest@sha256:0c3d36f317d6335831765546ece49b60ad35933250dc14f43f0fd1402450532e AS runtime
LABEL maintainer "DeadNews <aurczpbgr@mozmail.com>"

COPY --from=go-builder /app/dist/deadnews-template-go /usr/local/bin/deadnews-template-go

USER nonroot:nonroot
EXPOSE 1271
HEALTHCHECK NONE

CMD ["deadnews-template-go"]
