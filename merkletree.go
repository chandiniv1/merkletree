package main

import (
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
}

func buildMerkleNode(left *MerkleNode, right *MerkleNode, value []byte) *MerkleNode {

	newNode := MerkleNode{}

	if left == nil && right == nil {
		hash := sha256.Sum256(value)
		newNode.Value = hash[0:]
	} else {
		prevHashVal := append(left.Value, right.Value...)
		hash := sha256.Sum256(prevHashVal)
		newNode.Value = hash[0:]

	}

	newNode.Left = left
	newNode.Right = right

	return &newNode

}

func buildMerkleTree(values [][]byte) *MerkleTree {
	var newNodes []MerkleNode

	for _, v := range values {
		newNode := buildMerkleNode(nil, nil, v)
		newNodes = append(newNodes, *newNode)
	}

	for len(newNodes) > 1 {
		//checking whether the length is even or not
		if len(newNodes)%2 == 1 {
			newNodes = append(newNodes, newNodes[len(newNodes)-1])
		}

		//here we are changing every val in value into merkle node
		var parentHashes []MerkleNode
		for i := 0; i < len(newNodes); i += 2 {
			newNode := buildMerkleNode(&newNodes[i], &newNodes[i+1], nil)
			parentHashes = append(parentHashes, *newNode)
		}
		newNodes = parentHashes
	}

	//fmt.Println("new node", newNodes)

	tree := MerkleTree{&newNodes[0]}

	return &tree
}

func main() {

	values := []string{"GeeksforGeeks", "A", "Computer", "Science", "Portal", "For", "Geeks"}
	//values := []string{}
	data := [][]byte{}

	if len(values) == 0 {
		fmt.Println("NO transactions")
	} else {
		for k, v := range values {
			_ = k
			a := []byte(v)
			data = append(data, a)
		}
		rootHash := buildMerkleTree(data)
		fmt.Println("Root hash", hex.EncodeToString(rootHash.Root.Value))
	}

}
