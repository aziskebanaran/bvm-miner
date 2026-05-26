package miner

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type DiscoveryResponse struct {
	CoreAddress string `json:"core_address"`
	Latency     string `json:"latency"`
	Status      string `json:"status"`
}

type MinerClient struct {
	NexusURL string
}
// 🚀 TIPE DATA SATELLIT UNTUK MENANGKAP STATUS KERNEL CORE
type CoreStatusResponse struct {
        Height     int64  `json:"height"`
        LatestHash string `json:"latest_hash"`
        Status     string `json:"status"`
        Version    int    `json:"version"`
}

func NewMinerClient(nexusURL string) *MinerClient {
	return &MinerClient{NexusURL: nexusURL}
}


// 1. DiscoverCore: Cari lokasi Core Node via Nexus
func (c *MinerClient) DiscoverCore() (string, error) {
	resp, err := http.Get(c.NexusURL + "/api/discover-core")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var disco DiscoveryResponse
	if err := json.NewDecoder(resp.Body).Decode(&disco); err != nil {
		return "", err
	}
	return disco.CoreAddress, nil
}

// 2. GetMempoolTxs: Cek apakah ada antrean transaksi di Core Node
func (c *MinerClient) GetMempoolTxs(coreURL string) (int, error) {
	resp, err := http.Get(coreURL + "/api/mempool")
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var wrapper struct {
		Txs []interface{} `json:"txs"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&wrapper); err != nil {
		return 0, nil
	}
	return len(wrapper.Txs), nil
}

// 3. FetchWork: Ambil tugas blok baru (GetWork) dari Core L1
func (c *MinerClient) FetchWork(coreURL string, minerAddr string) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/api/getwork?address=%s&miner=BVM-EXTERNAL-PROV3", coreURL, minerAddr)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var blockTask map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&blockTask); err != nil {
		return nil, err
	}
	return blockTask, nil
}

// 4. SubmitBlock: Setor hasil pahatan blok yang berhasil dipecahkan
func (c *MinerClient) SubmitBlock(coreURL string, blockData map[string]interface{}) error {
	jsonData, _ := json.Marshal(blockData)
	url := fmt.Sprintf("%s/api/mine", coreURL)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf("node menolak blok, status: %d", resp.StatusCode)
	}
	return nil
}

// 5. ReportToNexus: Laporkan intelijen hashrate ke Dashboard Nexus
func (c *MinerClient) ReportToNexus(minerAddr string, hashes uint64, duration float64) error {
	report := struct {
		MinerAddr string  `json:"miner_addr"`
		Hashes    uint64  `json:"hashes"`
		Duration  float64 `json:"duration"`
	}{MinerAddr: minerAddr, Hashes: hashes, Duration: duration}

	jsonData, _ := json.Marshal(report)
	url := fmt.Sprintf("%s/api/stats/report-hash", c.NexusURL)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

// 6. SendHeartbeatToCore: Absensi dinamis ke Mempool Core
func (c *MinerClient) SendHeartbeatToCore(coreURL string, minerAddr string) error {
	pingData := struct{ Address string `json:"address"` }{Address: minerAddr}
	jsonData, _ := json.Marshal(pingData)
	url := fmt.Sprintf("%s/api/mempool/ping", coreURL)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

// GetNodeStatus: Menembak langsung api /api/status milik Core L1 Jenderal
func (c *MinerClient) GetNodeStatus(coreURL string, minerAddr string) (*CoreStatusResponse, error) {
	url := fmt.Sprintf("%s/api/status?address=%s", coreURL, minerAddr)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("node return status: %d", resp.StatusCode)
	}

	var status CoreStatusResponse
	if err := json.NewDecoder(resp.Body).Decode(&status); err != nil {
		return nil, err
	}
	return &status, nil
}
