package tesselator

import (
	"fmt"
	"testing"
)

// TestSimpleQuadTessellation 测试简单四边形的三角剖分
// 一个四边形应该被剖分为两个三角形
func TestSimpleQuadTessellation(t *testing.T) {
	// 定义一个简单的四边形顶点 (0,0), (1,0), (1,1), (0,1)
	quadVertices := []Vertex{
		{X: 0, Y: 0, Z: 0},
		{X: 1, Y: 0, Z: 0},
		{X: 1, Y: 1, Z: 0},
		{X: 0, Y: 1, Z: 0},
	}

	// 创建轮廓
	contour := Contour(quadVertices)
	contours := []Contour{contour}

	// 执行三角剖分
	indices, vertices, err := Tesselate(contours, WindingRuleOdd)

	// 输出调试信息
	fmt.Printf("TestSimpleQuadTessellation: 错误: %v\n", err)
	fmt.Printf("TestSimpleQuadTessellation: 顶点数量: %d\n", len(vertices))
	fmt.Printf("TestSimpleQuadTessellation: 索引数量: %d\n", len(indices))
	fmt.Printf("TestSimpleQuadTessellation: 顶点: %v\n", vertices)
	fmt.Printf("TestSimpleQuadTessellation: 索引: %v\n", indices)

	// 检查是否有错误
	if err != nil {
		t.Errorf("Tessellate returned an error: %v", err)
	}

	// 对于一个四边形，我们期望:
	// 1. 4个顶点 (输入的4个顶点)
	// 2. 6个索引 (2个三角形，每个三角形3个顶点)

	// 检查顶点数量
	if len(vertices) != 4 {
		t.Errorf("Expected 4 vertices, got %d", len(vertices))
	}

	// 检查索引数量
	if len(indices) != 6 {
		t.Errorf("Expected 6 indices (2 triangles), got %d", len(indices))
	}

	// 验证顶点坐标
	expectedVertices := []Vertex{
		{X: 0, Y: 0, Z: 0},
		{X: 1, Y: 0, Z: 0},
		{X: 1, Y: 1, Z: 0},
		{X: 0, Y: 1, Z: 0},
	}

	for i, expected := range expectedVertices {
		if i < len(vertices) {
			if vertices[i].X != expected.X || vertices[i].Y != expected.Y || vertices[i].Z != expected.Z {
				t.Errorf("Vertex %d: expected (%.1f, %.1f, %.1f), got (%.1f, %.1f, %.1f)",
					i, expected.X, expected.Y, expected.Z, vertices[i].X, vertices[i].Y, vertices[i].Z)
			}
		}
	}

	// 验证索引范围
	for i, idx := range indices {
		if idx < 0 || idx >= len(vertices) {
			t.Errorf("Index %d is out of range: %d (should be 0-%d)", i, idx, len(vertices)-1)
		}
	}

	// 验证三角形构成
	// 应该有两个三角形，例如 (0,1,2) 和 (0,2,3) 或其他有效组合
	if len(indices) >= 6 {
		fmt.Printf("Triangle 1: [%d, %d, %d]\n", indices[0], indices[1], indices[2])
		fmt.Printf("Triangle 2: [%d, %d, %d]\n", indices[3], indices[4], indices[5])
	}
}

// TestSimpleTriangleTessellation 测试简单三角形的三角剖分
func TestSimpleTriangleTessellation(t *testing.T) {
	// 定义一个简单的三角形顶点 (0,0), (1,0), (0,1)
	triangleVertices := []Vertex{
		{X: 0, Y: 0, Z: 0},
		{X: 1, Y: 0, Z: 0},
		{X: 0, Y: 1, Z: 0},
	}

	// 创建轮廓
	contour := Contour(triangleVertices)
	contours := []Contour{contour}

	// 执行三角剖分
	indices, vertices, err := Tesselate(contours, WindingRuleOdd)

	// 输出调试信息
	fmt.Printf("TestSimpleTriangleTessellation: 错误: %v\n", err)
	fmt.Printf("TestSimpleTriangleTessellation: 顶点数量: %d\n", len(vertices))
	fmt.Printf("TestSimpleTriangleTessellation: 索引数量: %d\n", len(indices))
	fmt.Printf("TestSimpleTriangleTessellation: 顶点: %v\n", vertices)
	fmt.Printf("TestSimpleTriangleTessellation: 索引: %v\n", indices)

	// 检查是否有错误
	if err != nil {
		t.Errorf("Tessellate returned an error: %v", err)
	}

	// 对于一个三角形，我们期望:
	// 1. 3个顶点 (输入的3个顶点)
	// 2. 3个索引 (1个三角形，每个三角形3个顶点)

	// 检查顶点数量
	if len(vertices) != 3 {
		t.Errorf("Expected 3 vertices, got %d", len(vertices))
	}

	// 检查索引数量
	if len(indices) != 3 {
		t.Errorf("Expected 3 indices (1 triangle), got %d", len(indices))
	}

	// 验证顶点坐标
	expectedVertices := []Vertex{
		{X: 0, Y: 0, Z: 0},
		{X: 1, Y: 0, Z: 0},
		{X: 0, Y: 1, Z: 0},
	}

	for i, expected := range expectedVertices {
		if i < len(vertices) {
			if vertices[i].X != expected.X || vertices[i].Y != expected.Y || vertices[i].Z != expected.Z {
				t.Errorf("Vertex %d: expected (%.1f, %.1f, %.1f), got (%.1f, %.1f, %.1f)",
					i, expected.X, expected.Y, expected.Z, vertices[i].X, vertices[i].Y, vertices[i].Z)
			}
		}
	}

	// 验证索引范围
	for i, idx := range indices {
		if idx < 0 || idx >= len(vertices) {
			t.Errorf("Index %d is out of range: %d (should be 0-%d)", i, idx, len(vertices)-1)
		}
	}

	// 验证三角形构成
	if len(indices) >= 3 {
		fmt.Printf("Triangle: [%d, %d, %d]\n", indices[0], indices[1], indices[2])
	}
}
