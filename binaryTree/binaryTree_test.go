package binaryTree

import (
	"math/rand"
	"testing"
)

var tree BinaryTree

func TestBinaryTree_Add(t *testing.T) {
	for i:=1;i<=100;i++ {
		r := rand.Intn(1000)
		node := &Node{Value:r}
		tree.Add(node)
	}
	go check(tree.root, t)
}

func check(root *Node, t *testing.T) {
	if root == nil {
		return
	}
	if root.right == nil && root.left == nil {
		return
	}
	if root.right == nil {
		if root.Compare(root.left) < 0 {
			t.Error("不符合二分搜索树，左节点必须小于根节点")
		}
		go check(root.left, t)
		return
	}
	if root.left == nil {
		if root.Compare(root.right) > 0 {
			t.Error("不符合二分搜索树，右节点必须小于根节点")
		}
		go check(root.right, t)
		return
	}
	if  root.Compare(root.left) < 0 || root.Compare(root.right) > 0 {
		t.Error("不符合二分搜索树，左节点必须小于右节点")
		return
	}
}

func TestBinaryTree_Middle(t *testing.T) {
	read := make([]*Node, 0)
	tree.Middle(&read)
	for i:=0;i<=len(read)-2;i++ {
		if read[i].Value > read[i+1].Value {
			t.Error("中序排序没有满足排序条件")
		}
	}
}