FROM --platform=${BUILDPLATFORM:-linux/amd64} golang:1.24.4 as builder

ARG TARGETPLATFORM
ARG BUILDPLATFORM
ARG TARGETOS
ARG TARGETARCH

WORKDIR /go

# Helps speeding up the builds when dependencies do not change
COPY go.mod go.sum .
RUN go mod download

COPY main.go .
COPY pkg ./pkg
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -ldflags "-s -w" -o markhor

# ----------------------------------

FROM --platform=${BUILDPLATFORM:-linux/amd64} gcr.io/distroless/static:nonroot

COPY --from=builder /go/markhor /

LABEL org.opencontainers.image.source=https://github.com/markhork8s/markhor

USER nonroot:nonroot

CMD ["/markhor"]
