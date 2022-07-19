package merkletree

import (
	"bytes"
	"crypto/md5"
	"errors"
	"fmt"
)

type Block interface {
	Hash() ([]byte, error)
	Equals(block Block) (bool, error)
}

type Node struct {
	Parent *Node
	Left   *Node
	Right  *Node
	Hash   []byte

	block Block
}

type MerkleTree struct {
	root  *Node
	leafs []*Node
}

func New(blocks []Block) (*MerkleTree, error) {
	root, leafs, err := buildWithBlocks(blocks)
	if err != nil {
		return nil, err
	}
	return &MerkleTree{
		root:  root,
		leafs: leafs,
	}, nil
}

func (mt *MerkleTree) Print() {
	fmt.Printf("merkle tree root: %x", mt.root.Hash)
	fmt.Printf("leafs hash: %+v", mt.leafs)
}

func (mt *MerkleTree) Print2() {
	queue := []*Node{mt.root}
	for len(queue) > 0 {
		size := len(queue)
		for i := 0; i < size; i++ {
			node := queue[i]
			fmt.Printf("hash(%x)   ", node.Hash)
			if node.Left != nil {
				queue = append(queue, node.Left)
			}
			if node.Right != nil {
				queue = append(queue, node.Right)
			}
		}
		queue = queue[size:]
		fmt.Print("\n")
	}
}

func (mt *MerkleTree) FindMerklePath(block Block) ([][]byte, []bool, error) {
	for _, current := range mt.leafs {
		ok, err := current.block.Equals(block)
		if err != nil {
			return nil, nil, err
		}
		if ok {
			var (
				merklePaths [][]byte
				isLeftNodes []bool
				parent      = current.Parent
			)
			for parent != nil {
				isLeftNodes = append(isLeftNodes, bytes.Equal(parent.Left.Hash, current.Hash))
				merklePaths = append(merklePaths, current.Hash)
				current = parent
				parent = current.Parent
			}

			return merklePaths, isLeftNodes, nil
		}
	}
	return nil, nil, errors.New("error: block is not in merkle")
}

func buildWithBlocks(blocks []Block) (*Node, []*Node, error) {
	if len(blocks) == 0 {
		return nil, nil, errors.New("error: cannot build merkle tree with no blocks")
	}

	leafs := make([]*Node, len(blocks))
	for i, block := range blocks {
		hash, err := block.Hash()
		if err != nil {
			return nil, nil, err
		}
		leafs[i] = &Node{
			Hash:  hash,
			block: block,
		}
	}

	root, err := buildNodes(leafs)
	if err != nil {
		return nil, nil, err
	}
	return root, leafs, nil
}

func buildNodes(ns []*Node) (*Node, error) {
	var nodes []*Node

	n := len(ns)
	if n == 1 {
		return ns[0], nil
	}
	isEven := n%2 == 0
	var size int
	if isEven {
		size = n
	} else {
		size = n - 1
	}

	for i := 0; i < size; i += 2 {
		left, right := i, i+1

		data := append(ns[left].Hash, ns[right].Hash...)
		hash := md5.Sum(data)

		node := &Node{
			Hash:  hash[:],
			Left:  ns[left],
			Right: ns[right],
		}

		ns[left].Parent = node
		ns[right].Parent = node

		nodes = append(nodes, node)
	}
	if !isEven {
		index := n - 1
		hash := md5.Sum(ns[index].Hash)
		node := &Node{
			Hash: hash[:],
			Left: ns[index],
		}

		ns[index].Parent = node
		nodes = append(nodes, node)
	}

	return buildNodes(nodes)
}
