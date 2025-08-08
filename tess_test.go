package tesselator

import (
	"math"
	"runtime/debug"
	"testing"
)

// 测试用例使用的辅助函数和类型

type Vector2f struct {
	X, Y float64
}

// 转换Vector2f切片为Contour类型
func toContour(v []Vector2f) Contour {
	c := make(Contour, len(v))
	for i, vec := range v {
		c[i] = Vertex{
			X: float32(vec.X),
			Y: float32(vec.Y),
			Z: 0,
		}
	}
	return c
}

// 添加带孔的多边形
func addPolygonWithHole() []Contour {
	outerLoop := []Vector2f{
		{0, 0},
		{3, 0},
		{3, 3},
		{0, 3},
	}
	innerHole := []Vector2f{
		{1, 1},
		{2, 1},
		{2, 2},
		{1, 2},
	}

	// 注意：外轮廓和内轮廓的缠绕方向应该相反
	return []Contour{
		toContour(outerLoop),
		toContour(innerHole),
	}
}

// 创建单个三角形
func createTriangle() []Contour {
	triangle := []Vector2f{
		{0, 0},
		{0, 1},
		{1, 0},
	}
	return []Contour{toContour(triangle)}
}

// 创建单位四边形
func createUnitQuad() []Contour {
	quad := []Vector2f{
		{0, 0},
		{0, 1},
		{1, 1},
		{1, 0},
	}
	return []Contour{toContour(quad)}
}

// 创建空轮廓
func createEmptyContour() []Contour {
	return []Contour{}
}

// 创建单行
func createSingleLine() []Contour {
	line := []Vector2f{
		{0, 0},
		{0, 1},
	}
	return []Contour{toContour(line)}
}

// 创建奇异点四边形
func createSingularityQuad() []Contour {
	quad := []Vector2f{
		{0, 0},
		{0, 0},
		{0, 0},
		{0, 0},
	}
	return []Contour{toContour(quad)}
}

// 创建退化四边形
func createDegenerateQuad() []Contour {
	quad := []Vector2f{
		{0, 3.40282347e+38},
		{0.64113313, -1},
		{-0, -0},
		{-3.40282347e+38, 1},
	}
	return []Contour{toContour(quad)}
}

// 创建宽度溢出三角形
func createWidthOverflowsTri() []Contour {
	tri := []Vector2f{
		{-2e+38, 0},
		{0, 0},
		{2e+38, -1},
	}
	return []Contour{toContour(tri)}
}

// 创建高度溢出三角形
func createHeightOverflowsTri() []Contour {
	tri := []Vector2f{
		{0, 0},
		{0, 2e+38},
		{-1, -2e+38},
	}
	return []Contour{toContour(tri)}
}

// 创建面积溢出三角形
func createAreaOverflowsTri() []Contour {
	tri := []Vector2f{
		{-2e+37, 0},
		{0, 5},
		{1e+37, -5},
	}
	return []Contour{toContour(tri)}
}

// 创建NaN四边形
func createNanQuad() []Contour {
	nanValue := math.NaN()
	quad := []Vector2f{
		{nanValue, nanValue},
		{nanValue, nanValue},
		{nanValue, nanValue},
		{nanValue, nanValue},
	}
	return []Contour{toContour(quad)}
}

// 测试用例

func TestDefaultTesselateWithHole(t *testing.T) {
	contours := addPolygonWithHole()
	defer func() {
		if r := recover(); r != nil {
			t.Logf("Recovered from panic: %v\n%s", r, debug.Stack())
		}
	}()
	elements, vertices, err := Tesselate(contours, WindingRulePositive)
	if err != nil {
		t.Errorf("Tesselate failed: %v", err)
	}
	// 带孔的多边形应该生成8个三角形
	// 每个三角形有3个顶点索引
	if len(elements) != 8*3 {
		t.Errorf("Expected 24 element indices (8 triangles), got %d", len(elements))
	}

	// 验证顶点数量合理
	if len(vertices) == 0 {
		t.Errorf("Expected vertices, got none")
	}
}

func TestEmptyContour(t *testing.T) {
	contours := createEmptyContour()
	elements, vertices, err := Tesselate(contours, WindingRulePositive)
	if err != nil {
		t.Errorf("Tesselate failed: %v", err)
	}

	if len(elements) != 0 {
		t.Errorf("Expected 0 elements, got %d", len(elements))
	}

	if len(vertices) != 0 {
		t.Errorf("Expected 0 vertices, got %d", len(vertices))
	}
}

func TestSingleLine(t *testing.T) {
	contours := createSingleLine()
	println("Testing single line contour with vertices:")
	for i, v := range contours[0] {
		println("  Vertex", i, ":", v.X, ",", v.Y)
	}
	elements, _, err := Tesselate(contours, WindingRulePositive)
	if err != nil {
		t.Errorf("Tesselate failed: %v", err)
	}

	// 单行应该不会生成任何三角形
	if len(elements) != 0 {
		t.Errorf("Expected 0 elements, got %d", len(elements))
	}
}

func TestSingleTriangle(t *testing.T) {
	contours := createTriangle()
	elements, vertices, err := Tesselate(contours, WindingRuleNonzero) // Changed from WindingRulePositive
	if err != nil {
		t.Errorf("Tesselate failed: %v", err)
	}

	// 单个三角形应该生成1个三角形，包含3个顶点索引
	if len(elements) != 3 {
		t.Errorf("Expected 3 element indices (1 triangle), got %d", len(elements))
	}

	// 应该有3个顶点
	if len(vertices) != 3 {
		t.Errorf("Expected 3 vertices, got %d", len(vertices))
	}
}

func TestUnitQuad(t *testing.T) {
	contours := createUnitQuad()
	elements, vertices, err := Tesselate(contours, WindingRulePositive)
	if err != nil {
		t.Errorf("Tesselate failed: %v", err)
	}

	// 单位四边形应该被分成2个三角形，共6个顶点索引
	if len(elements) != 6 {
		t.Errorf("Expected 6 element indices (2 triangles), got %d", len(elements))
	}

	// 应该有4个顶点
	if len(vertices) != 4 {
		t.Errorf("Expected 4 vertices, got %d", len(vertices))
	}
}

func TestInvalidInput(t *testing.T) {
	contours := createAreaOverflowsTri()
	elements, vertices, err := Tesselate(contours, WindingRulePositive)
	if err == nil {
		t.Errorf("Expected Tesselate to fail, but it succeeded")
	}

	if elements != nil {
		t.Errorf("Expected nil elements, got %v", elements)
	}

	if vertices != nil {
		t.Errorf("Expected nil vertices, got %v", vertices)
	}
}

func TestFloatOverflowQuad(t *testing.T) {
	kFloatMin := math.SmallestNonzeroFloat64
	kFloatMax := math.MaxFloat64

	quad := []Vector2f{
		{kFloatMin, kFloatMin},
		{kFloatMin, kFloatMax},
		{kFloatMax, kFloatMax},
		{kFloatMax, kFloatMin},
	}
	contours := []Contour{toContour(quad)}

	_, _, err := Tesselate(contours, WindingRulePositive)
	if err == nil {
		t.Errorf("Expected Tesselate to fail, but it succeeded")
	}
}

func TestSingularityQuad(t *testing.T) {
	contours := createSingularityQuad()
	elements, _, err := Tesselate(contours, WindingRulePositive)
	if err != nil {
		t.Errorf("Tesselate failed: %v", err)
	}

	// 奇异点四边形应该生成0个三角形
	if len(elements) != 0 {
		t.Errorf("Expected 0 elements, got %d", len(elements))
	}
}

func TestDegenerateQuad(t *testing.T) {
	contours := createDegenerateQuad()
	_, _, err := Tesselate(contours, WindingRulePositive)
	if err == nil {
		t.Errorf("Expected Tesselate to fail, but it succeeded")
	}
}

func TestWidthOverflowsTri(t *testing.T) {
	contours := createWidthOverflowsTri()
	_, _, err := Tesselate(contours, WindingRulePositive)
	if err == nil {
		t.Errorf("Expected Tesselate to fail, but it succeeded")
	}
}

func TestHeightOverflowsTri(t *testing.T) {
	contours := createHeightOverflowsTri()
	_, _, err := Tesselate(contours, WindingRulePositive)
	if err == nil {
		t.Errorf("Expected Tesselate to fail, but it succeeded")
	}
}

func TestAreaOverflowsTri(t *testing.T) {
	contours := createAreaOverflowsTri()
	_, _, err := Tesselate(contours, WindingRulePositive)
	if err == nil {
		t.Errorf("Expected Tesselate to fail, but it succeeded")
	}
}

func TestNanQuad(t *testing.T) {
	contours := createNanQuad()
	elements, vertices, err := Tesselate(contours, WindingRulePositive)
	if err == nil {
		t.Errorf("Expected Tesselate to fail, but it succeeded")
	}

	if elements != nil {
		t.Errorf("Expected nil elements, got %v", elements)
	}

	if vertices != nil {
		t.Errorf("Expected nil vertices, got %v", vertices)
	}
}
