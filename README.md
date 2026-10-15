# Fullstack Authentication System - PT Kuliner Kreasindo Inovasi (FTI)

Sistem autentikasi fullstack berbasis web yang dibangun menggunakan **Go**, **Svelte**, dan **PostgreSQL 17** untuk kebutuhan internal tim. Project ini diorkestrasi menggunakan **Docker Compose** dengan sistem migrasi basis data otomatis dan inisialisasi akun admin bawaan.

---

## 🚀 Gambaran Singkat & Fitur

Aplikasi ini menyediakan manajemen identitas dan kontrol akses berbasis peran (*Role-Based Access Control / RBAC*) dengan standar keamanan tinggi dan pembatasan sesi perangkat.

### Fitur Utama (Kebutuhan Minimum)
* **Pendaftaran (Register):** Pengguna baru dapat mendaftar dengan mengisi email, username, dan kata sandi.
* **Autentikasi (Login & Logout):** Login menggunakan email atau username dengan kata sandi, serta fitur logout untuk mengakhiri sesi.
* **Profil Pengguna (Protected Route):** Halaman profil terproteksi yang menampilkan data pengguna yang sedang login.
* **Role Management (User & Admin):** Halaman khusus admin untuk melihat daftar seluruh pengguna terdaftar.
* **Auto-Seeding Admin Pertama:** Akun admin pertama dibuat secara otomatis saat kontainer Docker pertama kali dijalankan tanpa langkah manual tambahan.
* **Pembatasan Maksimal 2 Sesi Aktif (Ketentuan Khusus):** Satu akun hanya dapat memiliki maksimal 2 sesi aktif secara bersamaan. Login pada perangkat/browser ketiga akan mengakhiri sesi paling awal secara otomatis.
* **Force Logout oleh Admin (Ketentuan Khusus):** Admin memiliki wewenang untuk memutus seluruh sesi aktif pengguna lain secara paksa.

### Fitur Tambahan
* **Health Check Endpoint:** Monitoring status koneksi aplikasi dan database (`/api/health`).
* **CORS & Security Headers:** Konfigurasi keamanan pada backend untuk mendukung komunikasi aman antar-domain.

---

## 🛠️ Teknologi yang Digunakan

* **Backend:** Go (Golang) + REST API
* **Frontend:** Svelte + Vite + Tailwind CSS (disajikan via Nginx Web Server)
* **Database:** PostgreSQL 17
* **Containerization:** Docker & Docker Compose

---

## ⚡ Cara Menjalankan Aplikasi

Aplikasi dirancang agar dapat dijalankan secara langsung dengan satu perintah Docker Compose.

### Persyaratan Sistem
* Docker Engine & Docker Compose v2+
* Git

### Langkah Jalankan Project
1. Clone repositori ini:
   git clone <URL_REPOSITORI_GITHUB_ANDA>
   cd <NAMA_FOLDER_REPOSITORI>

2. Salin berkas lingkungan (`.env.example` ke `.env`):
   cp .env.example .env

3. Jalankan aplikasi menggunakan Docker Compose:
   docker compose up -d --build

Setelah perintah selesai, aplikasi dapat diakses melalui:
* **Frontend:** `http://localhost:3000`
* **Backend:** `http://localhost:8080`

### Kredensial Admin Pertama (Default)
Akun administrator telah di-seed secara otomatis saat sistem pertama kali dinyalakan:
* **Identifier (Username/Email):** `admin` atau `admin@kuliner.co.id`
* **Password:** `AdminSecure123!`

---

## 🔌 Daftar Endpoint API

| Method | Endpoint | Proteksi | Role | Deskripsi |
| :--- | :--- | :--- | :--- | :--- |
| **GET** | `/api/health` | Publik | All | Cek status ketersediaan backend & database |
| **POST** | `/api/register` | Publik | All | Pendaftaran pengguna baru |
| **POST** | `/api/login` | Publik | All | Autentikasi & penerbitan token sesi |
| **POST** | `/api/logout` | Protected | User/Admin | Mengakhiri sesi aktif saat ini |
| **GET** | `/api/profile` | Protected | User/Admin | Mendapatkan data profil pengguna yang login |
| **GET** | `/api/users` | Protected | Admin | Mendapatkan daftar seluruh pengguna |
| **POST** | `/api/users/force-logout` | Protected | Admin | Memaksa logout pengguna lain seluruh sesi |

---

## 💡 Keputusan Teknis Penting

1. **Stateful Session Management di Database:**
   Meskipun autentikasi menggunakan token/bearer, setiap token sesi dicatat di dalam tabel database (`user_sessions`). Hal ini dipilih untuk memenuhi persyaratan bisnis secara presisi: melacak jumlah sesi aktif per pengguna (maksimal 2) dan memungkinkan instan *force logout* oleh admin.
2. **Auto Migration & Seeding via Go Startup:**
   Proses migrasi skema database dan pembuatan akun admin pertama (*auto-seeding*) dieksekusi langsung pada fungsi inisialisasi backend (`internal/database/database.go`) saat kontainer backend dinyalakan. Ini menjamin aplikasi siap pakai tanpa *step* manual di terminal.
3. **Multi-Stage Docker Build:**
   Proses build Svelte frontend dan kompilasi biner Go menggunakan *multi-stage build* untuk menghasilkan ukuran *image* Docker yang kecil, efisien, dan aman.

---

## ❓ Jawaban Catatan Implementasi

1. **Bagaimana kata sandi disimpan dalam sistem Anda, dan mengapa Anda memilih cara tersebut?**
   Kata sandi disimpan dalam bentuk hash terenkripsi menggunakan algoritma **bcrypt** dengan *cost factor* default (10) pada handler pendaftaran (`internal/handlers/auth.go`). Cara ini dipilih karena bcrypt secara otomatis menyertakan *salt* unik untuk mencegah *rainbow table attacks* dan merupakan standar industri yang aman serta komputasinya dapat disesuaikan untuk menahan serangan *brute-force*.

2. **Bagaimana sistem Anda mendefinisikan sebuah "sesi", dan pada berkas atau fungsi mana batas 2 sesi ditegakkan?**
   Sebuah "sesi" didefinisikan sebagai satu entitas unik di tabel `user_sessions` yang menyimpan `user_id`, `token` (UUID/JWT string), `created_at`, dan `expires_at`. Batas maksimal 2 sesi ditegakkan pada berkas `internal/handlers/auth.go` di dalam fungsi `Login()`. Saat proses login berhasil, sistem menghitung jumlah sesi aktif pengguna; jika sudah mencapai 2 sesi, fungsi akan menghapus (*delete*) sesi paling lama (`created_at` paling awal) sebelum memasukkan token sesi baru.

3. **Apabila token atau sesi seorang pengguna dicuri, apa yang terjadilah pada sistem Anda dan apa yang dapat dilakukan admin?**
   Penyerang yang mencuri token dapat mengakses endpoint terproteksi atas nama korban selama sesi tersebut belum habis masa berlakunya atau belum digantikan oleh login baru. Untuk menanggulangi hal tersebut, admin dapat memanggil fungsi *Force Logout* di halaman manajemen user (endpoint `/api/users/force-logout`), yang akan langsung menghapus seluruh rekaman sesi milik target pengguna dari basis data sehingga token yang dicuri seketika menjadi tidak valid (*revoked*).

4. **Sebutkan dua kekurangan atau celah yang Anda sadari masih terdapat pada hasil kerja Anda.**
   Pertama, pencatatan sesi di PostgreSQL memerlukan I/O database pada setiap pengecekan autentikasi yang dapat menjadi *bottleneck* jika trafik sangat tinggi (seharusnya menggunakan in-memory cache seperti Redis). Kedua, belum tersedianya fitur *rate-limiting* pada endpoint `/api/login` sehingga masih berpotensi terkena serangan *brute-force* atau *Denial of Service (DoS)*.

5. **Apabila diminta menambahkan fitur "pengguna dapat keluar dari seluruh perangkat sekaligus", bagian kode mana yang akan Anda ubah?**
   Saya akan menambahkan fungsi baru `LogoutAllSessions()` di `internal/handlers/auth.go` dan mendaftarkan endpoint baru (misalnya `POST /api/logout-all`). Di dalam handler tersebut, saya cukup mengambil `user_id` dari konteks sesi yang sedang aktif, lalu mengeksekusi query SQL `DELETE FROM user_sessions WHERE user_id = $1` untuk menghapus seluruh baris sesi milik pengguna tersebut.

---

## 🤖 Penggunaan AI

Dalam penyelesaian asesmen ini, alat AI yang digunakan adalah **Gemini / Claude AI**:
1. **Penyusunan Boilerplate & Troubleshooting Config:** Membantu dalam penataan konfigurasi `docker-compose.yml`, pembetulan plugin `@sveltejs/vite-plugin-svelte`, serta pemecahan masalah bentrok port PostgreSQL di lingkungan Windows.
2. **Diskusi Arsitektur Pembatasan Sesi:** Membantu memvalidasi logika pembersihan sesi paling lama (*FIFO session rotation*) pada handler Go.
3. **Pengembangan Antarmuka (UI/UX) & Refaktorisasi:** Claude membantu merancang serta merapikan struktur antarmuka pengguna (UI/UX) dan melakukan refactoring struktur direktori frontend menjadi komponen modular
(base component).

---

© 2026 PT Kuliner Kreasindo Inovasi - Divisi Future Technology and Innovation (FTI)