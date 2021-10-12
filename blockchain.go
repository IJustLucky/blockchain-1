package blockchain

import (
	"fmt"
	"log"

	badger "github.com/dgraph-io/badger"
)

type BlockChain struct {
	LastHash []byte
	Database *badger.DB
}

type BlockChainIterator struct {
	CurrentHash []byte
	Database    *badger.DB
}

func InitBlockChain() *BlockChain {
	var lastHash []byte

	db, err := badger.Open(badger.DefaultOptions("/blocks"))
	Handler(err)

	err = db.Update(func(txn *badger.Txn) error {
		if _, err := txn.Get([]byte("lh")); err == badger.ErrKeyNotFound {
			fmt.Println("No exisiting blockchain found")
			genesis := Genesis()
			err = txn.Set(genesis.Hash, genesis.Serialize())
			Handler(err)
			err = txn.Set([]byte("lh"), genesis.Hash)

			lastHash = genesis.Hash

			return err
		} else {
			item, err := txn.Get([]byte("lh"))
			Handler(err)
			lastHash, err = item.ValueCopy(nil)
			return err
		}
	})

	Handler(err)
	blockchain := BlockChain{lastHash, db}
	return &blockchain
}

func (chain *BlockChain) AddBlock(data string) {
	var lastHash []byte

	err := chain.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lh"))
		Handler(err)
		lastHash, err = item.ValueCopy(nil)

		return err
	})
	Handler(err)

	newBlock := CreateBlock(data, lastHash)

	err = chain.Database.Update(func(txn *badger.Txn) error {
		err := txn.Set(newBlock.Hash, newBlock.Serialize())
		Handler(err)
		err = txn.Set([]byte("lh"), newBlock.Hash)

		chain.LastHash = newBlock.Hash

		return err
	})
	Handler(err)
}

func (chain *BlockChain) Iterator() *BlockChainIterator {
	iter := &BlockChainIterator{chain.LastHash, chain.Database}

	return iter
}

func (iter *BlockChainIterator) Next() *Block {
	var block *Block

	err := iter.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get(iter.CurrentHash)
		Handler(err)
		encodedBlock, err := item.ValueCopy(nil)
		block = Deserialize(encodedBlock)

		return err
	})
	Handler(err)

	iter.CurrentHash = block.PrevHash

	return block
}

func Handler(err error) {
	if err != nil {
		log.Panic(err)
	}
}
