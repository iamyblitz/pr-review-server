FROM golang:1.22 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o pr-reviewer-service ./cmd/server

FROM alpine:3.20

WORKDIR /app

COPY --from=builder /app/pr-reviewer-service /app/pr-reviewer-service

EXPOSE 8080

ENTRYPOINT ["/app/pr-reviewer-service"]
