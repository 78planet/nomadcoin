package blockchain

import (
	"github.com/78planet/nomadcoin/db"
	"github.com/78planet/nomadcoin/utils"
	"sync"
)

const (
	defaultDifficulty   int = 2
	difficultyInterval  int = 5
	blockIntervalMinute int = 2
	allowedRange        int = 2
)

type blockchain struct {
	NewestHash        string `json:"newest_hash,omitempty"`
	Height            int    `json:"height,omitempty"`
	CurrentDifficulty int    `json:"current_difficulty,omitempty"`
}

func (b *blockchain) persist() {
	db.SaveCheckpoint(utils.ToBytes(b))
}

func (b *blockchain) AddBlock(data string) {
	block := createBlock(data, b.NewestHash, b.Height+1)
	b.NewestHash = block.Hash
	b.Height = block.Height
	b.CurrentDifficulty = block.Difficulty
	b.persist()
}

func (b *blockchain) Blocks() []*Block {
	var blocks []*Block
	hashCursor := b.NewestHash
	for {
		block, _ := FindBlock(hashCursor)
		blocks = append(blocks, block)
		if block.PrevHash != "" {
			hashCursor = block.PrevHash
		} else {
			break
		}
	}
	return blocks
}

func (b *blockchain) restore(data []byte) {
	utils.FromBytes(b, data)
}

func (b *blockchain) difficulty() int {
	if b.Height == 0 {
		return defaultDifficulty
	} else if b.Height%difficultyInterval == 0 {
		// recalculate difficulty
		return b.recalculateDifficulty()
	} else {
		return b.CurrentDifficulty
	}
}

func (b *blockchain) recalculateDifficulty() int {
	allBlocks := b.Blocks()
	newestBlock := allBlocks[0]
	lastRecalculateBlock := allBlocks[difficultyInterval-1]
	actualTime := (newestBlock.Timestamp / 60) - (lastRecalculateBlock.Timestamp / 60)
	expectedTime := difficultyInterval * blockIntervalMinute
	if actualTime <= (expectedTime - allowedRange) {
		return b.CurrentDifficulty + 1
	} else if actualTime >= (expectedTime + allowedRange) {
		return b.CurrentDifficulty - 1
	}
	return b.CurrentDifficulty
}

var b *blockchain
var once sync.Once

func Blockchain() *blockchain {
	if b == nil {
		once.Do(func() {
			b = &blockchain{
				Height: 0,
			}
			checkpoint := db.Checkpoint()
			if checkpoint == nil {
				b.AddBlock("Genesis Block")
			} else {
				b.restore(checkpoint)
			}
		})
	}
	return b
}
