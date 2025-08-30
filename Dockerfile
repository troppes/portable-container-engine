FROM golang:1.25-alpine AS builder

WORKDIR /app
COPY . .

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /pce ./cmd/pce

FROM alpine:latest

WORKDIR /
COPY --from=builder /pce /usr/local/bin/pce

ENTRYPOINT ["pce"]
