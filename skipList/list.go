package skipList

type Node interface{}

type List interface{
	Add(*node) int
	Remove(*node) bool
	Find(*node) int
	Contains(*node) bool
	Range(start , end int) []*node
}
