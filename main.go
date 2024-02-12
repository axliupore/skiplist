package main

import "fmt"

func main() {
	// 创建一个新的跳表实例
	skiplist := NewSkiplist()

	// 插入键值对
	skiplist.Put(3, 30)
	skiplist.Put(1, 10)
	skiplist.Put(2, 20)

	// 查找键为2的值
	val, found := skiplist.Get(2)
	if found {
		fmt.Println("Value for key 2:", val)
	} else {
		fmt.Println("Key 2 not found.")
	}

	// 删除键为1的节点
	skiplist.Del(1)

	// 查找键为1的值（已删除的节点）
	val, found = skiplist.Get(1)
	if found {
		fmt.Println("Value for key 1:", val)
	} else {
		fmt.Println("Key 1 not found.")
	}

	// 范围查询键值在范围 [2, 3] 内的节点
	rangeResult := skiplist.Range(2, 3)
	fmt.Println("Nodes in range [2, 3]:", rangeResult)

	// 查找大于等于键为2且最接近的节点
	ceilingResult, ceilingFound := skiplist.Ceiling(2)
	if ceilingFound {
		fmt.Println("Ceiling for key 2:", ceilingResult)
	} else {
		fmt.Println("Ceiling for key 2 not found.")
	}

	// 查找小于等于键为2且最接近的节点
	floorResult, floorFound := skiplist.Floor(2)
	if floorFound {
		fmt.Println("Floor for key 2:", floorResult)
	} else {
		fmt.Println("Floor for key 2 not found.")
	}
}
