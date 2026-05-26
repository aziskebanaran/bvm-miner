package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"github.com/aziskebanaran/bvm-lib/crypto" // 🚀 GUNAKAN KRYPTO OFF-LINE ANDALAN SULTAN
	"github.com/aziskebanaran/bvm-lib/logger"
	"github.com/aziskebanaran/bvm-miner/miner"
)

const NEXUS_URL = "http://localhost:9092"

func main() {
	// =========================================================================
	// 📡 KONFIGURASI FLAG UTAMA
	// =========================================================================
	nexusURLFlag := flag.String("nexus", NEXUS_URL, "Alamat Node Pusat Nexus Dashboard")
	walletFileFlag := flag.String("wallet", "node_wallet.json", "Lokasi file identitas dompet miner")
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		printUsage()
		return
	}

	command := args[0]

	switch command {
	// --- KASUS 1: MELAHIRKAN DOMPET BARU SECARA OFFLINE 100% ---
	case "create-wallet":
		fmt.Println("⚡ Memulai Ritual Pembuatan Mnemonic & Kunci Kriptografi Offline...")

		// 1. Lahirkan 12 Kata Rahasia (BIP39)
		mnemonic, err := crypto.GenerateMnemonic()
		if err != nil {
			logger.Error("CLI", "❌ Gagal merakit Entropy Mnemonic!")
			return
		}

		// 2. Turunkan Kunci ke Struktur BVMWallet (Index 0 untuk Miner Utama)
		newWallet, err := crypto.CreateFromMnemonic(mnemonic, 0)
		if err != nil {
			logger.Error("CLI", fmt.Sprintf("❌ Gagal menurunkan kunci dari mnemonic: %v", err))
			return
		}

		// 3. Konversi ke bentuk text JSON bersih
		walletBytes, err := json.MarshalIndent(newWallet, "", "    ")
		if err != nil {
			logger.Error("CLI", "❌ Gagal melakukan encoding JSON dompet!")
			return
		}

		// 4. Kunci ke dalam file disk secara lokal/offline
		err = ioutil.WriteFile(*walletFileFlag, walletBytes, 0600) // Hak akses 0600 agar hanya miner yang bisa baca
		if err != nil {
			logger.Error("CLI", fmt.Sprintf("❌ Gagal memahat file %s ke disk!", *walletFileFlag))
			return
		}

		logger.Success("CLI", "✨ Kunci Kriptografi Berhasil Dirakit Secara Offline!")
		fmt.Printf("   🔹 Alamat BVM (Reward) : %s\n", newWallet.Address)
		fmt.Printf("   🔹 Alamat Ethereum L2  : %s\n", newWallet.EthAddress)
		fmt.Printf("   💾 File Identitas      : %s\n", *walletFileFlag)
		logger.Warning("CLI", "⚠️ JANGAN BAGIKAN FILE JSON INI KARENA BERISI PRIVATE KEY SAKRAL ANDA!")

        // --- KASUS 2: MEMULAI EKSPEDISI TAMBANG ---
        case "start":
                fmt.Println("-------------------------------------------")
                fmt.Println("👷 BVM EXTERNAL MINER - SENTINEL MODULAR V3")
                fmt.Println("-------------------------------------------")

                // =========================================================================
                // 📡 SUNTIKAN RADAR SULTAN: AUTO-CALIBRATION PORT LINTAS JARINGAN
                // =========================================================================
                // 1. Cek apakah operator memaksa menyalakan Testnet lewat Environment
                if os.Getenv("BVM_CHAIN_ID") == "" {
                        os.Setenv("BVM_CHAIN_ID", "9999")
                        os.Setenv("BVM_NETWORK_NAME", "BVM Atomic Testnet")
                }

                // 2. KUNCI KOORDINAT PORT TARGET SECARA DINAMIS
                // Jika bendera --nexus TIDAK diubah manual oleh user lewat CLI, kita arahkan otomatis
                if *nexusURLFlag == NEXUS_URL { 
                        if os.Getenv("BVM_CHAIN_ID") == "9999" {
                                *nexusURLFlag = "http://localhost:9092" // 🚩 Belokkan otomatis ke Nexus Testnet!
                        } else {
                                *nexusURLFlag = "http://localhost:9094" // 🚩 Tetap di Nexus Mainnet
                        }
                }

                logger.Info("CLI", fmt.Sprintf("🔒 Frekuensi Jaringan Terkunci: %s (ChainID: %s)",
                        os.Getenv("BVM_CHAIN_ID"), os.Getenv("BVM_NETWORK_NAME")))
                logger.Info("CLI", fmt.Sprintf("📡 TARGET PORT NEXUS          : %s", *nexusURLFlag))
                // =========================================================================

                // 1. Muat Dompet Kerja dari Disk secara Mandiri
                walletBytes, err := ioutil.ReadFile(*walletFileFlag)
                if err != nil {
                        logger.Error("CLI", fmt.Sprintf("❌ Gagal membaca file %s! Silakan buat dompet terlebih dahulu.", *walletFileFlag))
                        return
                }

                var myWallet crypto.BVMWallet
                if err := json.Unmarshal(walletBytes, &myWallet); err != nil {
                        logger.Error("CLI", "❌ Struktur file JSON dompet rusak atau korup!")
                        return
                }

                logger.Info("CLI", fmt.Sprintf("👷 Alamat Reward Aktif: %s", myWallet.Address))
                logger.Info("CLI", fmt.Sprintf("📡 Hub Hubungan Nexus  : %s", *nexusURLFlag))

                // 2. Inisialisasi Klien Komunikasi Dan Mesin Modular
                client := miner.NewMinerClient(*nexusURLFlag)
                engine := miner.NewMiningEngine(client, myWallet.Address)

                // 3. Nyalakan Mesin Penambang!
                engine.StartWorkerLoop()


	default:
		logger.Error("CLI", fmt.Sprintf("❌ Perintah '%s' tidak dikenali oleh sistem.", command))
		printUsage()
	}
}

func printUsage() {
	fmt.Println("\n📋 Panduan Komando BVM-Miner:")
	fmt.Println("   bvm-miner [options] <command>")
	fmt.Println("\nPilihan Perintah (<command>):")
	fmt.Println("   start          : Menyalakan mesin tambang dan absensi ke L1 Core")
	fmt.Println("   create-wallet  : Melahirkan identitas dompet 100% offline (BIP39)")
	fmt.Println("\nPilihan Flags ([options]):")
	fmt.Println("   --nexus        : Mengubah URL target Nexus (Default: http://localhost:9092)")
	fmt.Println("   --wallet       : Mengubah file target identitas (Default: node_wallet.json)")
	fmt.Println("")
}
