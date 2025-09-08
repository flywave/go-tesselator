package tesselator

import (
	"fmt"
	"math"
	"math/rand"
	"os"
)

// 生成SVG可视化
func GenerateSVG(filename string, contours []Contour, vertices []Vertex, elements []int) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// 计算边界框
	minX, minY, maxX, maxY := computeBounds(contours)
	padding := float32(10.0)
	width := maxX - minX + 2*padding
	height := maxY - minY + 2*padding

	// SVG头部
	fmt.Fprintf(file, `<svg xmlns="http://www.w3.org/2000/svg" width="%f" height="%f" viewBox="%f %f %f %f">`,
		width, height, minX-padding, minY-padding, width, height)
	fmt.Fprintln(file)

	// 绘制三角形
	fmt.Fprintln(file, `  <g stroke="black" stroke-width="1">`)
	for i := 0; i < len(elements); i += 3 {
		if elements[i] == -1 || elements[i+1] == -1 || elements[i+2] == -1 {
			continue // 跳过无效三角形
		}

		v1 := vertices[elements[i]]
		v2 := vertices[elements[i+1]]
		v3 := vertices[elements[i+2]]

		// 随机颜色
		color := fmt.Sprintf("#%02x%02x%02x",
			rand.Intn(256), rand.Intn(256), rand.Intn(256))

		fmt.Fprintf(file, `    <polygon points="%f,%f %f,%f %f,%f" fill="%s" fill-opacity="0.6"/>`,
			v1.X, v1.Y, v2.X, v2.Y, v3.X, v3.Y, color)
		fmt.Fprintln(file)
	}
	fmt.Fprintln(file, "  </g>")

	// 绘制原始轮廓
	fmt.Fprintln(file, `  <g stroke="blue" stroke-width="2" fill="none">`)
	for _, contour := range contours {
		fmt.Fprint(file, `    <path d="`)
		for i, vertex := range contour {
			if i == 0 {
				fmt.Fprintf(file, "M %f %f ", vertex.X, vertex.Y)
			} else {
				fmt.Fprintf(file, "L %f %f ", vertex.X, vertex.Y)
			}
		}
		fmt.Fprintln(file, `Z"/>`)
	}
	fmt.Fprintln(file, "  </g>")

	// SVG尾部
	fmt.Fprintln(file, "</svg>")

	return nil
}

// 生成正多边形
func GenerateRegularPolygon(sides int, cx, cy, radius float32) Contour {
	contour := Contour{}
	for i := 0; i < sides; i++ {
		angle := 2 * math.Pi * float64(i) / float64(sides)
		x := cx + radius*float32(math.Cos(angle))
		y := cy + radius*float32(math.Sin(angle))
		contour = append(contour, Vertex{X: x, Y: y})
	}
	return contour
}

// 生成星形
func GenerateStar(points int, cx, cy, outerRadius, innerRadius float32) Contour {
	contour := Contour{}
	for i := 0; i < points*2; i++ {
		angle := math.Pi * float64(i) / float64(points)
		radius := outerRadius
		if i%2 == 1 {
			radius = innerRadius
		}
		x := cx + radius*float32(math.Cos(angle))
		y := cy + radius*float32(math.Sin(angle))
		contour = append(contour, Vertex{X: x, Y: y})
	}
	return contour
}

// 生成随机多边形
func GenerateRandomPolygon(points int, cx, cy, radius float32) Contour {
	contour := Contour{}
	angles := make([]float64, points)
	for i := 0; i < points; i++ {
		angles[i] = rand.Float64() * 2 * math.Pi
	}

	// 排序角度以确保多边形是简单的
	for i := 0; i < points-1; i++ {
		for j := i + 1; j < points; j++ {
			if angles[i] > angles[j] {
				angles[i], angles[j] = angles[j], angles[i]
			}
		}
	}

	for i := 0; i < points; i++ {
		r := radius * (0.8 + 0.4*rand.Float32()) // 添加一些随机性
		x := cx + r*float32(math.Cos(angles[i]))
		y := cy + r*float32(math.Sin(angles[i]))
		contour = append(contour, Vertex{X: x, Y: y})
	}
	return contour
}

// 计算边界框
func computeBounds(contours []Contour) (minX, minY, maxX, maxY float32) {
	minX, minY = float32(math.MaxFloat32), float32(math.MaxFloat32)
	maxX, maxY = float32(-math.MaxFloat32), float32(-math.MaxFloat32)

	for _, contour := range contours {
		for _, vertex := range contour {
			if vertex.X < minX {
				minX = vertex.X
			}
			if vertex.Y < minY {
				minY = vertex.Y
			}
			if vertex.X > maxX {
				maxX = vertex.X
			}
			if vertex.Y > maxY {
				maxY = vertex.Y
			}
		}
	}

	return minX, minY, maxX, maxY
}

func TessellateAndGenerateSVG(filename string, contours []Contour) error {
	// 运行三角剖分
	elements, vertices, err := Tesselate(contours, WindingRuleOdd)
	if err != nil {
		return err
	}

	// 生成SVG可视化
	svgFilename := fmt.Sprintf("test_%s.svg", filename)
	if err := GenerateSVG(svgFilename, contours, vertices, elements); err != nil {
		return err
	}
	return nil
}
