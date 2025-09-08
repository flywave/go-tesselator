package tesselator

import (
	"os"
	"testing"
)

// TestTessellateAndGenerateSVG 测试 TessellateAndGenerateSVG 函数
func TestTessellateAndGenerateSVG(t *testing.T) {
	// 创建一个简单的四边形轮廓
	quadVertices := []Vertex{
		{X: 0, Y: 0, Z: 0},
		{X: 100, Y: 0, Z: 0},
		{X: 100, Y: 100, Z: 0},
		{X: 0, Y: 100, Z: 0},
	}
	contour := Contour(quadVertices)
	contours := []Contour{contour}

	// 测试文件名
	testFilename := "quad_test"

	// 调用函数
	err := TessellateAndGenerateSVG(testFilename, contours)
	if err != nil {
		t.Fatalf("TessellateAndGenerateSVG failed: %v", err)
	}

	// 检查生成的SVG文件是否存在
	svgFilename := "test_" + testFilename + ".svg"
	if _, err := os.Stat(svgFilename); os.IsNotExist(err) {
		t.Errorf("SVG file was not created: %s", svgFilename)
	}

	// 保留生成的文件并移动到 assets 文件夹
	// os.Rename(svgFilename, "assets/"+svgFilename)
	// 注意：在实际测试中我们不移动文件，以免影响其他测试
	// 这里注释掉移动操作，但在 README 中展示时我们会手动移动文件
}

// TestTessellateAndGenerateSVGWithHole 测试带孔的多边形
func TestTessellateAndGenerateSVGWithHole(t *testing.T) {
	// 创建外部轮廓（正方形）
	outerVertices := []Vertex{
		{X: 0, Y: 0, Z: 0},
		{X: 100, Y: 0, Z: 0},
		{X: 100, Y: 100, Z: 0},
		{X: 0, Y: 100, Z: 0},
	}
	outerContour := Contour(outerVertices)

	// 创建内部孔洞（较小的正方形）
	innerVertices := []Vertex{
		{X: 25, Y: 25, Z: 0},
		{X: 75, Y: 25, Z: 0},
		{X: 75, Y: 75, Z: 0},
		{X: 25, Y: 75, Z: 0},
	}
	innerContour := Contour(innerVertices)

	contours := []Contour{outerContour, innerContour}

	// 测试文件名
	testFilename := "hole_test"

	// 调用函数
	err := TessellateAndGenerateSVG(testFilename, contours)
	if err != nil {
		t.Fatalf("TessellateAndGenerateSVG failed: %v", err)
	}

	// 检查生成的SVG文件是否存在
	svgFilename := "test_" + testFilename + ".svg"
	if _, err := os.Stat(svgFilename); os.IsNotExist(err) {
		t.Errorf("SVG file was not created: %s", svgFilename)
	}

	// 保留生成的文件并移动到 assets 文件夹
	// os.Rename(svgFilename, "assets/"+svgFilename)
}

// TestTessellateAndGenerateSVGWithMultipleContours 测试多个独立轮廓
func TestTessellateAndGenerateSVGWithMultipleContours(t *testing.T) {
	// 创建第一个轮廓（正方形）
	squareVertices := []Vertex{
		{X: 0, Y: 0, Z: 0},
		{X: 50, Y: 0, Z: 0},
		{X: 50, Y: 50, Z: 0},
		{X: 0, Y: 50, Z: 0},
	}
	squareContour := Contour(squareVertices)

	// 创建第二个轮廓（三角形）
	triangleVertices := []Vertex{
		{X: 60, Y: 0, Z: 0},
		{X: 110, Y: 0, Z: 0},
		{X: 85, Y: 50, Z: 0},
	}
	triangleContour := Contour(triangleVertices)

	contours := []Contour{squareContour, triangleContour}

	// 测试文件名
	testFilename := "multiple_contours_test"

	// 调用函数
	err := TessellateAndGenerateSVG(testFilename, contours)
	if err != nil {
		t.Fatalf("TessellateAndGenerateSVG failed: %v", err)
	}

	// 检查生成的SVG文件是否存在
	svgFilename := "test_" + testFilename + ".svg"
	if _, err := os.Stat(svgFilename); os.IsNotExist(err) {
		t.Errorf("SVG file was not created: %s", svgFilename)
	}

	// 保留生成的文件并移动到 assets 文件夹
	// os.Rename(svgFilename, "assets/"+svgFilename)
}

// TestTessellateAndGenerateSVGComplexPolygon 测试复杂多边形（五角星形）
func TestTessellateAndGenerateSVGComplexPolygon(t *testing.T) {
	// 创建一个复杂的星形轮廓
	starVertices := []Vertex{
		{X: 50, Y: 0, Z: 0},      // 顶点
		{X: 61.8, Y: 38.2, Z: 0}, // 右上点
		{X: 100, Y: 38.2, Z: 0},  // 右尖点
		{X: 69, Y: 61.8, Z: 0},   // 右下点
		{X: 80.9, Y: 100, Z: 0},  // 底部右点
		{X: 50, Y: 76.4, Z: 0},   // 底部中心点
		{X: 19.1, Y: 100, Z: 0},  // 底部左点
		{X: 31, Y: 61.8, Z: 0},   // 左下点
		{X: 0, Y: 38.2, Z: 0},    // 左尖点
		{X: 38.2, Y: 38.2, Z: 0}, // 左上点
	}
	starContour := Contour(starVertices)

	// 创建一个内部孔洞（小正方形）
	holeVertices := []Vertex{
		{X: 40, Y: 40, Z: 0},
		{X: 60, Y: 40, Z: 0},
		{X: 60, Y: 60, Z: 0},
		{X: 40, Y: 60, Z: 0},
	}
	holeContour := Contour(holeVertices)

	contours := []Contour{starContour, holeContour}

	// 测试文件名
	testFilename := "complex_star_test"

	// 调用函数
	err := TessellateAndGenerateSVG(testFilename, contours)
	if err != nil {
		t.Fatalf("TessellateAndGenerateSVG failed: %v", err)
	}

	// 检查生成的SVG文件是否存在
	svgFilename := "test_" + testFilename + ".svg"
	if _, err := os.Stat(svgFilename); os.IsNotExist(err) {
		t.Errorf("SVG file was not created: %s", svgFilename)
	}

	// 保留生成的文件并移动到 assets 文件夹
	// os.Rename(svgFilename, "assets/"+svgFilename)
}

// TestTessellateAndGenerateSVGRandomPolygon 测试随机生成的复杂多边形
func TestTessellateAndGenerateSVGRandomPolygon(t *testing.T) {
	// 设置随机种子以确保可重复性
	// rand.Seed(42) // 在Go 1.20中rand.Seed已弃用，使用随机种子

	// 生成一个随机多边形
	randomContour := GenerateRandomPolygon(12, 100, 100, 80) // 12个点，中心(100,100)，半径80

	// 生成一个内部孔洞
	holeContour := GenerateRandomPolygon(8, 100, 100, 30) // 8个点，中心(100,100)，半径30

	contours := []Contour{randomContour, holeContour}

	// 测试文件名
	testFilename := "random_polygon_test"

	// 调用函数
	err := TessellateAndGenerateSVG(testFilename, contours)
	if err != nil {
		t.Fatalf("TessellateAndGenerateSVG failed: %v", err)
	}

	// 检查生成的SVG文件是否存在
	svgFilename := "test_" + testFilename + ".svg"
	if _, err := os.Stat(svgFilename); os.IsNotExist(err) {
		t.Errorf("SVG file was not created: %s", svgFilename)
	}

	// 保留生成的文件并移动到 assets 文件夹
	// os.Rename(svgFilename, "assets/"+svgFilename)
}
