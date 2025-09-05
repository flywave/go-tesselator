package tesselator

import (
	"testing"
)

// TestTessMeshNewMesh 测试创建新网格
func TestTessMeshNewMesh(t *testing.T) {
	mesh := tessMeshNewMesh()
	if mesh == nil {
		t.Fatal("tessMeshNewMesh returned nil")
	}

	// 检查网格头部是否正确初始化
	if mesh.vHead.next != &mesh.vHead || mesh.vHead.prev != &mesh.vHead {
		t.Error("Vertex head not properly initialized")
	}
	if mesh.fHead.next != &mesh.fHead || mesh.fHead.prev != &mesh.fHead {
		t.Error("Face head not properly initialized")
	}
	if mesh.eHead.next != &mesh.eHead {
		t.Error("Edge head not properly initialized")
	}
	if mesh.eHead.Sym != &mesh.eHeadSym {
		t.Error("Edge head symmetry not properly initialized")
	}
}

// TestTessMeshMakeEdge 测试创建边
func TestTessMeshMakeEdge(t *testing.T) {
	mesh := tessMeshNewMesh()
	edge := tessMeshMakeEdge(mesh)

	if edge == nil {
		t.Fatal("tessMeshMakeEdge returned nil")
	}

	// 检查边是否正确创建
	if edge.Sym == nil {
		t.Error("Edge symmetry not created")
	}

	if edge.Org == nil {
		t.Error("Edge origin vertex not created")
	}

	if edge.Sym.Org == nil {
		t.Error("Edge symmetric origin vertex not created")
	}

	if edge.Lface == nil {
		t.Error("Edge left face not created")
	}

	// 检查网格中的元素数量
	vertexCount := 0
	for v := mesh.vHead.next; v != &mesh.vHead; v = v.next {
		vertexCount++
	}
	if vertexCount != 2 {
		t.Errorf("Expected 2 vertices, got %d", vertexCount)
	}

	faceCount := 0
	for f := mesh.fHead.next; f != &mesh.fHead; f = f.next {
		faceCount++
	}
	if faceCount != 1 {
		t.Errorf("Expected 1 face, got %d", faceCount)
	}
}

// TestTessMeshSplice 测试拼接边
func TestTessMeshSplice(t *testing.T) {
	mesh := tessMeshNewMesh()

	// 创建两个边
	e1 := tessMeshMakeEdge(mesh)
	e2 := tessMeshMakeEdge(mesh)

	// 确保两个边有不同的顶点
	if e1.Org == e2.Org {
		t.Error("Edges should have different origin vertices initially")
	}

	// 拼接两个边
	tessMeshSplice(mesh, e1, e2)

	// 检查拼接后的结果
	// 顶点应该被合并
	if e1.Org != e2.Org {
		t.Error("Origin vertices should be the same after splice")
	}
}

// TestTessMeshDelete 测试删除边
func TestTessMeshDelete(t *testing.T) {
	mesh := tessMeshNewMesh()

	// 创建一个边
	edge := tessMeshMakeEdge(mesh)

	// 删除边
	tessMeshDelete(mesh, edge)

	// 检查边是否被删除
	edgeCount := 0
	for e := mesh.eHead.next; e != &mesh.eHead; e = e.next {
		edgeCount++
	}
	// 应该没有边剩下
	if edgeCount != 0 {
		t.Errorf("Expected 0 edges after deletion, got %d", edgeCount)
	}
}

// TestTessMeshAddEdgeVertex 测试添加边顶点
func TestTessMeshAddEdgeVertex(t *testing.T) {
	mesh := tessMeshNewMesh()

	// 创建一个边
	eOrg := tessMeshMakeEdge(mesh)

	// 添加边顶点
	eNew := tessMeshAddEdgeVertex(mesh, eOrg)

	if eNew == nil {
		t.Fatal("tessMeshAddEdgeVertex returned nil")
	}

	// 检查新边是否正确创建
	if eNew.Sym == nil {
		t.Error("New edge symmetry not created")
	}

	// 检查新顶点是否创建
	if eNew.Sym.Org == nil {
		t.Error("New vertex not created")
	}

	// 检查边的关系 - 修复：eNew.Lnext 应该等于 eOrg.Lnext
	if eNew.Lnext != eOrg.Lnext {
		t.Error("New edge Lnext not properly set")
	}
}

// TestTessMeshSplitEdge 测试分割边
func TestTessMeshSplitEdge(t *testing.T) {
	mesh := tessMeshNewMesh()

	// 创建一个边
	eOrg := tessMeshMakeEdge(mesh)

	// 分割边
	eNew := tessMeshSplitEdge(mesh, eOrg)

	if eNew == nil {
		t.Fatal("tessMeshSplitEdge returned nil")
	}

	// 检查新边是否正确创建
	if eNew.Sym == nil {
		t.Error("New edge symmetry not created")
	}

	// 检查新顶点是否创建
	if eNew.Org == nil {
		t.Error("New vertex not created")
	}

	// 检查边的关系
	if eOrg.Lnext != eNew {
		t.Error("Original edge Lnext not properly set")
	}
}

// TestTessMeshConnect 测试连接边
func TestTessMeshConnect(t *testing.T) {
	mesh := tessMeshNewMesh()

	// 创建两个边
	eOrg := tessMeshMakeEdge(mesh)
	eDst := tessMeshMakeEdge(mesh)

	// 连接两个边
	eNew := tessMeshConnect(mesh, eOrg, eDst)

	if eNew == nil {
		t.Fatal("tessMeshConnect returned nil")
	}

	// 检查新边是否正确创建
	if eNew.Sym == nil {
		t.Error("New edge symmetry not created")
	}

	// 检查顶点连接
	if eNew.Org != eOrg.dst() {
		t.Error("New edge origin not properly set")
	}

	if eNew.Sym.Org != eDst.Org {
		t.Error("New edge symmetric origin not properly set")
	}
}

// TestTessMeshZapFace 测试删除面
func TestTessMeshZapFace(t *testing.T) {
	mesh := tessMeshNewMesh()

	// 创建一个边（带有一个面）
	edge := tessMeshMakeEdge(mesh)
	face := edge.Lface

	// 删除面
	tessMeshZapFace(mesh, face)

	// 检查面是否被删除
	faceCount := 0
	for f := mesh.fHead.next; f != &mesh.fHead; f = f.next {
		faceCount++
	}
	// 应该没有面剩下
	if faceCount != 0 {
		t.Errorf("Expected 0 faces after zapping, got %d", faceCount)
	}

	// 检查边的左面是否被设置为nil
	if edge.Lface != nil {
		t.Error("Edge left face should be nil after zapping face")
	}
}

// TestMeshNavigation 测试网格导航
func TestMeshNavigation(t *testing.T) {
	mesh := tessMeshNewMesh()

	// 创建一个边
	edge := tessMeshMakeEdge(mesh)

	// 测试对称边
	if edge.Sym.Sym != edge {
		t.Error("Symmetry of symmetry should be the original edge")
	}

	// 测试Onext - 修复：对于单边循环，Onext应该指向对称边
	if edge.Onext != edge.Sym {
		t.Error("Onext should point to symmetric edge for a single edge loop")
	}

	// 测试Lnext - 修复：对于单边循环，Lnext应该指向对称边
	if edge.Lnext != edge.Sym {
		t.Error("Lnext should point to symmetric edge for a single edge loop")
	}
}

// TestVertexOperations 测试顶点操作
func TestVertexOperations(t *testing.T) {
	mesh := tessMeshNewMesh()

	// 创建一个边
	edge := tessMeshMakeEdge(mesh)
	vertex := edge.Org

	// 添加调试信息
	println("TestVertexOperations: edge =", edge)
	println("TestVertexOperations: vertex =", vertex)
	if vertex != nil {
		println("TestVertexOperations: vertex.anEdge =", vertex.anEdge)
	}

	// 检查顶点的anEdge指针
	if vertex.anEdge != edge {
		t.Errorf("Vertex anEdge should point to the edge, but vertex.anEdge=%v and edge=%v", vertex.anEdge, edge)
	}

	// 检查顶点在网格中的链接
	if vertex.next.prev != vertex {
		t.Error("Vertex next.prev should point back to vertex")
	}
	if vertex.prev.next != vertex {
		t.Error("Vertex prev.next should point back to vertex")
	}
}

// TestFaceOperations 测试面操作
func TestFaceOperations(t *testing.T) {
	mesh := tessMeshNewMesh()

	// 创建一个边
	edge := tessMeshMakeEdge(mesh)
	face := edge.Lface

	// 检查面的anEdge指针
	if face.anEdge != edge {
		t.Error("Face anEdge should point to the edge")
	}

	// 检查面在网格中的链接
	if face.next.prev != face {
		t.Error("Face next.prev should point back to face")
	}
	if face.prev.next != face {
		t.Error("Face prev.next should point back to face")
	}
}

// TestEdgeOperations 测试边操作
func TestEdgeOperations(t *testing.T) {
	mesh := tessMeshNewMesh()

	// 创建一个边
	edge := tessMeshMakeEdge(mesh)

	// 检查边在网格中的链接
	if edge.next.Sym.next != edge.Sym {
		t.Error("Edge next.Sym.next should point back to edge.Sym")
	}
	if edge.Sym.next.Sym.next != edge {
		t.Error("Edge.Sym.next.Sym.next should point back to edge")
	}
}
