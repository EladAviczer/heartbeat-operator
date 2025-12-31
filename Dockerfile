# Build Stage
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
# Build the binary inside cmd/probe-operator
RUN go build -o probe-operator cmd/probe-operator/main.go

# Run Stage
FROM alpine:latest
# Install tools for 'exec' probes
RUN apk add --no-cache curl netcat-openbsd bash
WORKDIR /root/
COPY --from=builder /app/probe-operator .
CMD ["./probe-operator"]