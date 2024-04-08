FROM golang:1.22 as builder
COPY . .
RUN go build -ldflags "-linkmode external -extldflags -static" -o markhor

FROM scratch
COPY --from=builder /go/markhor .
CMD ["./markhor"]
