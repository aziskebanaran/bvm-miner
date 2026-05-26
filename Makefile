# BVM Miner Makefile
# Digunakan untuk membangun External Miner (Auto-Navigation)

BINARY_NAME = bvm-miner
SOURCE_FILE = ./cmd/bvm-miner/main.go

.PHONY: all build clean help

all: build

# Membangun Binary Miner
build:
	@echo "🔨 Membangun BVM External Miner..."
	go build -o $(BINARY_NAME) $(SOURCE_FILE)
	@echo "✅ Selesai: ./$(BINARY_NAME)"

# Membersihkan sisa log dan binary
clean:
	@echo "🧹 Membersihkan binary dan log penambangan..."
	rm -f $(BINARY_NAME) miner.log
	@echo "✨ Bersih!"

help:
	@echo "Perintah Miner:"
	@echo "  make build  - Membangun binary penambang"
	@echo "  make clean  - Menghapus binary dan file log"
