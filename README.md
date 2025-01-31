# RescueHub API

## Deskripsi
RescueHub adalah sistem manajemen bencana berbasis API yang dirancang untuk mengelola laporan bencana, pengungsi, koordinasi relawan, distribusi bantuan, dan jalur evakuasi secara efisien. API ini memungkinkan pengguna untuk melaporkan bencana, mengelola shelter, serta menyalurkan bantuan dan donasi secara transparan.

## Fitur-Fitur Utama
1. **Manajemen Pengguna**
   - Registrasi, login, dan autentikasi menggunakan JWT.
   - Dukungan **Two-Factor Authentication (2FA)** untuk keamanan login.
   - Hak akses berbasis peran (Admin, Donor, User).

2. **Manajemen Bencana**
   - Pelaporan bencana oleh admin atau relawan.
   - Pencatatan status bencana (Active, Resolved, Archived).
   - Monitoring shelter yang tersedia untuk pengungsi.

3. **Manajemen Pengungsi & Shelter**
   - Pendataan pengungsi berdasarkan bencana yang terjadi.
   - Pendaftaran dan pengelolaan shelter beserta kapasitasnya.
   
4. **Manajemen Relawan**
   - Pengguna dapat mendaftar sebagai relawan dan ditugaskan ke bencana tertentu.
   - Relawan dapat memiliki spesialisasi keterampilan seperti medis, logistik, atau penyelamatan.
   
5. **Manajemen Bantuan & Donasi**
   - Pengguna dapat memberikan donasi berupa barang atau uang untuk bencana tertentu.
   - Admin dapat mendistribusikan bantuan logistik ke shelter dan pengungsi.
   
6. **Manajemen Jalur Evakuasi**
   - Pendataan jalur evakuasi aman, termasuk status jalur (Safe, Risky, Blocked).
   - Pengguna dapat melihat jalur evakuasi berdasarkan bencana.
   
7. **Laporan Darurat**
   - Pengguna dapat mengirim laporan darurat terkait kebutuhan mendesak di lokasi bencana.
   
## Dokumentasi API
### **1️. Users**
| Method | Endpoint | Deskripsi | Hak Akses |
|--------|---------|-----------|------------|
| POST | `/users/` | Registrasi pengguna baru | Semua Pengguna |
| POST | `/users/login` | Login dan mendapatkan token JWT | Semua Pengguna |
| POST | `/users/verify-otp` | Verifikasi OTP untuk 2FA | Semua Pengguna |
| GET | `/users/` | Mendapatkan daftar semua pengguna | Admin, Donor |
| GET | `/users/:id` | Mendapatkan detail pengguna berdasarkan ID | Admin, Pemilik Akun |
| PUT | `/users/:id` | Mengedit informasi pengguna | Admin, Pemilik Akun |
| DELETE | `/users/:id` | Menghapus akun pengguna | Admin, Pemilik Akun |
| PUT | `/users/enable-2fa` | Mengaktifkan atau menonaktifkan 2FA | Pemilik Akun |
| PUT | `/users/:id/edit-info` | Mengedit info pengguna tanpa email | Admin, Pemilik Akun |
| GET | `/users/:id/emergency-reports` | Mendapatkan semua laporan darurat pengguna | Admin, Pemilik Akun |

### **2️. Disasters**
| Method | Endpoint | Deskripsi | Hak Akses |
|--------|---------|-----------|------------|
| POST | `/disasters/` | Membuat laporan bencana | Admin, Volunteer |
| GET | `/disasters/` | Mendapatkan semua laporan bencana | Semua Pengguna |
| GET | `/disasters/:id` | Mendapatkan detail laporan bencana | Semua Pengguna |
| PUT | `/disasters/:id` | Mengedit laporan bencana | Admin, Volunteer |
| DELETE | `/disasters/:id` | Menghapus laporan bencana | Admin |
| GET | `/disasters/:id/shelters` | Mendapatkan daftar shelter untuk bencana tertentu | Admin, Volunteer |
| GET | `/disasters/:id/volunteers` | Mendapatkan daftar relawan untuk bencana tertentu | Admin, Volunteer |
| GET | `/disasters/:id/logistics` | Mendapatkan daftar bantuan logistik | Admin, Volunteer |
| GET | `/disasters/:id/emergency-reports` | Mendapatkan daftar laporan darurat | Semua Pengguna |
| GET | `/disasters/:id/evacuation-routes` | Mendapatkan daftar jalur evakuasi | Semua Pengguna |

### **3️. Shelters**
| Method | Endpoint | Deskripsi | Hak Akses |
|--------|---------|-----------|------------|
| POST | `/shelters/` | Membuat shelter baru | Admin, Volunteer |
| GET | `/shelters/` | Mendapatkan semua shelter | Semua Pengguna |
| GET | `/shelters/:id` | Mendapatkan detail shelter | Semua Pengguna |
| PUT | `/shelters/:id` | Mengedit informasi shelter | Admin, Volunteer |
| DELETE | `/shelters/:id` | Menghapus shelter | Admin |
| GET | `/shelters/:id/refugees` | Mendapatkan daftar pengungsi di shelter tertentu | Semua Pengguna |
| GET | `/shelters/:id/logistics` | Mendapatkan daftar bantuan logistik di shelter tertentu | Semua Pengguna |

### **4️. Refugees**
| Method | Endpoint | Deskripsi | Hak Akses |
|--------|---------|-----------|------------|
| POST | `/refugees/` | Mendaftarkan pengungsi baru | Admin, Volunteer |
| GET | `/refugees/` | Mendapatkan semua pengungsi | Semua Pengguna |
| GET | `/refugees/:id` | Mendapatkan detail pengungsi | Semua Pengguna |
| PUT | `/refugees/:id` | Mengedit informasi pengungsi | Admin, Volunteer |
| DELETE | `/refugees/:id` | Menghapus data pengungsi | Admin |
| GET | `/refugees/:id/distribution-logs` | Mendapatkan log distribusi bantuan | Semua Pengguna |

### **5️. Logistics**
| Method | Endpoint | Deskripsi | Hak Akses |
|--------|---------|-----------|------------|
| GET | `/logistics/` | Mendapatkan semua bantuan logistik | Semua Pengguna |
| GET | `/logistics/:id` | Mendapatkan detail bantuan logistik | Semua Pengguna |
| POST | `/logistics/` | Menambahkan bantuan logistik baru | Admin, Volunteer |
| PUT | `/logistics/:id` | Mengedit informasi bantuan logistik | Admin, Volunteer |
| DELETE | `/logistics/:id` | Menghapus bantuan logistik | Admin |

### **6️. Distribution Logs**
| Method | Endpoint | Deskripsi | Hak Akses |
|--------|---------|-----------|------------|
| GET | `/distribution-logs/` | Mendapatkan semua log distribusi bantuan | Semua Pengguna |
| GET | `/distribution-logs/:id` | Mendapatkan detail log distribusi bantuan | Semua Pengguna |
| POST | `/distribution-logs/` | Membuat log distribusi bantuan baru | Admin, Volunteer |
| PUT | `/distribution-logs/:id` | Mengedit log distribusi bantuan | Admin, Volunteer |
| DELETE | `/distribution-logs/:id` | Menghapus log distribusi bantuan | Admin |

### **7️. Evacuation Routes**
| Method | Endpoint | Deskripsi | Hak Akses |
|--------|---------|-----------|------------|
| GET | `/evacuation-routes/` | Mendapatkan semua jalur evakuasi | Semua Pengguna |
| GET | `/evacuation-routes/:id` | Mendapatkan detail jalur evakuasi | Semua Pengguna |
| POST | `/evacuation-routes/` | Menambahkan jalur evakuasi baru | Admin, Volunteer |
| PUT | `/evacuation-routes/:id` | Mengedit jalur evakuasi | Admin, Volunteer |
| DELETE | `/evacuation-routes/:id` | Menghapus jalur evakuasi | Admin |

### **8️. Emergency Reports**
| Method | Endpoint | Deskripsi | Hak Akses |
|--------|---------|-----------|------------|
| GET | `/emergency-reports/` | Mendapatkan semua laporan darurat | Admin, Volunteer |
| GET | `/emergency-reports/:id` | Mendapatkan detail laporan darurat | Admin, Volunteer, Pemilik Akun |
| POST | `/emergency-reports/` | Mengirim laporan darurat | Semua Pengguna |
| PUT | `/emergency-reports/:id` | Mengedit laporan darurat | Admin, Pemilik Akun |
| DELETE | `/emergency-reports/:id` | Menghapus laporan darurat | Admin, Pemilik Akun |

### **9️. Donations**
| Method | Endpoint | Deskripsi | Hak Akses |
|--------|---------|-----------|------------|
| GET | `/donations/` | Mendapatkan semua donasi | Admin, Donor |
| GET | `/donations/:id` | Mendapatkan detail donasi | Admin, Donor, Pemilik Akun |
| POST | `/donations/` | Membuat donasi baru | Admin, Donor, Pemilik Akun |
| PUT | `/donations/:id` | Mengedit donasi | Admin, Donor, Pemilik Akun |
| DELETE | `/donations/:id` | Menghapus donasi | Admin, Donor, Pemilik Akun |

### **10. Volunteers**
| Method | Endpoint | Deskripsi | Hak Akses |
|--------|---------|-----------|------------|
| GET | `/volunteers/` | Mendapatkan semua relawan | Admin, Volunteer |
| GET | `/volunteers/:id` | Mendapatkan detail relawan | Admin, Volunteer, Pemilik Akun |
| POST | `/volunteers/` | Mendaftarkan relawan baru | Admin, Volunteer, Pemilik Akun |
| PUT | `/volunteers/:id` | Mengedit informasi relawan | Admin, Volunteer, Pemilik Akun |
| DELETE | `/volunteers/:id` | Menghapus relawan | Admin, Volunteer, Pemilik Akun |

## 🔑 Keamanan API & Role Access
### **1️⃣ Two-Factor Authentication (2FA)**
- **Admin**: Wajib menggunakan Two-Factor Authentication (2FA) - OTP via Email.
- **User**: Dapat mengaktifkan atau menonaktifkan 2FA melalui `/users/enable-2fa`.
- **Donor**: Tidak diwajibkan menggunakan 2FA, tetapi dapat mengaktifkannya.

### **2️⃣ Role-Based Access Control**
| Role | Hak Akses |
|------|----------|
| **Admin** | Mengelola semua fitur termasuk bencana, shelter, relawan, dan donasi. |
| **Donor** | Dapat melihat dan membuat donasi. |
| **User** | Dapat melaporkan bencana, menjadi relawan, dan melihat informasi umum. |
| **Volunteer** | Bisa mengelola shelter, mencatat pengungsi, distribusi bantuan, dan laporan jalur evakuasi. |

🚀 **Sistem ini sudah lengkap dan siap digunakan!** 🎉

