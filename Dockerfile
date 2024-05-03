FROM golang:1.22 as builder

# Helps speeding up the builds when dependencies do not change
COPY go.mod go.sum .
RUN go mod download

COPY . .
RUN go build -ldflags "-linkmode external -extldflags -static" -o markhor

# ----------------------------------

FROM scratch

COPY --from=builder /go/markhor .

CMD ["./markhor"]
