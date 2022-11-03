package db

import (
	"github.com/78planet/nomadcoin/utils"
	"github.com/boltdb/bolt"
)

const (
	dbName      = "blockchain.db"
	dataBucket  = "data"
	blockBucket = "blocks"
	checkpoint  = "checkpoint"
)

var db *bolt.DB

func DB() *bolt.DB {
	if db == nil {
		dbPointer, err := bolt.Open(dbName, 0600, nil)
		utils.HandleErr(err)
		db = dbPointer
		err = db.Update(func(tx *bolt.Tx) error {
			_, err := tx.CreateBucketIfNotExists([]byte(dataBucket))
			utils.HandleErr(err)
			_, err = tx.CreateBucketIfNotExists([]byte(blockBucket))
			return err
		})
		utils.HandleErr(err)
	}
	return db
}

func Close() {
	utils.HandleErr(DB().Close())
}

func SaveBlock(hash string, data []byte) {
	//fmt.Printf("Saving block %s\n Data: %b\n", hash, data)
	err := DB().Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blockBucket))
		err := bucket.Put([]byte(hash), data)
		return err
	})
	utils.HandleErr(err)
}

func SaveCheckpoint(data []byte) {
	err := DB().Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(dataBucket))
		err := bucket.Put([]byte(checkpoint), data)
		return err
	})
	utils.HandleErr(err)
}

func Checkpoint() []byte {
	var data []byte
	err := DB().View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(dataBucket))
		data = bucket.Get([]byte(checkpoint))
		return nil
	})
	utils.HandleErr(err)
	return data
}

func Block(hash string) []byte {
	var data []byte
	err := DB().View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blockBucket))
		data = bucket.Get([]byte(hash))
		return nil
	})
	utils.HandleErr(err)
	return data
}
