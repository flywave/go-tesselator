package tesselator

import (
	"runtime"
)

// The mesh operations below have three motivations: completeness,
// convenience, and efficiency.  The basic mesh operations are MakeEdge,
// Splice, and Delete.  All the other edge operations can be implemented
// in terms of these.  The other operations are provided for convenience
// and/or efficiency.
//
// When a face is split or a vertex is added, they are inserted into the
// global list //before// the existing vertex or face (ie. e.Org or e.Lface).
// This makes it easier to process all vertices or faces in the global lists
// without worrying about processing the same data twice.  As a convenience,
// when a face is split, the "inside" flag is copied from the old face.
// Other internal data (v.data, v.activeRegion, f.data, f.marked,
// f.trail, e.winding) is set to zero.
//
// ********************** Basic Edge Operations **************************
//
// tessMeshMakeEdge( mesh ) creates one edge, two vertices, and a loop.
// The loop (face) consists of the two new half-edges.
//
// tessMeshSplice( eOrg, eDst ) is the basic operation for changing the
// mesh connectivity and topology.  It changes the mesh so that
//  eOrg.Onext <- OLD( eDst.Onext )
//  eDst.Onext <- OLD( eOrg.Onext )
// where OLD(...) means the value before the meshSplice operation.
//
// This can have two effects on the vertex structure:
//  - if eOrg.Org != eDst.Org, the two vertices are merged together
//  - if eOrg.Org == eDst.Org, the origin is split into two vertices
// In both cases, eDst.Org is changed and eOrg.Org is untouched.
//
// Similarly (and independently) for the face structure,
//  - if eOrg.Lface == eDst.Lface, one loop is split into two
//  - if eOrg.Lface != eDst.Lface, two distinct loops are joined into one
// In both cases, eDst.Lface is changed and eOrg.Lface is unaffected.
//
// tessMeshDelete( eDel ) removes the edge eDel.  There are several cases:
// if (eDel.Lface != eDel.Rface), we join two loops into one; the loop
// eDel.Lface is deleted.  Otherwise, we are splitting one loop into two;
// the newly created loop will contain eDel.Dst.  If the deletion of eDel
// would create isolated vertices, those are deleted as well.
//
// ********************** Other Edge Operations **************************
//
// tessMeshAddEdgeVertex( eOrg ) creates a new edge eNew such that
// eNew == eOrg.Lnext, and eNew.Dst is a newly created vertex.
// eOrg and eNew will have the same left face.
//
// tessMeshSplitEdge( eOrg ) splits eOrg into two edges eOrg and eNew,
// such that eNew == eOrg.Lnext.  The new vertex is eOrg.Dst == eNew.Org.
// eOrg and eNew will have the same left face.
//
// tessMeshConnect( eOrg, eDst ) creates a new edge from eOrg.Dst
// to eDst.Org, and returns the corresponding half-edge eNew.
// If eOrg.Lface == eDst.Lface, this splits one loop into two,
// and the newly created loop is eNew.Lface.  Otherwise, two disjoint
// loops are merged into one, and the loop eDst.Lface is destroyed.
//
// ************************ Other Operations *****************************
//
// tessMeshNewMesh() creates a new mesh with no edges, no vertices,
// and no loops (what we usually call a "face").
//
// tessMeshUnion( mesh1, mesh2 ) forms the union of all structures in
// both meshes, and returns the new mesh (the old meshes are destroyed).
//
// tessMeshDeleteMesh( mesh ) will free all storage for any valid mesh.
//
// tessMeshZapFace( fZap ) destroys a face and removes it from the
// global face list.  All edges of fZap will have a nil pointer as their
// left face.  Any edges which also have a nil pointer as their right face
// are deleted entirely (along with any isolated vertices this produces).
// An entire mesh can be deleted by zapping its faces, one at a time,
// in any order.  Zapped faces cannot be used in further mesh operations!
//
// tessMeshCheckMesh( mesh ) checks a mesh for self-consistency.

var vertexIDCounter int

type vertex struct {
	next   *vertex   // next vertex (never nil)
	prev   *vertex   // previous vertex (never nil)
	anEdge *halfEdge // a half-edge with this origin
	id     int       // unique vertex identifier

	// Internal data (keep hidden)

	coords   [3]float // vertex location in 3D
	s, t     float    // projection onto the sweep plane
	pqHandle *vertex  // to allow deletion from priority queue
	n        index    // to allow identify unique vertices
	idx      index    // to allow map result to original verts
}

type face struct {
	next   *face     // next face (never nil)
	prev   *face     // previous face (never nil)
	anEdge *halfEdge // a half edge with this left face

	// Internal data (keep hidden)

	trail  *face // "stack" for conversion to strips
	n      index // to allow identiy unique faces
	marked bool  // flag for conversion to strips
	inside bool  // this face is in the polygon interior
}

type halfEdge struct {
	next  *halfEdge // doubly-linked list (prev==Sym.next)
	Sym   *halfEdge // same edge, opposite direction
	Onext *halfEdge // next edge CCW around origin
	Lnext *halfEdge // next edge CCW around left face
	Org   *vertex   // origin vertex (Overtex too long)
	Lface *face     // left face

	// Internal data (keep hidden)

	activeRegion *activeRegion // a region with this upper edge

	// change in winding number when crossing
	// from the right face to the left face
	winding int
}

// The mesh structure is similar in spirit, notation, and operations
// to the "quad-edge" structure (see L. Guibas and J. Stolfi, Primitives
// for the manipulation of general subdivisions and the computation of
// Voronoi diagrams, ACM Transactions on Graphics, 4(2):74-123, April 1985).
// For a simplified description, see the course notes for CS348a,
// "Mathematical Foundations of Computer Graphics", available at the
// Stanford bookstore (and taught during the fall quarter).
// The implementation also borrows a tiny subset of the graph-based approach
// use in Mantyla's Geometric Work Bench (see M. Mantyla, An Introduction
// to Sold Modeling, Computer Science Press, Rockville, Maryland, 1988).
//
// The fundamental data structure is the "half-edge".  Two half-edges
// go together to make an edge, but they point in opposite directions.
// Each half-edge has a pointer to its mate (the "symmetric" half-edge Sym),
// its origin vertex (Org), the face on its left side (Lface), and the
// adjacent half-edges in the CCW direction around the origin vertex
// (Onext) and around the left face (Lnext).  There is also a "next"
// pointer for the global edge list (see below).
//
// The notation used for mesh navigation:
//
//	Sym   = the mate of a half-edge (same edge, but opposite direction)
//	Onext = edge CCW around origin vertex (keep same origin)
//	Dnext = edge CCW around destination vertex (keep same dest)
//	Lnext = edge CCW around left face (dest becomes new origin)
//	Rnext = edge CCW around right face (origin becomes new dest)
//
// "prev" means to substitute CW for CCW in the definitions above.
//
// The mesh keeps global lists of all vertices, faces, and edges,
// stored as doubly-linked circular lists with a dummy header node.
// The mesh stores pointers to these dummy headers (vHead, fHead, eHead).
//
// The circular edge list is special; since half-edges always occur
// in pairs (e and e.Sym), each half-edge stores a pointer in only
// one direction.  Starting at eHead and following the e.next pointers
// will visit each //edge// once (ie. e or e.Sym, but not both).
// e.Sym stores a pointer in the opposite direction, thus it is
// always true that e.Sym.next.Sym.next == e.
//
// Each vertex has a pointer to next and previous vertices in the
// circular list, and a pointer to a half-edge with this vertex as
// the origin (nil if this is the dummy header).  There is also a
// field "data" for client data.
//
// Each face has a pointer to the next and previous faces in the
// circular list, and a pointer to a half-edge with this face as
// the left face (nil if this is the dummy header).  There is also
// a field "data" for client data.
//
// Note that what we call a "face" is really a loop; faces may consist
// of more than one loop (ie. not simply connected), but there is no
// record of this in the data structure.  The mesh may consist of
// several disconnected regions, so it may not be possible to visit
// the entire mesh by starting at a half-edge and traversing the edge
// structure.
//
// The mesh does NOT support isolated vertices; a vertex is deleted along
// with its last edge.  Similarly when two faces are merged, one of the
// faces is deleted (see tessMeshDelete below).  For mesh operations,
// all face (loop) and vertex pointers must not be nil.  However, once
// mesh manipulation is finished, TESSmeshZapFace can be used to delete
// faces of the mesh, one at a time.  All external faces can be "zapped"
// before the mesh is returned to the client; then a nil face indicates
// a region which is not part of the output polygon.
type mesh struct {
	vHead    vertex   // dummy header for vertex list
	fHead    face     // dummy header for face list
	eHead    halfEdge // dummy header for edge list
	eHeadSym halfEdge // and its symmetric counterpart
}

// makeEdge creates a new pair of half-edges which form their own loop.
// No vertex or face structures are allocated, but these must be assigned
// before the current edge operation is completed.
func makeEdge(m *mesh, eNext *halfEdge) *halfEdge {
	e := &halfEdge{}
	eSym := &halfEdge{}

	// Make sure eNext points to the first edge of the edge pair
	if eNext.Sym != eNext {
		eNext = eNext.Sym
	}

	// Insert in circular doubly-linked list before eNext.
	// Note that the prev pointer is stored in Sym.next.
	ePrev := eNext.Sym.next
	eSym.next = ePrev
	ePrev.Sym.next = e
	e.next = eNext
	eNext.Sym.next = eSym

	e.Sym = eSym
	e.Onext = e
	e.Lnext = eSym
	e.Org = nil
	e.Lface = nil
	e.winding = 0
	e.activeRegion = nil

	eSym.Sym = e
	eSym.Onext = eSym
	eSym.Lnext = e
	eSym.Org = nil
	eSym.Lface = nil
	eSym.winding = 0
	eSym.activeRegion = nil

	return e
}

// splice is best described by the Guibas/Stolfi paper or the
// CS348a notes (see mesh.h).  Basically it modifies the mesh so that
// a.Onext and b.Onext are exchanged.  This can have various effects
// depending on whether a and b belong to different face or vertex rings.
// For more explanation see tessMeshSplice() below.
func splice(a, b *halfEdge) {
	aOnext := a.Onext
	bOnext := b.Onext

	aOnext.Sym.Lnext = b
	bOnext.Sym.Lnext = a
	a.Onext = bOnext
	b.Onext = aOnext
}

// makeVertex attaches a new vertex and makes it the
// origin of all edges in the vertex loop to which eOrig belongs. "vNext" gives
// a place to insert the new vertex in the global vertex list.  We insert
// the new vertex *before* vNext so that algorithms which walk the vertex
// list will not see the newly created vertices.
func makeVertex(newVertex *vertex, eOrig *halfEdge, vNext *vertex) {
	if newVertex.id == 0 {
		vertexIDCounter++
		newVertex.id = vertexIDCounter
	}
	println("makeVertex: creating vertex ID", newVertex.id, "at address", newVertex)
	// Special tracking for vertices 11 and 13
	if newVertex.id == 11 || newVertex.id == 13 {
		println("WARNING: Creating vertex with ID 11 or 13!")
		println("  Call stack:")
		var buf [4096]byte
		n := runtime.Stack(buf[:], false)
		println(string(buf[:n]))
	}
	if eOrig.Org != nil {
		println("Edge origin before update:", eOrig.Org, "(ID:", eOrig.Org.id, ")")
	} else {
		println("Warning: eOrig.Org is nil")
	}
	if vNext != nil {
		println("vNext:", vNext, "(ID:", vNext.id, ")")
	} else {
		println("Warning: vNext is nil")
	}
	vNew := newVertex

	assert(vNew != nil)

	// insert in circular doubly-linked list before vNext
	vPrev := vNext.prev
	vNew.prev = vPrev
	vPrev.next = vNew
	vNew.next = vNext
	vNext.prev = vNew

	println("makeVertex: setting anEdge for vertex ID", vNew.id, "to edge", eOrig)
	vNew.anEdge = eOrig
	// leave coords, s, t undefined

	// fix other edges on this vertex loop
	e := eOrig
	println("makeVertex: updating edges for vertex ID", vNew.id)
	for {
		prevID := -1
		if e.Org != nil {
			prevID = e.Org.id
		}
		e.Org = vNew
		println("  Updated edge", e, "from vertex", prevID, "to vertex", vNew.id)
		e = e.Onext
		if e == eOrig {
			break
		}
	}
}

// makeFace attaches a new face and makes it the left
// face of all edges in the face loop to which eOrig belongs.  "fNext" gives
// a place to insert the new face in the global face list.  We insert
// the new face *before* fNext so that algorithms which walk the face
// list will not see the newly created faces.
func makeFace(newFace *face, eOrig *halfEdge, fNext *face) {
	if newFace == nil {
		println("WARNING: makeFace called with nil newFace")
		return
	}
	if eOrig == nil {
		println("WARNING: makeFace called with nil eOrig")
		return
	}
	if fNext == nil {
		println("WARNING: makeFace called with nil fNext")
		return
	}

	fNew := newFace

	// insert in circular doubly-linked list before fNext
	fPrev := fNext.prev
	if fPrev == nil {
		println("WARNING: fNext.prev is nil in makeFace")
		return
	}

	fNew.prev = fPrev
	fPrev.next = fNew
	fNew.next = fNext
	fNext.prev = fNew

	fNew.anEdge = eOrig
	fNew.trail = nil
	fNew.marked = false

	// The new face is marked "inside" if the old one was.  This is a
	// convenience for the common case where a face has been split in two.
	fNew.inside = fNext.inside

	// fix other edges on this face loop
	e := eOrig
	for {
		e.Lface = fNew
		e = e.Lnext
		if e == nil || e == eOrig {
			break
		}
	}
}

// killEdge destroys an edge (the half-edges eDel and eDel.Sym),
// and removes from the global edge list.
func killEdge(m *mesh, eDel *halfEdge) {
	// Half-edges are allocated in pairs, see EdgePair above
	if eDel.Sym != eDel {
		eDel = eDel.Sym
	}

	// delete from circular doubly-linked list
	eNext := eDel.next
	ePrev := eDel.Sym.next
	eNext.Sym.next = ePrev
	ePrev.Sym.next = eNext
}

// killVertex destroys a vertex and removes it from the global
// vertex list.  It updates the vertex loop to point to a given new vertex.
func killVertex(m *mesh, vDel *vertex, newOrg *vertex) {
	if vDel == nil {
		return // Can't delete nil vertex
	}
	if newOrg == nil {
		return // Need a valid newOrg to reassign edges
	}
	eStart := vDel.anEdge

	// change the origin of all affected edges
	if eStart != nil {
		e := eStart
		for {
			e.Org = newOrg
			e = e.Onext
			if e == eStart || e == nil {
				break
			}
		}
	}

	// delete from circular doubly-linked list
	vPrev := vDel.prev
	vNext := vDel.next
	println("killVertex: deleting vertex ID", vDel.id)
	println("  vPrev ID:", func() int {
		if vPrev != nil {
			return vPrev.id
		} else {
			return -1
		}
	}())
	println("  vNext ID:", func() int {
		if vNext != nil {
			return vNext.id
		} else {
			return -1
		}
	}())
	if vPrev != nil && vNext != nil {
		if vPrev != vNext {
			println("  Updating both vPrev and vNext")
			vPrev.next = vNext
			vNext.prev = vPrev
		} else {
			// vPrev and vNext are the same, which means this was the only vertex
			println("  Only vertex in list, clearing pointers")
			vPrev.next = nil
			vPrev.prev = nil
		}
	} else if vPrev != nil {
		println("  Only updating vPrev")
		vPrev.next = nil
	} else if vNext != nil {
		println("  Only updating vNext")
		vNext.prev = nil
	}

	// Update the anEdge pointer of the newOrg vertex
	newOrg.anEdge = eStart
	println("  Updated newOrg(ID:", newOrg.id, ") anEdge to", eStart)

	// Clear vDel's pointers to help with garbage collection
	vDel.anEdge = nil
	vDel.next = nil
	vDel.prev = nil
}

// killFace destroys a face and removes it from the global face
// list.  It updates the face loop to point to a given new face.
func killFace(m *mesh, fDel *face, newLface *face) {
	eStart := fDel.anEdge

	// change the left face of all affected edges
	e := eStart
	for {
		e.Lface = newLface
		e = e.Lnext
		if e == eStart {
			break
		}
	}

	// delete from circular doubly-linked list
	fPrev := fDel.prev
	fNext := fDel.next
	fNext.prev = fPrev
	fPrev.next = fNext
}

// tessMeshMakeEdge creates one edge, two vertices, and a loop (face).
// The loop consists of the two new half-edges.
func tessMeshMakeEdge(m *mesh) *halfEdge {
	newVertex1 := &vertex{}
	newVertex2 := &vertex{}
	newFace := &face{}

	e := makeEdge(m, &m.eHead)

	// Set edge origins before calling makeVertex
	e.Org = newVertex1
	e.Sym.Org = newVertex2

	makeVertex(newVertex1, e, &m.vHead)
	makeVertex(newVertex2, e.Sym, &m.vHead)
	makeFace(newFace, e, &m.fHead)
	return e
}

// tessMeshSplice is the basic operation for changing the
// mesh connectivity and topology.  It changes the mesh so that
//
//	eOrg.Onext <- OLD( eDst.Onext )
//	eDst.Onext <- OLD( eOrg.Onext )
//
// where OLD(...) means the value before the meshSplice operation.
//
// This can have two effects on the vertex structure:
//   - if eOrg.Org != eDst.Org, the two vertices are merged together
//   - if eOrg.Org == eDst.Org, the origin is split into two vertices
//
// In both cases, eDst.Org is changed and eOrg.Org is untouched.
//
// Similarly (and independently) for the face structure,
//   - if eOrg.Lface == eDst.Lface, one loop is split into two
//   - if eOrg.Lface != eDst.Lface, two distinct loops are joined into one
//
// In both cases, eDst.Lface is changed and eOrg.Lface is unaffected.
//
// Some special cases:
// If eDst == eOrg, the operation has no effect.
// If eDst == eOrg.Lnext, the new face will have a single edge.
// If eDst == eOrg.Lprev, the old face will have a single edge.
// If eDst == eOrg.Onext, the new vertex will have a single edge.
// If eDst == eOrg.Oprev, the old vertex will have a single edge.
func tessMeshSplice(m *mesh, eOrg *halfEdge, eDst *halfEdge) {
	joiningLoops := false
	joiningVertices := false

	if eOrg == eDst {
		return
	}

	if eDst.Org != eOrg.Org {
		// We are merging two disjoint vertices -- destroy eDst.Org
		joiningVertices = true
		killVertex(m, eDst.Org, eOrg.Org)
	}
	if eDst.Lface != eOrg.Lface {
		// We are connecting two disjoint loops -- destroy eDst.Lface
		joiningLoops = true
		killFace(m, eDst.Lface, eOrg.Lface)
	}

	// Change the edge structure
	splice(eDst, eOrg)

	if !joiningVertices {
		newVertex := &vertex{}

		// We split one vertex into two -- the new vertex is eDst.Org.
		// Make sure the old vertex points to a valid half-edge.
		makeVertex(newVertex, eDst, eOrg.Org)
		eOrg.Org.anEdge = eOrg
	}
	if !joiningLoops {
		newFace := &face{}

		// We split one loop into two -- the new loop is eDst.Lface.
		// Make sure the old face points to a valid half-edge.
		if eOrg.Lface != nil {
			makeFace(newFace, eDst, eOrg.Lface)
			eOrg.Lface.anEdge = eOrg
		} else {
			println("WARNING: eOrg.Lface is nil in tessMeshSplice, skipping makeFace")
			// Add the new face to the mesh's face list directly
			newFace.anEdge = eDst
			newFace.trail = nil
			newFace.marked = false
			newFace.inside = false

			// Insert the new face at the beginning of the face list
			if m.fHead.next != nil {
				newFace.next = m.fHead.next
				newFace.prev = &m.fHead
				m.fHead.next.prev = newFace
				m.fHead.next = newFace
			} else {
				m.fHead.next = newFace
				m.fHead.prev = newFace
				newFace.next = &m.fHead
				newFace.prev = &m.fHead
			}
		}
	}
}

// tessMeshDelete removes the edge eDel.  There are several cases:
// if (eDel.Lface != eDel.Rface), we join two loops into one; the loop
// eDel.Lface is deleted.  Otherwise, we are splitting one loop into two;
// the newly created loop will contain eDel.Dst.  If the deletion of eDel
// would create isolated vertices, those are deleted as well.
//
// This function could be implemented as two calls to tessMeshSplice
// plus a few calls to memFree, but this would allocate and delete
// unnecessary vertices and faces.
func tessMeshDelete(m *mesh, eDel *halfEdge) {
	eDelSym := eDel.Sym
	joiningLoops := false

	// First step: disconnect the origin vertex eDel.Org.  We make all
	// changes to get a consistent mesh in this "intermediate" state.
	if eDel.Lface != eDel.rFace() {
		// We are joining two loops into one -- remove the left face
		joiningLoops = true
		killFace(m, eDel.Lface, eDel.rFace())
	}

	if eDel.Onext == eDel {
		killVertex(m, eDel.Org, nil)
	} else {
		// Make sure that eDel.Org and eDel.Rface point to valid half-edges
		eDel.rFace().anEdge = eDel.oPrev()
		eDel.Org.anEdge = eDel.Onext

		splice(eDel, eDel.oPrev())
		if !joiningLoops {
			newFace := &face{}

			// We are splitting one loop into two -- create a new loop for eDel.
			makeFace(newFace, eDel, eDel.Lface)
		}
	}

	// Claim: the mesh is now in a consistent state, except that eDel.Org
	// may have been deleted.  Now we disconnect eDel.Dst.
	if eDelSym.Onext == eDelSym {
		killVertex(m, eDelSym.Org, nil)
		killFace(m, eDelSym.Lface, nil)
	} else {
		// Make sure that eDel.Dst and eDel.Lface point to valid half-edges
		eDel.Lface.anEdge = eDelSym.oPrev()
		eDelSym.Org.anEdge = eDelSym.Onext
		splice(eDelSym, eDelSym.oPrev())
	}

	// Any isolated vertices or faces have already been freed.
	killEdge(m, eDel)
}

// All these routines can be implemented with the basic edge
// operations above.  They are provided for convenience and efficiency.

// tessMeshAddEdgeVertex creates a new edge eNew such that
// eNew == eOrg.Lnext, and eNew.Dst is a newly created vertex.
// eOrg and eNew will have the same left face.
func tessMeshAddEdgeVertex(m *mesh, eOrg *halfEdge) *halfEdge {
	eNew := makeEdge(m, eOrg)
	eNewSym := eNew.Sym

	// Connect the new edge appropriately
	splice(eNew, eOrg.Lnext)

	// Set the vertex and face information
	eNew.Org = eOrg.dst()
	vertexIDCounter++
	newVertex := &vertex{id: vertexIDCounter}
	println("Creating new vertex in tessMeshAddEdgeVertex with ID:", newVertex.id)
	// Fix: Set eNewSym.Org to the new vertex
	eNewSym.Org = newVertex
	makeVertex(newVertex, eNewSym, eNew.Org)
	eNew.Lface = eOrg.Lface
	eNewSym.Lface = eOrg.Lface

	return eNew
}

// tessMeshSplitEdge splits eOrg into two edges eOrg and eNew,
// such that eNew == eOrg.Lnext.  The new vertex is eOrg.Dst == eNew.Org.
// eOrg and eNew will have the same left face.
func tessMeshSplitEdge(m *mesh, eOrg *halfEdge) *halfEdge {
	tempHalfEdge := tessMeshAddEdgeVertex(m, eOrg)

	eNew := tempHalfEdge.Sym

	// Disconnect eOrg from eOrg.Dst and connect it to eNew.Org
	splice(eOrg.Sym, eOrg.Sym.oPrev())
	splice(eOrg.Sym, eNew)

	// Set the vertex and face information
	eOrg.setDst(eNew.Org)
	eNew.dst().anEdge = eNew.Sym // may have pointed to eOrg.Sym
	eNew.setRFace(eOrg.rFace())
	eNew.winding = eOrg.winding // copy old winding information
	eNew.Sym.winding = eOrg.Sym.winding

	return eNew
}

// tessMeshConnect creates a new edge from eOrg.Dst
// to eDst.Org, and returns the corresponding half-edge eNew.
// If eOrg.Lface == eDst.Lface, this splits one loop into two,
// and the newly created loop is eNew.Lface.  Otherwise, two disjoint
// loops are merged into one, and the loop eDst.Lface is destroyed.
//
// If (eOrg == eDst), the new face will have only two edges.
// If (eOrg.Lnext == eDst), the old face is reduced to a single edge.
// If (eOrg.Lnext.Lnext == eDst), the old face is reduced to two edges.
func tessMeshConnect(m *mesh, eOrg *halfEdge, eDst *halfEdge) *halfEdge {
	joiningLoops := false
	eNew := makeEdge(m, eOrg)
	eNewSym := eNew.Sym

	if eDst.Lface != eOrg.Lface {
		// We are connecting two disjoint loops -- destroy eDst.Lface
		joiningLoops = true
		killFace(m, eDst.Lface, eOrg.Lface)
	}

	// Connect the new edge appropriately
	splice(eNew, eOrg.Lnext)
	splice(eNewSym, eDst)

	// Set the vertex and face information
	eNew.Org = eOrg.dst()
	eNewSym.Org = eDst.Org
	eNew.Lface = eOrg.Lface
	eNewSym.Lface = eOrg.Lface

	// Make sure the old face points to a valid half-edge
	eOrg.Lface.anEdge = eNewSym

	if !joiningLoops {
		newFace := &face{}

		// We split one loop into two -- the new loop is eNew.Lface
		makeFace(newFace, eNew, eOrg.Lface)
	}
	return eNew
}

// tessMeshZapFace destroys a face and removes it from the
// global face list.  All edges of fZap will have a nil pointer as their
// left face.  Any edges which also have a nil pointer as their right face
// are deleted entirely (along with any isolated vertices this produces).
// An entire mesh can be deleted by zapping its faces, one at a time,
// in any order.  Zapped faces cannot be used in further mesh operations!
func tessMeshZapFace(m *mesh, fZap *face) {
	eStart := fZap.anEdge

	// walk around face, deleting edges whose right face is also nil
	eNext := eStart.Lnext
	for {
		e := eNext
		eNext = e.Lnext

		e.Lface = nil
		if e.rFace() == nil {
			// delete the edge -- see TESSmeshDelete above

			if e.Onext == e {
				killVertex(m, e.Org, nil)
			} else {
				// Make sure that e.Org points to a valid half-edge
				e.Org.anEdge = e.Onext
				splice(e, e.oPrev())
			}
			eSym := e.Sym
			if eSym.Onext == eSym {
				killVertex(m, eSym.Org, nil)
			} else {
				// Make sure that eSym.Org points to a valid half-edge
				eSym.Org.anEdge = eSym.Onext
				splice(eSym, eSym.oPrev())
			}
			killEdge(m, e)
		}
		if e == eStart {
			break
		}
	}

	// delete from circular doubly-linked list
	fPrev := fZap.prev
	fNext := fZap.next
	fNext.prev = fPrev
	fPrev.next = fNext
}

// tessMeshNewMesh creates a new mesh with no edges, no vertices,
// and no loops (what we usually call a "face").
func tessMeshNewMesh() *mesh {
	mesh := &mesh{}

	v := &mesh.vHead
	f := &mesh.fHead
	e := &mesh.eHead
	eSym := &mesh.eHeadSym

	v.next = v
	v.prev = v
	v.anEdge = nil

	f.next = f
	f.prev = f
	f.anEdge = nil
	f.trail = nil
	f.marked = false
	f.inside = false

	e.next = e
	e.Sym = eSym
	e.Onext = nil
	e.Lnext = nil
	e.Org = nil
	e.Lface = nil
	e.winding = 0
	e.activeRegion = nil

	eSym.next = eSym
	eSym.Sym = e
	eSym.Onext = nil
	eSym.Lnext = nil
	eSym.Org = nil
	eSym.Lface = nil
	eSym.winding = 0
	eSym.activeRegion = nil

	return mesh
}

// tessMeshUnion forms the union of all structures in
// both meshes, and returns the new mesh (the old meshes are destroyed).
//
// TODO: This function is not used anywhere?
func tessMeshUnion(mesh1, mesh2 *mesh) *mesh {
	f1 := &mesh1.fHead
	v1 := &mesh1.vHead
	e1 := &mesh1.eHead
	f2 := &mesh2.fHead
	v2 := &mesh2.vHead
	e2 := &mesh2.eHead

	// Add the faces, vertices, and edges of mesh2 to those of mesh1
	if f2.next != f2 {
		f1.prev.next = f2.next
		f2.next.prev = f1.prev
		f2.prev.next = f1
		f1.prev = f2.prev
	}

	if v2.next != v2 {
		v1.prev.next = v2.next
		v2.next.prev = v1.prev
		v2.prev.next = v1
		v1.prev = v2.prev
	}

	if e2.next != e2 {
		e1.Sym.next.Sym.next = e2.next
		e2.next.Sym.next = e1.Sym.next
		e2.Sym.next.Sym.next = e1
		e1.Sym.next = e2.Sym.next
	}
	return mesh1
}

func countFaceVerts(f *face) int {
	eCur := f.anEdge
	n := 0
	for {
		n++
		eCur = eCur.Lnext
		if eCur == f.anEdge {
			break
		}
	}
	return n
}

// EdgeIsInternal returns true if the edge is internal (not on the boundary)
func EdgeIsInternal(e *halfEdge) bool {
	return e.Sym != nil && e.Sym.Lface != nil
}

func tessMeshMergeConvexFaces(m *mesh, maxVertsPerFace int) {
	// Process edges instead of faces to avoid redundant work
	processed := make(map[*halfEdge]bool)

	for e := m.eHead.next; e != &m.eHead; e = e.next {
		// Skip if this edge or its symmetric has already been processed
		if processed[e] || processed[e.Sym] || e.Org == nil {
			continue
		}

		// Only consider internal edges between two inside faces
		if !EdgeIsInternal(e) || !e.Lface.inside || !e.Sym.Lface.inside {
			continue
		}

		f1 := e.Lface
		f2 := e.Sym.Lface

		// Check if merging these faces would exceed the maximum vertex count
		curNv := countFaceVerts(f1)
		symNv := countFaceVerts(f2)
		if curNv+symNv-2 > maxVertsPerFace {
			continue
		}

		// Check if the merged polygon would be convex
		// This is a simplified check - in production code you might want to
		// check all vertices of the merged polygon
		if vertCCW(e.lPrev().Org, e.Org, e.Sym.Lnext.Lnext.Org) && vertCCW(e.Sym.lPrev().Org, e.Sym.Org, e.Lnext.Lnext.Org) {
			// Mark both edges as processed
			processed[e] = true
			processed[e.Sym] = true

			// Merge the faces by deleting the common edge
			tessMeshDelete(m, e)
		}
	}
}

// tessMeshFlipEdge flips an edge, recomputing the face topology.
// Returns true if the edge was flipped, false otherwise.
func tessMeshFlipEdge(m *mesh, e *halfEdge) bool {
	// Check if the edge can be flipped
	if e.Lface == e.Sym.Lface || !EdgeIsInternal(e) {
		return false
	}

	// Get the vertices involved in the flip
	a := e.Org
	b := e.dst()
	c := e.Lnext.Org
	d := e.Sym.Lnext.Org

	// Check if the flip would create a degenerate face
	if a == c || b == d || a == d || b == c {
		return false
	}

	// Check if the quadrilateral is convex
	if !vertCCW(a, b, c) || !vertCCW(b, c, d) || !vertCCW(c, d, a) || !vertCCW(d, a, b) {
		return false
	}

	// Remember the faces
	f1 := e.Lface
	f2 := e.Sym.Lface

	// Disconnect the edge from its current neighbors
	eNext := e.Lnext
	eSymNext := e.Sym.Lnext

	// Create a new edge between c and d
	eNew := makeEdge(m, e)
	eNewSym := eNew.Sym

	// Connect the new edge
	splice(eNew, eNext)
	splice(eNewSym, eSymNext)

	// Set the vertices for the new edge
	eNew.Org = c
	eNewSym.Org = d

	// Set the faces for the new edge
	eNew.Lface = f1
	eNewSym.Lface = f2

	// Update the face pointers
	f1.anEdge = eNew
	f2.anEdge = eNewSym

	// Delete the old edge
	tessMeshDelete(m, e)

	// Check if the new faces are valid
	tessMeshCheckMesh(m)

	return true
}

// tessMeshCheckMesh checks a mesh for self-consistency.
func tessMeshCheckMesh(m *mesh) {
	fHead := &m.fHead
	vHead := &m.vHead
	eHead := &m.eHead

	var f *face
	fPrev := fHead
	for {
		f = fPrev.next
		if f == fHead {
			break
		}
		assert(f.prev == fPrev)
		e := f.anEdge
		for {
			assert(e.Sym != e)
			assert(e.Sym.Sym == e)
			assert(e.Lnext.Onext.Sym == e)
			assert(e.Onext.Sym.Lnext == e)
			assert(e.Lface == f)
			e = e.Lnext
			if e == f.anEdge {
				break
			}
		}
		fPrev = f
	}
	assert(f.prev == fPrev && f.anEdge == nil)

	var v *vertex
	vPrev := vHead
	for {
		v = vPrev.next
		if v == vHead {
			break
		}
		assert(v.prev == vPrev)
		e := v.anEdge
		for {
			assert(e.Sym != e)
			assert(e.Sym.Sym == e)
			assert(e.Lnext.Onext.Sym == e)
			assert(e.Onext.Sym.Lnext == e)
			if e.Org != v {
				println("Assertion failed: e.Org != v")
				println("Vertex address:", v, "(ID:", v.id, ")")
				println("Vertex s,t:", v.s, ",", v.t)
				println("Edge address:", e)
				if e.Org != nil {
					println("Edge Org address:", e.Org, "(ID:", e.Org.id, ")")
					println("Edge Org s,t:", e.Org.s, ",", e.Org.t)
				} else {
					println("Edge Org is nil")
				}
				if e.Sym != nil && e.Sym.Org != nil {
					println("Edge Sym Org s,t:", e.Sym.Org.s, ",", e.Sym.Org.t)
				} else {
					println("Edge Sym or Edge Sym Org is nil")
				}
				println("Vertex loop starting at:", v.anEdge)
				println("Traversing vertex loop to find incorrect edge...")
				// Traverse the vertex loop to find where the error occurs
				startE := v.anEdge
				tempE := startE
				count := 0
				for {
					count++
					if count > 100 { // Prevent infinite loop
						println("Warning: Loop detected, exiting after 100 iterations")
						break
					}
					// Check if tempE.Org is nil before accessing id
					if tempE.Org == nil {
						println("ERROR: tempE.Org is nil in vertex loop")
						println("Current edge:", tempE)
						panic("nil tempE.Org in tessMeshCheckMesh loop")
					} else {
						println("  Edge", tempE, "Org ID:", tempE.Org.id)
					}
					if tempE.Org.id != v.id {
						println("Found incorrect edge in loop!")
					}
					tempE = tempE.Onext
					if tempE == startE {
						break
					}
				}
				// Add nil check before assertion to prevent crash
				if e.Org == nil {
					println("ERROR: e.Org is nil in tessMeshCheckMesh")
					println("Current edge:", e)
					println("Current vertex:", v)
					panic("nil e.Org in tessMeshCheckMesh")
				}
				assert(e.Org == v)
			}
			e = e.Onext
			if e == v.anEdge {
				break
			}
		}
		vPrev = v
	}
	assert(v.prev == vPrev && v.anEdge == nil)

	var e *halfEdge
	ePrev := eHead
	for {
		e = ePrev.next
		if e == eHead {
			break
		}
		assert(e.Sym.next == ePrev.Sym)
		assert(e.Sym != e)
		assert(e.Sym.Sym == e)
		assert(e.Org != nil)
		assert(e.dst() != nil)
		assert(e.Lnext.Onext.Sym == e)
		assert(e.Onext.Sym.Lnext == e)
		ePrev = e
	}
	assert(e.Sym.next == ePrev.Sym)
	assert(e.Sym == &m.eHeadSym)
	assert(e.Sym.Sym == e)
	assert(e.Org == nil && e.dst() == nil)
	assert(e.Lface == nil && e.rFace() == nil)
}
