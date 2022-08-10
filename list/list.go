package list

// 定义基本的数据结构
// redis中list 是双链表格式的
//链表节点
type ListNode struct {
	pre   *ListNode
	next  *ListNode
	value string // 这边后期考虑传入interface{} 暂时使用字符串表示
}

// 定义链表结构
type List struct {
	head *ListNode
	tail *ListNode
	len  int
}

//创建节点
func NewListNode(value string) *ListNode {
	return &ListNode{
		value: value,
	}
}

//获取当前节点的前一个节点
func (n *ListNode) Prev() *ListNode {
	prev := n.pre
	return prev
}

//获取当前节点的后一个节点
func (n *ListNode) Next() *ListNode {
	next := n.next
	return next
}

// 获取当前节点值
func (n *ListNode) GetValue() string {
	if n == nil {
		return ""
	}
	value := n.value
	return value
}

//创建空链表
func NewList() *List {
	return &List{}
}

//返回表头节点
func (l *List) Head() *ListNode {
	head := l.head
	return head
}

// 返回表尾节点
func (l *List) Tail() *ListNode {
	tail := l.tail
	return tail
}

// 返回链表长度
func (l *List) LenList() int {
	len := l.len
	return len
}

// 在链表右边插入数据
func (l *List) RPush(value string) {
	node := NewListNode(value)

	if l.len == 0 {
		l.head = node
		l.tail = node
	} else {
		tail := l.tail
		tail.next = node
		node.pre = tail

		l.tail = node
	}

	//更新长度
	l.len += 1
}

// 从链表左边取出数据
func (l *List) LPop() *ListNode {
	if l.len == 0 {
		return &ListNode{}
	}
	node := l.head
	if node.next == nil {
		// 链表为空
		l.head = nil
		l.tail = nil
	} else {
		l.head = node.next

	}
	l.len -= 1
	return node
}

// 通过索引查找节点
func (l *List) Index(index int) *ListNode {
	if index < 0 {
		index = (-index) - 1
		node := l.tail
		for {
			if node == nil {
				return nil
			}

			if index == 0 {
				return node
			}
			node = node.pre
			index--
		}
	} else {
		node := l.head
		for ; index > 0 && node != nil; index-- {
			node = node.next
		}
		return node
	}

}

// 返回指定区间的元素
func (l *List) RangeList(start, stop int) []*ListNode {
	nodes := make([]*ListNode, 0)
	if start < 0 {
		start = l.len + start
		if start < 0 {
			start = 0
		}
	}
	if start < 0 {
		start = l.len + start
		if start < 0 {
			start = 0
		}
	}

	rangeLen := stop - start + 1
	if rangeLen < 0 {
		return nodes
	}

	startNode := l.Index(start)
	for i := 0; i < rangeLen; i++ {
		if startNode == nil {
			break
		}
		nodes = append(nodes, startNode)
		startNode = startNode.next
	}
	return nodes

}
