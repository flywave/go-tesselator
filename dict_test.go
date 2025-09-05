package tesselator

import (
	"testing"
)

// 模拟tesselator结构体，用于测试
func newMockTesselator() *tesselator {
	t := &tesselator{
		mesh:  tessMeshNewMesh(),
		pq:    newPriorityQ(),
		event: newMockVertex(0, 0, 0), // 初始化event字段
	}
	t.dict = newDict(t)
	return t
}

// 模拟创建activeRegion
func newMockActiveRegion(eUp *halfEdge) *activeRegion {
	return &activeRegion{
		eUp:           eUp,
		windingNumber: 0,
		inside:        false,
		sentinel:      false,
		dirty:         false,
		fixUpperEdge:  false,
	}
}

// 模拟创建halfEdge，并确保满足edgeEval的断言条件
func newMockHalfEdge(org, dst *vertex) *halfEdge {
	e := &halfEdge{
		Org: org,
	}
	e.Sym = &halfEdge{
		Org: dst,
		Sym: e,
	}
	// 确保e.dst()函数可以正常工作
	e.Sym.Org = dst
	return e
}

// 模拟创建vertex
func newMockVertex(id int, s, t float) *vertex {
	v := &vertex{
		id: id,
		s:  s,
		t:  t,
	}
	v.next = v
	v.prev = v
	return v
}

// 测试dict的基本功能，使用简化的测试数据
func TestDictBasicOperations(t *testing.T) {
	tess := newMockTesselator()
	d := tess.dict

	// 创建测试用的顶点和边
	v1 := newMockVertex(1, 0, 0)
	v2 := newMockVertex(2, 1, 1)

	// 设置tess.event为v1
	tess.event = v1

	// 创建边，确保方向正确
	e1 := newMockHalfEdge(v1, v2) // 从v1到v2 (s从0到1)

	// 创建测试用的activeRegion
	r1 := newMockActiveRegion(e1)

	// 测试插入
	n1 := d.insert(r1)

	// 验证插入结果
	if dictKey(n1) != r1 {
		t.Errorf("Insert failed: expected keys do not match")
	}

	// 测试min
	minNode := d.min()
	if dictKey(minNode) != r1 {
		t.Errorf("min should return the only node")
	}

	// 测试删除
	dictDelete(n1)

	// 验证删除后字典为空
	if d.min() != &d.head {
		t.Errorf("Delete failed: dictionary should be empty")
	}
}

// 测试dict的insertBefore功能
func TestDictInsertBefore(t *testing.T) {
	tess := newMockTesselator()
	d := newDict(tess)

	// 创建测试用的activeRegion
	r1 := newMockActiveRegion(newMockHalfEdge(newMockVertex(1, 0, 0), newMockVertex(2, 1, 1)))
	r2 := newMockActiveRegion(newMockHalfEdge(newMockVertex(2, 1, 1), newMockVertex(3, 2, 2)))
	r3 := newMockActiveRegion(newMockHalfEdge(newMockVertex(3, 2, 2), newMockVertex(1, 0, 0)))

	// 先插入r1和r3
	n1 := d.insert(r1)
	d.insert(r3)

	// 在n1之前插入r2
	n2 := d.insertBefore(n1, r2)

	// 验证插入位置正确
	if n2.next != n1 {
		t.Errorf("insertBefore failed: new node's next is not the specified node")
	}
	if n1.prev != n2 {
		t.Errorf("insertBefore failed: specified node's prev is not the new node")
	}
}

// 测试dict的迭代功能
func TestDictIteration(t *testing.T) {
	tess := newMockTesselator()
	d := newDict(tess)

	// 创建多个activeRegion并插入
	regions := []*activeRegion{}
	for i := 0; i < 5; i++ {
		v1 := newMockVertex(i, float(i), float(i))
		v2 := newMockVertex(i+1, float(i+1), float(i+1))
		e := newMockHalfEdge(v1, v2)
		r := newMockActiveRegion(e)
		regions = append(regions, r)
		d.insert(r)
	}

	// 遍历dict并收集所有节点
	collected := []*activeRegion{}
	current := d.head.next
	for current != &d.head {
		collected = append(collected, dictKey(current))
		current = current.next
	}

	// 验证所有插入的节点都被遍历到
	if len(collected) != len(regions) {
		t.Errorf("Iteration failed: expected %d nodes, got %d", len(regions), len(collected))
	}

	// 检查所有节点都在收集的列表中
	for _, r := range regions {
		found := false
		for _, c := range collected {
			if r == c {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Iteration failed: missing region")
		}
	}
}

// 测试dict的search功能
func TestDictSearch(t *testing.T) {
	tess := newMockTesselator()
	d := newDict(tess)

	// 创建测试用的activeRegion
	r1 := newMockActiveRegion(newMockHalfEdge(newMockVertex(1, 0, 0), newMockVertex(2, 1, 1)))
	r2 := newMockActiveRegion(newMockHalfEdge(newMockVertex(2, 1, 1), newMockVertex(3, 2, 2)))
	r3 := newMockActiveRegion(newMockHalfEdge(newMockVertex(3, 2, 2), newMockVertex(4, 3, 3)))

	// 插入节点
	d.insert(r1)
	d.insert(r2)
	d.insert(r3)

	// 测试搜索存在的key (应该返回具有相同key的节点)
	result := d.search(r2)
	// 由于我们无法完全控制排序，我们只验证返回的节点不是head节点
	if result == &d.head {
		t.Errorf("Search failed: should not return head node for existing key")
	}

	// 测试空字典的搜索
	emptyDict := newDict(tess)
	result = emptyDict.search(r1)
	if dictKey(result) != nil {
		t.Errorf("Search failed: expected to find nil in empty dict")
	}
}

// 测试dictSucc和dictPred功能
func TestDictSuccAndPred(t *testing.T) {
	tess := newMockTesselator()
	d := newDict(tess)

	// 创建测试用的activeRegion
	r1 := newMockActiveRegion(newMockHalfEdge(newMockVertex(1, 0, 0), newMockVertex(2, 1, 1)))
	r2 := newMockActiveRegion(newMockHalfEdge(newMockVertex(2, 1, 1), newMockVertex(3, 2, 2)))
	r3 := newMockActiveRegion(newMockHalfEdge(newMockVertex(3, 2, 2), newMockVertex(4, 3, 3)))

	// 按顺序插入节点
	n1 := d.insert(r1)
	n2 := d.insert(r2)
	n3 := d.insert(r3)

	// 测试dictSucc
	if dictSucc(n1) != n2 {
		t.Errorf("dictSucc failed: n1's successor should be n2")
	}
	if dictSucc(n2) != n3 {
		t.Errorf("dictSucc failed: n2's successor should be n3")
	}
	// head的前驱应该是最后一个节点
	if dictSucc(&d.head) != n1 {
		t.Errorf("dictSucc failed: head's successor should be n1")
	}

	// 测试dictPred
	if dictPred(n3) != n2 {
		t.Errorf("dictPred failed: n3's predecessor should be n2")
	}
	if dictPred(n2) != n1 {
		t.Errorf("dictPred failed: n2's predecessor should be n1")
	}
	// head的前驱应该是最后一个节点
	if dictPred(&d.head) != n3 {
		t.Errorf("dictPred failed: head's predecessor should be n3")
	}
}

// 测试dict的max功能
func TestDictMax(t *testing.T) {
	tess := newMockTesselator()
	d := newDict(tess)

	// 测试空字典的max
	maxNode := d.max()
	if maxNode != &d.head {
		t.Errorf("max failed: empty dict should return head")
	}

	// 创建测试用的activeRegion
	r1 := newMockActiveRegion(newMockHalfEdge(newMockVertex(1, 0, 0), newMockVertex(2, 1, 1)))
	r2 := newMockActiveRegion(newMockHalfEdge(newMockVertex(2, 1, 1), newMockVertex(3, 2, 2)))
	r3 := newMockActiveRegion(newMockHalfEdge(newMockVertex(3, 2, 2), newMockVertex(4, 3, 3)))

	// 按顺序插入节点
	d.insert(r1)
	d.insert(r2)
	n3 := d.insert(r3)

	// 测试非空字典的max
	maxNode = d.max()
	if maxNode != n3 {
		t.Errorf("max failed: should return the last inserted node")
	}
}
