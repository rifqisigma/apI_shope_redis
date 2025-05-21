# Golang REST API with Redis Integration

Project ini merupakan implementasi REST API menggunakan **Golang** dengan kombinasi **Redis** sebagai datastore dan antrian tugas. Didesain dengan **clean architecture** dan memiliki fitur caching & worker system yang optimal.

---

## 🚀 Studi Kasus

REST API ini dibangun dengan studi kasus toko online sederhana, dengan fitur:

- **Autentikasi pengguna** menggunakan JWT.
- **Penyimpanan data utama** menggunakan MySQL via GORM.
- **Integrasi Redis** sebagai datastore dan worker queue.
- **Worker Redis** untuk mengirim email melalui task queue.
- **Strategi caching**:
  - Lazy loading pada store.
  - Write-around caching pada product.
  - Write-through caching pada cart item.

cmi