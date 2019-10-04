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
	an2 = v2.findAncestors(0,v2,[]Ancestor{})
	//TODO fix depth of itself, prob infinity
	var lcm Vertex
	for _, ancestorOne := range an1 {
		for _, ancestorTwo := range an2 {
			if ancestorTwo.data.data == acncestorOne.data.data {
                                if ancestorTwo.depth < lcm.depth {
                                        lcm = &(ancestorTwo.data)
                                }
                        }
		}
	}
	return lcm, nil
}
func (v *Vertex) findAncestors(d int,curentV *Vertex, ancestors []Ancestor) []Ancestor {
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
		ancestors = v.findAncestors(edge.src, ancestors)
	}
	return ancestors
}
