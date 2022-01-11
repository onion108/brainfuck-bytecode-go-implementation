package main

type VMStack struct {
	content   []int
	stackSize int
}

func VMStackInit() *VMStack {
	return &VMStack{[]int{}, 0}
}

func (stack *VMStack) Push(data int) {
	if stack == nil {
		return
	}
	stack.content = append(stack.content, data)
	stack.stackSize++
}

func (stack *VMStack) Pop() int {
	if stack == nil {
		return -114514
	}
	r := stack.content[len(stack.content)-1]
	stack.content = stack.content[:len(stack.content)-1]
	stack.stackSize--
	return r
}

func (stack *VMStack) Empty() bool {
	if stack == nil {
		return true
	}
	return stack.stackSize == 0
}

func (stack *VMStack) Seek() int {
	return stack.content[len(stack.content)-1]
}
