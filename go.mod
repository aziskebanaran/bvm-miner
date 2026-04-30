module bvm-miner

go 1.26.2

replace github.com/aziskebanaran/bvm-core => ../bvm-core

require github.com/aziskebanaran/bvm-core v1.1.8

require (
	github.com/aziskebanaran/bvm-lib v0.0.0-00010101000000-000000000000 // indirect
	github.com/cbergoon/merkletree v0.2.0 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/syndtr/goleveldb v1.0.1-0.20210819022825-2ae1ddf74ef7 // indirect
	github.com/tyler-smith/go-bip39 v1.1.0 // indirect
	github.com/vmihailenco/msgpack/v5 v5.4.1 // indirect
	github.com/vmihailenco/tagparser/v2 v2.0.0 // indirect
	golang.org/x/crypto v0.0.0-20220518034528-6f7dac969898 // indirect
)

replace github.com/aziskebanaran/bvm-lib => /data/data/com.termux/files/home/bvm-lib
