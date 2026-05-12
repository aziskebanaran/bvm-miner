#!/bin/bash

# Target Operasi: BVM External Miner
cd ~/bvm-miner

echo "👷 Memulai sinkronisasi BVM Miner ke Awan..."

# 1. Ambil data terbaru dari Awan untuk mencegah konflik (Penambahan Taktis)
echo "📡 Mengunduh data terbaru dari GitHub..."
git pull origin main --rebase

# 2. Menambahkan semua perubahan
git add .

# 3. Perhitungan Versi
latest_tag=$(git tag --list 'v*' | sort -V | tail -n 1)
if [ -z "$latest_tag" ]; then
    next_tag="v1.0.0"
else
    base_version=$(echo $latest_tag | cut -d. -f1,2)
    patch_version=$(echo $latest_tag | cut -d. -f3)
    next_tag="$base_version.$((patch_version + 1))"
fi

echo "📟 Versi Miner terakhir : $latest_tag"
echo "🆕 Menyiapkan versi baru  : $next_tag"

# 4. Komando Jenderal
echo "📝 Apa pesan untuk versi $next_tag ini, Jenderal?"
read message

if [ -z "$message" ]; then
    message="Update Miner: Auto-discovery stabilization"
fi

# 5. Eksekusi Git
git commit -m "$next_tag: $message"
git tag -a "$next_tag" -m "$message"

echo "📡 Mengirim armada Miner ke GitHub..."
# Gunakan branch spesifik agar lebih aman
git push origin main
git push origin "$next_tag"

echo "---------------------------------------"
echo "✅ MISI MINER SELESAI!"
echo "📍 Versi  : $next_tag"
echo "🚀 Status : Worker & Auto-Navigation aman di Awan."
