package main

import (
	"errors"
	"fmt"
	"reflect"
)
type Node1 struct {
	Value   interface{}
}

func (n *Node1) Compare(node *Node) (int, error) {
	if reflect.TypeOf(n.Value) != reflect.TypeOf(node) {
		return 0, errors.New("类型不同无法比较")
	}
	switch node.Value.(type) {
	case int,int8,int16,int32,int64,uint,uint8,uint16,uint32,uint64:
		if n.Value.(int) > node.Value.(int) {
			return 1, nil
		}
		if n.Value.(int) < node.Value.(int) {
			return 0, nil
		}
		if n.Value.(int) > node.Value.(int) {
			return -1, nil
		}
	case string:
		if n.Value.(string) > node.Value.(string) {
			return 1, nil
		}
		if n.Value.(int) < node.Value.(int) {
			return 0, nil
		}
		if n.Value.(int) > node.Value.(int) {
			return -1, nil
		}
	}
	return 0, errors.New("不支持相关类型的比较")
}

func main() {
	a := make([]int, 1)
	func(s *[]int) {
		*s = append(*s, 1)
	}(&a)
	fmt.Println(a)
	return

}


