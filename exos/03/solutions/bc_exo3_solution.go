package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"time"

	couchdb "github.com/zemirco/couchdb"
)

//-----------------------------------------------------------------
// Blockchain models and api
//-----------------------------------------------------------------
type Block struct {
	ID        string `json:"_id,omitempty"`
	Rev       string `json:"_rev,omitempty"`
	LastID    string `json:"_LastId,omitempty"`
	CreatedAt string `json:"_CreatedAt,omitempty"`
	Data      Data   `json:"Data"`
}

// GetID returns document id
func (d *Block) GetID() string {
	return d.ID
}

// GetRev returns document revision
func (d *Block) GetRev() string {
	return d.Rev
}

// Data structure which contains the transaction data.
type Data struct {
	Reference string  `json:"reference,omitempty"`
	Quantity  int     `json:"quantity,omitempty"`
	Price     float32 `json:"price,omitempty"`
}

// Chain structure.
type Chain struct {
	lastKey string
	blocks  map[string]Block
}

// Enrich chain with a add block feature
func (chain *Chain) addBlock(b Block) {
	lastID := chain.last().ID
	b.ID = hashBlock(b)
	b.LastID = lastID
	chain.lastKey = b.ID
	chain.blocks[b.ID] = b
}

/**
 * Enrich chain with a last block function.
 */
func (chain *Chain) last() Block {
	return chain.blocks[chain.lastKey]
}

// Create the default chain structure.
func genesis() Chain {
	return Chain{
		lastKey: "",
		blocks:  make(map[string]Block),
	}
}

func hashBlock(b Block) string {
	convertedBlock, err := json.Marshal(b)
	if err != nil {
		panic(err)
	}
	//fmt.Println(string(convertedBlock))
	return hash([]byte(string(convertedBlock)))
}

func hash(obj []byte) string {
	fmt.Println(obj)
	h := sha256.New()
	_, err := h.Write(obj)
	if err != nil {
		panic(err)
	}
	return hex.EncodeToString(h.Sum(nil))
}

// Create a new default block with data.
func generateBlock(reference string, quantity int, price float32) Block {
	return Block{
		LastID:    "",
		CreatedAt: time.Now().Format(time.RFC3339),
		Data:      Data{Reference: reference, Price: price, Quantity: quantity},
	}
}

//-----------------------------------------------------------------
// Db models and api
//-----------------------------------------------------------------

type ClientFactory struct {
	DbName         string
	Uri            string
	RemoteClient   *couchdb.Client
	RemoteDatabase couchdb.DatabaseService
}

/*
* Connect on the couch db.
 */
func (factory *ClientFactory) connect() (string, error) {
	var STATUS = "CONNECTED"
	if factory.Uri == "" {
		return "NOT " + STATUS, errors.New("Uri must be initialized !!!")
	}

	u, err := url.Parse("http://" + factory.Uri)
	if err != nil {
		STATUS = "NOT " + STATUS
		fmt.Println("Can't parse URL: ", err)
	}

	client, err := couchdb.NewClient(u)
	if err != nil {
		STATUS = "NOT " + STATUS
		fmt.Println("Error when connecting on database: ", err)
	}
	factory.RemoteClient = client
	return STATUS, err
}

func (factory *ClientFactory) info() {
	info, err := factory.RemoteClient.Info()
	if err != nil {
		panic(err)
	}
	fmt.Println(info)
}

func (factory *ClientFactory) isDbExists() bool {

	if factory.RemoteClient == nil {
		fmt.Println("factory.RemoteClient is nil")
	}
	if info, err := factory.RemoteClient.Get(factory.DbName); err != nil && info == nil {
		return false
	}
	return true
}

func (factory *ClientFactory) createDb() {
	fmt.Println("factory db name: ", factory.DbName)
	_, err := factory.RemoteClient.Create(factory.DbName)
	if err != nil {
		panic(err)
	}
	factory.RemoteClient.Use(factory.DbName)
}

func (factory *ClientFactory) execute() couchdb.DatabaseService {
	factory.RemoteDatabase = factory.RemoteClient.Use(factory.DbName)
	return factory.RemoteDatabase
}

func main() {

	fmt.Println("Simple Block Chain Creation")
	var blockChain = genesis()
	// Add three transactions:
	blockChain.addBlock(generateBlock("croissants", 5, 1.2))
	blockChain.addBlock(generateBlock("pains", 2, 2.3))
	blockChain.addBlock(generateBlock("croissants", 4, 1.77))

	// Display blockchain
	for m, v := range blockChain.blocks {
		fmt.Println("Block: ", m, v)
	}

	var factory = ClientFactory{Uri: "127.0.0.1:5984", DbName: "blockchain"}
	_, err := factory.connect()
	if err != nil {
		panic(err)
	}
	factory.info()
	if !factory.isDbExists() {
		factory.createDb()
	}

	for k, b := range blockChain.blocks {
		fmt.Println(k)
		fmt.Println(b.Data)
		_, err = factory.execute().Post(&b)
		if err != nil {
			panic(err)
		}
	}
}
