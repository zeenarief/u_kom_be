FROM ubuntu:latest
LABEL authors="zeen"

ENTRYPOINT ["top", "-b"]

# Tahap 1: Build
FROM golang:1.22-alpine AS builder

# Set direktori kerja di dalam container
WORKDIR /app

# Copy file go.mod dan go.sum
COPY go.mod go.sum ./

# Download dependency
RUN go mod download

# Copy semua kode sumber ke dalam container
COPY . .

# Build binary
RUN go build -o server ./cmd/server/main.go

# Tahap 2: Runtime
FROM alpine:latest

# Set direktori kerja
WORKDIR /app

# Copy binary dari tahap build
COPY --from=builder /app/server .

# Copy file .env jika dibutuhkan
COPY .env .env

# Port yang digunakan oleh aplikasi
EXPOSE 8080

# Jalankan aplikasi
CMD ["./server"]
