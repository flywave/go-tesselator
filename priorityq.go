package tesselator

import (
	"container/heap"
)

type pq struct {
	data []*vertex
	// 为了提高删除效率，添加一个映射来存储顶点到索引的映射
	index map[*vertex]int
	// 自由列表，用于重用已删除的槽位
	freeList []int
}

func (p pq) Len() int {
	return len(p.data)
}

func (p pq) Less(i, j int) bool {
	return vertLeq(p.data[i], p.data[j])
}

func (p pq) Swap(i, j int) {
	p.data[i], p.data[j] = p.data[j], p.data[i]
	// 更新索引映射
	p.index[p.data[i]] = i
	p.index[p.data[j]] = j
}

// 清理无效的自由索引
func (p *pq) cleanupFreeList() {
	validFreeList := []int{}
	for _, idx := range p.freeList {
		if idx < len(p.data) {
			validFreeList = append(validFreeList, idx)
		}
	}
	p.freeList = validFreeList
}

func (p *pq) Push(x interface{}) {
	v := x.(*vertex)
	if len(p.freeList) > 0 {
		// 清理无效索引
		p.cleanupFreeList()
		if len(p.freeList) > 0 {
			// 重用自由列表中的槽位
			idx := p.freeList[len(p.freeList)-1]
			p.freeList = p.freeList[:len(p.freeList)-1]
			p.data[idx] = v
			p.index[v] = idx
			return
		}
	}
	// 没有可用的自由槽位，直接追加
	p.index[v] = len(p.data)
	p.data = append(p.data, v)
}

func (p *pq) Pop() interface{} {
	if len(p.data) == 0 {
		return nil
	}
	old := p.data
	x := old[len(old)-1]
	p.data = old[:len(old)-1]
	// 从索引映射中删除
	delete(p.index, x)
	// 将槽位添加到自由列表
	p.freeList = append(p.freeList, len(old)-1)
	return x
}

func newPriorityQ() *pq {
	p := &pq{
		data:     make([]*vertex, 0),
		index:    make(map[*vertex]int),
		freeList: make([]int, 0),
	}
	heap.Init(p)
	return p
}

// 批量插入多个顶点并初始化堆
func (p *pq) batchInsert(keys []*vertex) {
	for _, key := range keys {
		p.data = append(p.data, key)
		p.index[key] = len(p.data) - 1
	}
	heap.Init(p)
}

func (p *pq) insert(key *vertex) *vertex {
	heap.Push(p, key)
	return key
}

func (p *pq) extractMin() *vertex {
	if len(p.data) == 0 {
		return nil
	}
	return heap.Pop(p).(*vertex)
}

func (p *pq) delete(key *vertex) {
	// 检查key是否在队列中
	idx, exists := p.index[key]
	if !exists {
		return
	}
	heap.Remove(p, idx)
	// 从索引映射中删除
	delete(p.index, key)
	// 将槽位添加到自由列表
	p.freeList = append(p.freeList, idx)
}

func (p *pq) minimum() *vertex {
	if len(p.data) == 0 {
		return nil
	}
	return p.data[0]
}

// 添加isEmpty函数，检查队列是否为空
func (p *pq) isEmpty() bool {
	return len(p.data) == 0
}

// 添加clear函数，清空队列
func (p *pq) clear() {
	p.data = make([]*vertex, 0)
	p.index = make(map[*vertex]int)
	p.freeList = make([]int, 0)
	heap.Init(p)
}
