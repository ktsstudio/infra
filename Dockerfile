FROM --platform=$BUILDPLATFORM golang:1.27.3 AS builder
ARG TARGETARCH
ARG TARGETOS=linux
ARG BUILDVERSION=99.99.99999
ARG TELEMETRY_WRITE_KEY
WORKDIR /go/src/github.com/infrahq/infra

# get deps first so it's cached
COPY go.mod .
COPY go.sum .
RUN --mount=type=cache,id=gomod,target=/go/pkg/mod \
    --mount=type=cache,id=gobuild,target=/root/.cache/go-build \
    go mod download

COPY . .

RUN --mount=type=cache,id=gomod,target=/go/pkg/mod \
    --mount=type=cache,id=gobuild,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH \
    go build -trimpath \
    -ldflags "-s -X github.com/infrahq/infra/internal.Version=$BUILDVERSION -X github.com/infrahq/infra/internal.TelemetryWriteKey=$TELEMETRY_WRITE_KEY" \
    .

FROM alpine:3.24.3@sha256:be2656061b536c1883c4034790e21d0a09c5aaaf51968b89368a3c0371ea2c9a
RUN apk add --no-cache ca-certificates
COPY --from=builder /go/src/github.com/infrahq/infra/infra /bin/infra
EXPOSE 80
EXPOSE 443
ARG USER=infra
ARG GROUP=infra
RUN addgroup -g 1000 $GROUP && adduser -u 1000 -DG $GROUP $USER
USER $USER:$GROUP
ENTRYPOINT ["/bin/infra"]
CMD ["server"]
