package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
    "time"

    "github.com/aziskebanaran/bvm-lib/logger"
    "github.com/aziskebanaran/bvm-core/pkg/wallet"
    "github.com/aziskebanaran/bvm-core/pkg/miner" // Pastikan ini diimport
)

const NEXUS_URL = "http://localhost:9092"

// DiscoveryResponse untuk navigasi otomatis
type DiscoveryResponse struct {
    CoreAddress string `json:"core_address"`
    Latency     string `json:"latency"`
    Status      string `json:"status"`
}

// reportToNexus mengirim intelijen hashrate ke Dashboard Nexus
func reportToNexus(minerAddr string, hashes uint64, duration float64) {
    report := struct {
        MinerAddr string  `json:"miner_addr"`
        Hashes    uint64  `json:"hashes"`
        Duration  float64 `json:"duration"`
    }{
        MinerAddr: minerAddr,
        Hashes:    hashes,
        Duration:  duration,
    }

    jsonData, _ := json.Marshal(report)
    url := fmt.Sprintf("%s/api/stats/report-hash", NEXUS_URL)
    
    resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
    if err != nil {
        return // Nexus sedang sibuk, abaikan laporan ini
    }
    defer resp.Body.Close()
}

func main() {
    fmt.Println("-------------------------------------------")
    fmt.Println("👷 BVM EXTERNAL MINER - SENTINEL MODE 2.0")
    fmt.Println("-------------------------------------------")

    // 1. Muat Identitas
    myWallet, err := wallet.LoadWallet("node_wallet.json")
    if err != nil {
        logger.Error("MINER", "❌ Gagal muat wallet!")
        return
    }

    logger.Info("MINER", fmt.Sprintf("👷 Alamat Reward: %s", myWallet.Address))

    for {
        // 2. NAVIGASI: Cari Koordinat
        resp, err := http.Get(NEXUS_URL + "/api/discover-core")
        if err != nil {
            time.Sleep(5 * time.Second)
            continue
        }
        var disco DiscoveryResponse
        json.NewDecoder(resp.Body).Decode(&disco)
        resp.Body.Close()

        if disco.CoreAddress == "" {
            time.Sleep(5 * time.Second)
            continue
        }

        // 3. OPERASI PENAMBANGAN DI BACKGROUND
        // Menggunakan 'go' agar StartMining tidak mengunci perulangan
        logger.Info("MINER", "🔥 Memulai operasi penambangan di background...")
        go miner.StartMining(myWallet.Address, disco.CoreAddress)

        // 4. LOOP PELAPORAN (Reporting Loop)
        // Lapor setiap 10 detik tanpa menghentikan Miner
        for {
            startTime := time.Now()
            time.Sleep(10 * time.Second) // Interval laporan

            // Gunakan angka benchmark Jenderal (1.5jt hash per batch)
            hashesDone := uint64(1545775) 
            duration := time.Since(startTime).Seconds()

            // Kirim intelijen ke Nexus Dashboard
            reportToNexus(myWallet.Address, hashesDone, duration)
            logger.Info("MINER", "📡 Intelijen Hashrate terkirim ke Nexus.")

            // Opsional: Re-discovery setiap 5 menit untuk update koordinat
            if time.Since(startTime) > 5*time.Minute {
                break 
            }
        }
    }
}
