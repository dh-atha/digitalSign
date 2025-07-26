# Digital Signature & Steganography Application

Aplikasi web Go untuk digital signature dan steganografi menggunakan algoritma RSA dan metode LSB (Least Significant Bit).

## Fitur

### 1. Digital Signature

- **Generate RSA Key Pair**: Membuat pasangan kunci RSA (1024, 2048, atau 4096 bit)
- **Sign Document**: Menandatangani dokumen menggunakan private key RSA
- **Verify Signature**: Memverifikasi tanda tangan digital menggunakan public key RSA

### 2. Steganography

- **Hide Message**: Menyembunyikan pesan dalam gambar menggunakan metode LSB
- **Extract Message**: Mengekstrak pesan tersembunyi dari gambar

### 3. Combined Operations

- **Sign and Hide**: Menandatangani dokumen dan menyembunyikannya dalam gambar
- **Extract and Verify**: Mengekstrak dokumen dari gambar dan memverifikasi tanda tangannya

## Teknologi yang Digunakan

- **Go**: Bahasa pemrograman utama
- **RSA Cryptography**: Untuk digital signature
- **LSB Steganography**: Untuk menyembunyikan pesan dalam gambar
- **Gorilla Mux**: Router untuk web server
- **HTML/CSS/JavaScript**: Frontend interface

## Algoritma

### RSA Digital Signature

1. **Key Generation**: Menggunakan `crypto/rsa` package untuk membuat pasangan kunci
2. **Signing**: SHA-256 hash dokumen kemudian ditandatangani dengan RSA PKCS#1 v1.5
3. **Verification**: Memverifikasi tanda tangan menggunakan public key

### LSB Steganography

1. **Hiding**: Pesan disembunyikan dalam bit terakhir (LSB) dari setiap channel RGB
2. **Extraction**: Mengekstrak bit LSB untuk merekonstruksi pesan asli
3. **Format**: Mendukung gambar PNG dan JPEG

## Struktur Project

```
digitalSign/
├── main.go                 # Entry point aplikasi
├── go.mod                 # Go module definition
├── crypto/
│   └── rsa.go            # RSA cryptography utilities
├── steganography/
│   └── lsb.go            # LSB steganography implementation
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
- `POST /api/sign-document` - Sign a document
- `POST /api/verify-signature` - Verify signature

### Steganography

- `POST /api/hide-message` - Hide message in image
- `POST /api/extract-message` - Extract message from image

### Combined Operations

- `POST /api/sign-and-hide` - Sign document and hide in image
- `POST /api/extract-and-verify` - Extract and verify signed document

## Web Pages

- `/` - Homepage dengan navigasi utama
- `/digital-signature` - Halaman digital signature
- `/steganography` - Halaman steganografi
- `/combined` - Halaman operasi gabungan

## Penjelasan Metode

### Digital Signature dengan RSA

Digital signature menggunakan RSA untuk memastikan:

- **Authenticity**: Dokumen berasal dari pemilik private key
- **Integrity**: Dokumen tidak diubah setelah ditandatangani
- **Non-repudiation**: Pengirim tidak dapat menyangkal telah menandatangani

### Steganografi LSB

Steganografi LSB menyembunyikan pesan dengan:

- Mengubah bit terakhir (LSB) dari setiap channel RGB pixel
- Menyimpan panjang pesan di awal untuk ekstraksi
- Menggunakan end marker untuk menandai akhir pesan

### Operasi Gabungan

Menggabungkan kedua metode untuk:

- Melindungi integritas dengan digital signature
- Menyembunyikan dokumen dan signature dalam gambar
- Memberikan perlindungan berlapis untuk dokumen penting

## Keamanan

- RSA key minimum 1024 bit (disarankan 2048 bit)
- SHA-256 untuk hashing dokumen
- End marker untuk validasi ekstraksi pesan
- Error handling untuk input yang tidak valid

## Limitasi

- Ukuran pesan terbatas oleh ukuran gambar
- Format gambar yang didukung: PNG dan JPEG
- Kualitas gambar JPEG dapat menurun setelah steganografi

## Penggunaan untuk Tugas

Aplikasi ini dapat digunakan untuk memenuhi requirements:

1. ✅ Aplikasi kriptografi dengan bahasa pemrograman (Go)
2. ✅ Metode enkripsi (RSA digital signature)
3. ✅ Perancangan dan implementasi algoritma
4. ✅ Deklerasi dan algoritma yang jelas
5. ✅ Format laporan dalam bentuk web application

Untuk keperluan akademik, aplikasi ini mendemonstrasikan penggunaan kriptografi RSA dan steganografi dalam konteks praktis.
