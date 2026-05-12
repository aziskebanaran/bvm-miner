module github.com/aziskebanaran/bvm-miner

go 1.26.2

replace github.com/aziskebanaran/bvm-core => ../bvm-core

require (
	github.com/aziskebanaran/bvm-core v1.1.27
	github.com/aziskebanaran/bvm-lib v1.0.3
)

require (
	github.com/cbergoon/merkletree v0.2.0 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/syndtr/goleveldb v1.0.1-0.20210819022825-2ae1ddf74ef7 // indirect
	github.com/tyler-smith/go-bip39 v1.1.0 // indirect
	github.com/vmihailenco/msgpack/v5 v5.4.1 // indirect
	github.com/vmihailenco/tagparser/v2 v2.0.0 // indirect
	golang.org/x/crypto v0.50.0 // indirect
)

replace github.com/aziskebanaran/bvm-lib => ../bvm-lib
