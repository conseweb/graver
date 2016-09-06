package merkle

import "sync"

func NewMerkleNode() *MerkleTreeNode {
	return &MerkleTreeNode{
		MerkleTreeLeaf: MerkleTreeLeaf{
			Value:     "",
			Container: nil,
		},
		children: map[int]TreeNode{},
	}
}
func NewMerkleLeaf() *MerkleTreeLeaf {
	return &MerkleTreeLeaf{
		Value:     "",
		Container: nil,
	}
}

type TreeNode interface {
	IsLeaf() bool
	AddChild(child TreeNode)
	RemoveChild(child TreeNode)
	GetContainer() interface{}
	GetValue() string
	SetValue(v string)
	GetChildren() []TreeNode
}
type MerkleTreeLeaf struct {
	Value     string      //hash value for this node
	Container interface{} //the container which holds the read data
	Next      TreeNode    //pointer of the next leaf
	sync.Mutex
}

func (leaf *MerkleTreeLeaf) IsLeaf() bool {
	return true
}
func (leaf *MerkleTreeLeaf) AddChild(child TreeNode) {
	panic("Unsupported method")
}
func (leaf *MerkleTreeLeaf) RemoveChild(child TreeNode) {
	panic("Unsupported method")
}
func (leaf *MerkleTreeLeaf) GetChildren() []TreeNode {
	panic("Unsupported method")
}
func (leaf *MerkleTreeLeaf) GetContainer() interface{} {
	return leaf.Container
}
func (leaf *MerkleTreeLeaf) GetValue() string {
	return leaf.Value
}
func (leaf *MerkleTreeLeaf) SetValue(v string) {
	leaf.Lock()
	leaf.Value = v
	leaf.Unlock()
}

type MerkleTreeNode struct {
	MerkleTreeLeaf
	children map[int]TreeNode //the children belong to this node
}

func (node *MerkleTreeNode) IsLeaf() bool {
	return false
}
func (node *MerkleTreeNode) AddChild(child TreeNode) {
	node.Lock()
	childCount := len(node.children)
	node.children[childCount] = child
	node.Unlock()
}
func (node *MerkleTreeNode) RemoveChild(child TreeNode) {
	node.Lock()
	var target int = -1
	for key, value := range node.children {
		if value == child {
			target = key
			break
		}
	}
	if target != -1 {
		delete(node.children, target)
	}
}
func (node *MerkleTreeNode) GetChildren() []TreeNode {
	cs := []TreeNode{}
	node.Lock()
	for _, c := range node.children {
		cs = append(cs, c)
	}
	node.Unlock()
	return cs
}
