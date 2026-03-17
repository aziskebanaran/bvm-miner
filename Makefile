# Makefile untuk BVM Miner Pro
BINARY_NAME=bvm-mining
MASTER_URL=http://127.0.0.1:8080
ADDRESS=bvmf110f63ad9e15ed26e9e

.PHONY: run build clean update

# Perintah utama: make run
run:
	@echo "⛏️ Memulai BVM Worker..."
	@go run main.go --master $(MASTER_URL) --address $(ADDRESS)

# Perintah untuk update kode dari GitHub Sultan
update:
	@echo "🔄 Mengambil update terbaru dari GitHub..."
	@git pull origin main
	@go mod tidy

# Perintah untuk build jadi file mentah (agar lebih cepat jalan)
build:
	@echo "🛠️ Compiling binary..."
	@go build -o $(BINARY_NAME) main.go

clean:
	@rm -f $(BINARY_NAME)
