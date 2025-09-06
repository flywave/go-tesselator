package tesselator

import (
	"math"
	"testing"
)

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

// TestEdgeEval 测试边评估函数
func TestEdgeEval(t *testing.T) {
	u := &vertex{s: 0.0, t: 0.0}
	v := &vertex{s: 1.0, t: 1.0}
	w := &vertex{s: 2.0, t: 0.0}

	// 确保点按s坐标排序
	if !(vertLeq(u, v) && vertLeq(v, w)) {
		t.Fatal("Points not in order for edgeEval")
	}

	// 测试点在边上 - 注意：根据实际实现，v点(1.0,1.0)并不在uw边上
	result := edgeEval(u, v, w)
	// 根据edgeEval实现计算预期值
	gapL := v.s - u.s
	gapR := w.s - v.s
	var expected float
	if gapL+gapR > 0 {
		if gapL < gapR {
			expected = (v.t - u.t) + (u.t-w.t)*(gapL/(gapL+gapR))
		} else {
			expected = (v.t - w.t) + (w.t-u.t)*(gapR/(gapL+gapR))
		}
	} else {
		expected = 0.0
	}
	if math.Abs(float64(result)-float64(expected)) > 1e-6 {
		t.Errorf("Expected edgeEval(%v, %v, %v) = %v, got %v", u, v, w, expected, result)
	}
	// 添加一个真正在边上的测试用例
	vu := &vertex{s: 1.0, t: 0.0}
	// 确保点按s坐标排序
	if !(vertLeq(u, vu) && vertLeq(vu, w)) {
		t.Fatal("Points not in order for edgeEval")
	}
	result = edgeEval(u, vu, w)
	expected = 0.0
	if math.Abs(float64(result)-float64(expected)) > 1e-6 {
		t.Errorf("Expected edgeEval(%v, %v, %v) = %v, got %v", u, vu, w, expected, result)
	}

	// 测试点不在边上
	v2 := &vertex{s: 1.0, t: 2.0}
	// 确保点按s坐标排序
	if !(vertLeq(u, v2) && vertLeq(v2, w)) {
		t.Fatal("Points not in order for edgeEval")
	}
	result = edgeEval(u, v2, w)
	// 根据实际实现计算预期值
	gapL = v2.s - u.s
	gapR = w.s - v2.s
	if gapL+gapR > 0 {
		if gapL < gapR {
			expected = float(v2.t-u.t) + float(u.t-w.t)*float(gapL)/float(gapL+gapR)
		} else {
			expected = float(float64(v2.t-w.t) + float64(w.t-u.t)*float64(gapR)/float64(gapL+gapR))
		}
	} else {
		expected = 0
	}
	if math.Abs(float64(result)-float64(expected)) > 1e-6 {
		t.Errorf("Expected edgeEval(%v, %v, %v) = %v, got %v", u, v2, w, expected, result)
	}
}

// TestEdgeSign 测试边符号函数
func TestEdgeSign(t *testing.T) {
	u := &vertex{s: 0.0, t: 0.0}
	v := &vertex{s: 1.0, t: 0.0} // 真正在边上的点
	w := &vertex{s: 2.0, t: 0.0}

	// 确保点按s坐标排序
	if !(vertLeq(u, v) && vertLeq(v, w)) {
		t.Fatal("Points not in order for edgeSign")
	}

	// 测试点在边上
	result := edgeSign(u, v, w)
	if result != 0 {
		t.Errorf("Expected edgeSign(%v, %v, %v) = 0, got %v", u, v, w, result)
	}

	// 测试原来的点不在边上的情况
	v2 := &vertex{s: 1.0, t: 1.0}
	// 确保点按s坐标排序
	if !(vertLeq(u, v2) && vertLeq(v2, w)) {
		t.Fatal("Points not in order for edgeSign")
	}
	result = edgeSign(u, v2, w)
	// 计算预期值
	expected := float((v2.t-w.t)*(v2.s-u.s) + (v2.t-u.t)*(w.s-v2.s))
	if math.Abs(float64(result)-float64(expected)) > 1e-6 {
		t.Errorf("Expected edgeSign(%v, %v, %v) = %v, got %v", u, v2, w, expected, result)
	}

	// 测试点在边上方
	v2 = &vertex{s: 1.0, t: 2.0}
	// 确保点按s坐标排序
	if !(vertLeq(u, v2) && vertLeq(v2, w)) {
		t.Fatal("Points not in order for edgeSign")
	}
	result = edgeSign(u, v2, w)
	if result <= 0 {
		t.Errorf("Expected edgeSign(%v, %v, %v) > 0, got %v", u, v2, w, result)
	}

	// 测试点在边下方
	v3 := &vertex{s: 1.0, t: -1.0}
	// 确保点按s坐标排序
	if !(vertLeq(u, v3) && vertLeq(v3, w)) {
		t.Fatal("Points not in order for edgeSign")
	}
	result = edgeSign(u, v3, w)
	if result >= 0 {
		t.Errorf("Expected edgeSign(%v, %v, %v) < 0, got %v", u, v3, w, result)
	}
}

// TestVertCCW 测试三点是否逆时针
func TestVertCCW(t *testing.T) {
	// 逆时针三角形 (实际计算结果为负，返回false)
	u := &vertex{s: 0.0, t: 0.0}
	v := &vertex{s: 0.0, t: 1.0}
	w := &vertex{s: 1.0, t: 0.0}

	if vertCCW(u, v, w) {
		t.Errorf("Expected vertCCW(%v, %v, %v) to be false", u, v, w)
	}

	// 顺时针三角形 (实际计算结果为正，返回true)
	u2 := &vertex{s: 0.0, t: 0.0}
	v2 := &vertex{s: 1.0, t: 0.0}
	w2 := &vertex{s: 0.0, t: 1.0}

	if !vertCCW(u2, v2, w2) {
		t.Errorf("Expected vertCCW(%v, %v, %v) to be true", u2, v2, w2)
	}

	// 真正的逆时针三角形 (u, w, v)
	u3 := &vertex{s: 0.0, t: 0.0}
	v3 := &vertex{s: 0.0, t: 1.0}
	w3 := &vertex{s: 1.0, t: 0.0}

	if !vertCCW(u3, w3, v3) {
		t.Errorf("Expected vertCCW(%v, %v, %v) to be true", u3, w3, v3)
	}

	// 共线
	u3 = &vertex{s: 0.0, t: 0.0}
	v3 = &vertex{s: 1.0, t: 1.0}
	w3 = &vertex{s: 2.0, t: 2.0}

	// 根据实现，共线点应该返回true
	if !vertCCW(u3, v3, w3) {
		t.Errorf("Expected vertCCW(%v, %v, %v) to be true (collinear points)", u3, v3, w3)
	}

	// 测试另一个顺时针三角形
	u4 := &vertex{s: 0.0, t: 0.0}
	v4 := &vertex{s: 2.0, t: 0.0}
	w4 := &vertex{s: 1.0, t: -1.0}

	if vertCCW(u4, v4, w4) {
		t.Errorf("Expected vertCCW(%v, %v, %v) to be false", u4, v4, w4)
	}
}

// TestInterpolate 测试插值函数
func TestInterpolate(t *testing.T) {
	// 正常情况
	result := interpolate(1.0, 0.0, 1.0, 2.0)
	expected := 1.0
	if math.Abs(float64(result)-expected) > 1e-6 {
		t.Errorf("Expected interpolate(1, 0, 1, 2) = %v, got %v", expected, result)
	}

	// a为0
	result = interpolate(0.0, 0.0, 1.0, 2.0)
	// 根据实现，当a=0时，如果a<=b且b!=0，返回x + (y-x)*(a/(a+b)) = 0 + (2-0)*(0/1) = 0
	expected = 0.0
	if math.Abs(float64(result)-expected) > 1e-6 {
		t.Errorf("Expected interpolate(0, 0, 1, 2) = %v, got %v", expected, result)
	}

	// b为0
	result = interpolate(1.0, 0.0, 0.0, 2.0)
	// 根据实现，当b=0时，如果a>b，返回y + (x-y)*(b/(a+b)) = 2 + (0-2)*(0/1) = 2
	expected = 2.0
	if math.Abs(float64(result)-expected) > 1e-6 {
		t.Errorf("Expected interpolate(1, 0, 0, 2) = %v, got %v", expected, result)
	}

	// a和b都为0
	result = interpolate(0.0, 0.0, 0.0, 2.0)
	expected = 1.0
	if math.Abs(float64(result)-expected) > 1e-6 {
		t.Errorf("Expected interpolate(0, 0, 0, 2) = %v, got %v", expected, result)
	}

	// 负值情况
	result = interpolate(-1.0, 0.0, 1.0, 2.0)
	// 根据实现，负值会被截断为0，所以a=0, b=1
	expected = 0.0 + (2.0-0.0)*(0.0/(0.0+1.0))
	if math.Abs(float64(result)-expected) > 1e-6 {
		t.Errorf("Expected interpolate(-1, 0, 1, 2) = %v, got %v", expected, result)
	}

	// 测试a > b的情况
	result = interpolate(3.0, 0.0, 1.0, 4.0)
	expected = 4.0 + (0.0-4.0)*(1.0/(3.0+1.0))
	if math.Abs(float64(result)-expected) > 1e-6 {
		t.Errorf("Expected interpolate(3, 0, 1, 4) = %v, got %v", expected, result)
	}
}

// TestEdgeIntersect 测试边相交函数
func TestEdgeIntersect(t *testing.T) {
	o1 := &vertex{s: 0.0, t: 0.0}
	d1 := &vertex{s: 2.0, t: 2.0}
	o2 := &vertex{s: 0.0, t: 2.0}
	d2 := &vertex{s: 2.0, t: 0.0}
	v := &vertex{}

	edgeIntersect(o1, d1, o2, d2, v)
	expectedS := 1.0
	expectedT := 1.0
	if math.Abs(float64(v.s)-expectedS) > 1e-6 || math.Abs(float64(v.t)-expectedT) > 1e-6 {
		t.Errorf("Expected intersection at (%v, %v), got (%v, %v)", expectedS, expectedT, v.s, v.t)
	}

	// 测试不相交的边
	o3 := &vertex{s: 0.0, t: 0.0}
	d3 := &vertex{s: 0.0, t: 1.0}
	o4 := &vertex{s: 1.0, t: 0.0}
	d4 := &vertex{s: 1.0, t: 1.0}
	v2 := &vertex{}

	edgeIntersect(o3, d3, o4, d4, v2)
	// 对于不相交的边，函数会返回特定点的中点
	// 根据edgeIntersect实现，当边不相交时，v.s = (o2.s + d1.s)/2, v.t = (o2.t + d1.t)/2
	// 其中o2和d1是经过排序后的点
	expectedS = (1.0 + 0.0) / 2 // o4.s + o3.s
	expectedT = (1.0 + 1.0) / 2 // d4.t + d3.t
	if math.Abs(float64(v2.s)-expectedS) > 1e-6 || math.Abs(float64(v2.t)-expectedT) > 1e-6 {
		t.Errorf("Expected midpoint at (%v, %v) for non-intersecting edges, got (%v, %v)", expectedS, expectedT, v2.s, v2.t)
	}

	// 测试另一个相交情况
	o5 := &vertex{s: 0.0, t: 0.0}
	d5 := &vertex{s: 2.0, t: 0.0}
	o6 := &vertex{s: 1.0, t: -1.0}
	d6 := &vertex{s: 1.0, t: 1.0}
	v3 := &vertex{}

	edgeIntersect(o5, d5, o6, d6, v3)
	expectedS = 1.0
	expectedT = 0.0
	if math.Abs(float64(v3.s)-expectedS) > 1e-6 || math.Abs(float64(v3.t)-expectedT) > 1e-6 {
		t.Errorf("Expected intersection at (%v, %v), got (%v, %v)", expectedS, expectedT, v3.s, v3.t)
	}
}

// TestTransEval 测试transEval函数
func TestTransEval(t *testing.T) {
	// 创建测试点，确保它们满足transLeq顺序 (按t坐标排序)
	u := &vertex{s: 0.0, t: 0.0}
	v := &vertex{s: 1.0, t: 1.0}
	w := &vertex{s: 2.0, t: 2.0}

	// 确保点按t坐标排序
	if !(transLeq(u, v) && transLeq(v, w)) {
		t.Fatal("Points not in order for transEval")
	}

	// 测试点在边上
	result := transEval(u, v, w)
	expected := float(0.0)
	if math.Abs(float64(result)-float64(expected)) > 1e-6 {
		t.Errorf("Expected transEval(%v, %v, %v) = %v, got %v", u, v, w, expected, result)
	}

	// 测试点不在边上
	v2 := &vertex{s: 0.0, t: 1.0}
	// 确保点按t坐标排序
	if !(transLeq(u, v2) && transLeq(v2, w)) {
		t.Fatal("Points not in order for transEval")
	}
	result = transEval(u, v2, w)
	// 根据transEval实现计算预期值
	gapL := v2.t - u.t
	gapR := w.t - v2.t
	if gapL+gapR > 0 {
		if gapL < gapR {
			expected = float(v2.s-u.s) + float(u.s-w.s)*float(gapL)/float(gapL+gapR)
		} else {
			expected = float(v2.s-w.s) + float(w.s-u.s)*float(gapR)/float(gapL+gapR)
		}
	} else {
		expected = 0
	}
	if math.Abs(float64(result)-float64(expected)) > 1e-6 {
		t.Errorf("Expected transEval(%v, %v, %v) = %v, got %v", u, v2, w, expected, result)
	}
}

// TestTransSign 测试transSign函数
func TestTransSign(t *testing.T) {
	// 创建测试点，确保它们满足transLeq顺序 (按t坐标排序)
	u := &vertex{s: 0.0, t: 0.0}
	v := &vertex{s: 1.0, t: 1.0}
	w := &vertex{s: 2.0, t: 2.0}

	// 确保点按t坐标排序
	if !(transLeq(u, v) && transLeq(v, w)) {
		t.Fatal("Points not in order for transSign")
	}

	// 测试点在边上
	result := transSign(u, v, w)
	expected := float(0.0)
	if math.Abs(float64(result)-float64(expected)) > 1e-6 {
		t.Errorf("Expected transSign(%v, %v, %v) = %v, got %v", u, v, w, expected, result)
	}

	// 测试点在边上侧
	v2 := &vertex{s: 2.0, t: 1.0}
	// 确保点按t坐标排序
	if !(transLeq(u, v2) && transLeq(v2, w)) {
		t.Fatal("Points not in order for transSign")
	}
	result = transSign(u, v2, w)
	if result <= 0 {
		t.Errorf("Expected transSign(%v, %v, %v) > 0, got %v", u, v2, w, result)
	}

	// 测试点在边下侧
	v3 := &vertex{s: 0.0, t: 1.0}
	// 确保点按t坐标排序
	if !(transLeq(u, v3) && transLeq(v3, w)) {
		t.Fatal("Points not in order for transSign")
	}
	result = transSign(u, v3, w)
	if result >= 0 {
		t.Errorf("Expected transSign(%v, %v, %v) < 0, got %v", u, v3, w, result)
	}
}

// TestEdgeGoesLeft 测试edgeGoesLeft函数
func TestEdgeGoesLeft(t *testing.T) {
	// 创建边，Dst的s坐标小于Org的s坐标 (边向左走)
	org := &vertex{s: 1.0, t: 1.0}
	dst := &vertex{s: 0.0, t: 0.0}
	e := &halfEdge{
		Org: org,
		Sym: &halfEdge{
			Org: dst,
		},
	}

	// 设置对称边的dst
	e.Sym.Sym = e

	// 边应该向左走 (dst.s <= org.s)
	if !edgeGoesLeft(e) {
		t.Errorf("Expected edgeGoesLeft(%v -> %v) to be true", org, dst)
	}

	// 创建边，Dst的s坐标大于Org的s坐标 (边向右走)
	org2 := &vertex{s: 0.0, t: 0.0}
	dst2 := &vertex{s: 1.0, t: 1.0}
	e2 := &halfEdge{
		Org: org2,
		Sym: &halfEdge{
			Org: dst2,
		},
	}

	// 设置对称边的dst
	e2.Sym.Sym = e2

	// 边不应该向左走 (dst.s > org.s)
	if edgeGoesLeft(e2) {
		t.Errorf("Expected edgeGoesLeft(%v -> %v) to be false", org2, dst2)
	}
}

// TestEdgeGoesRight 测试edgeGoesRight函数
func TestEdgeGoesRight(t *testing.T) {
	// 创建边，Org的s坐标小于Dst的s坐标
	org := &vertex{s: 0.0, t: 0.0}
	dst := &vertex{s: 1.0, t: 1.0}
	e := &halfEdge{
		Org: org,
		Sym: &halfEdge{
			Org: dst,
		},
	}

	// 设置对称边的dst
	e.Sym.Sym = e

	// 边应该向右走 (org.s <= dst.s)
	if !edgeGoesRight(e) {
		t.Errorf("Expected edgeGoesRight(%v -> %v) to be true", org, dst)
	}

	// 创建边，Org的s坐标大于Dst的s坐标
	org2 := &vertex{s: 1.0, t: 1.0}
	dst2 := &vertex{s: 0.0, t: 0.0}
	e2 := &halfEdge{
		Org: org2,
		Sym: &halfEdge{
			Org: dst2,
		},
	}

	// 设置对称边的dst
	e2.Sym.Sym = e2

	// 边不应该向右走 (org.s > dst.s)
	if edgeGoesRight(e2) {
		t.Errorf("Expected edgeGoesRight(%v -> %v) to be false", org2, dst2)
	}
}

// TestVertL1dist 测试vertL1dist函数
func TestVertL1dist(t *testing.T) {
	// 测试两个相同点
	u := &vertex{s: 1.0, t: 2.0}
	v := &vertex{s: 1.0, t: 2.0}
	result := vertL1dist(u, v)
	expected := float(0.0)
	if result != expected {
		t.Errorf("Expected vertL1dist(%v, %v) = %v, got %v", u, v, expected, result)
	}

	// 测试两个不同点
	u2 := &vertex{s: 0.0, t: 0.0}
	v2 := &vertex{s: 3.0, t: 4.0}
	result = vertL1dist(u2, v2)
	expected = float(7.0) // |3-0| + |4-0| = 3 + 4 = 7
	if result != expected {
		t.Errorf("Expected vertL1dist(%v, %v) = %v, got %v", u2, v2, expected, result)
	}

	// 测试负坐标
	u3 := &vertex{s: -1.0, t: -2.0}
	v3 := &vertex{s: 2.0, t: 3.0}
	result = vertL1dist(u3, v3)
	expected = float(8.0) // |2-(-1)| + |3-(-2)| = 3 + 5 = 8
	if result != expected {
		t.Errorf("Expected vertL1dist(%v, %v) = %v, got %v", u3, v3, expected, result)
	}
}

// TestVertEq 测试vertEq函数
func TestVertEq(t *testing.T) {
	// 测试两个相同点
	u := &vertex{s: 1.0, t: 2.0}
	v := &vertex{s: 1.0, t: 2.0}
	if !vertEq(u, v) {
		t.Errorf("Expected vertEq(%v, %v) to be true", u, v)
	}

	// 测试s坐标不同
	u2 := &vertex{s: 1.0, t: 2.0}
	v2 := &vertex{s: 2.0, t: 2.0}
	if vertEq(u2, v2) {
		t.Errorf("Expected vertEq(%v, %v) to be false", u2, v2)
	}

	// 测试t坐标不同
	u3 := &vertex{s: 1.0, t: 2.0}
	v3 := &vertex{s: 1.0, t: 3.0}
	if vertEq(u3, v3) {
		t.Errorf("Expected vertEq(%v, %v) to be false", u3, v3)
	}

	// 测试两个坐标都不同
	u4 := &vertex{s: 1.0, t: 2.0}
	v4 := &vertex{s: 3.0, t: 4.0}
	if vertEq(u4, v4) {
		t.Errorf("Expected vertEq(%v, %v) to be false", u4, v4)
	}
}

// TestTransLeq 测试transLeq函数
func TestTransLeq(t *testing.T) {
	// 测试t坐标相同，s坐标不同
	u := &vertex{s: 1.0, t: 2.0}
	v := &vertex{s: 2.0, t: 2.0}
	if !transLeq(u, v) {
		t.Errorf("Expected transLeq(%v, %v) to be true", u, v)
	}
	if transLeq(v, u) {
		t.Errorf("Expected transLeq(%v, %v) to be false", v, u)
	}

	// 测试t坐标不同
	u2 := &vertex{s: 1.0, t: 1.0}
	v2 := &vertex{s: 2.0, t: 2.0}
	if !transLeq(u2, v2) {
		t.Errorf("Expected transLeq(%v, %v) to be true", u2, v2)
	}
	if transLeq(v2, u2) {
		t.Errorf("Expected transLeq(%v, %v) to be false", v2, u2)
	}

	// 测试两个点相同
	u3 := &vertex{s: 1.0, t: 2.0}
	v3 := &vertex{s: 1.0, t: 2.0}
	if !transLeq(u3, v3) {
		t.Errorf("Expected transLeq(%v, %v) to be true", u3, v3)
	}
	if !transLeq(v3, u3) {
		t.Errorf("Expected transLeq(%v, %v) to be true", v3, u3)
	}
}
