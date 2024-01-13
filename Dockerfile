# Build the application from source.
FROM golang:1.21.6-alpine@sha256:fd78f2fb1e49bcf343079bbbb851c936a18fc694df993cbddaa24ace0cc724c5 AS go-builder

WORKDIR /app
COPY go.mod go.sum cmd ./
RUN go build -o /app/dist/deadnews-template-go ./...

# Deploy the application binary into a lean image.
FROM gcr.io/distroless/static-debian12:latest@sha256:4a2c1a51ae5e10ec4758a0f981be3ce5d6ac55445828463fce8dff3a355e0b75 AS runtime
LABEL maintainer "DeadNews <aurczpbgr@mozmail.com>"

ENV PORT=1271

COPY --from=go-builder /app/dist/deadnews-template-go /usr/local/bin/deadnews-template-go

USER nonroot:nonroot
EXPOSE ${PORT}
HEALTHCHECK NONE

CMD ["deadnews-template-go"]
