# Stage 1: Build the Go application
FROM golang:alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o main .

# Stage 2: Create the final image with the executable
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/main .

CMD ["./main"]
