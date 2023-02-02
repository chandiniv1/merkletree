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

func AddNode(values []string) []string {
	var newnode string
	fmt.Println("enter the node you want to add")
	fmt.Scanln(&newnode)
	values = append(values, newnode)
	return values

}

func DeleteNode(Values []string) []string {
	var deletednode string
	flag := 0
	fmt.Println("enter the node you want to delete")
	fmt.Scanln(&deletednode)
	for i := 0; i < len(Values); i++ {
		if Values[i] == deletednode {
			Values = append(Values[:i], Values[i+1:]...)
			flag = 1
		}
	}
	if flag == 0 {
		fmt.Println("transaction not found")
	}
	return Values
}

func (root *MerkleNode) verify(data string) bool {
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

	if bytes.Equal(root.Value, target) {
		return true
	}
	var left, right bool
	if root.Left != nil {
		left = VerifyNode(root.Left, target)
	}
	if root.Right != nil {
		right = VerifyNode(root.Right, target)
	}
	return left || right
}

func main() {
	values := []string{"GeeksforGeeks", "A", "Computer", "Science", "Portal", "For", "Geeks"}
	//values := []string{}
	updatedvalues := AddNode(values)
	DeleteNode(updatedvalues)
	data := [][]byte{}

	if len(updatedvalues) == 0 {
		fmt.Println("NO transactions")
	} else {
		for k, v := range updatedvalues {
			_ = k
			a := []byte(v)
			data = append(data, a)
		}
		rootHash := buildMerkleTree(data)
		fmt.Println("Root hash", hex.EncodeToString(rootHash.Root.Value))
	}
	c := &MerkleNode{}
	fmt.Println("node is present/not: ", c.verify("xyz"))

}
