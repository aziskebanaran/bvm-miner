#!/bin/bash

# Target Operasi: BVM External Miner
cd ~/bvm-miner

echo "👷 Memulai sinkronisasi BVM Miner ke Awan..."

# 1. AMANKAN PERUBAHAN LOKAL DULU (Pindahkan git add ke sini)
git add .

# 2. PERHITUNGAN VERSI
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

# 3. KOMANDO JENDERAL
echo "📝 Apa pesan untuk versi $next_tag ini, Jenderal?"
read message
if [ -z "$message" ]; then
    message="Update Miner: v$next_tag"
fi

# 4. COMMIT LOKAL (Wajib sebelum Pull/Rebase)
git commit -m "$next_tag: $message"
git tag -a "$next_tag" -m "$message"

# 5. SINKRONISASI DENGAN AWAN (Setelah Commit Aman)
echo "📡 Menyelaraskan data dengan GitHub..."
git pull origin main --rebase --no-edit

# 6. PUSH ARMADA KE GITHUB
echo "🚀 Mengirim armada ke GitHub..."
git push origin main
git push origin "$next_tag"

echo "---------------------------------------"
echo "✅ MISI MINER SELESAI!"
echo "📍 Versi  : $next_tag"
echo "🚀 Status : Worker & Auto-Navigation aman di Awan."
