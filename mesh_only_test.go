package tesselator

import (
	"testing"
)

// TestMeshOnly tests only the mesh functionality without the rest of the tesselator
func TestMeshOnly(t *testing.T) {
	// Test basic mesh creation
	mesh := tessMeshNewMesh()
	if mesh == nil {
		t.Fatal("tessMeshNewMesh returned nil")
	}

	// Test edge creation
	edge := tessMeshMakeEdge(mesh)
	if edge == nil {
		t.Fatal("tessMeshMakeEdge returned nil")
	}

	// Test that the edge has the correct properties
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

	// Test vertex count
	vertexCount := 0
	for v := mesh.vHead.next; v != &mesh.vHead; v = v.next {
		vertexCount++
	}
	if vertexCount != 2 {
		t.Errorf("Expected 2 vertices, got %d", vertexCount)
	}

	// Test face count
	faceCount := 0
	for f := mesh.fHead.next; f != &mesh.fHead; f = f.next {
		faceCount++
	}
	if faceCount != 1 {
		t.Errorf("Expected 1 face, got %d", faceCount)
	}

	// Test navigation
	if edge.Sym.Sym != edge {
		t.Error("Symmetry of symmetry should be the original edge")
	}

	// Test Onext - for a single edge loop, Onext should point to the symmetric edge
	if edge.Onext != edge.Sym {
		t.Error("Onext should point to symmetric edge for a single edge loop")
	}

	// Test Lnext - for a single edge loop, Lnext should point to the symmetric edge
	if edge.Lnext != edge.Sym {
		t.Error("Lnext should point to symmetric edge for a single edge loop")
	}

	// Test vertex operations
	// edge.Org should be the first vertex
	vertex1 := edge.Org
	if vertex1.anEdge != edge {
		t.Errorf("Vertex1 anEdge should point to the edge, but vertex1.anEdge=%v and edge=%v", vertex1.anEdge, edge)
	}

	// edge.Sym.Org should be the second vertex
	vertex2 := edge.Sym.Org
	if vertex2.anEdge != edge.Sym {
		t.Errorf("Vertex2 anEdge should point to the symmetric edge, but vertex2.anEdge=%v and edge.Sym=%v", vertex2.anEdge, edge.Sym)
	}

	// Test face operations
	face := edge.Lface
	if face.anEdge != edge {
		t.Error("Face anEdge should point to the edge")
	}
}
