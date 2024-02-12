package main

import "math/rand"

type Skiplist struct {
	head *node
}

type node struct {
	nexts    []*node // 长度为当前节点的高度
	key, val int
}

// NewSkiplist 初始化 Skiplist
func NewSkiplist() *Skiplist {
	return &Skiplist{
		head: &node{
			nexts: make([]*node, 1),
		},
	}
}

// Get 根据 key 读取 val，第二个 bool flag 反映 key 在 Skiplist 中是否存在
func (s *Skiplist) Get(key int) (int, bool) {
	if n := s.search(key); n != nil {
		return n.val, true
	}
	return -1, false
}

// 在 Skiplist 中检索 key 对应的 node
func (s *Skiplist) search(key int) *node {
	// 每次检索从头部出发
	move := s.head
	// 每次检索从最大高度出发，直到来到首层
	for level := len(s.head.nexts) - 1; level >= 0; level-- {
		// 在每一层中持续向右遍历，直到下一个节点不存在或者 key 值大于等于 key
		for move.nexts[level] != nil && move.nexts[level].key < key {
			move = move.nexts[level]
		}
		// 如果 key 值相等，则找到了目标直接返回
		if move.nexts[level] != nil && move.nexts[level].key == key {
			return move.nexts[level]
		}
		// 当前层没找到目标，则层数减 1，继续向下
	}
	return nil
}

// 随机值，决定一个待插入的新节点在 Skiplist 中最高层对应的 index
func (s *Skiplist) roll() int {
	const m = 1 << 30 // 设置一个较大的上限值
	var level int
	// 生成一个范围在 [0, max) 内的随机整数
	for rand.Intn(m) < m/2 {
		level++
	}
	return level
}

// Put 将 key, val 加入 Skiplist
func (s *Skiplist) Put(key, val int) {
	// 假如 kv 对已存在，则直接对值进行更新并返回
	if n := s.search(key); n != nil {
		n.val = val
		return
	}

	// roll 出新节点的高度
	level := s.roll()

	// 新节点高度超出跳表最大高度，则需要对高度进行补齐
	for len(s.head.nexts)-1 < level {
		s.head.nexts = append(s.head.nexts, nil)
	}

	// 创建出新的节点
	newNode := node{
		key:   key,
		val:   val,
		nexts: make([]*node, level+1),
	}

	// 从头节点的最高层出发
	move := s.head
	for level := level; level >= 0; level-- {
		// 向右遍历，直到右侧节点不存在或者 key 值大于 key
		for move.nexts[level] != nil && move.nexts[level].key < key {
			move = move.nexts[level]
		}

		// 调整指针关系，完成新节点的插入
		newNode.nexts[level] = move.nexts[level]
		move.nexts[level] = &newNode
	}
}

// Del 根据 key 从跳表中删除对应的节点

func (s *Skiplist) Del(key int) {
	// 如果 kv 对不存在，则无需删除直接返回
	if n := s.search(key); n == nil {
		return
	}

	// 从头节点的最高层出发
	move := s.head
	for level := len(s.head.nexts) - 1; level > 0; level-- {
		// 向右遍历，直到右侧节点不存在或者 key 值大于等于 key
		for move.nexts[level] != nil && move.nexts[level].key < key {
			move = move.nexts[level]
		}

		// 右侧节点不存在或者 key 值大于 target，则直接跳过
		if move.nexts[level] == nil || move.nexts[level].key > key {
			continue
		}

		// 走到此处意味着右侧节点的 key 值必然等于 key，则调整指针引用
		move.nexts[level] = move.nexts[level].nexts[level]
	}

	// 对跳表的最大高度进行更新
	var dif int
	// 倘若某一层已经不存在数据节点，高度需要递减
	for level := len(s.head.nexts) - 1; level > 0 && s.head.nexts[level] == nil; level-- {
		dif++
	}
	s.head.nexts = s.head.nexts[:len(s.head.nexts)-dif]
}

// 找到 key 值大于等于 target 且 key 值最接近于 target 的节点
func (s *Skiplist) ceiling(target int) *node {
	move := s.head

	// 自上而下，找到 key 值小于 target 且最接近 target 的 kv 对
	for level := len(s.head.nexts) - 1; level >= 0; level-- {
		for move.nexts[level] != nil && move.nexts[level].key < target {
			move = move.nexts[level]
		}
		// 如果 key 值等于 target 的 kv 对存在，则直接返回
		if move.nexts[level] != nil && move.nexts[level].key == target {
			return move.nexts[level]
		}
	}

	// 此时 move 已经对应于在首层 key 值小于 key 且最接近于 key 的节点，其右侧第一个节点即为所寻找的目标节点
	return move.nexts[0]
}

// Range 找到 Skiplist 当中 ≥ start，且 ≤ end 的 kv 对
func (s *Skiplist) Range(start, end int) [][2]int {
	// 首先通过 ceiling 方法，找到 Skiplist 中 key 值大于等于 start 且最接近于 start 的节点 ceilNode
	ceilNode := s.ceiling(start)
	// 如果不存在，直接返回
	if ceilNode == nil {
		return [][2]int{}
	}

	// 从 ceilNode 首层出发向右遍历，把所有位于 [start,end] 区间内的节点统统返回
	var res [][2]int
	for move := ceilNode; move != nil && move.key <= end; move = move.nexts[0] {
		res = append(res, [2]int{move.key, move.val})
	}
	return res
}

// Ceiling 找到 Skiplist 中，key 值大于等于 target 且最接近于 target 的 key-value 对
func (s *Skiplist) Ceiling(target int) ([2]int, bool) {
	if ceilNode := s.ceiling(target); ceilNode != nil {
		return [2]int{ceilNode.key, ceilNode.val}, true
	}

	return [2]int{}, false
}

// 找到 key 值小于等于 target 且 key 值最接近于 target 的节点
func (s *Skiplist) floor(target int) *node {
	move := s.head

	// 自上而下，找到 key 值小于 target 且最接近 target 的 kv 对
	for level := len(s.head.nexts) - 1; level >= 0; level-- {
		for move.nexts[level] != nil && move.nexts[level].key < target {
			move = move.nexts[level]
		}
		// 如果 key 值等于 target 的 kv对存在，则直接返回
		if move.nexts[level] != nil && move.nexts[level].key == target {
			return move.nexts[level]
		}
	}

	// move 是首层中 key 值小于 target 且最接近 target 的节点，直接返回 move 即可
	return move
}

// Floor 找到 skiplist 中，key 值小于等于 target 且最接近于 target 的 key-value 对
func (s *Skiplist) Floor(target int) ([2]int, bool) {
	// 引用 floor 方法，取 floorNode 值进行返回
	if floorNode := s.floor(target); floorNode != nil {
		return [2]int{floorNode.key, floorNode.val}, true
	}

	return [2]int{}, false
}
