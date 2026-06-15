# fogus

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](./LICENSE)

`fogus` adalah terminal tool untuk interstitial journaling.

Tujuannya sederhana: membantu Anda cepat menangkap apa yang sedang dikerjakan, apa yang menginterupsi fokus, dan bagaimana hari Anda benar-benar berjalan.

Tool ini cocok untuk orang yang sering mengalami context switching, lompat antar task, atau butuh external memory yang ringan selama bekerja.

## Kenapa fogus

Interstitial journaling berguna karena Anda tidak perlu menunggu sampai akhir hari untuk mengingat apa yang terjadi. Anda cukup mencatat potongan kecil aktivitas saat itu juga.

`fogus` mencoba menjaga proses itu tetap ringan:

1. terminal-first
2. local-first
3. cepat dipakai
4. format file tetap manusiawi
5. tidak mengunci data Anda di format yang sulit dibaca

## Status saat ini

Versi sekarang adalah MVP yang sudah usable untuk kebutuhan harian pribadi.

Fokusnya belum ke fitur banyak, tapi ke flow dasar yang harus cepat dan stabil.

## Fitur yang sudah ada

1. Mencatat log aktivitas dengan timestamp otomatis
2. Menambahkan log untuk waktu tertentu dengan `--at HH:MM`
3. Menambahkan log untuk tanggal tertentu dengan `--date YYYY-MM-DD`
4. Menulis log satu baris atau multi-line
5. Menulis log lewat editor dengan `--editor`
6. Menyimpan journal harian otomatis ke Markdown
7. Menampilkan timeline hari tertentu dengan `today`
8. Menampilkan summary netral hari tertentu dengan `summary`
9. Menyisipkan entry secara kronologis, bukan sekadar append di bawah

## Quickstart

### Setup dari nol

1. Install Go:

```bash
brew install go
```

2. Clone repo:

```bash
git clone <repo-url>
cd fogus
```

3. Install `fogus` ke `~/.local/bin`:

```bash
make install
```

4. Verifikasi command:

```bash
fogus help
```

5. Coba langsung:

```bash
fogus log "Mulai kerja"
fogus today
fogus summary
```

### Setup singkat

Kalau Go sudah ada dan Anda sudah di repo:

```bash
make install
fogus log "Mulai kerja"
fogus today
fogus summary
```

## Instalasi

Binary akan dipasang ke:

```text
~/.local/bin/fogus
```

Direktori itu harus ada di `PATH` shell Anda.

## Lokasi penyimpanan

Semua log disimpan otomatis di:

```text
~/.fogus/journal
```

Format file per hari:

```text
~/.fogus/journal/YYYY-MM-DD.md
```

Contoh:

```text
~/.fogus/journal/2026-06-12.md
```

Kalau folder belum ada, `fogus` akan membuatnya otomatis.

## Command

```bash
fogus log [--date YYYY-MM-DD] [--at HH:MM] [--editor] [text...]
fogus today [--date YYYY-MM-DD]
fogus summary [--date YYYY-MM-DD]
```

## Penggunaan dasar

### Quick capture

```bash
fogus log "Mulai PBI data retention."
```

`fogus` akan otomatis:

1. mengambil waktu saat ini
2. menentukan file hari ini
3. menyimpan log ke `~/.fogus/journal`

### Multi-line dari terminal

```bash
fogus log
```

Lalu tulis isi log, dan akhiri dengan `Ctrl+D`.

Contoh:

```text
09:47 - write log, end with Ctrl+D
Kepikiran cek AWS billing.
Tahan dulu. Lanjut retention.
```

### Multi-line lewat editor

```bash
export EDITOR=nvim
fogus log --editor
```

`fogus` akan membuka editor dari environment variable `$EDITOR`.

Kalau `$EDITOR` belum diset, command ini akan gagal.

### Backfill waktu tertentu

Kalau Anda lupa log tepat waktu:

```bash
fogus log --at 09:45 "Selesai analisa transactions.parquet."
```

Entry akan disisipkan ke posisi waktu yang benar.

### Backfill tanggal tertentu

Kalau Anda ingin mencatat ke hari sebelumnya:

```bash
fogus log --date 2026-06-12 --at 11:20 "Stuck di query finder."
```

### Lihat timeline

```bash
fogus today
fogus today --date 2026-06-12
```

Jika belum ada log:

```text
No entries for 2026-06-13
```

### Lihat summary

```bash
fogus summary
fogus summary --date 2026-06-12
```

Summary versi sekarang bersifat netral. Belum ada grouping seperti `focus`, `break`, atau `blocker`.

## Contoh workflow nyata

Contoh ini cukup dekat dengan pola kerja yang sering terjadi saat fokus pecah, ada distraksi kecil, lalu harus kembali ke task utama.

### Mulai kerja

```bash
fogus log "Mulai PBI data retention."
```

### Selesai satu langkah analisa

```bash
fogus log --at 09:45 "Selesai analisa transactions.parquet."
```

### Ada distraksi mendadak

```bash
fogus log
```

Lalu isi:

```text
Kepikiran cek AWS billing.
Tahan dulu. Lanjut retention.
```

### Ambil break singkat

```bash
fogus log "Break kopi."
```

### Kembali ke task utama

```bash
fogus log "Balik ke script batch delete."
```

### Stuck dan butuh referensi

```bash
fogus log
```

Lalu isi:

```text
Stuck di query finder.
Cari referensi schema.
```

### Review hari berjalan

```bash
fogus today
fogus summary
```

Dengan flow seperti ini, Anda tidak perlu mengingat semuanya di kepala. `fogus` menjadi jejak kecil yang membantu Anda:

1. melihat progres nyata
2. menangkap interupsi sebelum hilang
3. kembali ke konteks yang benar setelah terdistraksi
4. meninjau pola kerja di akhir hari

## Contoh output

### Timeline

```text
2026-06-12

09:00
Mulai PBI data retention.

09:45
Selesai analisa transactions.parquet.

09:47
Kepikiran cek AWS billing.
Tahan dulu. Lanjut retention.
```

### Summary

```text
2026-06-12
6 entries
First log: 09:00
Last log: 11:20

Entries
- 09:00 Mulai PBI data retention.
- 09:45 Selesai analisa transactions.parquet.
- 09:47 Kepikiran cek AWS billing. (+1 line)
```

## Format file journal

Contoh file hasil simpan:

```md
# 2026-06-12

09:00
Mulai PBI data retention.

09:45
Selesai analisa transactions.parquet.

09:47
Kepikiran cek AWS billing.
Tahan dulu. Lanjut retention.
```

Format ini sengaja sederhana supaya:

1. mudah dibaca langsung
2. mudah di-backup
3. mudah dipindahkan ke Obsidian atau Git
4. tetap berguna walau tanpa `fogus`

## Perilaku penting

1. Entry dengan timestamp yang lebih awal akan disisipkan ke tengah file jika perlu
2. Entry dengan timestamp yang sama akan mempertahankan urutan input
3. File journal ditulis ulang dengan format yang konsisten setiap ada perubahan
4. `today` dan `summary` membaca sumber Markdown yang sama

## Batasan versi sekarang

1. Belum ada tagging seperti `#break` atau `#blocker`
2. Belum ada grouping summary otomatis
3. Belum ada weekly review
4. Belum ada sync ke Obsidian vault tertentu
5. Belum ada mode TUI full-screen

## Roadmap

### Near term

1. Grouping summary sederhana seperti `focus`, `break`, `blocker`, dan `interruption`
2. Tagging ringan tanpa membuat input terasa berat
3. Weekly review untuk melihat pola kerja mingguan
4. Perintah install dan distribusi yang lebih rapi

### Mid term

1. Integrasi opsional ke Obsidian vault tertentu
2. Export atau format output yang lebih cocok untuk refleksi harian
3. Pencarian dan filter berdasarkan tanggal atau kata kunci
4. Prompt atau shortcut untuk quick capture yang lebih cepat

### Longer term

1. TUI mode untuk journaling dan review tanpa meninggalkan terminal
2. Template harian atau review terstruktur
3. Insight ringan dari pola journaling, tanpa membuat tool terasa berat
4. Positioning yang lebih matang untuk open source dan distribusi berbayar

## Untuk open source atau produk

Kalau `fogus` nanti dibuka ke publik atau dijual, positioning yang menurut saya kuat adalah:

1. journaling tool untuk orang yang sering context-switching
2. terminal-native capture tool untuk pekerja knowledge work
3. local-first daily work log yang tetap manusiawi dibaca

Nilai jual utamanya bukan sekadar "CLI note app", tapi tool yang membantu orang kembali ke konteks kerja saat perhatian mereka pecah.

## Development

```bash
make build
make test
make install
make clean
```
