# Nutribox API Docker Setup

Dokumen ini berisi petunjuk untuk menjalankan Nutribox API menggunakan Docker di lingkungan production.

## Persyaratan

- Docker dan Docker Compose terinstal di sistem Anda
- Git (untuk mengkloning repositori)

## Struktur Docker

Setup Docker ini terdiri dari:

1. **Dockerfile** - Konfigurasi untuk membangun image aplikasi Go
2. **docker-compose.yml** - Konfigurasi untuk menjalankan container aplikasi dan PostgreSQL
3. **run-docker.sh** - Script untuk mempermudah pengelolaan container

## Cara Penggunaan

### 1. Persiapan

Pastikan Anda telah mengkloning repositori dan berada di direktori root proyek:

```bash
git clone https://github.com/TheValeHack/nutribox-api.git
cd nutribox-api
```

### 2. Konfigurasi Environment

Salin file `.env.example` menjadi `.env` dan sesuaikan konfigurasi sesuai kebutuhan:

```bash
cp .env.example .env
```

Edit file `.env` untuk mengubah konfigurasi seperti port, kredensial database, dan pengaturan lainnya.

### 3. Menjalankan Aplikasi

Gunakan script `run-docker.sh` untuk menjalankan aplikasi:

```bash
./run-docker.sh
```

Script ini akan:
- Memeriksa keberadaan file `.env`
- Membangun dan menjalankan container Docker
- Menampilkan URL untuk mengakses API

### Perintah yang Tersedia

Script `run-docker.sh` mendukung beberapa perintah:

- `./run-docker.sh start` - Menjalankan container
- `./run-docker.sh stop` - Menghentikan container
- `./run-docker.sh restart` - Me-restart container
- `./run-docker.sh rebuild` - Membangun ulang dan me-restart container
- `./run-docker.sh logs` - Menampilkan log container
- `./run-docker.sh help` - Menampilkan bantuan

## Struktur Container

### 1. Container Aplikasi (nutribox-api)

- Image: Dibangun dari Dockerfile
- Port: 8097 (default, dapat diubah di .env)
- Volume: `./uploads:/app/uploads` untuk menyimpan file upload

### 2. Container Database (nutribox-postgres)

- Image: postgres:15-alpine
- Port: 5432 (default, dapat diubah di .env)
- Volume: `postgres_data` untuk persistensi data

## Pengelolaan Data

Data PostgreSQL disimpan dalam volume Docker bernama `postgres_data` untuk memastikan data tetap ada meskipun container dihapus atau di-restart.

## Keamanan

- Aplikasi berjalan sebagai pengguna non-root (appuser) di dalam container
- Kredensial database disimpan dalam file .env dan tidak di-hardcode
- Container terhubung melalui jaringan Docker internal (nutribox-network)

## Pemecahan Masalah

1. **Masalah Koneksi Database**
   - Pastikan kredensial database di file .env sudah benar
   - Periksa log container dengan `./run-docker.sh logs`

2. **Port Sudah Digunakan**
   - Ubah port di file .env jika port default sudah digunakan

3. **Masalah Izin**
   - Pastikan direktori uploads memiliki izin yang benar
   - Jalankan `chmod -R 777 uploads` jika diperlukan