# Build the application from source.
FROM golang:1.21.2-alpine@sha256:a76f153cff6a59112777c071b0cde1b6e4691ddc7f172be424228da1bfb7bbda AS go-builder

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
