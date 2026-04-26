# Technical Requirements Document (TRD) - Project Pantau Gizi V2

## 1. Ringkasan Eksekutif
Platform **Pantau Gizi** adalah sistem crowdsourcing independen untuk memantau kualitas program Makan Bergizi Gratis (MBG). Fokus pengembangan: efisiensi biaya, skalabilitas, dan integritas data publik.

## 2. Arsitektur Infrastruktur (Hybrid Model)
Sistem menggunakan pendekatan hybrid untuk optimasi performa vs biaya:
- **Backend API:** Go (Golang) 1.22+ (Dockerized on VPS).
- **Frontend:** Next.js (PWA enabled) on Vercel/Cloudflare Pages.
- **Database Utama:** PostgreSQL (Managed/Serverless - Neon.tech/Supabase).
- **Cache & Rate Limiting:** Redis (Upstash/Self-hosted).
- **Object Storage:** Cloudflare R2 (S3-compatible) untuk foto.

## 3. Spesifikasi Teknis Backend (Go)
- **Framework:** Gin atau Echo.
- **Database Handler:** GORM atau sqlx dengan connection pooling.
- **Logging:** Structured logging (zerolog/zap).
- **Middleware:** CORS, Recovery, dan Redis-based Rate Limiting.

## 4. Strategi Pengelolaan Data & Media (Kritis)
### A. Mekanisme Presigned URL
Server Go **TIDAK DIPERBOLEHKAN** memproses upload file biner secara langsung.
1. Client `POST` ke `/v1/upload/request`.
2. Backend generate **S3 Presigned URL** (Cloudflare R2, expired 5-10 mnt).
3. Client upload biner langsung ke R2 via `PUT`.
4. Client mengirim metadata (URL, rating, koordinat) ke API utama.

### B. Kompresi Sisi Client (Frontend)
Wajib dilakukan sebelum upload:
- **Max File Size:** 500 KB.
- **Max Resolution:** 1280px (long side).
- **Format:** WebP/JPEG (Quality 0.7 - 0.8).

## 5. Implementasi Caching (Redis)
- Data agregat (rating nasional/wilayah) disimpan di Redis (TTL 15-30 mnt).
- Query ke PostgreSQL hanya dilakukan jika cache miss.

## 6. Skema Database (Ringkasan)
| Tabel | Kolom Utama | Indeks Kritis |
| :--- | :--- | :--- |
| **users** | id, email, name, role | email |
| **schools** | id, npsn, name, city_code, lat, long | npsn, city_code |
| **reports** | id, user_id, school_id, photo_url, rating | school_id, created_at |

## 7. Keamanan
- **Captcha:** Cloudflare Turnstile pada setiap laporan.
- **Validation:** Backend memvalidasi domain `photo_url` (harus R2 bucket).
- **Sanitization:** Proteksi terhadap SQL Injection & XSS.

## 8. Panduan Deployment
Wajib menyertakan `Dockerfile` dan `docker-compose.yml` mencakup:
- Go API container.
- Redis container.
- Nginx Reverse Proxy (SSL).
