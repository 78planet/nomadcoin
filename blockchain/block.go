package blockchain

import (
	"errors"
	"fmt"
	"github.com/78planet/nomadcoin/db"
	"github.com/78planet/nomadcoin/utils"
	"strings"
	"time"
)

type Block struct {
	Data       string `json:"data,omitempty"`
	Hash       string `json:"hash,omitempty"`
	PrevHash   string `json:"prev_hash,omitempty"`
	Height     int    `json:"height,omitempty"`
	Difficulty int    `json:"difficulty,omitempty"`
	Nonce      int    `json:"nonce"`
	Timestamp  int    `json:"timestamp,omitempty"`
}

func (b *Block) persist() {
	db.SaveBlock(b.Hash, utils.ToBytes(b))
}

func (b *Block) restore(data []byte) {
	utils.FromBytes(b, data)
}

func (b *Block) mine() {
	target := strings.Repeat("0", b.Difficulty)
	for {
		b.Timestamp = int(time.Now().Unix())
		hash := utils.Hash(b)
		fmt.Printf("Target:%s\n Hash:%s\n Nonce:%d\n\n\n", target, hash, b.Nonce)
		if strings.HasPrefix(hash, target) {
			b.Hash = hash
			break
		} else {
			b.Nonce++
		}
	}
}

var ErrNotFound = errors.New("Block Not Found")

func FindBlock(hash string) (*Block, error) {
	blockBytes := db.Block(hash)
	if blockBytes == nil {
		return nil, ErrNotFound
	}
	block := &Block{}
	block.restore(blockBytes)
	return block, nil
}

func createBlock(data string, prevHash string, height int) *Block {
	block := Block{
		Data:       data,
		Hash:       "",
		PrevHash:   prevHash,
		Height:     height,
		Difficulty: Blockchain().difficulty(),
		Nonce:      0,
	}
	block.mine()
	block.persist()
	return &block
}
