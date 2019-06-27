package heap

import (
	"errors"
	"fmt"
)

type Node interface {
	//如果目标节点比自己小，则返回正数；比自己大则返回负数；相等则返回 0
	//如果比较的目标和自己的数据类型不同则报错
	Compare(Node) (int, error)
}

//二叉堆
//			  65
//		    /   \
//		  41     30
//		 /  \    / \
//      28   16 20  13
//     / \   /
//    19 17 15
//	完全二叉树
// 		叶子节点之间的深度差只能是0或者1，叶子节点只能在最深的那层和倒数第一层
//		最后一排的叶子节点，是从左到右排列的，也就是添加的时候是一行一行添加的
//	堆
// 		首先是一颗完全二叉树
//		每个节点都比它的子节点大，最大的就是跟节点
//		由于是完全二叉树，所以可以用数组来存储
//		这里还是每个节点必须是可以被比较的
type heapTree struct {
	heap []Node
}

func New() heapTree {
	tree := heapTree{}
	tree.heap = make([]Node, 0)
	return tree
}

//获取堆中的元素数量
func (h *heapTree) GetSize() int {
	return len(h.heap)
}

//把一个数组变成堆
//如果循环一个个执行Add()方法，是 O(N*logN)
//具体方法：
//	这里是整个数组从后往前找，找第一个非叶子节点的索引值，
//这个索引以后的节点只是叶子节点，是不需要 siftDown() 操作的
//所以省了一半的节点操作，综合下来时间复杂度是 O(n)
func (h *heapTree) Heapify(s []Node) *heapTree {
	tree := heapTree{}
	copy(tree.heap, s)
	//找到第一个不是叶子节点的节点，也就是整个数组最后一个节点的父节点
	downIndex, err := h.getParent(len(tree.heap) - 1)
	if err != nil {
		return &tree
	}
	for i := 0; i <= downIndex; i++ {
		h.siftDown(i)
	}
	return &tree
}

//添加元素
func (h *heapTree) Add(node Node) {
	h.heap = append(h.heap, node)
	h.siftUp(len(h.heap) - 1)
}

//返回堆中的最大元素，其实就是根节点
func (h *heapTree) GetMax() Node {
	if len(h.heap) > 0 {
		root := h.heap[0]
		//把最后一个节点拿到当前缺失的跟节点上
		h.heap[0] = h.heap[len(h.heap)-1]
		//删除最后一个元素
		h.heap = h.heap[:len(h.heap)-1]
		h.siftDown(0)
		return root
	}
	return nil
}

//添加一个node，并且返回这个堆的最大元素
//如果先调用 Add 再调用 GetMax 就是两次 O(logN) 的操作
//所以单独封装一下，变成一次 O(logN) 操作
func (h *heapTree) Replace(node Node) Node {
	if len(h.heap) <= 0 {
		return node
	}
	//先取出跟节点，也就是最大的元素
	max := h.heap[0]
	//把当前要插入的元素当作当前的跟节点
	h.heap[0] = node
	//下沉操作
	h.siftDown(0)
	return max
}

//下沉操作
func (h *heapTree) siftDown(index int) error {
	var err error
	var compareResult int
	for {
		right := h.getRight(index)
		left := h.getLeft(index)
		//右子节点已经不存在了，说明已经下沉到底了
		if right > len(h.heap)-1 {
			return nil
		}
		//先取出两个子节点
		compareResult, err = h.heap[right].Compare(h.heap[left])
		//如果传入的节点类型不正确
		if err != nil {
			return err
		}
		//取两个子节点中最大的
		var maxIndex int
		if compareResult > 0 {
			maxIndex = right
		} else {
			maxIndex = left
		}
		//如果两个子节点中最大的节点比父节点大，则调换位置
		compareResult, err = h.heap[maxIndex].Compare(h.heap[index])
		if err != nil {
			return err
		}
		if compareResult > 0 {
			h.heap[index], h.heap[maxIndex] = h.heap[maxIndex], h.heap[index]
		}
		index = maxIndex
	}
}

//根据一个子节点的索引，查询父节点的索引
func (h *heapTree) getParent(index int) (int, error) {
	//这里只能利用 err 来表示，因为 0 被用来代表了根节点
	if index <= 0 {
		return 0, errors.New("当前节点为根节点")
	}
	//防止被除数为 0
	if index == 1 || index == 2 {
		return 0, nil
	}
	return (index - 1) / 2, nil
}

//查询一个给定节点索引的左子节点索引
func (h *heapTree) getLeft(index int) int {
	return index*2 + 1
}

//查询一个给定节点索引的左子节点索引
func (h *heapTree) getRight(index int) int {
	return index*2 + 2
}

//上浮操作
func (h *heapTree) siftUp(index int) error {
	var compareResult int
	var err error
	i := index
	for {
		//已经到达跟节点
		if i <= 0 {
			return err
		}
		//如果父节点大于当前节点，说明已经上浮到此就足够了
		parentIndex, _ := h.getParent(i)
		compareResult, err = h.heap[i].Compare(h.heap[parentIndex])
		if err != nil {
			return err
		}
		if compareResult < 0 {
			return nil
		}
		h.heap[i], h.heap[parentIndex] = h.heap[parentIndex], h.heap[i]
		i = parentIndex
	}
}

func (h *heapTree) Debug() []Node {
	for _, v := range h.heap {
		fmt.Printf(" %d", v)
	}
	fmt.Println()
	return h.heap
}
