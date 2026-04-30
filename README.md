# 👷 BVM External Miner (Auto-Navigation)

BVM Miner adalah unit pekerja dalam ekosistem **BVM (Bitcoin Virtual Machine)** yang bertugas melakukan komputasi blok. Miner ini dilengkapi dengan fitur **Auto-Navigation** yang secara otomatis meminta koordinat Core aktif kepada **Nexus Gateway**.

## ⚙️ Cara Kerja
1. **Discovery**: Menghubungi Nexus di `http://localhost:9092/api/discover-core`.
2. **Handshake**: Mendapatkan alamat Core dengan latensi terendah.
3. **Mining**: Melakukan proses penambangan dan mengirimkan bukti kerja ke Core.
4. **Auto-Reconnect**: Jika koneksi terputus, miner akan melakukan re-discovery secara otomatis.

## 🏗️ Prasyarat
Miner ini bergantung pada library lokal `bvm-core`. Pastikan struktur folder Anda:
```text
/home/
  ├── bvm-core/
  └── bvm-miner/

🚀 Memulai Penambangan

1. Siapkan Dompet Reward:
Salin file dompet dari unit wallet ke sini:

cp ../bvm-wallet/wallet.json ./node_wallet.json

2. Build & Run:

make build
./bvm-miner

📊 Monitoring

Aktivitas penambangan dapat dipantau melalui log:

tail -f miner.log

