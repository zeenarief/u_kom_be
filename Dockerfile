# Tahap 1: Build
FROM golang:1.25-alpine AS builder

# Set direktori kerja di dalam container
WORKDIR /app

# Copy file go.mod dan go.sum
COPY go.mod go.sum ./

# Download dependency
RUN go mod download

# Copy semua kode sumber ke dalam container
COPY . .

# Build binary yang optimal untuk produksi:
# CGO_ENABLED=0: Membuat binary statis (tidak butuh library C di runtime)
# -ldflags '-s -w': Menghapus simbol debug, membuat file lebih kecil
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags '-s -w' -o server ./cmd/server

# Tahap 2: Runtime
FROM alpine:latest

# Instal paket yang umum dibutuhkan di produksi:
# ca-certificates: Untuk koneksi HTTPS ke service lain
# tzdata: Untuk data timezone (agar `TZ=UTC` di .env.prod berfungsi)
RUN apk add --no-cache ca-certificates tzdata

# Set direktori kerja
WORKDIR /app

# Copy binary dari tahap build
COPY --from=builder /app/server .

# Copy folder migrasi ke dalam image
COPY ./migrations ./migrations

# Port yang digunakan oleh aplikasi
EXPOSE 8080

# Jalankan aplikasi
CMD ["./server"]