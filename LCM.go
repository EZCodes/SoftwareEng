package LCA

import "fmt"

type Graph struct {
	edges []*Edge
}
type Vertex struct {
	data int
}
type Edge struct {
	dest *Vertex
	src  *Vertex
}
type Ancestor struct {
	data     *Vertex
	distance int
}

// LowestCommonAncestor find the lowest common ancestor of two nodes for Directed Acyclic Graph,
// if found, the node with it is returned, if two nodes dont have a common ancestor
// empty vertex (value 0) is returned
// does not work for graphs with cycles
func (g *Graph) LowestCommonAncestor(v1, v2 *Vertex) (Vertex, error) {
	var an1 []Ancestor
	var an2 []Ancestor
	if v1 == nil || v2 == nil {
		return Vertex{}, fmt.Errorf("One of the vertices is nil. v1: %v, v2: %v", v1, v2)
	}
	an1 = v1.findAncestors(0, v1, []Ancestor{}, g.edges)
	an1 = an1[1:] // last element is itself
	an2 = v2.findAncestors(0, v2, []Ancestor{}, g.edges)
	an2 = an2[1:]        // last element is itsef, we dont count node itself as own ancestor
	an1 = distanceToDepth(an1) // only need to do this with one list since if they have LCA, max depth will be same
	lcm := Ancestor{
		distance: 0,
		data:     &Vertex{data: 0},
	}

	for _, ancestorOne := range an1 {
		for _, ancestorTwo := range an2 {
			if ancestorTwo.data.data == ancestorOne.data.data {
				if ancestorOne.distance > lcm.distance {
					lcm = ancestorOne
				}
			}
		}
	}
	return *(lcm.data), nil
}
func (v *Vertex) findAncestors(d int, currentV *Vertex, ancestors []Ancestor, edges []*Edge) []Ancestor {
	currAncestor := Ancestor{
		data:     currentV,
		distance: d,
	}
	isContained, ancestorIndex := contains(currAncestor, ancestors)
	if !isContained {
		ancestors = append(ancestors, currAncestor)
	} else if ancestors[ancestorIndex].distance < currAncestor.distance {
		ancestors[ancestorIndex].distance = currAncestor.distance
	}
	directAncestors := currentV.findSources(edges)
	if len(directAncestors) == 0 {
		return ancestors
	}
	d++
	for _, vertex := range directAncestors {
		ancestors = v.findAncestors(d, vertex, ancestors, edges)
	}
	return ancestors
}

func (v *Vertex) findSources(edges []*Edge) []*Vertex {
	var sources []*Vertex
	for _, edge := range edges {
		if v.data == edge.dest.data {
			sources = append(sources, edge.src)
		}
	}
	return sources
}

// convert distance from node of interest into depth from the "root"
// basically a reverse
func distanceToDepth(ancestors []Ancestor) []Ancestor {
	currentDepth := 1
	var result []Ancestor
	maxDepth := findMax(ancestors)
	for maxDepth > 0 {
		for _, ancestor := range ancestors {
			if ancestor.distance == maxDepth {
				newAncestor := ancestor			
				newAncestor.distance = currentDepth
				result = append(result, newAncestor)
			}
		}
		currentDepth++
		maxDepth--
	}
	return result
}

func contains(ancestor Ancestor, ancestors []Ancestor) (bool, int) {
	isContained := false
	var containedAncestor int
	for index, other := range ancestors {
		if ancestor.data.data == other.data.data {
			isContained = true
			containedAncestor = index
		}
	}
	return isContained, containedAncestor
}

func findMax(ancestors []Ancestor) int {
	maxAncestor := Ancestor{distance: 0}
	for _, ancestor := range ancestors {
		if ancestor.distance > maxAncestor.distance {
			maxAncestor = ancestor
		}
	}
	return maxAncestor.distance
}
