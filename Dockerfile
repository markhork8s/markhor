FROM golang:1.22 as builder
COPY . .
RUN go build -ldflags "-linkmode external -extldflags -static" -o sops_k8s

FROM scratch
COPY --from=builder /go/sops_k8s .
CMD ["./sops_k8s"]
