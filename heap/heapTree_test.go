package heap

import (
	"errors"
	"testing"
)

//声明一个节点类型，实现 Node 接口
type NodeInt int
func (n NodeInt) Compare (node Node) (int, error) {
	if nodeTmp, ok := node.(NodeInt); ok {
		return int(n) - int(nodeTmp), nil
	}
	return 0, errors.New("数据类型不正确")
}

func TestHeapTree(t *testing.T) {
	var heap heapTree
	heap = New()
	//判断元素是否都添加进去了
	insert := []int{34, 34, 13, 42, 125, 32, 7, 2, 1, 32, 32, 32}
	for _, i := range insert {
		node := NodeInt(i)
		heap.Add(node)
	}
	//heap.Debug()
	if heap.GetSize() != len(insert) {
		t.Error("当前堆，添加元素算法错误")
	}
	//判断元素是否都取出来了
	var node Node
	var nodes []Node
	nodes = make([]Node, 0)
	for {
		node = heap.GetMax()
		//heap.Debug()
		if node == nil {
			break
		}
		nodes = append(nodes, node)
	}
	if heap.GetSize() != 0 {
		t.Error("当前堆，取出最大元素算法错误")
	}
	//判断是否满足堆的条件，根节点永远最大
	for i := 0; i < len(nodes)-1; i++ {
		if r,err := nodes[i].Compare(nodes[i+1]);err == nil && r < 0 {
			t.Error("当前堆，父节点大于子节点")
		}
	}
}
