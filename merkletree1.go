package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

type MerkleTree struct {
	Root *MerkleNode
}

type MerkleNode struct {
	Left  *MerkleNode
	Right *MerkleNode
	Value []byte
	Data  []byte
}

func buildMerkleNode(left *MerkleNode, right *MerkleNode, value []byte) *MerkleNode {

	newNode := MerkleNode{}

	// If node is a leafnode
	if left == nil && right == nil {
		hash := sha256.Sum256(value)
		newNode.Value = hash[0:]
	} else { // If node is not a leaf node
		prevHashVal := append(left.Value, right.Value...)
		hash := sha256.Sum256(prevHashVal)
		newNode.Value = hash[0:]

	}

	// Update left, right nodes and data
	newNode.Left = left
	newNode.Right = right
	newNode.Data = value

	return &newNode
}

func buildMerkleTree(values [][]byte) *MerkleTree {
	var newNodes []MerkleNode
	// Create MerkleNode for each transaction
	for _, v := range values {
		newNode := buildMerkleNode(nil, nil, v)
		newNodes = append(newNodes, *newNode)
	}

	for len(newNodes) > 1 {
		// If odd number of transactions, duplicate last value
		if len(newNodes)%2 == 1 {
			newNodes = append(newNodes, newNodes[len(newNodes)-1])
		}

		// Find root hashes for each level
		var parentHashes []MerkleNode
		for i := 0; i < len(newNodes); i += 2 {
			newNode := buildMerkleNode(&newNodes[i], &newNodes[i+1], nil)
			parentHashes = append(parentHashes, *newNode)
		}
		newNodes = parentHashes
	}

	// First node is the root node of merkle tree
	tree := MerkleTree{&newNodes[0]}

	return &tree
}

func addNode(mTree *MerkleTree, data []byte) *MerkleTree {
	newData := [][]byte{data}
	newTree := buildMerkleTree(newData)

	prevLevel := []*MerkleNode{mTree.Root, newTree.Root}

	for len(prevLevel) > 1 {
		newLevel := []*MerkleNode{}
		// Calculate news hashes for each level
		for i := 0; i < len(prevLevel); i += 2 {
			node := buildMerkleNode(prevLevel[i], prevLevel[i+1], nil)
			newLevel = append(newLevel, node)
		}

		prevLevel = newLevel
	}
	// First node is the root node of Merkle Tree
	mTree.Root = prevLevel[0]
	return mTree
}

func deleteNode(mTree *MerkleTree, data []byte) *MerkleTree {
	// get all the data in the Merkle Tree
	allData := mTree.getAllData()

	// remove the data to be deleted
	for i, d := range allData {
		if string(d) == string(data) {
			allData = append(allData[:i], allData[i+1:]...)
			break
		}
	}

	// reconstructing Merkle Tree tree
	mTree = buildMerkleTree(allData)
	return mTree
}

func (mTree *MerkleTree) getAllData() [][]byte {
	allData := [][]byte{}
	getAllTransactions(mTree.Root, &allData)
	return allData
}

func getAllTransactions(node *MerkleNode, allData *[][]byte) {
	// data of node
	if node.Left == nil && node.Right == nil {
		*allData = append(*allData, node.Value)
		return
	}

	// Left subtree
	getAllTransactions(node.Left, allData)
	// Right subtree
	getAllTransactions(node.Right, allData)
}

func verify(root *MerkleNode, data string) bool {
	var hash []byte
	bytedata := []byte(data)
	hash32 := sha256.Sum256(bytedata)
	hash = append(hash, hash32[0:]...)
	return VerifyNode(root, hash)
}

func VerifyNode(root *MerkleNode, target []byte) bool {
	if root == nil {
		return false
	}

	// Return true if root value matches target
	if bytes.Equal(root.Value, target) {
		return true
	}

	var left, right bool
	// If left mode exist
	if root.Left != nil {
		left = VerifyNode(root.Left, target)
	}
	// If right node exist
	if root.Right != nil {
		right = VerifyNode(root.Right, target)
	}
	return left || right
}

func main() {
	transactions := []string{"AB", "CDEF", "G", "HIJKL", "MNO", "PQRSTU", "VW"}

	data := [][]byte{}

	if len(transactions) == 0 {
		fmt.Println("NO transactions")
	} else {
		for i := 0; i < len(transactions); i++ {
			data = append(data, []byte(transactions[i]))
		}

		merkleRoot := buildMerkleTree(data)
		fmt.Println("Root hash", hex.EncodeToString(merkleRoot.Root.Value))

		if verify(merkleRoot.Root, "AB") {
			fmt.Println("Found")
		} else {
			fmt.Println("Not Found")
		}

		merkleRoot = addNode(merkleRoot, []byte("XYZ"))
		fmt.Println("Root hash", hex.EncodeToString(merkleRoot.Root.Value))

		if verify(merkleRoot.Root, "XYZ") {
			fmt.Println("Found")
		} else {
			fmt.Println("Not Found")
		}

		merkleRoot = deleteNode(merkleRoot, []byte("XYZ"))
		fmt.Println("Root hash", hex.EncodeToString(merkleRoot.Root.Value))
		if verify(merkleRoot.Root, "XYZ") {
			fmt.Println("Found")
		} else {
			fmt.Println("Not Found")
		}
	}

}
