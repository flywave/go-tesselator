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

	// 测试6: 自由列表功能
	t.Run("FreeListFunctionality", func(t *testing.T) {
		pq := newPriorityQ()

		// 插入一些元素
		pq.insert(v1)
		pq.insert(v2)
		pq.insert(v3)

		// 检查自由列表为空
		if len(pq.freeList) != 0 {
			t.Errorf("Expected freeList to be empty, got %d elements", len(pq.freeList))
		}

		// 删除一个元素，这应该会向自由列表添加一个索引
		pq.delete(v2)

		// 检查自由列表不为空
		if len(pq.freeList) == 0 {
			t.Error("Expected freeList to contain elements after deletion")
		}

		// 插入新元素，应该重用自由列表中的槽位
		pq.insert(v4)

		// 验证队列状态
		vertices := make([]*vertex, 0, 2)
		for !pq.isEmpty() {
			vertices = append(vertices, pq.extractMin())
		}

		// 验证我们得到了正确的元素数量
		if len(vertices) != 2 {
			t.Errorf("Expected 2 vertices, got %d", len(vertices))
		}
	})

	// 测试7: 索引映射功能
	t.Run("IndexMapFunctionality", func(t *testing.T) {
		pq := newPriorityQ()

		// 插入元素
		pq.insert(v1)
		pq.insert(v2)
		pq.insert(v3)

		// 验证索引映射存在
		if _, exists := pq.index[v1]; !exists {
			t.Error("Expected v1 to be in index map")
		}
		if _, exists := pq.index[v2]; !exists {
			t.Error("Expected v2 to be in index map")
		}
		if _, exists := pq.index[v3]; !exists {
			t.Error("Expected v3 to be in index map")
		}

		// 删除元素后，验证索引映射被移除
		pq.delete(v2)
		if _, exists := pq.index[v2]; exists {
			t.Error("Expected v2 to be removed from index map")
		}

		// 其他元素应该仍然在索引映射中
		if _, exists := pq.index[v1]; !exists {
			t.Error("Expected v1 to still be in index map")
		}
		if _, exists := pq.index[v3]; !exists {
			t.Error("Expected v3 to still be in index map")
		}
	})

	// 测试8: 批量插入功能
	t.Run("BatchInsertFunctionality", func(t *testing.T) {
		pq := newPriorityQ()

		// 批量插入
		vertices := []*vertex{v1, v2, v3, v4, v5}
		pq.batchInsert(vertices)

		// 验证所有元素都被插入
		expectedOrder := []*vertex{v5, v3, v1, v4, v2}
		for i := 0; i < len(expectedOrder); i++ {
			min := pq.extractMin()
			if min != expectedOrder[i] {
				t.Errorf("Expected vertex %v at position %d, got %v", expectedOrder[i], i, min)
			}
		}
	})

	// 测试9: 空队列操作
	t.Run("EmptyQueueOperations", func(t *testing.T) {
		pq := newPriorityQ()

		// 在空队列上操作
		if min := pq.extractMin(); min != nil {
			t.Errorf("Expected nil from extractMin on empty queue, got %v", min)
		}

		if min := pq.minimum(); min != nil {
			t.Errorf("Expected nil from minimum on empty queue, got %v", min)
		}

		if !pq.isEmpty() {
			t.Error("Expected empty queue to return true for isEmpty")
		}

		// 删除不存在的元素应该不会出错
		pq.delete(v1)
	})

	// 测试10: 重复元素插入
	t.Run("DuplicateElementInsertion", func(t *testing.T) {
		pq := newPriorityQ()

		// 插入相同元素多次
		pq.insert(v1)
		pq.insert(v1) // 插入相同元素

		// 应该有两个v1在队列中
		min1 := pq.extractMin()
		min2 := pq.extractMin()

		if min1 != v1 || min2 != v1 {
			t.Error("Expected both elements to be v1")
		}

		// 队列应该为空
		if !pq.isEmpty() {
			t.Error("Expected queue to be empty after extracting all elements")
		}
	})
}
