FROM golang:1.22 as builder

WORKDIR /workspace
COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download

COPY main.go main.go
COPY internal internal

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o webhook .

FROM cgr.dev/chainguard/static:latest
WORKDIR /
COPY --from=builder /workspace/webhook /

ENTRYPOINT ["/webhook"]
