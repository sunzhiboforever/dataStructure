package skipList

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

var skipListTest skipList

func init() {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	skipListTest = New()
}

func TestSkipList_getLevel(t *testing.T) {
	return
	skipList := New()
	level := skipList.getLevel()
	fmt.Println(level)
	return
	if level > MAX_LEVEl {
		t.Error("跳跃表层数比设定的最大层数还要大")
	}
}

func TestSkipList_Add(t *testing.T) {
	skipListTest.Add("test1", 99)
	fmt.Println("添加完成！")
	fmt.Println(skipListTest.String())
	fmt.Println()


	skipListTest.Add("test2", 64)
	fmt.Println("添加完成！")
	fmt.Println(skipListTest.String())
	fmt.Println()

	skipListTest.Add("test3", 78)
	fmt.Println("添加完成！")
	fmt.Println(skipListTest.String())
	fmt.Println()

	skipListTest.Add("test4", 33)
	fmt.Println("添加完成！")
	fmt.Println(skipListTest.String())
	fmt.Println()

	skipListTest.Add("test5", 45)
	fmt.Println("添加完成！")
	fmt.Println(skipListTest.String())
	fmt.Println()

	skipListTest.Add("test6", 54)
	fmt.Println("添加完成！")
	fmt.Println(skipListTest.String())
}

func TestSkipList_RemoveIndex(t *testing.T) {
	skipListTest.RemoveIndex(5)
	fmt.Println("删除第 5 个节点")
	fmt.Println(skipListTest.String())
}

func TestSkipList_FindIndex(t *testing.T) {
	var data interface{}

	data = skipListTest.FindIndex(0)
	fmt.Printf("查找索引为0的元素值:%q\n", data)
	fmt.Println("----------------------------------------")

	data = skipListTest.FindIndex(1)
	fmt.Printf("查找索引为1的元素值:%q\n", data)
	fmt.Println("----------------------------------------")

	data = skipListTest.FindIndex(2)
	fmt.Printf("查找索引为2的元素值:%q\n", data)
	fmt.Println("----------------------------------------")

	data = skipListTest.FindIndex(3)
	fmt.Printf("查找索引为3的元素值:%q\n", data)
	fmt.Println("----------------------------------------")

	data = skipListTest.FindIndex(4)
	fmt.Printf("查找索引为4的元素值:%q\n", data)
	fmt.Println("----------------------------------------")

	data = skipListTest.FindIndex(5)
	fmt.Printf("查找索引为5的元素值:%q\n", data)
	fmt.Println("----------------------------------------")

	data = skipListTest.FindIndex(6)
	fmt.Printf("查找索引为6的元素值:%q\n", data)
	fmt.Println("----------------------------------------")

	data = skipListTest.FindIndex(7)
	fmt.Printf("查找索引为7的元素值:%q\n", data)
}

func TestSkipList_FindRange(t *testing.T) {
	return
	var data []interface{}
	data = skipListTest.FindRange(2, -1)
	for _, v := range data {
		fmt.Printf(" %s\n", v.(string))
	}
}
