FROM golang:1.24.7-bookworm AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/main

FROM alpine:latest

WORKDIR /app

# Essential files for the app to run
COPY --from=builder /app/main .
COPY --from=builder /app/mnist_dataset ./mnist_dataset
COPY --from=builder /app/data ./data

CMD ["./main"]