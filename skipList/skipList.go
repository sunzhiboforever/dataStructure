package skipList

import (
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"time"
)

const MAX_LEVEl = 10

//跳跃表结构
type skipList struct {
	size  int
	lists *[]*list
	count int
	rand  *rand.Rand
	w     sync.RWMutex
}

//链表结构
type list struct {
	head *node
	tail *node
}

//节点结构
type node struct {
	span  int
	score int
	data  interface{}
	next  *node
	prev  *node
	down  *node
	up    *node
}

func New() skipList {
	skipList := skipList{}
	a := make([]*list, 0)
	skipList.lists = &a
	*skipList.lists = append(*skipList.lists, &list{head: nil})
	skipList.rand = rand.New(rand.NewSource(time.Now().UnixNano()))
	skipList.w = sync.RWMutex{}
	return skipList
}

//掷硬币决定要添加到第几层
//最底层是第零层
func (s *skipList) getLevel() int {
	var level int
	for {
		//投硬币是反面   0：反面 1：正面
		if s.rand.Intn(2) == 0 {
			return level
		}
		//已经达到最大层数+1层了
		if level >= len(*s.lists) {
			return level
		}
		//如果比规定的最大层数还要大
		if level >= MAX_LEVEl {
			return level
		}
		level = level + 1
	}
	return level
}

// 获取当前元素个数
func (s *skipList) GetSize() int {
	return s.size
}

// 按照索引查找范围
// start 代表开始的索引，length代表要返回几个值
// length 为 -1 时，代表返回 start 后的所有节点
func (s *skipList) FindRange(start, length int) []interface{} {
	s.w.RLock()
	defer s.w.RUnlock()
	if length == 0 {
		return []interface{}{}
	}

	var startNode *node
	var result []interface{}
	var currentNode *node
	var currentData interface{}

	//想取出排在前几名的元素，直接从头节点开始查找就可以了
	if start == 1 {
		startNode = (*s.lists)[0].head
	} else {
		// 这里拿到的节点，可能不是最底层的节点
		startNode, _ = s.findNode(start)
		for startNode.down != nil {
			startNode = startNode.down
		}
	}
	currentNode = startNode
	dataCount := 1

	for {
		if currentNode == nil {
			break
		}
		if dataCount > length && length > 0 {
			break
		}
		fmt.Printf("当前元素: %s\n", currentNode.data.(string))
		currentData = (*currentNode).data
		result = append(result, currentData)
		currentNode = currentNode.next
		dataCount = dataCount + 1
	}
	return result
}

// 按照索引查找元素值
func (s *skipList) FindIndex(index int) interface{} {
	s.w.RLock()
	defer s.w.RUnlock()
	node, _ := s.findNode(index)
	if node == nil {
		return nil
	}
	return node.data
}

// 删除元素
func (s *skipList) RemoveIndex(index int) bool {
	s.w.Lock()
	defer s.w.Unlock()
	// 当前要删除的元素（最上层链表）
	var removeNode *node

	// 当前删除元素所在的层数
	var currentLevel int

	// 当前要删除的元素指针
	var currentNode *node

	removeNode, currentLevel = s.findNode(index)
	fmt.Printf("找到要删除的节点：%s，此节点层级:%d\n", removeNode.data.(string), currentLevel)
	// 此链表中没有找到相应的节点
	if removeNode == nil {
		return false
	}
	// 从要删除的节点往上循环维护 span
	fmt.Println("开始维护上层节点span")
	currentNode = removeNode
	for currentNode != nil{
		fmt.Printf("当前节点:%s\n", currentNode.data.(string))
		if currentNode.up != nil {
			currentNode.up.span = currentNode.up.span - 1
			currentNode = currentNode.up
		} else {
			currentNode = currentNode.next
		}
	}

	fmt.Println("开始删除节点")
	// 从要删除的节点层级往下循环删除
	currentNode = removeNode
	for currentNode != nil{
		fmt.Printf("当前节点:%s\n", currentNode.data.(string))
		//情况1：当前节点为头节点，也是尾节点，也就是这层链表就它一个元素，则当前所在链表层数必然为最高层
		if currentNode.prev == nil && currentNode.next == nil {
			//去掉最高层的链表
			list := *s.lists
			list = list[0:len(list) - 1]
			s.lists = &list

			//情况2：当前节点为头节点
		} else if currentNode.prev == nil {
			currentNode.next.span = currentNode.next.span + currentNode.span - 1
			(*s.lists)[currentLevel].head = currentNode.next
			currentNode.next.prev = nil

			//情况3：当前节点为尾节点
		} else if currentNode.next == nil {
			currentNode.prev.next = nil

			//情况4：当前节点为中间节点
		} else if currentNode.prev != nil && currentNode.next != nil {
			currentNode.next.span = currentNode.next.span + currentNode.span - 1
			currentNode.prev.next = currentNode.next
			currentNode.next.prev = currentNode.prev
		}
		currentNode.up = nil
		currentNode.prev = nil
		currentNode.next = nil
		currentNode = currentNode.down
		currentLevel--
	}
	s.size--
	return true
}

// 按照索引查找节点
// 注意：此方法返回的可能不是最下层链表的节点
func (s *skipList) findNode(index int) (*node, int) {
	// 当前作为比较的节点
	var currentNode *node
	// 当前节点的索引值
	var currentNodeIndex int
	// 当前的链表层数
	var currentLevel int

	// 初始化 比较的节点为最上层链表的头指针
	currentNode = (*s.lists)[len(*s.lists)-1].head
	// 初始化 索引值为最上层链表头节点的索引值
	currentNodeIndex = (*s.lists)[len(*s.lists)-1].head.span
	// 初始化 比较的层数为最上层链表
	currentLevel = len(*s.lists) - 1
	// 空链表判断
	if currentNode == nil {
		return nil, 0
	}

	for {
		fmt.Printf("当前节点：%s, 当前节点索引%d, 当前层级:%d\n", currentNode.data.(string), currentNodeIndex, currentLevel)
		// 情况1：要查找的索引比当前链表的头节点索引还要小
		if currentNode.prev == nil && currentNodeIndex > index {
			// 如果此时是最底层链表，说明索引值的元素并不存在
			if currentLevel <= 0 {
				return nil, 0
			} else { // 把当前要比较的节点置为下层链表的头节点
				currentNode = (*s.lists)[currentLevel-1].head
				currentNodeIndex = currentNode.span
				currentLevel = currentLevel - 1
			}
			// 情况2：要查找的索引比当前链表的尾节点索引还要大
		} else if currentNode.next == nil && currentNodeIndex < index {
			// 如果此时是最底层链表，说明索引值的元素并不存在
			if currentLevel <= 0 {
				return nil, 0
			} else { // 把当前要比较的节点置为下层链表的对应节点
				currentNode = currentNode.down
				currentLevel = currentLevel - 1
			}
			// 情况3：要查找的索引，正好比当前节点索引大，比当前节点的后一个节点的索引小
		} else if currentNodeIndex < index && currentNodeIndex + currentNode.next.span > index {
			// 把当前要比较的节点置为下层链表的对应节点
			currentNode = currentNode.down
			currentLevel = currentLevel - 1

			// 情况4：要查找的索引，比当前节点索引大，比当前节点的后一个节点的索引也大
		} else if currentNodeIndex < index {
			// 把当前要比较的节点置为当前节点的下一个节点
			currentNode = currentNode.next
			currentNodeIndex = currentNodeIndex + currentNode.span

			// 情况5：当前的节点的索引，等于要找到的索引
		} else if currentNodeIndex == index {
			return currentNode, currentLevel
		}
	}
	return currentNode, currentLevel
}

//插入元素
//如果 score 分数有重复则按照时间顺序排序
func (s *skipList) Add(data interface{}, score int) {
	s.w.Lock()
	defer s.w.Unlock()

	fmt.Printf("当前添加元素 data:%s, score: %d \n", data, score)

	// 抛硬币，确定要添加的层数
	level := s.getLevel()

	// 如果当前的层数不够，则要再加一层
	if level > len(*s.lists)-1 {
		*s.lists = append(*s.lists, &list{head: nil})
	}
	fmt.Printf("当前共有链表层数：%d\n", len(*s.lists))

	// 要插入的节点
	var insertNode *node

	// 循环链表的时候，记录当前节点
	var currentNode *node

	// 查找时，当前层找到的节点，要作为下层链表开始查找的节点
	var nextListHead *node

	// 插入时，存入上层链表的插入节点，以便调整 down 指针到下层节点
	var upperNode *node

	//循环要从最上层链表开始往下循环
	for currentLevel := len(*s.lists) - 1; currentLevel >= 0; currentLevel-- {
		currentList := (*s.lists)[currentLevel]
		fmt.Printf("当前循环第：%d 层链表\n", currentLevel)

		//判断开始查找的头节点
		if nextListHead != nil {
			currentNode = nextListHead
			fmt.Printf("来自上层节点的查找结果，当前从 %s，开始往后找\n", currentNode.data.(string))
		} else {
			currentNode = currentList.head
			fmt.Printf("没有上层节点的结果，当前从头节点，开始往后找\n")
		}
		insertNode = &node{
			score: score,
			data:  data,
			span:  0,
			next:  nil,
			prev:  nil,
			down:  nil,
			up:    nil,
		}
		isInsert := false
		// 循环当前层的链表
		for {
			// 如果节点需要插入到当前链表中
			if currentLevel <= level {
				// 说明此链表是新添加出来的，当前节点为头节点，一个节点都没有
				if currentNode == nil {
					fmt.Println("当前节点为头节")
					currentList.head = insertNode
					currentList.tail = insertNode
					isInsert = true
					// 当前节点(头节点) < 要插入的节点
				} else if currentNode.score < score && currentNode == currentList.head {
					fmt.Println("当前为头节点插入")
					currentList.head = insertNode
					insertNode.next = currentNode
					currentNode.prev = insertNode
					isInsert = true
					// 当前节点(尾节点) >= 要插入的节点
				} else if currentNode.score >= score && currentNode == currentList.tail {
					fmt.Println("当前为尾节点插入")
					currentList.tail = insertNode
					currentNode.next = insertNode
					insertNode.prev = currentNode
					nextListHead = currentNode.down
					isInsert = true
					// 当前节点 >= 要插入的节点 > 当前节点的下一个节点
				} else if currentNode.score > score && score >= currentNode.next.score {
					fmt.Println("当前为中间节点插入")
					insertNode.next = currentNode.next
					insertNode.next.prev = insertNode
					insertNode.prev = currentNode
					currentNode.next = insertNode
					nextListHead = currentNode.down
					isInsert = true
				} else {
					fmt.Println("没有找到相应节点，要跳转到下一个节点")
					currentNode = currentNode.next
					continue
				}

				// 插入成功了
				if isInsert == true {
					// 把上层节点的 down 指针指向当前节点，当前节点的 up 指针指向上层节点
					if upperNode != nil {
						upperNode.down = insertNode
						insertNode.up = upperNode
					}
					// 把当前节点作为为下一层链表的上层节点
					upperNode = insertNode
					break
				}
			}
			// 当前链表的层数大于要插入的层数，不需要插入，只执行查找操作
			if currentLevel > level {
				// 当前节点(头节点) < 要插入的节点
				if currentNode.score < score && currentNode == currentList.head {
					nextListHead = nil
					fmt.Println("要插入的节点比当前链表的头节点还要小")
					break
					// 当前节点(尾节点) >= 要插入的节点
				} else if currentNode.score >= score && currentNode == currentList.tail {
					nextListHead = currentNode.down
					fmt.Printf("当前层查找到节点: %s，作为下层链表的头节点\n", currentNode.data.(string))
					break
					// 当前节点 >= 要插入的节点 > 当前节点的下一个节点 （后一个节点为空的情况在上面判断过了，不会走到这一步，所以不会报空指针）
				} else if currentNode.score > score && score >= currentNode.next.score {
					nextListHead = currentNode.down
					fmt.Printf("当前层查找到节点: %s，作为下层链表的头节点，它的下层节点：%s\n", currentNode.data.(string), currentNode.down.data.(string))
					break
				}
				currentNode = currentNode.next
				continue
			}
		}
	}

	// 插入完成后，要完成 span 的维护
	insertNode.span = 1

	//维护当前插入的节点，以及它的上节点们的span
	currentNode = insertNode

	// 每次循环更新的是上节点的 span 和上节点的后节点的 span
	for {
		if currentNode.up == nil {
			break
		}
		// 开始水平循环，计算当前节点和前一个节点中间的跨度
		lineNode := currentNode
		span := currentNode.span
		for {
			// 上层节点为头节点
			if currentNode.up.prev == nil && lineNode.prev == nil {
				break
			}
			// 上层节点不为头节点
			if currentNode.up.prev != nil && lineNode.prev == currentNode.up.prev.down {
				break
			}
			span = span + lineNode.prev.span
			lineNode = lineNode.prev
		}
		currentNode.up.span = span

		// 更新上节点的后节点的span，最后再加上 1 是因为这个添加方法有一个新添加的元素
		if currentNode.up.next != nil {
			currentNode.up.next.span = currentNode.up.next.span - currentNode.up.span + 1
		}
		currentNode = currentNode.up
	}

	flag := false
	for {
		if flag == true {
			break
		}
		lineNode := currentNode
		for {
			if lineNode == nil {
				flag = true
				break
			}
			if lineNode.up != nil {
				lineNode.up.span += 1
				currentNode = lineNode.up
				break
			}
			lineNode = lineNode.next
		}
	}
	s.size ++
}

//实现接口，可以直接 fmt.Println()
func (s *skipList) String() string {
	builder := strings.Builder{}
	for k, l := range *s.lists {
		list := l
		level := k
		builder.WriteString(fmt.Sprintf("level: %d\n", level))
		current := list.head
		for {
			builder.WriteString(fmt.Sprintf("- %s(%d)(%d) ", current.data, current.score, current.span))
			if current.next == nil {
				break
			}
			current = current.next
		}
		builder.WriteString("\n")

	}
	return builder.String()
}
