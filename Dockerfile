FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY *.go ./
RUN go build -o sentinel .

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/sentinel .
CMD ["./sentinel"]