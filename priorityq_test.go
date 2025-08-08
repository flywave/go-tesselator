package tesselator

import (
	"container/heap"
	"testing"
)

// 测试priorityq.go中的功能
func TestPriorityQ(t *testing.T) {
	// 创建测试顶点并正确初始化
	v1 := &vertex{s: 1.0, t: 2.0}
	v2 := &vertex{s: 3.0, t: 4.0}
	v3 := &vertex{s: 1.0, t: 1.0}
	v4 := &vertex{s: 2.0, t: 3.0}
	v5 := &vertex{s: 0.0, t: 5.0}

	// 初始化next和prev指针以避免nil引用
	v1.next, v1.prev = v1, v1
	v2.next, v2.prev = v2, v2
	v3.next, v3.prev = v3, v3
	v4.next, v4.prev = v4, v4
	v5.next, v5.prev = v5, v5

	// 测试1: 基本插入和提取最小元素
	t.Run("BasicInsertAndExtractMin", func(t *testing.T) {
		pq := newPriorityQ()

		// 插入顶点
		pq.insert(v1)
		pq.insert(v2)
		pq.insert(v3)
		pq.insert(v4)
		pq.insert(v5)

		// 验证提取顺序是否正确
		// 预期顺序: v5 (s=0), v3 (s=1,t=1), v1 (s=1,t=2), v4 (s=2), v2 (s=3)
		order := []*vertex{v5, v3, v1, v4, v2}

		for i := 0; i < len(order); i++ {
			min := pq.extractMin()
			if min != order[i] {
				t.Errorf("Expected vertex %v at position %d, got %v", order[i], i, min)
			}
		}

		// 队列应该为空
		if !pq.isEmpty() {
			t.Error("Expected queue to be empty after extracting all elements")
		}
	})

	// 测试2: 删除元素
	t.Run("DeleteElement", func(t *testing.T) {
		pq := newPriorityQ()

		// 插入顶点
		pq.insert(v1)
		pq.insert(v2)
		pq.insert(v3)
		pq.insert(v4)
		pq.insert(v5)

		// 删除v1
		pq.delete(v1)

		// 预期顺序: v5, v3, v4, v2
		order := []*vertex{v5, v3, v4, v2}

		for i := 0; i < len(order); i++ {
			min := pq.extractMin()
			if min != order[i] {
				t.Errorf("Expected vertex %v at position %d, got %v", order[i], i, min)
			}
		}
	})

	// 测试3: 批量插入
	t.Run("BatchInsert", func(t *testing.T) {
		pq := &pq{
			data:     make([]*vertex, 0),
			index:    make(map[*vertex]int),
			freeList: make([]int, 0),
		}

		// 批量插入
		vertices := []*vertex{v1, v2, v3, v4, v5}
		pq.batchInsert(vertices)
		heap.Init(pq)

		// 验证提取顺序
		order := []*vertex{v5, v3, v1, v4, v2}

		for i := 0; i < len(order); i++ {
			min := pq.extractMin()
			if min != order[i] {
				t.Errorf("Expected vertex %v at position %d, got %v", order[i], i, min)
			}
		}
	})

	// 测试4: 清空队列
	t.Run("ClearQueue", func(t *testing.T) {
		pq := newPriorityQ()

		pq.insert(v1)
		pq.insert(v2)
		pq.clear()

		if !pq.isEmpty() {
			t.Error("Expected queue to be empty after clear")
		}

		// 确保清空后可以继续使用
		pq.insert(v3)
		min := pq.extractMin()
		if min != v3 {
			t.Errorf("Expected vertex %v, got %v", v3, min)
		}
	})

	// 测试5: 检查最小值
	t.Run("CheckMinimum", func(t *testing.T) {
		pq := newPriorityQ()

		// 先插入v2，检查最小值是否是v2
		pq.insert(v2)
		min := pq.minimum()
		if min != v2 {
			t.Errorf("Expected minimum vertex %v, got %v", v2, min)
		}

		// 插入v5，检查最小值是否变为v5
		pq.insert(v5)
		min = pq.minimum()
		if min != v5 {
			t.Errorf("Expected minimum vertex %v, got %v", v5, min)
		}

		// 插入v3，检查最小值是否仍然是v5
		pq.insert(v3)
		min = pq.minimum()
		if min != v5 {
			t.Errorf("Expected minimum vertex %v, got %v", v5, min)
		}
	})
}

// 测试vertLeq函数
func TestVertLeq(t *testing.T) {
	v1 := &vertex{s: 1.0, t: 2.0}
	v2 := &vertex{s: 1.0, t: 2.0}
	v3 := &vertex{s: 1.0, t: 3.0}
	v4 := &vertex{s: 2.0, t: 1.0}

	if !vertLeq(v1, v2) {
		t.Error("Expected v1 <= v2")
	}

	if !vertLeq(v1, v3) {
		t.Error("Expected v1 <= v3")
	}

	if vertLeq(v3, v1) {
		t.Error("Expected v3 > v1")
	}

	if !vertLeq(v1, v4) {
		t.Error("Expected v1 <= v4")
	}

	if vertLeq(v4, v1) {
		t.Error("Expected v4 > v1")
	}
}
