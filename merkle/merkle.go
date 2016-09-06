package merkle

import (
	"os"
	"path"
	"time"

	"github.com/conseweb/stonemason/task"
)

type FileContainer struct {
	Path      string
	Timestamp time.Time //the last modification time of this file
}

/** build a merkle tree based a directory of a file system */
type FileSystemMerkleTree struct {
	root      TreeNode        //the root of this merkle tree
	leaf      *MerkleTreeLeaf //the first leaf of this merkle tree
	executor  *task.TaskExecutor
	exclusion []string //exclusion list
}

func NewMerkleTreeFromPath(path string, exclusion []string) (*FileSystemMerkleTree, error) {
	fi, err := os.Stat(path)
	if os.IsNotExist(err) {
		return nil, err
	}
	fc := &FileContainer{
		Path:      path,
		Timestamp: fi.ModTime(),
	}
	var node TreeNode = nil
	if fi.IsDir() {
		merkleNode := NewMerkleNode()
		merkleNode.Container = fc
		node = merkleNode
	} else {
		node = NewMerkleLeaf()
	}
	return &FileSystemMerkleTree{
		root:      node,
		leaf:      nil,
		executor:  task.NewTaskExecutor(32),
		exclusion: exclusion,
	}, nil
}

//build the merkle tree from root
func (tree *FileSystemMerkleTree) build(node TreeNode) {
	if node.IsLeaf() {
		tmpLeaf := node.(*MerkleTreeLeaf)
		tmpLeaf.Lock()
		tmpLeaf.Next = tree.leaf
		tree.leaf = tmpLeaf
		tmpLeaf.Unlock()
	} else {
		tree.buildChildren(node)
		merkleNode, _ := node.(*MerkleTreeNode)
		for _, c := range merkleNode.children {
			if c == nil {
				continue
			}
			tree.build(c)
		}
	}
	tree.buildValue(node)
}
func (tree *FileSystemMerkleTree) buildChildren(node TreeNode) {
	fc, ok := node.GetContainer().(*FileContainer)
	if !ok || fc == nil {
		return
	}
	file, err := os.Open(fc.Path)
	if err != nil {
		return
	}
	subfiles, err := file.Readdir(-1)
	if err != nil {
		return
	}
	for _, sf := range subfiles {
		tmpFC := &FileContainer{
			Path:      path.Join(fc.Path, sf.Name()),
			Timestamp: sf.ModTime(),
		}
		if sf.IsDir() {
			merkleNode := NewMerkleNode()
			merkleNode.Container = tmpFC
			node.AddChild(merkleNode)
		} else {
			merkleLeaf := NewMerkleLeaf()
			merkleLeaf.Container = tmpFC
			node.AddChild(merkleLeaf)
		}
	}
}
func (tree *FileSystemMerkleTree) buildValue(node TreeNode) {
	t := &MerkleValueCalTask{
		node: node,
	}
	waiter := tree.executor.Submit(t)
	waiter.Done()
}
