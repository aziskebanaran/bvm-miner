package miner

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/aziskebanaran/bvm-lib/logger"
	coreTypes "github.com/aziskebanaran/bvm-core/x/bvm/types"
)

type MiningEngine struct {
	Client    *MinerClient
	WalletHex string
}

func NewMiningEngine(client *MinerClient, walletHex string) *MiningEngine {
	// Inisialisasi generator angka acak berbasis waktu agar seed unik antar-terminal Termux
	rand.Seed(time.Now().UnixNano())
	return &MiningEngine{
		Client:    client,
		WalletHex: walletHex,
	}
}

func (e *MiningEngine) StartWorkerLoop() {
	logger.Success("MINER-ENGINE", "⚡ Modular V3 Active Hashing Engine (Mode Anti-Spam Multi-Miner) Aktif!")

	if os.Getenv("BVM_CHAIN_ID") == "" {
		os.Setenv("BVM_CHAIN_ID", "9999")
	}

	// Buat pencatatan internal untuk mencegah penambangan berulang blok usang
	var lastProcessedHeight int64 = 0

	for {
		// 1. Navigasi koordinat Core L1 lewat Nexus
		coreURL, err := e.Client.DiscoverCore()
		if err != nil || coreURL == "" {
			logger.Warning("MINER-ENGINE", "🌐 Gagal navigasi koordinat Core via Nexus, menanti 5 detik...")
			time.Sleep(5 * time.Second)
			continue
		}

		// SENSOR MEMPOOL RADAR: Jika kosong, jangan ngegas (Hemat CPU)
		txCount, err := e.Client.GetMempoolTxs(coreURL)
		if err != nil {
			logger.Error("MINER-ENGINE", "⚠️ Gagal cek Mempool. Core mungkin offline. Re-sync 5s...")
			time.Sleep(5 * time.Second)
			continue
		}

		if txCount == 0 {
			_ = e.Client.SendHeartbeatToCore(coreURL, e.WalletHex)
			time.Sleep(3 * time.Second)
			continue
		}

		// 2. MINTA PEKERJAAN RIIL KARENA TRANSAKSI TERDETEKSI
		task, err := e.Client.FetchWork(coreURL, e.WalletHex)
		if err != nil {
			time.Sleep(2 * time.Second)
			continue
		}

		// 🔄 TRANSFER DATA SECARA TOTAL KE CETAKAN ASLI CORE BLOK SULTAN
		taskBytes, _ := json.Marshal(task)
		var coreBlock coreTypes.Block
		if err := json.Unmarshal(taskBytes, &coreBlock); err != nil {
			time.Sleep(2 * time.Second)
			continue
		}

		// =========================================================================
		// 🛡️ BARIKADE 1: BLOCK HEIGHT GUARD (PENCEGAH PENAMBANGAN BLOK USANG)
		// =========================================================================
		// Tanya ke Core status ketinggian blok detik ini juga
		status, err := e.Client.GetNodeStatus(coreURL, e.WalletHex)
		if err == nil && status != nil {
			// Jika indeks blok yang kita bawa ternyata sudah tertinggal atau sama dengan tinggi Core sekarang
			if coreBlock.Index <= status.Height {
				logger.Warning("MINER-ENGINE", fmt.Sprintf("⏩ Blok #%d sudah dipahat miner lain (Tinggi Core: %d). Abort!", coreBlock.Index, status.Height))
				// Ambil napas acak sejenak agar tidak langsung membombardir HTTP GetWork kembali
				jitter := rand.Intn(3) + 1
				time.Sleep(time.Duration(jitter) * time.Second)
				continue
			}
		}

		// Jika tinggi blok baru ini sama dengan yang barusan kita kerjakan dan gagal, beri jeda pengaman
		if coreBlock.Index == lastProcessedHeight {
			jitter := rand.Intn(2) + 1
			time.Sleep(time.Duration(jitter) * time.Second)
		}

		logger.Info("MINER-ENGINE", fmt.Sprintf("🔥 Radar Sinkron! Menambang Blok #%d | Diff: %d", coreBlock.Index, coreBlock.Difficulty))

		// 3. PROSES PROOF OF WORK MENGGUNAKAN MESIN ORIGINAL CORE BINARY
		startTime := time.Now()
		targetPrefix := strings.Repeat("0", int(coreBlock.Difficulty))

		var nonce int64 = 0 
		var finalHash string
		found := false

		// Jalankan siklus pencarian dengan interupsi pengaman berkala
		for i := 0; i < 3000000; i++ { 
			nonce++
			coreBlock.Nonce = int32(nonce) 
			hashStr := coreBlock.CalculateBlockHash()

			if strings.HasPrefix(hashStr, targetPrefix) {
				finalHash = hashStr
				coreBlock.Hash = hashStr
				found = true
				break
			}
		}

		duration := time.Since(startTime).Seconds()
		lastProcessedHeight = coreBlock.Index

		if found {
			logger.Success("MINER-ENGINE", fmt.Sprintf("💎 SUKSES MEMECAHKAN BLOK #%d! Hash: %s", coreBlock.Index, finalHash))

			var finalTask map[string]interface{}
			finalBytes, _ := json.Marshal(coreBlock)
			json.Unmarshal(finalBytes, &finalTask)

			// 4. SETOR BLOK KEMENANGAN KE CORE NODE JERAL
			if err := e.Client.SubmitBlock(coreURL, finalTask); err != nil {
				// =========================================================================
				// 💤 STRATEGI SULTAN SLEEP: JEDA ACAK KETIKA DITOLAK JELANG BALAPAN
				// =========================================================================
				// Jika ditolak karena didahului miner lain (406/500/409), aktivasi jeda acak 2-5 detik!
				randomSleepSec := rand.Intn(4) + 2 
				logger.Error("MINER-ENGINE", fmt.Sprintf("❌ Blok ditolak/didahului. Aktivasi Protokol Sultan Sleep %d detik...", randomSleepSec))
				time.Sleep(time.Duration(randomSleepSec) * time.Second)
			} else {
				logger.Success("MINER-ENGINE", fmt.Sprintf("🧱 BERHASIL! Blok #%d Resmi Dipatenkan di Blockchain Jenderal!", coreBlock.Index))
				// Sukses memahat, istirahat 1 detik bersih untuk sinkronisasi state database database L1
				time.Sleep(1 * time.Second)
			}
		} else {
			// Jika dalam 3 juta putaran belum ketemu dan transaksi belum habis, ambil jeda mikro 500ms
			time.Sleep(500 * time.Millisecond)
		}

		_ = e.Client.ReportToNexus(e.WalletHex, uint64(nonce), duration)
		_ = e.Client.SendHeartbeatToCore(coreURL, e.WalletHex)
	}
}
