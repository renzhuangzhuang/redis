package list

import "testing"

func TestListNode(t *testing.T) {
	// 这块用于测试 链表节点的性能

	//初始化一个 节点
	node := NewListNode("redis")

	node1 := NewListNode("hello")

	node2 := NewListNode("ListNode")

	//构建双链表用于下面功能测试
	node.next = node1
	node1.pre = node
	node1.next = node2
	node2.pre = node1

	// 测试
	if node.Next().value != "hello" {
		t.Error(`获取后一个节点功能失败`)
	}

	if node1.Prev().value != "redis" {
		t.Error(`获取前一个节点功能失败`)
	}

	if node.GetValue() != "redis" {
		t.Error(`获取节点值功能失败`)
	}

}

func TestList(t *testing.T) {
	// 测试链表功能

	list := NewList()
	test_string := []string{"hello redis", "hello golang", "hello world"}
	for _, value := range test_string {
		list.RPush(value)
	}
	if list.LenList() != 3 {
		t.Error(`右边插入数据功能失败`)
	}

	node := list.LPop()

	if node.GetValue() != "hello redis" {
		t.Error(`左边获取数据失败`)
	}

	//插入功能还需完善 1 代表的是除头结点以后的节点开始算的
	node1 := list.Index(1)

	if node1.GetValue() != "hello world" {
		t.Error(`查询功能失败`)
	}

	// 测试这块功能时候 需要单独测试
	node2 := list.RangeList(0, 2)

	for i := 0; i < len(node2); i++ {
		if node2[i].GetValue() != test_string[i] {
			t.Error(`返回失败`)
		}
	}

}
