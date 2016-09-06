package merkle

import (
	"crypto/sha1"
	"fmt"
	"io"
	"os"
)

type MerkleValueCalTask struct {
	node TreeNode
}

func (t *MerkleValueCalTask) Run() {

	if t.node.IsLeaf() {
		t.calLeafValue()
	} else {
		t.calNodeValue()
	}
}

func (t *MerkleValueCalTask) calLeafValue() {
	fc, ok := t.node.GetContainer().(*FileContainer)
	if !ok || fc == nil {
		return
	}
	file, err := os.Open(fc.Path)
	defer file.Close()
	if err != nil {
		return
	}
	h := sha1.New()
	io.Copy(h, file)
	sha1Byte := h.Sum(nil)
	sha1Value := fmt.Sprintf("%x", sha1Byte)
	t.node.SetValue(sha1Value)
}

func (t *MerkleValueCalTask) calNodeValue() {
	children := t.node.GetChildren()
	h := sha1.New()
	for _, c := range children {
		io.WriteString(h, c.GetValue())
	}
	sha1Byte := h.Sum(nil)
	sha1Value := fmt.Sprintf("%x", sha1Byte)
	t.node.SetValue(sha1Value)
}
