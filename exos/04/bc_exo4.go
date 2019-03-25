package main

import (
	"encoding/hex"
	"fmt"
	"time"
)

// Block Data
type Data struct {
	//Tree *MerkleTree
}

type Transaction struct {
	User      string
	CreatedAt string
	Reference string
	Quantity  int
	Price     float32
}

func (d Transaction) CalculateHash() ([]byte, error) {
	return HashStructBytes(d)
}

func (d Transaction) Equals(other Content) (bool, error) {
	otherData := other.(Transaction)
	return d.CreatedAt == otherData.CreatedAt && d.User == otherData.User, nil
}

func NewTx(u string, todayPlus int, r string, q int, p float32) Transaction {
	return Transaction{User: u, CreatedAt: time.Now().Add(time.Duration(todayPlus) * time.Millisecond).Format(time.RFC3339), Reference: r, Quantity: q, Price: p}
}

func PathToString(p [][]byte) []string {
	var result []string
	for _, v := range p {
		result = append(result, hex.EncodeToString(v))
	}
	return result
}

var table = struct {
	txs []Content
}{
	txs: []Content{
		NewTx("A", 1, "Bread", 1, 0.8),
		NewTx("B", 2, "Croissant", 3, 1.2),
		NewTx("C", 3, "Bread", 5, 0.8),
		NewTx("D", 4, "Croissant", 3, 1.2),
	},
}

func main() {

	fmt.Println("Exo 4 - Merkle Tree")
	fmt.Println("\n--Data--\n")

	for _, v := range table.txs {
		fmt.Println(v, "	 -> ", hashStruct(v))
	}

	fmt.Println("\n--Tree--\n")

	tree, err := NewTree(table.txs)
	if err != nil {
		fmt.Println("error: unexpected error:  ", err)
	}

	fmt.Println("root: ", hash(tree.MerkleRoot()))

	//	fmt.Println("tree: ", tree)

	for _, v := range table.txs {
		p, i, _ := tree.GetMerklePath(v)
		fmt.Println(PathToString(p), " -> ", i)

	}

	// tuple {0, 1} -> {Left, Right}
	h0L, _ := tree.Leafs[0].calculateNodeHash()
	h0R, _ := tree.Leafs[1].calculateNodeHash()

	fmt.Println(hash(append(h0L, h0R...)))

	// tuple {2, 3} -> {Left, Right}
	h1L, _ := tree.Leafs[2].calculateNodeHash()
	h1R, _ := tree.Leafs[3].calculateNodeHash()

	fmt.Println(hash(append(h1L, h1R...)))

	// route
	fmt.Println(hash(append(hashBytes(append(h0R, h0L...)), hashBytes(append(h1R, h1L...))...)))
	fmt.Println(hash(append(hashBytes(append(h1L, h1R...)), hashBytes(append(h0L, h0R...))...)))

}
