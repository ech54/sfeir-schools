package main

// Block defines a node in blockchain structure.
type Block struct {
	ID        string
	LastID    string
	CreatedAt string
	Data      Data
}

/* Data structure which contains the transaction data.
type Data struct {
	Reference string
	Quantity  int
	Price     float32
}*/

// Chain structure.
type Chain struct {
	LastKey string
	Blocks  map[string]Block
}

// Enrich chain with a add block feature

func (chain *Chain) addBlock(b Block) {
	b.LastID = chain.last().ID
	b.ID = hashStruct(b)
	chain.LastKey = b.ID
	chain.Blocks[b.ID] = b
}

/**
 * Enrich chain with a last block function.
 */
func (chain *Chain) last() Block {
	return chain.Blocks[chain.LastKey]
}

// Create the default chain structure.
func genesis() Chain {
	return Chain{
		LastKey: hash([]byte{0}),
		Blocks:  make(map[string]Block),
	}
}

/* Create a new default block with data.
func generateBlock(reference string, quantity int, price float32) Block {
	return Block{
		LastID: time.Now().Format(time.RFC3339), CreatedAt: time.Now().Format(time.RFC3339),
		Data: Data{Reference: reference, Quantity: quantity, Price: price},
	}
}
*/
