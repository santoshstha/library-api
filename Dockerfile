FROM golang:1.23-alpine AS builder
WORKDIR /app

# Install git
RUN apk add --no-cache git

# Initialize the Go module
RUN go mod init library-api

# Fetch dependencies
RUN go get -u gorm.io/gorm
RUN go get -u gorm.io/driver/mysql
RUN go get -u github.com/gorilla/mux
RUN go get -u github.com/dgrijalva/jwt-go
RUN go get -u golang.org/x/crypto/bcrypt

# Copy all source files
COPY . .

# Build the app
RUN go build -o main .

# Final lightweight image
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/main .
EXPOSE 8080
CMD ["./main"]