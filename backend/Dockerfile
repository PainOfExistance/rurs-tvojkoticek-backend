# Stage 1: Build
FROM golang:1.20 as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Compile with static linking for Linux
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Stage 2: Run
FROM alpine:latest

# Install necessary runtime libraries
RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/main .
COPY .env /root/.env

EXPOSE 8080

CMD ["./main"]
