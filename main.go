package main

import (
	"bvm-core/modules/wallet"
	"bvm-core/pkg/client"
	"bvm-core/pkg/logger"
	"bvm-core/x/bvm/types" // Import types untuk struct Block
	"fmt"
	"time"
        "net/http"
        "encoding/json" // <--- MASUKKAN LAGI INI
)

const CORE_URL = "http://127.0.0.1:8080"

func main() {
    bvm := client.NewBVMClient(CORE_URL)

    myWallet, err := wallet.LoadWallet("node_wallet.json")
    if err != nil {
        logger.Error("WORKER", "Gagal memuat wallet. Pastikan node_wallet.json ada.")
        return
    }

    logger.Info("WORKER", "BVM Worker Pro Aktif")
    logger.Info("WORKER", "Mining untuk Address: "+myWallet.Address)

    // --- TAMBAHKAN INI: JANTUNG OTOMATIS (Background Heartbeat) ---
    // Ini memastikan radar Kernel tetap mendeteksi HP ini meskipun mining sedang berat
    go func() {
        for {
            // Lapor ke Kernel setiap 30 detik
            bvm.GetNodeStatus(myWallet.Address) 
            time.Sleep(30 * time.Second)
        }
    }()

    // Mulai siklus penambangan utama
    startMining(bvm, myWallet.Address)
}

func startMining(bvm *client.BVMClient, minerAddr string) {
    for {
        status, err := bvm.GetNodeStatus(minerAddr)
        if err != nil {
            fmt.Println("❌ Gagal mengambil status node, mencoba lagi...")
            time.Sleep(5 * time.Second)
            continue
        }

        // --- 1. AMBIL TX DARI MEMPOOL ---
        resp, err := http.Get(CORE_URL + "/api/mempool")
        var txs []types.Transaction // Inisialisasi awal

        if err == nil {
            // Gunakan struktur yang persis sama dengan yang dikirim Kernel
            var data struct {
                Count int                 `json:"count"`
                Txs   []types.Transaction `json:"txs"`
            }
            
            errDecode := json.NewDecoder(resp.Body).Decode(&data)
            resp.Body.Close()

            if errDecode == nil && data.Txs != nil {
                txs = data.Txs
            }
        }

        // --- 2. LOG VALIDASI ---
        if len(txs) > 0 {
            fmt.Printf("\n🔥 MENDETEKSI %d TRANSAKSI! Mengangkut ke Blok #%d...", len(txs), status.Height)
            fmt.Printf("\n📝 Detail: %s -> %s (%.2f BVM)", txs[0].From[:10], txs[0].To[:10], txs[0].Amount)
        }

        // --- 3. RAKIT BLOK ---
        newBlock := types.Block{
            Index:        status.Height,
            Timestamp:    time.Now().Unix(),
            PrevHash:     status.LastHash,
            Difficulty:   status.TargetDifficulty,
            Miner:        minerAddr,
            Transactions: txs, // DIPASTIKAN TERISI SEKARANG
            Data:         fmt.Sprintf("BVM Sultan Node | %d Txs", len(txs)),
        }

        // --- 4. POW & SUBMIT ---
        // PENTING: Gunakan method CalculateHash dari types
        newBlock.Hash, newBlock.Nonce = performProofOfWork(newBlock)
        
        // Pastikan SDK client Sultan mendukung SubmitBlock(types.Block)
        errSubmit := bvm.SubmitBlock(newBlock)
        if errSubmit != nil {
            fmt.Printf("\n❌ Gagal setor blok: %v", errSubmit)
        } else {
            fmt.Print(" ✅ BLOK DITERIMA KERNEL!")
        }

        time.Sleep(1 * time.Second)
    }
}


func performProofOfWork(b types.Block) (string, int) {
	nonce := 0
	start := time.Now()

	for {
		b.Nonce = nonce
		hash := b.CalculateHash() // Pastikan types.Block punya method CalculateHash()

		// Cek apakah hash diawali dengan nol sebanyak Difficulty
		if b.HasValidTarget(hash) { 
			duration := time.Since(start)
			fmt.Printf("\n✨ Blok Ditemukan!\n🔗 Hash: %s\n⏱️ Waktu: %v | Nonce: %d\n", hash, duration, nonce)
			return hash, nonce
		}

		nonce++
		if nonce%1000000 == 0 {
			fmt.Print("⚡") // Indikator 1 juta percobaan
		}
	}
}
