package LCA

import "fmt"

type Graph struct {
	vertices []*Vertex
	edges    []*Edge
}
type Vertex struct {
	data   int
	outEdg []*Edge
	incEdg []*Edge
}
type Edge struct {
	dest *Vertex
	src  *Vertex
}
type Ancestor struct {
	ancOf *Vertex
	data *Vertex
	distance int
}

// LowestCommonAncestor find the lowest common ancestor of two nodes
// if found, the node with it is returned, if two nodes dont have a common ancestor
// empty node is returned
func (g *Graph) LowestCommonAncestor(v1, v2 *Vertex) (Vertex, error) {
	var an1 []Ancestor
	var an2 []Ancestor
	if v1 == nil || v2 == nil {
		return Vertex{}, fmt.Errorf("One of the vertices is nil. v1: %v, v2: %v", v1, v2)
	}
	an1 = v1.findAncestors(0,v1,[]Ancestor{})
	an1 = an1[1:] // last element is itself
	an2 = v2.findAncestors(0,v2,[]Ancestor{})
	an2 = an2[1:] // last element is itsef, we dont count node itself as own ancestor
	distanceToDepth(an1) // only need to do this with one list since if they have LCA, max depth will be same
	lcm :=  Ancestor{distance: 0} // 32bits should be enough for this scale
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
func (v *Vertex) findAncestors(d int, currentV *Vertex, ancestors []Ancestor) []Ancestor {
	ancestors = append(ancestors, Ancestor{
		ancOf: v,
		data : currentV,
		distance: d,
	})	
	if len(currentV.incEdg) == 0 {
		return ancestors
	}
	d++
	for _, edge := range currentV.incEdg {
		ancestors = v.findAncestors(d, edge.src, ancestors)
	}
	return ancestors
}

// convert distance from node of interest into depth from the "root"
// basically a reverse
func distanceToDepth (ancestors []Ancestor) []Ancestor {
	currentDepth := 1
	maxDepth := findMax(ancestors)
	for ;maxDepth > 0; {
		for _, ancestor := range ancestors {
			if ancestor.distance == maxDepth {
				ancestor.distance = currentDepth
			}
		}
		currentDepth++
		maxDepth--
	}
	return ancestors
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
