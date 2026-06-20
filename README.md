# Shared JSON / Text Share

Aplikasi sederhana untuk berbagi teks, JSON, atau konten apa pun antar perangkat saat bekerja dengan banyak perangkat.

## Tujuan

- Memudahkan berbagi teks dan JSON secara real-time.
- Memfasilitasi transfer konten antar device tanpa perlu email atau chat.
- Mendukung pengiriman teks bebas, bukan hanya JSON.

## Cara pakai

1. Jalankan server:
   ```bash
   go run main.go
   ```
2. Buka browser ke `http://localhost:8080`.
3. Tempel teks, JSON, atau konten lain ke textarea.
4. Tekan `Enter` untuk mengirim dan membagikannya ke semua perangkat yang terhubung.
5. Klik tombol `Copy` pada setiap pesan untuk menyalinnya ke clipboard.

## Fitur

- Berbagi konten real-time melalui WebSocket.
- Salin ke clipboard dengan fallback untuk browser yang tidak mendukung API Clipboard modern.
- Pesan terbaru ditampilkan di bagian atas.

## Catatan

Aplikasi ini cocok untuk bekerja dengan beberapa perangkat, misalnya saat berpindah dari laptop ke tablet atau dari PC kantor ke komputer rumah, dan ingin menyalin teks/JSON dengan cepat.
