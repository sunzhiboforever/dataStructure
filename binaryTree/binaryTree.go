package binaryTree

import (
	"errors"
	"fmt"
)

//二分搜索树
//满足条件：
//	首先是一颗二叉树
//	对于树中的每个节点，它的右节点一定比自己大，左节点一定比自己小
//	树中的每个节点必须是可以比较的
//特点：
//	具有顺序性，能快速拿到最大值最小值，也能指定一个元素快速拿到它的前驱和后继
type BinaryTree struct {
	//树的跟节点指针
	root *Node
	//树中的元素个数
	size int
}

type Node struct {
	right *Node
	left  *Node
	Value int
}

//node节点必须是可以比较的
//@TODO 实现多种数据类型
func (t *Node) Compare(arm *Node) int {
	if t.Value == arm.Value {
		return 0
	}
	if t.Value > arm.Value {
		return 1
	}
	return -1
}

//获取当前树元素个数
func (t *BinaryTree) GetSize() int {
	return t.size
}

//往树中添加元素，如果元素已存在则不做操作
func (t *BinaryTree) Add(node *Node) {
	t.root = t.add(t.root, node)
	t.size++
}

//往以root节点为根节点的树中，添加元素node
//返回添加node之后的树的跟节点
func (t *BinaryTree) add(root, node *Node) *Node {
	//如果当前树的跟节点为空节点（空节点也可以理解为一棵树的跟节点），说明这就是应该插入的位置
	if root == nil {
		return node
	}
	//如果插入的节点比当前树的根节点小，继续往左子树搜索，返回的值为添加后的子树的跟节点
	if root.Compare(node) > 0 {
		root.left = t.add(root.left, node)
		return root
	}
	//如果插入的节点比当前树的根节点大，继续往右子树搜索，返回的值为添加后的子树的跟节点
	if root.Compare(node) < 0 {
		root.right = t.add(root.right, node)
		return root
	}
	//如果节点已经存在，则不做任何操作
	if root.Compare(node) == 0 {
		return nil
	}
	return root
}

//广度优先遍历
func (t *BinaryTree) Foreach(read []int) ([]int){
	if t.size == 0 {
		return read
	}
	//这里队列仍然借用slice简陋的实现一下
	list := make([]*Node, 0)
	//push
	list = append(list, t.root)
	var node *Node
	for {
		if len(list) == 0 {
			break
		}
		//pop
		node = list[0]
		list = list[1:]
		read = append(read, node.Value)
		if node.left != nil {
			list = append(list, node.left)
		}
		if node.right != nil {
			list = append(list, node.right)
		}
	}
	return read
}

//深度优先遍历-利用栈来实现非递归遍历
func (t *BinaryTree) FrontStack() {
	if t.size == 0 {
		return
	}
	//利用栈实现，这里用slice临时代替一下，这里len分配0，要不然都初始化为nil了
	stack := make([]*Node, 0)
	stack = append(stack, t.root)
	var node *Node
	for {
		//栈空间为空则返回
		if len(stack) == 0 {
			return
		}
		//取出栈内的最后一个元素
		node = stack[len(stack)-1]
		//重置空间（从头开始，截取到刚才那个出栈的元素那里）
		stack = stack[:len(stack)-1]

		fmt.Printf(" %d", node.Value)

		//栈是先进后出，所以要先放右节点，再放左节点，确定左节点可以优先访问
		if node.right != nil {
			stack = append(stack, node.right)
		}
		if node.left != nil {
			stack = append(stack, node.left)
		}
	}
}

//深度优先遍历-前序遍历
//第一次访问到当前节点的时候打印值
func (t *BinaryTree) Front() {
	t.front(t.root)
}

func (t *BinaryTree) front(root *Node) {
	if root == nil {
		return
	}
	fmt.Printf(" %d", root.Value)
	t.front(root.left)
	t.front(root.right)
}

//深度优先遍历-中序遍历
//按照二分搜索树的顺序遍历
//第二次访问到当前节点的时候存储到read里面
func (t *BinaryTree) Middle(read *[]*Node) {
	t.middle(t.root, read)
}

func (t *BinaryTree) middle(root *Node, read *[]*Node) {
	if root == nil {
		return
	}
	t.middle(root.left, read)
	*read = append(*read, root)
	t.middle(root.right, read)
}

//深度优先遍历-后续遍历
//第三次访问当前节点的时候打印值
func (t *BinaryTree) Back() {
	t.back(t.root)
}

func (t *BinaryTree) back(root *Node) {
	if root == nil {
		return
	}
	t.back(root.left)
	t.back(root.right)
	fmt.Printf(" %d", root.Value)
}

//删除一个给定的节点node
func (t *BinaryTree) Remove(node *Node) {
	t.root = t.remove(t.root, node)
}

//在以root为跟节点的树中，删除node节点，并返回删除完成后的树的跟节点
func (t *BinaryTree) remove(root, node *Node) *Node {
	//搜索到跟节点为nil的树时，都找到节点node，说明节点不存在，返回的nil不会有任何影响，因为本来就是nil
	if root == nil {
		return nil
	}

	//如果插入的节点比当前树的根节点小，继续往左子树搜索，返回的值为删除后的子树的跟节点
	if root.Compare(node) > 0 {
		root.left = t.remove(root.left, node)
		return root
	}
	//如果插入的节点比当前树的根节点大，继续往右子树搜索，返回的值为删除后的子树的跟节点
	if root.Compare(node) < 0 {
		root.right = t.remove(root.right, node)
		return root
	}

	//找到了待删除的节点
	if root.Compare(node) == 0 {
		//一方子树为nil的时候，只需要把另外一方子树暂存并返回即可（兼容叶子节点的情况）
		if root.right == nil {
			left := root.left
			root.left = nil
			t.size--
			return left
		}
		if root.left == nil {
			right := root.right
			root.right = nil
			t.size--
			return right
		}

		//两方子树都不为nil的时候，应该取后继（右子树中的最小值）或前驱（左子树中的最大值）
		//来替代当前删除的节点，并把当前的左右子树套接在当前替代的节点上
		//最后完成删除后的节点为当前树的根节点返回回去
		maxNode := t.getMax(root.left)
		maxNode.right = root.right
		maxNode.left = t.removeMax(root.left)
		//被删除的节点与当前树解绑
		root.right = nil
		root.left = nil
		return maxNode
	}
	return root
}

//删除树中最大的节点
func (t *BinaryTree) RemoveMax() (*Node, error) {
	var max *Node
	var err error
	if max, err = t.GetMax(); err != nil {
		return nil, err
	}
	t.root = t.removeMax(t.root)
	return max, nil
}

//递归的方式查找并删除以root为根节点的树中最大的节点
//其实删除最大的节点，和单链表一样，一直找到右子树为nil的节点就可以了
//这里完全可以用while循环来做
func (t *BinaryTree) removeMax(root *Node) *Node {
	if root.right == nil {
		left := root.left
		root.left = nil
		t.size--
		return left
	}
	root.right = t.removeMax(root.right)
	return root
}

//获取树中最大的节点
func (t *BinaryTree) GetMax() (*Node, error) {
	if t.size == 0 {
		return nil, errors.New("当前树为空树")
	}
	return t.getMax(t.root), nil
}

//递归的方式查找以root为根的树的最大节点
//和删除最大节点一样，单链表操作，完全可以用while循环实现
func (t *BinaryTree) getMax(root *Node) *Node {
	if root.right == nil {
		return root
	}
	return t.getMax(root.right)
}
