package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/aziskebanaran/bvm-core/pkg/logger"
	"github.com/aziskebanaran/bvm-core/pkg/miner"
	"github.com/aziskebanaran/bvm-core/pkg/wallet"
)

const NEXUS_URL = "http://localhost:9092"

// Struktur yang sama dengan yang ada di Nexus Discovery
type DiscoveryResponse struct {
	CoreAddress string `json:"core_address"`
	Latency     string `json:"latency"`
	Status      string `json:"status"`
}

func main() {
	fmt.Println("-------------------------------------------")
	fmt.Println("👷 BVM EXTERNAL MINER - AUTO-NAVIGATION")
	fmt.Println("-------------------------------------------")

	// 1. Muat Identitas Penambang
	myWallet, err := wallet.LoadWallet("node_wallet.json")
	if err != nil {
		logger.Error("MINER", "❌ Gagal muat wallet! Jalankan 'cp ../bvm-wallet/wallet.json ./node_wallet.json'")
		return
	}

	logger.Info("MINER", fmt.Sprintf("👷 Alamat Reward: %s", myWallet.Address))

	for {
		logger.Info("MINER", "📡 Menghubungi Nexus untuk koordinat Core...")

		// 2. TANYA NEXUS: "Di mana Core-nya?"
		resp, err := http.Get(NEXUS_URL + "/api/discover-core")
		if err != nil {
			logger.Error("MINER", "❌ Nexus Offline. Mencoba lagi dalam 5 detik...")
			time.Sleep(5 * time.Second)
			continue
		}

		var disco DiscoveryResponse
		json.NewDecoder(resp.Body).Decode(&disco)
		resp.Body.Close()

		logger.Info("MINER", fmt.Sprintf("✅ Koordinat Diterima: %s (Status: %s)", disco.CoreAddress, disco.Status))

		// 3. MULAI MENAMBANG ke alamat yang diberikan Nexus
		// Kita gunakan disco.CoreAddress sebagai target kerja
		miner.StartMining(myWallet.Address, disco.CoreAddress)

		logger.Error("MINER", "⚠️ Koneksi terputus. Melakukan re-discovery...")
		time.Sleep(5 * time.Second)
	}
}
