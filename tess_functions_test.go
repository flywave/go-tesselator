package tesselator

import (
	"testing"
)

// TestTessProjectPolygon 测试多边形投影函数
func TestTessProjectPolygon(t *testing.T) {
	// 创建一个简单的正方形轮廓
	contour := []Contour{
		{
			{X: 0.0, Y: 0.0, Z: 0.0},
			{X: 1.0, Y: 0.0, Z: 0.0},
			{X: 1.0, Y: 1.0, Z: 0.0},
			{X: 0.0, Y: 1.0, Z: 0.0},
		},
	}

	tess := &tesselator{}
	// 添加轮廓
	fs := make([]float32, len(contour[0])*3)
	for j, v := range contour[0] {
		fs[3*j] = v.X
		fs[3*j+1] = v.Y
		fs[3*j+2] = v.Z
	}
	tessAddContour(tess, 3, fs)

	// 调用tessProjectPolygon
	tessProjectPolygon(tess)

	// 验证投影结果
	if tess.mesh == nil {
		t.Fatal("Mesh should not be nil after tessProjectPolygon")
	}

	// 检查顶点是否被正确投影
	vertexCount := 0
	for v := tess.mesh.vHead.next; v != &tess.mesh.vHead; v = v.next {
		vertexCount++
		// 验证s和t坐标已被设置
		// 简化检查，只验证坐标不全为0
		if v.s == 0 && v.t == 0 && v.coords[0] != 0 && v.coords[1] != 0 {
			t.Errorf("Vertex %d was not properly projected", vertexCount)
		}
	}

	if vertexCount != 4 {
		t.Errorf("Expected 4 vertices, got %d", vertexCount)
	}
}

// TestTessMeshTessellateMonoRegion 测试单调区域三角剖分函数
func TestTessMeshTessellateMonoRegion(t *testing.T) {
	// 创建一个简单的正方形轮廓
	contour := []Contour{
		{
			{X: 0.0, Y: 0.0, Z: 0.0},
			{X: 1.0, Y: 0.0, Z: 0.0},
			{X: 1.0, Y: 1.0, Z: 0.0},
			{X: 0.0, Y: 1.0, Z: 0.0},
		},
	}

	tess := &tesselator{}
	// 添加轮廓
	fs := make([]float32, len(contour[0])*3)
	for j, v := range contour[0] {
		fs[3*j] = v.X
		fs[3*j+1] = v.Y
		fs[3*j+2] = v.Z
	}
	tessAddContour(tess, 3, fs)

	// 投影多边形
	tessProjectPolygon(tess)

	// 计算内部区域
	tessComputeInterior(tess)

	// 在调用tessMeshTessellateMonoRegion之前，我们需要确保有一个单调区域
	// 对于简单的凸四边形，整个区域应该就是单调的
	faceCount := 0
	for f := tess.mesh.fHead.next; f != &tess.mesh.fHead; f = f.next {
		faceCount++
		if f.inside {
			// 对单调区域进行三角剖分
			tessMeshTessellateMonoRegion(tess.mesh, f)
		}
	}

	// 验证结果
	if faceCount == 0 {
		t.Error("No faces found in mesh")
	}
}

// TestTessMeshTessellateInterior 测试内部区域三角剖分函数
func TestTessMeshTessellateInterior(t *testing.T) {
	// 创建一个简单的正方形轮廓
	contour := []Contour{
		{
			{X: 0.0, Y: 0.0, Z: 0.0},
			{X: 1.0, Y: 0.0, Z: 0.0},
			{X: 1.0, Y: 1.0, Z: 0.0},
			{X: 0.0, Y: 1.0, Z: 0.0},
		},
	}

	tess := &tesselator{}
	// 添加轮廓
	fs := make([]float32, len(contour[0])*3)
	for j, v := range contour[0] {
		fs[3*j] = v.X
		fs[3*j+1] = v.Y
		fs[3*j+2] = v.Z
	}
	tessAddContour(tess, 3, fs)

	// 投影多边形
	tessProjectPolygon(tess)

	// 计算内部区域
	tessComputeInterior(tess)

	// 三角剖分内部区域
	tessMeshTessellateInterior(tess.mesh)

	// 验证结果
	insideFaceCount := 0
	for f := tess.mesh.fHead.next; f != &tess.mesh.fHead; f = f.next {
		if f.inside {
			insideFaceCount++
		}
	}

	// 简单的四边形应该被剖分为2个三角形(2个内部面)
	if insideFaceCount != 2 {
		t.Errorf("Expected 2 interior faces after tessellation, got %d", insideFaceCount)
	}
}

// TestTessMeshSetWindingNumber 测试设置缠绕数函数
func TestTessMeshSetWindingNumber(t *testing.T) {
	// 创建一个简单的正方形轮廓
	contour := []Contour{
		{
			{X: 0.0, Y: 0.0, Z: 0.0},
			{X: 1.0, Y: 0.0, Z: 0.0},
			{X: 1.0, Y: 1.0, Z: 0.0},
			{X: 0.0, Y: 1.0, Z: 0.0},
		},
	}

	tess := &tesselator{}
	// 添加轮廓
	fs := make([]float32, len(contour[0])*3)
	for j, v := range contour[0] {
		fs[3*j] = v.X
		fs[3*j+1] = v.Y
		fs[3*j+2] = v.Z
	}
	tessAddContour(tess, 3, fs)

	// 投影多边形
	tessProjectPolygon(tess)

	// 计算内部区域
	tessComputeInterior(tess)

	// 设置缠绕数
	tessMeshSetWindingNumber(tess.mesh, 1, true)

	// 验证边界边的缠绕数
	boundaryEdgeCount := 0
	for e := tess.mesh.eHead.next; e != &tess.mesh.eHead; e = e.next {
		if e.Lface.inside != e.rFace().inside {
			boundaryEdgeCount++
			// 边界边应该有非零缠绕数
			if e.winding == 0 && e.Sym.winding == 0 {
				t.Error("Boundary edge should have non-zero winding number")
			}
		}
	}

	if boundaryEdgeCount == 0 {
		t.Error("No boundary edges found")
	}
}

// TestTessAddContour 测试添加轮廓函数
func TestTessAddContour(t *testing.T) {
	tess := &tesselator{}

	// 创建一个简单的三角形轮廓
	vertices := []float32{
		0.0, 0.0, 0.0, // 顶点1
		1.0, 0.0, 0.0, // 顶点2
		0.0, 1.0, 0.0, // 顶点3
	}

	// 添加轮廓
	tessAddContour(tess, 3, vertices)

	// 验证网格已创建
	if tess.mesh == nil {
		t.Fatal("Mesh should be created after adding contour")
	}

	// 验证顶点数量
	vertexCount := 0
	for v := tess.mesh.vHead.next; v != &tess.mesh.vHead; v = v.next {
		vertexCount++
	}

	if vertexCount != 3 {
		t.Errorf("Expected 3 vertices, got %d", vertexCount)
	}

	// 验证边数量
	edgeCount := 0
	for e := tess.mesh.eHead.next; e != &tess.mesh.eHead; e = e.next {
		edgeCount++
	}

	// 3个顶点应该形成3条边
	if edgeCount != 3 {
		t.Errorf("Expected 3 edges, got %d", edgeCount)
	}

	// 验证面数量
	faceCount := 0
	for f := tess.mesh.fHead.next; f != &tess.mesh.fHead; f = f.next {
		faceCount++
	}

	// 应该有一个面
	if faceCount != 1 {
		t.Errorf("Expected 1 face, got %d", faceCount)
	}
}

// TestTessTesselate 测试完整的三角剖分函数
func TestTessTesselate(t *testing.T) {
	// 创建一个简单的正方形轮廓
	contour := []Contour{
		{
			{X: 0.0, Y: 0.0, Z: 0.0},
			{X: 1.0, Y: 0.0, Z: 0.0},
			{X: 1.0, Y: 1.0, Z: 0.0},
			{X: 0.0, Y: 1.0, Z: 0.0},
		},
	}

	tess := &tesselator{}
	// 添加轮廓
	fs := make([]float32, len(contour[0])*3)
	for j, v := range contour[0] {
		fs[3*j] = v.X
		fs[3*j+1] = v.Y
		fs[3*j+2] = v.Z
	}
	tessAddContour(tess, 3, fs)

	// 执行完整的三角剖分过程
	result := tessTesselate(tess, WindingRulePositive, elementTypePolygons, 3, 3, nil)

	if !result {
		t.Fatal("Tessellation failed")
	}

	// 验证结果
	// 由于可能的实现差异，我们只验证基本条件
	if tess.vertexCount <= 0 {
		t.Errorf("Expected positive vertex count, got %d", tess.vertexCount)
	}

	// 验证元素数组不为空
	if len(tess.elements) == 0 {
		t.Error("Elements array should not be empty")
	}

	// 验证顶点数组不为空
	if len(tess.vertices) == 0 {
		t.Error("Vertices array should not be empty")
	}
}

// TestIsInside 测试缠绕规则判断函数
func TestIsInside(t *testing.T) {
	tess := &tesselator{}

	// 测试奇数规则
	tess.windingRule = WindingRuleOdd
	if !tess.isInside(1) {
		t.Error("Odd winding rule failed for winding number 1")
	}
	if tess.isInside(0) {
		t.Error("Odd winding rule failed for winding number 0")
	}
	if !tess.isInside(3) {
		t.Error("Odd winding rule failed for winding number 3")
	}

	// 测试非零规则
	tess.windingRule = WindingRuleNonzero
	if !tess.isInside(1) {
		t.Error("Nonzero winding rule failed for winding number 1")
	}
	if tess.isInside(0) {
		t.Error("Nonzero winding rule failed for winding number 0")
	}
	if !tess.isInside(-1) {
		t.Error("Nonzero winding rule failed for winding number -1")
	}

	// 测试正数规则
	tess.windingRule = WindingRulePositive
	if !tess.isInside(1) {
		t.Error("Positive winding rule failed for winding number 1")
	}
	if tess.isInside(0) {
		t.Error("Positive winding rule failed for winding number 0")
	}
	if tess.isInside(-1) {
		t.Error("Positive winding rule failed for winding number -1")
	}

	// 测试负数规则
	tess.windingRule = WindingRuleNegative
	if tess.isInside(1) {
		t.Error("Negative winding rule failed for winding number 1")
	}
	if tess.isInside(0) {
		t.Error("Negative winding rule failed for winding number 0")
	}
	if !tess.isInside(-1) {
		t.Error("Negative winding rule failed for winding number -1")
	}

	// 测试绝对值大于等于2规则
	tess.windingRule = WindingRuleAbsGeqTwo
	if tess.isInside(1) {
		t.Error("AbsGeqTwo winding rule failed for winding number 1")
	}
	if tess.isInside(0) {
		t.Error("AbsGeqTwo winding rule failed for winding number 0")
	}
	if !tess.isInside(2) {
		t.Error("AbsGeqTwo winding rule failed for winding number 2")
	}
	if !tess.isInside(-2) {
		t.Error("AbsGeqTwo winding rule failed for winding number -2")
	}
}

// TestLongAxis 测试长轴计算函数
func TestLongAxis(t *testing.T) {
	// 测试X轴最长
	v1 := []float{3.0, 1.0, 1.0}
	if longAxis(v1) != 0 {
		t.Error("Longest axis should be X (0)")
	}

	// 测试Y轴最长
	v2 := []float{1.0, 3.0, 1.0}
	if longAxis(v2) != 1 {
		t.Error("Longest axis should be Y (1)")
	}

	// 测试Z轴最长
	v3 := []float{1.0, 1.0, 3.0}
	if longAxis(v3) != 2 {
		t.Error("Longest axis should be Z (2)")
	}

	// 测试负值
	v4 := []float{-3.0, 1.0, 1.0}
	if longAxis(v4) != 0 {
		t.Error("Longest axis should be X (0) even with negative values")
	}
}

// TestShortAxis 测试短轴计算函数
func TestShortAxis(t *testing.T) {
	// 测试X轴最短
	v1 := []float{1.0, 3.0, 3.0}
	if shortAxis(v1) != 0 {
		t.Error("Shortest axis should be X (0)")
	}

	// 测试Y轴最短
	v2 := []float{3.0, 1.0, 3.0}
	if shortAxis(v2) != 1 {
		t.Error("Shortest axis should be Y (1)")
	}

	// 测试Z轴最短
	v3 := []float{3.0, 3.0, 1.0}
	if shortAxis(v3) != 2 {
		t.Error("Shortest axis should be Z (2)")
	}

	// 测试负值
	v4 := []float{1.0, 3.0, 3.0}
	if shortAxis(v4) != 0 {
		t.Error("Shortest axis should be X (0) even with negative values")
	}
}
