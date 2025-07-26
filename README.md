# Digital Signature Application

Aplikasi web Go untuk digital signature menggunakan algoritma RSA.

## Fitur

### Digital Signature

- **Generate RSA Key Pair**: Membuat pasangan kunci RSA (1024, 2048, atau 4096 bit)
- **Sign Document**: Menandatangani dokumen menggunakan private key RSA dan mengunduh signature sebagai file .rsasig
- **Verify Signature**: Memverifikasi tanda tangan digital menggunakan public key RSA dan file signature

## Teknologi yang Digunakan

- **Go**: Bahasa pemrograman utama
- **RSA Cryptography**: Untuk digital signature
- **Gorilla Mux**: Router untuk web server
- **HTML/CSS/JavaScript**: Frontend interface

## Algoritma

### RSA Digital Signature

1. **Key Generation**: Menggunakan `crypto/rsa` package untuk membuat pasangan kunci
2. **Signing**: SHA-256 hash dokumen kemudian ditandatangani dengan RSA PKCS#1 v1.5
3. **Verification**: Memverifikasi tanda tangan menggunakan public key
4. **File Format**: Signature disimpan dalam format Base64 dengan ekstensi .rsasig

## Struktur Project

```
digitalSign/
├── main.go                 # Entry point aplikasi
├── go.mod                 # Go module definition
├── crypto/
│   └── rsa.go            # RSA cryptography utilities
├── handlers/
│   ├── main.go           # API handlers
│   └── pages.go          # Web page handlers
├── middleware/
│   └── cors.go           # CORS middleware
└── static/               # Static files directory
```

## Cara Menjalankan

1. Install Go (version 1.21 atau lebih baru)
2. Clone atau download project ini
3. Jalankan perintah:
   ```bash
   go mod tidy
   go run main.go
   ```
4. Buka browser dan akses `http://localhost:8080`

## API Endpoints

### Digital Signature

- `POST /api/generate-keys` - Generate RSA key pair
- `POST /api/sign-document` - Sign a document (returns .rsasig file)
- `POST /api/verify-signature` - Verify signature

## Web Pages

- `/` - Homepage dengan navigasi utama
- `/digital-signature` - Halaman digital signature

## Cara Penggunaan

### 1. Generate RSA Key Pair

- Pilih ukuran kunci (1024, 2048, atau 4096 bit)
- Klik "Generate Key Pair"
- Copy private key untuk signing
- Copy public key untuk verification

### 2. Sign Document

- Upload dokumen yang akan ditandatangani
- Paste private key pada textarea
- Klik "Sign Document"
- File signature (\*.rsasig) akan otomatis terunduh

### 3. Verify Signature

- Upload dokumen asli yang sama
- Upload file signature (\*.rsasig)
- Paste public key pada textarea
- Klik "Verify Signature"
- Hasil verifikasi akan ditampilkan

## Penjelasan Metode

### Digital Signature dengan RSA

Digital signature menggunakan RSA untuk memastikan:

- **Authenticity**: Dokumen berasal dari pemilik private key
- **Integrity**: Dokumen tidak diubah setelah ditandatangani
- **Non-repudiation**: Pengirim tidak dapat menyangkal telah menandatangani

### Proses Signing

1. Dokumen di-hash menggunakan SHA-256
2. Hash ditandatangani menggunakan RSA private key dengan PKCS#1 v1.5
3. Signature di-encode dalam Base64
4. Disimpan dalam file dengan ekstensi .rsasig

### Proses Verification

1. Dokumen di-hash menggunakan SHA-256
2. Signature file dibaca dan di-decode dari Base64
3. Signature diverifikasi menggunakan RSA public key
4. Hasil verifikasi menunjukkan valid/invalid

## Keamanan

- RSA key minimum 1024 bit (disarankan 2048 bit atau lebih)
- SHA-256 untuk hashing dokumen (cryptographically secure)
- RSA PKCS#1 v1.5 untuk signing scheme
- Base64 encoding untuk representasi signature
- Error handling untuk input yang tidak valid

## Format File Signature

File signature (.rsasig) berisi:

- Base64-encoded RSA signature
- Hasil dari signing SHA-256 hash dokumen
- Dapat dibaca sebagai plain text
- Ukuran tergantung pada ukuran RSA key

## Penggunaan untuk Tugas

Aplikasi ini dapat digunakan untuk memenuhi requirements:

1. ✅ Aplikasi kriptografi dengan bahasa pemrograman (Go)
2. ✅ Metode enkripsi (RSA digital signature)
3. ✅ Perancangan dan implementasi algoritma
4. ✅ Deklarasi dan algoritma yang jelas
5. ✅ Format laporan dalam bentuk web application

Untuk keperluan akademik, aplikasi ini mendemonstrasikan penggunaan kriptografi RSA dalam konteks digital signature yang praktis dan aman.
