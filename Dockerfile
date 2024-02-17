FROM golang:1.21 as builder

WORKDIR /workspace
COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download

COPY main.go main.go
COPY internal internal

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o webhook .

FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY --from=builder /workspace/webhook /

ENTRYPOINT ["/webhook"]
