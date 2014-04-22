package stack

type Stack struct {
	top  *Node
	size int
}

type Node struct {
	next *Node
	data interface{}
}

func (stack *Stack) Size() int {
	return stack.size
}

func (stack *Stack) Push(data interface{}) {
	top := stack.top
	node := &Node{top, data}
	stack.top = node
	stack.size++
}

func (stack *Stack) Pop() interface{} {
	if stack.size > 0 {
		top := stack.top
		next := top.next
		stack.top = next
		stack.size--
		return top.data
	}
	return nil
}

func (stack *Stack) Peek() interface{} {
	if stack.size > 0 {
		return stack.top.data
	} else {
		return nil
	}
}
