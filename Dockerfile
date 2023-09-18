# Build the application from source.
FROM golang:1.21-alpine@sha256:96634e55b363cb93d39f78fb18aa64abc7f96d372c176660d7b8b6118939d97b AS go-builder

WORKDIR /tmp/app
COPY go.mod go.sum cmd ./
RUN go build -o /tmp/deadnews-template-go ./...

# Deploy the application binary into a lean image.
FROM gcr.io/distroless/static-debian12:latest@sha256:98e138282ba524ff4f5124fec603f82ee2331df4ba981d169b3ded8bcd83ca52 AS final
LABEL maintainer "DeadNews <aurczpbgr@mozmail.com>"

WORKDIR /
COPY --from=go-builder /tmp/deadnews-template-go /bin/deadnews-template-go

USER nonroot:nonroot
EXPOSE 1271
HEALTHCHECK --interval=60s --timeout=3s CMD curl --fail http://127.0.0.1:1271/health || exit 1

CMD ["/bin/deadnews-template-go"]
