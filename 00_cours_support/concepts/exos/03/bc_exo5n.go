package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

//
func toString(structure Tx) string {
	bytesStructure, err := json.Marshal(structure)
	if err != nil {
		panic(err)
	}
	return string(bytesStructure)
}

//
func hash(obj []byte) string {
	return encodeHex(hashBytes(obj))
}

func encodeHex(obj []byte) string {
	return hex.EncodeToString(obj)
}

//
func hashBytes(obj []byte) []byte {
	h := sha256.New()
	_, err := h.Write(obj)
	if err != nil {
		panic(err)
	}
	return h.Sum(nil)
}

//
type Tree struct {
	Root       *Node
	merkleRoot []byte
	Leafs      []*Node
}

//
type Node struct {
	Parent *Node
	Left   *Node
	Right  *Node
	leaf   bool
	dup    bool
	Hash   []byte
	C      Tx
}

func (n *Node) calculateNodeHash() ([]byte, error) {
	if n.leaf {
		return n.C.CalculateHash(), nil
	}

	h := sha256.New()
	if _, err := h.Write(append(n.Left.Hash, n.Right.Hash...)); err != nil {
		return nil, err
	}

	return h.Sum(nil), nil
}

func NewTree(cs []Tx) (*Tree, error) {
	root, leafs, err := buildWithContent(cs)
	if err != nil {
		return nil, err
	}
	t := &Tree{
		Root:       root,
		merkleRoot: root.Hash,
		Leafs:      leafs,
	}
	return t, nil
}

func (m *Tree) GetMerklePath(content Tx) ([][]byte, []int64, error) {
	for _, current := range m.Leafs {
		ok, err := current.C.Equals(content)
		if err != nil {
			return nil, nil, err
		}
		var currentParent *Node
		if ok {
			currentParent = current.Parent
			var merklePath [][]byte
			var index []int64
			for currentParent != nil {
				if bytes.Equal(currentParent.Left.Hash, current.Hash) {
					merklePath = append(merklePath, currentParent.Right.Hash)
					index = append(index, 1) // right leaf
				} else {
					merklePath = append(merklePath, currentParent.Left.Hash)
					index = append(index, 0) // left leaf
				}
				current = currentParent
				currentParent = currentParent.Parent
			}
			return merklePath, index, nil
		}
	}
	return nil, nil, nil
}
func buildWithContent(cs []Tx) (*Node, []*Node, error) {
	if len(cs) == 0 {
		return nil, nil, errors.New("error: cannot construct tree with no content")
	}
	var leafs []*Node
	for _, c := range cs {
		leafs = append(leafs, &Node{
			Hash: c.CalculateHash(),
			C:    c,
			leaf: true,
		})
	}
	if len(leafs)%2 == 1 {
		duplicate := &Node{
			Hash: leafs[len(leafs)-1].Hash,
			C:    leafs[len(leafs)-1].C,
			leaf: true,
			dup:  true,
		}
		leafs = append(leafs, duplicate)
	}
	root, err := buildIntermediate(leafs)
	if err != nil {
		return nil, nil, err
	}

	return root, leafs, nil
}

func buildIntermediate(nl []*Node) (*Node, error) {
	var nodes []*Node
	for i := 0; i < len(nl); i += 2 {
		h := sha256.New()
		var left, right int = i, i + 1
		if i+1 == len(nl) {
			right = i
		}
		chash := append(nl[left].Hash, nl[right].Hash...)
		if _, err := h.Write(chash); err != nil {
			return nil, err
		}
		n := &Node{
			Left:  nl[left],
			Right: nl[right],
			Hash:  h.Sum(nil),
		}
		nodes = append(nodes, n)
		nl[left].Parent = n
		nl[right].Parent = n
		if len(nl) == 2 {
			return n, nil
		}
	}
	return buildIntermediate(nodes)
}

func (m *Tree) MerkleRoot() []byte {
	return m.merkleRoot
}

type Tx struct {
	User      string
	CreatedAt string
	Reference string
	Quantity  int
	Price     float32
}

func (d Tx) CalculateHash() []byte {
	return hashBytes([]byte(toString(d)))
}

func (d Tx) Equals(other Tx) (bool, error) {
	return d.CreatedAt == other.CreatedAt && d.User == other.User, nil
}

func NewTx(u string, todayPlus int, r string, q int, p float32) Tx {
	return Tx{User: u, CreatedAt: time.Now().Add(time.Duration(todayPlus) * time.Millisecond).Format(time.RFC3339), Reference: r, Quantity: q, Price: p}
}

func PathToString(p [][]byte) []string {
	var result []string
	for _, v := range p {
		result = append(result, hex.EncodeToString(v))
	}
	return result
}

var table = struct {
	txs []Tx
}{
	txs: []Tx{
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
		r := v.CalculateHash()
		fmt.Println(v, "	 -> ", encodeHex(r))
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
		r := v.CalculateHash()
		fmt.Println("tx -> ", v, " hash -> ", encodeHex(r), "\n  -> path: ", PathToString(p), " -> ", i)

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
