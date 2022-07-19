package merkletree

import (
	"bytes"
	"crypto/md5"
	"testing"

	"github.com/stretchr/testify/assert"
)

type myBlock string

func newMyBlock(s string) myBlock {
	return myBlock(s)
}

func (s myBlock) Hash() ([]byte, error) {
	hash := md5.Sum([]byte(s))
	return hash[:], nil
}
func (s myBlock) Equals(block Block) (bool, error) {
	hash1, err := block.Hash()
	if err != nil {
		return false, err
	}
	hash2, err := s.Hash()
	if err != nil {
		return false, err
	}
	return bytes.Equal(hash1, hash2), nil
}

func TestMerkleTree(t *testing.T) {
	blocks := []Block{
		newMyBlock("1"),
		newMyBlock("2"),
		newMyBlock("3"),
		newMyBlock("4"),
		newMyBlock("5"),
		// newMyBlock("6"),
	}
	mt, err := New(blocks)
	assert.Nil(t, err)
	mt.Print()
	path, isLefts, err := mt.FindMerklePath(newMyBlock("4"))
	assert.Nil(t, err)
	t.Logf("path: %+v\n", path)
	t.Logf("isLefts: %+v\n", isLefts)

	mt.Print2()
}
