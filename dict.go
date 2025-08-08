package tesselator

// activeRegion:
// For each pair of adjacent edges crossing the sweep line, there is
// an ActiveRegion to represent the region between them.  The active
// regions are kept in sorted order in a dynamic dictionary.  As the
// sweep line crosses each vertex, we update the affected regions.
type activeRegion struct {
	// upper edge, directed right to left
	eUp *halfEdge

	// dictionary node corresponding to eUp
	nodeUp *dictNode

	// used to determine which regions are
	// inside the polygon
	windingNumber int

	// is this region inside the polygon?
	inside bool

	// marks fake edges at t = +/-infinity
	sentinel bool

	// marks regions where the upper or lower
	// edge has changed, but we haven't checked
	// whether they intersect yet
	dirty bool

	// marks temporary edges introduced when
	// we process a "right vertex" (one without
	// any edges leaving to the right)
	fixUpperEdge bool
}

type dictNode struct {
	key  *activeRegion
	prev *dictNode
	next *dictNode
}

type dict struct {
	head  dictNode
	frame *tesselator
}

func newDict(frame *tesselator) *dict {
	d := &dict{
		frame: frame,
	}
	d.head.next = &d.head
	d.head.prev = &d.head
	return d
}

func (d *dict) insertBefore(n *dictNode, key *activeRegion) *dictNode {
	// Handle nil n by starting at head
	if n == nil {
		n = &d.head
	}
	for {
		n = n.prev
		if n.key == nil || edgeLeq(d.frame, n.key, key) {
			break
		}
	}

	nn := &dictNode{
		key:  key,
		next: n.next,
		prev: n,
	}
	n.next.prev = nn
	n.next = nn

	return nn
}

func dictDelete(n *dictNode) {
	n.next.prev = n.prev
	n.prev.next = n.next
}

// search returns the node with the smallest key greater than or equal
// to the given key.  If there is no such key, returns a node whose
// key is NULL.  Similarly, Succ(Max(d)) has a NULL key, etc.
func (d *dict) search(key *activeRegion) *dictNode {
	n := &d.head
	for {
		n = n.next
		if n.key == nil || edgeLeq(d.frame, key, n.key) {
			break
		}
	}
	return n
}

// dictKey returns the key of a dictionary node.
func dictKey(n *dictNode) *activeRegion {
	return n.key
}

func dictSucc(n *dictNode) *dictNode {
	return n.next
}

func dictPred(n *dictNode) *dictNode {
	return n.prev
}

func (d *dict) min() *dictNode {
	return d.head.next
}

func (d *dict) max() *dictNode {
	return d.head.prev
}

func (d *dict) insert(key *activeRegion) *dictNode {
	return d.insertBefore(&d.head, key)
}
